// Copyright 2019 PingCAP, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// See the License for the specific language governing permissions and
// limitations under the License.

package v1

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	"github.com/pingcap/tidb-operator/pkg/apis/pingcap/v1alpha1"
)

type CommonModel struct {
	AcrossK8s     bool   // same as tc.spec.acrossK8s
	ClusterDomain string // same as tc.spec.clusterDomain
}

func (c CommonModel) FormatClusterDomain() string {
	if len(c.ClusterDomain) > 0 {
		return "." + c.ClusterDomain
	}
	return ""
}

// TODO(aylei): it is hard to maintain script in go literal, we should figure out a better solution
// tidbStartScriptTpl is the template string of tidb start script
// Note: changing this will cause a rolling-update of tidb-servers
var tidbStartScriptTpl = template.Must(template.New("tidb-start-script").Parse(`#!/bin/sh

# This script is used to start tidb containers in kubernetes cluster

# Use DownwardAPIVolumeFiles to store informations of the cluster:
# https://kubernetes.io/docs/tasks/inject-data-application/downward-api-volume-expose-pod-information/#the-downward-api
#
#   runmode="normal/debug"
#
set -uo pipefail

ANNOTATIONS="/etc/podinfo/annotations"

if [[ ! -f "${ANNOTATIONS}" ]]
then
    echo "${ANNOTATIONS} does't exist, exiting."
    exit 1
fi
source ${ANNOTATIONS} 2>/dev/null
runmode=${runmode:-normal}
if [[ X${runmode} == Xdebug ]]
then
    echo "entering debug mode."
    tail -f /dev/null
fi

# Use HOSTNAME if POD_NAME is unset for backward compatibility.
POD_NAME=${POD_NAME:-$HOSTNAME}{{ if .AcrossK8s }}
pd_url="{{ .Path }}"
encoded_domain_url=$(echo $pd_url | base64 | tr "\n" " " | sed "s/ //g")
discovery_url="${CLUSTER_NAME}-discovery.${NAMESPACE}:10261"
until result=$(wget -qO- -T 3 http://${discovery_url}/verify/${encoded_domain_url} 2>/dev/null | sed 's/http:\/\///g'); do
echo "waiting for the verification of PD endpoints ..."
sleep $((RANDOM % 5))
done

ARGS="--store=tikv \
--advertise-address=${POD_NAME}.${HEADLESS_SERVICE_NAME}.${NAMESPACE}.svc{{ .FormatClusterDomain }} \
--host=0.0.0.0 \
--path=${result} \
{{ else }}
ARGS="--store=tikv \
--advertise-address=${POD_NAME}.${HEADLESS_SERVICE_NAME}.${NAMESPACE}.svc{{ .FormatClusterDomain }} \
--host=0.0.0.0 \
--path={{ .Path }} \{{ end }}
--config=/etc/tidb/tidb.toml
"

if [[ X${BINLOG_ENABLED:-} == Xtrue ]]
then
    ARGS="${ARGS} --enable-binlog=true"
fi

SLOW_LOG_FILE=${SLOW_LOG_FILE:-""}
if [[ ! -z "${SLOW_LOG_FILE}" ]]
then
    ARGS="${ARGS} --log-slow-query=${SLOW_LOG_FILE:-}"
fi

{{- if .EnablePlugin }}
ARGS="${ARGS}  --plugin-dir  {{ .PluginDirectory  }} --plugin-load {{ .PluginList }}  "
{{- end }}

echo "start tidb-server ..."
echo "/tidb-server ${ARGS}"
exec /tidb-server ${ARGS}
`))

type TidbStartScriptModel struct {
	CommonModel

	EnablePlugin    bool
	PluginDirectory string
	PluginList      string
	Path            string
}

// pdStartScriptTpl is the pd start script
// Note: changing this will cause a rolling-update of pd cluster
var pdStartScriptTplText = `#!/bin/sh

# This script is used to start pd containers in kubernetes cluster

# Use DownwardAPIVolumeFiles to store informations of the cluster:
# https://kubernetes.io/docs/tasks/inject-data-application/downward-api-volume-expose-pod-information/#the-downward-api
#
#   runmode="normal/debug"
#

set -uo pipefail

ANNOTATIONS="/etc/podinfo/annotations"

if [[ ! -f "${ANNOTATIONS}" ]]
then
    echo "${ANNOTATIONS} does't exist, exiting."
    exit 1
fi
source ${ANNOTATIONS} 2>/dev/null

runmode=${runmode:-normal}
if [[ X${runmode} == Xdebug ]]
then
    echo "entering debug mode."
    tail -f /dev/null
fi

# Use HOSTNAME if POD_NAME is unset for backward compatibility.
POD_NAME=${POD_NAME:-$HOSTNAME}
# the general form of variable PEER_SERVICE_NAME is: "<clusterName>-pd-peer"
cluster_name=` + "`" + `echo ${PEER_SERVICE_NAME} | sed 's/-pd-peer//'` + "`" +
	`
domain="${POD_NAME}.${PEER_SERVICE_NAME}.${NAMESPACE}.svc{{ .FormatClusterDomain }}"
discovery_url="${cluster_name}-discovery.${NAMESPACE}.svc:10261"
encoded_domain_url=` + "`" + `echo ${domain}:2380 | base64 | tr "\n" " " | sed "s/ //g"` + "`" +
	`
elapseTime=0
period=1
threshold=30
while true; do
sleep ${period}
elapseTime=$(( elapseTime+period ))

if [[ ${elapseTime} -ge ${threshold} ]]
then
echo "waiting for pd cluster ready timeout" >&2
exit 1
fi
{{ if eq .CheckDomainScript ""}}
if nslookup ${domain} 2>/dev/null
then
echo "nslookup domain ${domain}.svc success"
break
else
echo "nslookup domain ${domain} failed" >&2
fi {{- else}}{{.CheckDomainScript}}{{end}}
done

ARGS="--data-dir={{ .DataDir }} \
--name={{- if or .AcrossK8s .ClusterDomain }}${domain}{{- else }}${POD_NAME}{{- end }} \
--peer-urls={{ .Scheme }}://0.0.0.0:2380 \
--advertise-peer-urls={{ .Scheme }}://${domain}:2380 \
--client-urls={{ .Scheme }}://0.0.0.0:2379 \
--advertise-client-urls={{ .Scheme }}://${domain}:2379 \
--config=/etc/pd/pd.toml \
"

if [[ -f {{ .DataDir }}/join ]]
then
# The content of the join file is:
#   demo-pd-0=http://demo-pd-0.demo-pd-peer.demo.svc:2380,demo-pd-1=http://demo-pd-1.demo-pd-peer.demo.svc:2380
# The --join args must be:
#   --join=http://demo-pd-0.demo-pd-peer.demo.svc:2380,http://demo-pd-1.demo-pd-peer.demo.svc:2380
join=` + "`" + `cat {{ .DataDir }}/join | tr "," "\n" | awk -F'=' '{print $2}' | tr "\n" ","` + "`" + `
join=${join%,}
ARGS="${ARGS} --join=${join}"
elif [[ ! -d {{ .DataDir }}/member/wal ]]
then
until result=$(wget -qO- -T 3 http://${discovery_url}/new/${encoded_domain_url} 2>/dev/null); do
echo "waiting for discovery service to return start args ..."
sleep $((RANDOM % 5))
done
ARGS="${ARGS}${result}"
fi

echo "starting pd-server ..."
sleep $((RANDOM % 10))
echo "/pd-server ${ARGS}"
exec /pd-server ${ARGS}
`

func replacePDStartScriptCustomPorts(startScript string) string {
	// `DefaultPDClientPort`/`DefaultPDPeerPort` may be changed when building the binary
	if v1alpha1.DefaultPDClientPort != 2379 {
		startScript = strings.ReplaceAll(startScript, ":2379", fmt.Sprintf(":%d", v1alpha1.DefaultPDClientPort))
	}
	if v1alpha1.DefaultPDPeerPort != 2380 {
		startScript = strings.ReplaceAll(startScript, ":2380", fmt.Sprintf(":%d", v1alpha1.DefaultPDPeerPort))
	}
	return startScript
}

var pdStartScriptTpl = template.Must(template.New("pd-start-script").Parse(replacePDStartScriptCustomPorts(pdStartScriptTplText)))

var checkDNSV1 string = `
digRes=$(dig ${domain} A ${domain} AAAA +search +short)
if [ $? -ne 0  ]; then
  echo "$digRes"
  echo "domain resolve ${domain} failed"
  continue
fi

if [ -z "${digRes}" ]
then
  echo "domain resolve ${domain} no record return"
else
  echo "domain resolve ${domain} success"
  echo "$digRes"
  break
fi
`

type PDStartScriptModel struct {
	CommonModel

	Scheme            string
	DataDir           string
	CheckDomainScript string
}

var tikvStartScriptTplText = `#!/bin/sh

# This script is used to start tikv containers in kubernetes cluster

# Use DownwardAPIVolumeFiles to store informations of the cluster:
# https://kubernetes.io/docs/tasks/inject-data-application/downward-api-volume-expose-pod-information/#the-downward-api
#
#   runmode="normal/debug"
#

set -uo pipefail

ANNOTATIONS="/etc/podinfo/annotations"

if [[ ! -f "${ANNOTATIONS}" ]]
then
    echo "${ANNOTATIONS} does't exist, exiting."
    exit 1
fi
source ${ANNOTATIONS} 2>/dev/null

runmode=${runmode:-normal}
if [[ X${runmode} == Xdebug ]]
then
	echo "entering debug mode."
	tail -f /dev/null
fi

# Use HOSTNAME if POD_NAME is unset for backward compatibility.
POD_NAME=${POD_NAME:-$HOSTNAME}{{ if .AcrossK8s }}
pd_url="{{ .PDAddress }}"
encoded_domain_url=$(echo $pd_url | base64 | tr "\n" " " | sed "s/ //g")
discovery_url="${CLUSTER_NAME}-discovery.${NAMESPACE}:10261"

until result=$(wget -qO- -T 3 http://${discovery_url}/verify/${encoded_domain_url} 2>/dev/null); do
echo "waiting for the verification of PD endpoints ..."
sleep $((RANDOM % 5))
done

ARGS="--pd=${result} \
{{ else }}
ARGS="--pd={{ .PDAddress }} \{{ end }}
--advertise-addr=${POD_NAME}.${HEADLESS_SERVICE_NAME}.${NAMESPACE}.svc{{ .FormatClusterDomain }}:20160 \
--addr={{ .Addr }} \
--status-addr={{ .StatusAddr }} \{{if .EnableAdvertiseStatusAddr }}
--advertise-status-addr={{ .AdvertiseStatusAddr }}:20180 \{{end}}
--data-dir={{ .DataDir }} \
--capacity=${CAPACITY} \
--config=/etc/tikv/tikv.toml
"

if [ ! -z "${STORE_LABELS:-}" ]; then
  LABELS=" --labels ${STORE_LABELS} "
  ARGS="${ARGS}${LABELS}"
fi

echo "starting tikv-server ..."
echo "/tikv-server ${ARGS}"
exec /tikv-server ${ARGS}
`

func replaceTiKVStartScriptCustomPorts(startScript string) string {
	// `DefaultTiKVServerPort`/`DefaultTiKVStatusPort` may be changed when building the binary
	if v1alpha1.DefaultTiKVServerPort != 20160 {
		startScript = strings.ReplaceAll(startScript, ":20160", fmt.Sprintf(":%d", v1alpha1.DefaultTiKVServerPort))
	}
	if v1alpha1.DefaultTiKVStatusPort != 20180 {
		startScript = strings.ReplaceAll(startScript, ":20180", fmt.Sprintf(":%d", v1alpha1.DefaultTiKVStatusPort))
	}
	return startScript
}

var tikvStartScriptTpl = template.Must(template.New("tikv-start-script").Parse(replaceTiKVStartScriptCustomPorts(tikvStartScriptTplText)))

type TiKVStartScriptModel struct {
	CommonModel

	EnableAdvertiseStatusAddr bool
	AdvertiseStatusAddr       string
	DataDir                   string
	PDAddress                 string
	Addr                      string
	StatusAddr                string
}

// pumpStartScriptTpl is the template string of pump start script
// Note: changing this will cause a rolling-update of pump cluster
var pumpStartScriptTplText = `{{ if .AcrossK8s }}
pd_url="{{ .PDAddr }}"
encoded_domain_url=$(echo $pd_url | base64 | tr "\n" " " | sed "s/ //g")
discovery_url="{{ .ClusterName }}-discovery.{{ .Namespace }}:10261"
until result=$(wget -qO- -T 3 http://${discovery_url}/verify/${encoded_domain_url} 2>/dev/null); do
echo "waiting for the verification of PD endpoints ..."
sleep $((RANDOM % 5))
done

pd_url=$result

set -euo pipefail

/pump \
-pd-urls=$pd_url \{{ else }}set -euo pipefail

/pump \
-pd-urls={{ .PDAddr }} \{{ end }}
-L={{ .LogLevel }} \
-advertise-addr=` + "`" + `echo ${HOSTNAME}` + "`" + `.{{ .ClusterName }}-pump{{ .FormatPumpZone }}:8250 \
-config=/etc/pump/pump.toml \
-data-dir=/data \
-log-file=

if [ $? == 0 ]; then
    echo $(date -u +"[%Y/%m/%d %H:%M:%S.%3N %:z]") "pump offline, please delete my pod"
    tail -f /dev/null
fi`

func replacePumpStartScriptCustomPorts(startScript string) string {
	// `DefaultPumpPort` may be changed when building the binary
	if v1alpha1.DefaultPumpPort != 8250 {
		startScript = strings.ReplaceAll(startScript, ":8250", fmt.Sprintf(":%d", v1alpha1.DefaultPumpPort))
	}
	return startScript
}

var pumpStartScriptTpl = template.Must(template.New("pump-start-script").Parse(replacePumpStartScriptCustomPorts(pumpStartScriptTplText)))

type PumpStartScriptModel struct {
	CommonModel

	Scheme      string
	ClusterName string
	PDAddr      string
	LogLevel    string
	Namespace   string
}

func (pssm *PumpStartScriptModel) FormatPumpZone() string {
	if pssm.ClusterDomain != "" {
		return fmt.Sprintf(".%s.svc.%s", pssm.Namespace, pssm.ClusterDomain)
	}
	if pssm.ClusterDomain == "" && pssm.AcrossK8s {
		return fmt.Sprintf(".%s.svc", pssm.Namespace)
	}
	return ""
}

// tidbInitStartScriptTpl is the template string of tidb initializer start script
var tidbInitStartScriptTpl = template.Must(template.New("tidb-init-start-script").Parse(`import os, sys, time, MySQLdb
host = '{{ .ClusterName }}-tidb'
permit_host = '{{ .PermitHost }}'
port = {{ .TiDBServicePort }}
retry_count = 0
for i in range(0, 10):
    try:
{{- if and .TLS .SkipCA }}
        conn = MySQLdb.connect(host=host, port=port, user='root', charset='utf8mb4',connect_timeout=5, ssl={'cert': '{{ .CertPath }}', 'key': '{{ .KeyPath }}'})
{{- else if .TLS }}
        conn = MySQLdb.connect(host=host, port=port, user='root', charset='utf8mb4',connect_timeout=5, ssl={'ca': '{{ .CAPath }}', 'cert': '{{ .CertPath }}', 'key': '{{ .KeyPath }}'})
{{- else }}
        conn = MySQLdb.connect(host=host, port=port, user='root', connect_timeout=5, charset='utf8mb4')
{{- end }}
    except MySQLdb.OperationalError as e:
        print(e)
        retry_count += 1
        time.sleep(1)
        continue
    break
if retry_count == 10:
    sys.exit(1)

{{- if .PasswordSet }}
password_dir = '/etc/tidb/password'
for file in os.listdir(password_dir):
    if file.startswith('.'):
        continue
    user = file
    with open(os.path.join(password_dir, file), 'r') as f:
        lines = f.read().splitlines()
        password = lines[0] if len(lines) > 0 else ""
    if user == 'root':
        conn.cursor().execute("set password for 'root'@'%%' = %s;", (password,))
    else:
        conn.cursor().execute("create user %s@%s identified by %s;", (user, permit_host, password,))
{{- end }}
{{- if .InitSQL }}
with open('/data/init.sql', 'r') as sql:
    for line in sql.readlines():
        conn.cursor().execute(line)
        conn.commit()
{{- end }}
if permit_host != '%%':
    conn.cursor().execute("update mysql.user set Host=%s where User='root';", (permit_host,))
conn.cursor().execute("flush privileges;")
conn.commit()
conn.close()
`))

type TiDBInitStartScriptModel struct {
	ClusterName     string
	PermitHost      string
	PasswordSet     bool
	InitSQL         bool
	TLS             bool
	SkipCA          bool
	CAPath          string
	CertPath        string
	KeyPath         string
	TiDBServicePort int32
}

func RenderTiDBInitStartScript(model *TiDBInitStartScriptModel) (string, error) {
	return renderTemplateFunc(tidbInitStartScriptTpl, model)
}

// tidbInitInitStartScriptTpl is the template string of tidb initializer init container start script
var tidbInitInitStartScriptTpl = template.Must(template.New("tidb-init-init-start-script").Parse(`trap exit TERM
host={{ .ClusterName }}-tidb
port={{ .TiDBServicePort }}
while true; do
  nc -zv -w 3 $host $port
  if [ $? -eq 0 ]; then
	break
  else
	echo "info: failed to connect to $host:$port, sleep 1 second then retry"
	sleep 1
  fi
done
echo "info: successfully connected to $host:$port, able to initialize TiDB now"
`))

type TiDBInitInitStartScriptModel struct {
	ClusterName     string
	TiDBServicePort int32
}

func RenderTiDBInitInitStartScript(model *TiDBInitInitStartScriptModel) (string, error) {
	return renderTemplateFunc(tidbInitInitStartScriptTpl, model)
}

func renderTemplateFunc(tpl *template.Template, model interface{}) (string, error) {
	buff := new(bytes.Buffer)
	err := tpl.Execute(buff, model)
	if err != nil {
		return "", err
	}
	return buff.String(), nil
}

// dmMasterStartScriptTpl is the dm-master start script
// Note: changing this will cause a rolling-update of dm-master cluster
var dmMasterStartScriptTpl = template.Must(template.New("dm-master-start-script").Parse(`#!/bin/sh

# This script is used to start dm-master containers in kubernetes cluster

# Use DownwardAPIVolumeFiles to store informations of the cluster:
# https://kubernetes.io/docs/tasks/inject-data-application/downward-api-volume-expose-pod-information/#the-downward-api
#
#   runmode="normal/debug"
#

set -uo pipefail

ANNOTATIONS="/etc/podinfo/annotations"

if [[ ! -f "${ANNOTATIONS}" ]]
then
    echo "${ANNOTATIONS} does't exist, exiting."
    exit 1
fi
source ${ANNOTATIONS} 2>/dev/null

runmode=${runmode:-normal}
if [[ X${runmode} == Xdebug ]]
then
    echo "entering debug mode."
    tail -f /dev/null
fi

# Use HOSTNAME if POD_NAME is unset for backward compatibility.
POD_NAME=${POD_NAME:-$HOSTNAME}
# the general form of variable PEER_SERVICE_NAME is: "<clusterName>-dm-master-peer"
cluster_name=` + "`" + `echo ${PEER_SERVICE_NAME} | sed 's/-dm-master-peer//'` + "`" +
	`
domain="${POD_NAME}.${PEER_SERVICE_NAME}"
discovery_url="${cluster_name}-dm-discovery.${NAMESPACE}:10261"
encoded_domain_url=` + "`" + `echo ${domain}:8291 | base64 | tr "\n" " " | sed "s/ //g"` + "`" +
	`
elapseTime=0
period=1
threshold=30
while true; do
sleep ${period}
elapseTime=$(( elapseTime+period ))

if [[ ${elapseTime} -ge ${threshold} ]]
then
echo "waiting for dm-master cluster ready timeout" >&2
exit 1
fi
{{ if eq .CheckDomainScript ""}}
if nslookup ${domain} 2>/dev/null
then
echo "nslookup domain ${domain} success"
break
else
echo "nslookup domain ${domain} failed" >&2
fi {{- else}}{{.CheckDomainScript}}{{end}}
done

ARGS="--data-dir={{ .DataDir }} \
--name=${POD_NAME} \
--peer-urls={{ .Scheme }}://0.0.0.0:8291 \
--advertise-peer-urls={{ .Scheme }}://${domain}:8291 \
--master-addr=:8261 \
--advertise-addr=${domain}:8261 \
--config=/etc/dm-master/dm-master.toml \
"

if [[ -f {{ .DataDir }}/join ]]
then
# The content of the join file is:
#   demo-dm-master-0=http://demo-dm-master-0.demo-dm-master-peer.demo.svc:8291,demo-dm-master-1=http://demo-dm-master-1.demo-dm-master-peer.demo.svc:8291
# The --join args must be:
#   --join=http://demo-dm-master-0.demo-dm-master-peer.demo.svc:8261,http://demo-dm-master-1.demo-dm-master-peer.demo.svc:8261
join=` + "`" + `cat {{ .DataDir }}/join | sed -e 's/8291/8261/g' | tr "," "\n" | awk -F'=' '{print $2}' | tr "\n" ","` + "`" + `
join=${join%,}
ARGS="${ARGS} --join=${join}"
elif [[ ! -d {{ .DataDir }}/member/wal ]]
then
until result=$(wget -qO- -T 3 ${discovery_url}/new/${encoded_domain_url}/dm 2>/dev/null); do
echo "waiting for discovery service to return start args ..."
sleep $((RANDOM % 5))
done
ARGS="${ARGS}${result}"
fi

echo "starting dm-master ..."
sleep $((RANDOM % 10))
echo "/dm-master ${ARGS}"
exec /dm-master ${ARGS}
`))

// TODO: refactor to confine the checking script within the package
var DMMasterCheckDNSV1 string = `
digRes=$(dig ${domain} A ${domain} AAAA +search +short 2>/dev/null)
if [ -z "${digRes}" ]
then
echo "dig domain ${domain} failed" >&2
else
echo "dig domain ${domain} success"
break
fi`

type DMMasterStartScriptModel struct {
	Scheme            string
	DataDir           string
	CheckDomainScript string
}

func RenderDMMasterStartScript(model *DMMasterStartScriptModel) (string, error) {
	return renderTemplateFunc(dmMasterStartScriptTpl, model)
}

// dmWorkerStartScriptTpl is the dm-worker start script
// Note: changing this will cause a rolling-update of dm-worker cluster
var dmWorkerStartScriptTpl = template.Must(template.New("dm-worker-start-script").Parse(`#!/bin/sh

# This script is used to start dm-worker containers in kubernetes cluster

# Use DownwardAPIVolumeFiles to store informations of the cluster:
# https://kubernetes.io/docs/tasks/inject-data-application/downward-api-volume-expose-pod-information/#the-downward-api
#
#   runmode="normal/debug"
#

set -uo pipefail

ANNOTATIONS="/etc/podinfo/annotations"

if [[ ! -f "${ANNOTATIONS}" ]]
then
    echo "${ANNOTATIONS} does't exist, exiting."
    exit 1
fi
source ${ANNOTATIONS} 2>/dev/null

runmode=${runmode:-normal}
if [[ X${runmode} == Xdebug ]]
then
    echo "entering debug mode."
    tail -f /dev/null
fi


# Use HOSTNAME if POD_NAME is unset for backward compatibility.
POD_NAME=${POD_NAME:-$HOSTNAME}
# TODO: dm-worker will support data-dir in the future
ARGS="--name=${POD_NAME} \
--join={{ .MasterAddress }} \
--advertise-addr=${POD_NAME}.${HEADLESS_SERVICE_NAME}:8262 \
--worker-addr=0.0.0.0:8262 \
--config=/etc/dm-worker/dm-worker.toml
"

if [ ! -z "${STORE_LABELS:-}" ]; then
  LABELS=" --labels ${STORE_LABELS} "
  ARGS="${ARGS}${LABELS}"
fi

echo "starting dm-worker ..."
echo "/dm-worker ${ARGS}"
exec /dm-worker ${ARGS}
`))

type DMWorkerStartScriptModel struct {
	DataDir       string
	MasterAddress string
}

func RenderDMWorkerStartScript(model *DMWorkerStartScriptModel) (string, error) {
	return renderTemplateFunc(dmWorkerStartScriptTpl, model)
}
