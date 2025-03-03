//
// Jenkins pipeline for e2e kind job.
//
// We uses ghprb plugin to build pull requests and report results. Some special
// environment variables will be available for jobs that are triggered by GitHub
// Pull Request events.
//
// - ghprbActualCommit
//
// For more information about this plugin, please check out https://plugins.jenkins.io/ghprb/.
//

// Able to override default values in Jenkins job via environment variables.
env.DEFAULT_GIT_REF = env.DEFAULT_GIT_REF ?: 'master'
env.DEFAULT_GINKGO_NODES = env.DEFAULT_GINKGO_NODES ?: '8'
env.DEFAULT_E2E_ARGS = env.DEFAULT_E2E_ARGS ?: ''
env.DEFAULT_DELETE_NAMESPACE_ON_FAILURE = env.DEFAULT_DELETE_NAMESPACE_ON_FAILURE ?: 'true'

properties([
    parameters([
        string(name: 'GIT_URL', defaultValue: 'https://github.com/pingcap/tidb-operator', description: 'git repo url'),
        string(name: 'GIT_REF', defaultValue: env.DEFAULT_GIT_REF, description: 'git ref spec to checkout, e.g. master, release-1.1'),
        string(name: 'RELEASE_VER', defaultValue: '', description: "the version string in released tarball"),
        string(name: 'PR_ID', defaultValue: '', description: 'pull request ID, this will override GIT_REF if set, e.g. 1889'),
        string(name: 'GINKGO_NODES', defaultValue: env.DEFAULT_GINKGO_NODES, description: 'the number of ginkgo nodes'),
        string(name: 'E2E_ARGS', defaultValue: env.DEFAULT_E2E_ARGS, description: "e2e args, e.g. --ginkgo.focus='\\[Stability\\]'"),
        string(name: 'DELETE_NAMESPACE_ON_FAILURE', defaultValue: env.DEFAULT_DELETE_NAMESPACE_ON_FAILURE, description: 'delete ns after test case fails'),

        string(name: 'CUSTOM_PORT_TIDB_SERVER', defaultValue: '', description: 'custom component port: tidb server'),
        string(name: 'CUSTOM_PORT_TIDB_STATUS', defaultValue: '', description: 'custom component port: tidb status'),
        string(name: 'CUSTOM_PORT_PD_CLIENT', defaultValue: '', description: 'custom component port: pd client'),
        string(name: 'CUSTOM_PORT_PD_PEER', defaultValue: '', description: 'custom component port: pd peer'),
        string(name: 'CUSTOM_PORT_TIKV_SERVER', defaultValue: '', description: 'custom component port: tikv server'),
        string(name: 'CUSTOM_PORT_TIKV_STATUS', defaultValue: '', description: 'custom component port: tikv status'),
        string(name: 'CUSTOM_PORT_TIFLASH_TCP', defaultValue: '', description: 'custom component port: tiflash tcp'),
        string(name: 'CUSTOM_PORT_TIFLASH_HTTP', defaultValue: '', description: 'custom component port: tiflash http'),
        string(name: 'CUSTOM_PORT_TIFLASH_FLASH', defaultValue: '', description: 'custom component port: tiflash flash'),
        string(name: 'CUSTOM_PORT_TIFLASH_PROXY', defaultValue: '', description: 'custom component port: tiflash proxy'),
        string(name: 'CUSTOM_PORT_TIFLASH_METRICS', defaultValue: '', description: 'custom component port: tiflash metrics'),
        string(name: 'CUSTOM_PORT_TIFLASH_PROXY_STATUS', defaultValue: '', description: 'custom component port: tiflash proxy status'),
        string(name: 'CUSTOM_PORT_TIFLASH_INTERNAL_STATUS', defaultValue: '', description: 'custom component port: tiflash internal status'),
        string(name: 'CUSTOM_PORT_PUMP', defaultValue: '', description: 'custom component port: pump'),
        string(name: 'CUSTOM_PORT_DRAINER', defaultValue: '', description: 'custom component port: drainer'),
        string(name: 'CUSTOM_PORT_TICDC', defaultValue: '', description: 'custom component port: ticdc')
    ])
])

podYAML = '''\
apiVersion: v1
kind: Pod
metadata:
  labels:
    app: tidb-operator-e2e
spec:
  containers:
  - name: main
    image: hub-new.pingcap.net/tidb-operator/kubekins-e2e:v5-go1.19
    command:
    - runner.sh
    - exec
    - bash
    - -c
    - |
      sleep 1d & wait
    # we need privileged mode in order to do docker in docker
    securityContext:
      privileged: true
    env:
    - name: DOCKER_IN_DOCKER_ENABLED
      value: "true"
<% if (resources && (resources.requests || resources.limits)) { %>
    resources:
    <% if (resources.requests) { %>
      requests:
        cpu: <%= resources.requests.cpu %>
        memory: <%= resources.requests.memory %>
        ephemeral-storage: <%= resources.requests.storage %>
    <% } %>
    <% if (resources.limits) { %>
      limits:
        cpu: <%= resources.limits.cpu %>
        memory: <%= resources.limits.memory %>
    <% } %>
<% } %>
    # kind needs /lib/modules and cgroups from the host
    volumeMounts:
    - mountPath: /lib/modules
      name: modules
      readOnly: true
    - mountPath: /sys/fs/cgroup
      name: cgroup
    # dind expects /var/lib/docker to be volume
    - name: docker-root
      mountPath: /var/lib/docker
    # legacy docker path for cr.io/k8s-testimages/kubekins-e2e
    - name: docker-graph
      mountPath: /docker-graph
    # use memory storage for etcd hostpath in kind cluster
    - name: kind-data-dir
      mountPath: /kind-data
    - name: etcd-data-dir
      mountPath: /mnt/tmpfs/etcd
  volumes:
  - name: modules
    hostPath:
      path: /lib/modules
      type: Directory
  - name: cgroup
    hostPath:
      path: /sys/fs/cgroup
      type: Directory
  - name: docker-root
    emptyDir: {}
  - name: docker-graph
    emptyDir: {}
  - name: kind-data-dir
    emptyDir: {}
  - name: etcd-data-dir
    emptyDir: {}
  tolerations:
  - effect: NoSchedule
    key: tidb-operator
    operator: Exists
  affinity:
    podAntiAffinity:
      preferredDuringSchedulingIgnoredDuringExecution:
      - weight: 100
        podAffinityTerm:
          labelSelector:
            matchExpressions:
            - key: app
              operator: In
              values:
              - tidb-operator-e2e
          topologyKey: kubernetes.io/hostname
'''

String buildPodYAML(Map m = [:]) {
    m.putIfAbsent("resources", [:])
    m.putIfAbsent("any", false)
    def engine = new groovy.text.SimpleTemplateEngine()
    def template = engine.createTemplate(podYAML).make(m)
    return template.toString()
}

e2ePodResources = [
    requests: [
        cpu: "8",
        memory: "16Gi",
        storage: "250Gi"
    ],
    limits: [
        cpu: "16",
        memory: "32Gi",
        storage: "250Gi"
    ],
]

def build(String name, String code, Map resources = e2ePodResources) {
    podTemplate(yaml: buildPodYAML(resources: resources), namespace: "jenkins-tidb-operator", cloud: "kubernetes-ng") {
        node(POD_LABEL) {
            container('main') {
                def WORKSPACE = pwd()
                def ARTIFACTS = "${WORKSPACE}/go/src/github.com/pingcap/tidb-operator/_artifacts"
                try {
                    dir("${WORKSPACE}/go/src/github.com/pingcap/tidb-operator") {
                        unstash 'tidb-operator'
                        stage("Debug Info") {
                            sh """
                            echo "====== shell env ======"
                            echo "pwd: \$(pwd)"
                            env
                            echo "====== go env ======"
                            go env
                            echo "====== docker version ======"
                            docker version
                            """
                        }
                        stage('Run') {
                            sh """#!/bin/bash
                            export GOPATH=${WORKSPACE}/go
                            export ARTIFACTS=${ARTIFACTS}
                            export RUNNER_SUITE_NAME=${name}

                            echo "info: create local path for data and coverage"
                            mount --make-rshared /
                            mkdir -p /kind-data/control-plane/coverage
                            mkdir -p /kind-data/worker1/coverage
                            mkdir -p /kind-data/worker2/coverage
                            mkdir -p /kind-data/worker3/coverage
                            ${code}
                            """
                        }
                        stage('Coverage') {
                            withCredentials([
                                string(credentialsId: "tp-codecov-token", variable: 'CODECOV_TOKEN')
                            ]) {
                                sh """#!/bin/bash
                                echo "info: list all coverage files"
                                ls -dla /kind-data/control-plane/coverage/*
                                ls -dla /kind-data/worker1/coverage/*
                                ls -dla /kind-data/worker2/coverage/*
                                ls -dla /kind-data/worker3/coverage/*
                                echo "info: merging coverage files"
                                cp /kind-data/control-plane/coverage/*.cov /tmp
                                cp /kind-data/worker1/coverage/*.cov /tmp
                                cp /kind-data/worker2/coverage/*.cov /tmp
                                cp /kind-data/worker3/coverage/*.cov /tmp
                                ./bin/gocovmerge /tmp/*.cov > /tmp/coverage.txt
                                source EXPORT_GIT_COMMIT
                                echo "info: uploading coverage to codecov"
                                bash <(curl -s https://codecov.io/bash) -t ${CODECOV_TOKEN} -F e2e -n tidb-operator -f /tmp/coverage.txt
                                """
                            }
                        }
                    }
                } finally {
                    stage('Artifacts') {
                        dir(ARTIFACTS) {
                            sh """#!/bin/bash
                            echo "info: change ownerships for jenkins"
                            chown -R 1000:1000 .
                            echo "info: print total size of artifacts"
                            du -sh .
                            echo "info: list all files"
                            find .
                            echo "info: moving all artifacts into a sub-directory"
                            shopt -s extglob
                            mkdir ${name}
                            mv !(${name}) ${name}/
                            """
                            archiveArtifacts artifacts: "${name}/**", allowEmptyArchive: true
                            junit testResults: "${name}/*.xml", allowEmptyResults: true, keepLongStdio: true
                        }
                    }
                }
            }
        }
    }
}


try {

    def GITHASH
    def IMAGE_TAG

    def PROJECT_DIR = "go/src/github.com/pingcap/tidb-operator"

    // Git ref to checkout
    def GIT_REF = params.GIT_REF
    if (params.PR_ID != "") {
        GIT_REF = "refs/remotes/origin/pull/${params.PR_ID}/head"
    } else if (env.ghprbActualCommit) {
        // for PR jobs triggered by ghprb plugin
        GIT_REF = env.ghprbActualCommit
    }

    def CUSTOM_PORT_TIDB_SERVER = params.CUSTOM_PORT_TIDB_SERVER
    def CUSTOM_PORT_TIDB_STATUS = params.CUSTOM_PORT_TIDB_STATUS
    def CUSTOM_PORT_PD_CLIENT = params.CUSTOM_PORT_PD_CLIENT
    def CUSTOM_PORT_PD_PEER = params.CUSTOM_PORT_PD_PEER
    def CUSTOM_PORT_TIKV_SERVER = params.CUSTOM_PORT_TIKV_SERVER
    def CUSTOM_PORT_TIKV_STATUS = params.CUSTOM_PORT_TIKV_STATUS
    def CUSTOM_PORT_TIFLASH_TCP = params.CUSTOM_PORT_TIFLASH_TCP
    def CUSTOM_PORT_TIFLASH_HTTP = params.CUSTOM_PORT_TIFLASH_HTTP
    def CUSTOM_PORT_TIFLASH_FLASH = params.CUSTOM_PORT_TIFLASH_FLASH
    def CUSTOM_PORT_TIFLASH_PROXY = params.CUSTOM_PORT_TIFLASH_PROXY
    def CUSTOM_PORT_TIFLASH_METRICS = params.CUSTOM_PORT_TIFLASH_METRICS
    def CUSTOM_PORT_TIFLASH_PROXY_STATUS = params.CUSTOM_PORT_TIFLASH_PROXY_STATUS
    def CUSTOM_PORT_TIFLASH_INTERNAL_STATUS = params.CUSTOM_PORT_TIFLASH_INTERNAL_STATUS
    def CUSTOM_PORT_PUMP = params.CUSTOM_PORT_PUMP
    def CUSTOM_PORT_DRAINER = params.CUSTOM_PORT_DRAINER
    def CUSTOM_PORT_TICDC = params.CUSTOM_PORT_TICDC

    timeout (time: 2, unit: 'HOURS') {
        // use fixed label, so we can reuse previous workers
        // increase version in pod label when we update pod template
        def buildPodLabel = "tidb-operator-build-v5-pingcap-docker-mirror"
        def resources = [
            requests: [
                cpu: "4",
                memory: "10Gi",
                storage: "50Gi"
            ],
            limits: [
                cpu: "8",
                memory: "32Gi",
                storage: "50Gi"
            ],
        ]
        podTemplate(
            cloud: "kubernetes-ng",
            namespace: "jenkins-tidb-operator",
            label: buildPodLabel,
            yaml: buildPodYAML(resources: resources, any: true),
            // We allow this pod to remain active for a while, later jobs can
            // reuse cache in previous created nodes.
            idleMinutes: 30,
        ) {
        node(buildPodLabel) {
            container("main") {
                dir("${PROJECT_DIR}") {

                    stage('Checkout') {
                        sh """
                        echo "info: change ownerships for jenkins"
                        # we run as root in our pods, this is required
                        # otherwise jenkins agent will fail because of the lack of permission
                        chown -R 1000:1000 .
                        git config --global --add safe.directory '*'
                        """

                        // clean stale files because we may reuse previous created nodes
                        deleteDir()

                        try {
                            checkout changelog: false, poll: false, scm: [
                                    $class: 'GitSCM',
                                    branches: [[name: "${GIT_REF}"]],
                                    userRemoteConfigs: [[
                                            refspec: '+refs/heads/*:refs/remotes/origin/* +refs/pull/*:refs/remotes/origin/pull/*',
                                            url: "${params.GIT_URL}",
                                    ]]
                            ]
                        } catch (info) {
                            retry(3) {
                                echo "checkout failed, retry.."
                                sleep 10
                                checkout changelog: false, poll: false, scm: [
                                        $class: 'GitSCM',
                                        branches: [[name: "${GIT_REF}"]],
                                        userRemoteConfigs: [[
                                                refspec: '+refs/heads/*:refs/remotes/origin/* +refs/pull/*:refs/remotes/origin/pull/*',
                                                url: "${params.GIT_URL}",
                                        ]]
                                ]
                            }
                        }



                        GITHASH = sh(returnStdout: true, script: "git rev-parse HEAD").trim()
                        IMAGE_TAG = env.JOB_NAME + "-" + GITHASH.substring(0, 6)
                    }

                    stage("Build") {
                        withCredentials([
                            string(credentialsId: "tp-codecov-token", variable: 'CODECOV_TOKEN')
                        ]) {
                            sh """#!/bin/bash
                            set -eu
                            echo "info: building"
                            echo "info: patch charts and golang code to enable coverage profile"
                            ./hack/e2e-patch-codecov.sh
                            export CUSTOM_PORT_TIDB_SERVER=${CUSTOM_PORT_TIDB_SERVER}
                            export CUSTOM_PORT_TIDB_STATUS=${CUSTOM_PORT_TIDB_STATUS}
                            export CUSTOM_PORT_PD_CLIENT=${CUSTOM_PORT_PD_CLIENT}
                            export CUSTOM_PORT_PD_PEER=${CUSTOM_PORT_PD_PEER}
                            export CUSTOM_PORT_TIKV_SERVER=${CUSTOM_PORT_TIKV_SERVER}
                            export CUSTOM_PORT_TIKV_STATUS=${CUSTOM_PORT_TIKV_STATUS}
                            export CUSTOM_PORT_TIFLASH_TCP=${CUSTOM_PORT_TIFLASH_TCP}
                            export CUSTOM_PORT_TIFLASH_HTTP=${CUSTOM_PORT_TIFLASH_HTTP}
                            export CUSTOM_PORT_TIFLASH_FLASH=${CUSTOM_PORT_TIFLASH_FLASH}
                            export CUSTOM_PORT_TIFLASH_PROXY=${CUSTOM_PORT_TIFLASH_PROXY}
                            export CUSTOM_PORT_TIFLASH_METRICS=${CUSTOM_PORT_TIFLASH_METRICS}
                            export CUSTOM_PORT_TIFLASH_PROXY_STATUS=${CUSTOM_PORT_TIFLASH_PROXY_STATUS}
                            export CUSTOM_PORT_TIFLASH_INTERNAL_STATUS=${CUSTOM_PORT_TIFLASH_INTERNAL_STATUS}
                            export CUSTOM_PORT_PUMP=${CUSTOM_PORT_PUMP}
                            export CUSTOM_PORT_DRAINER=${CUSTOM_PORT_DRAINER}
                            export CUSTOM_PORT_TICDC=${CUSTOM_PORT_TICDC}
                            E2E=y make build e2e-build
                            make gocovmerge
                            """
                        }
                    }

                    stage("Prepare for e2e") {
                        withCredentials([usernamePassword(credentialsId: 'TIDB_OPERATOR_HUB_AUTH', usernameVariable: 'USERNAME', passwordVariable: 'PASSWORD')]) {
                            sh """#!/bin/bash
                            set -eu
                            echo "save GTI_COMMIT export script into file"
                            echo "export GIT_COMMIT=\$(git rev-parse HEAD)" > EXPORT_GIT_COMMIT
                            echo "info: logging into hub.pingcap.net"
                            docker login -u \$USERNAME --password-stdin hub.pingcap.net <<< \$PASSWORD
                            echo "info: build and push images for e2e"
                            echo "test: show docker daemon config file"
                            cat /etc/docker/daemon.json || true
                            E2E=y NO_BUILD=y DOCKER_REPO=hub.pingcap.net/tidb-operator-e2e IMAGE_TAG=${IMAGE_TAG} make docker-push e2e-docker-push
                            echo "info: download binaries for e2e"
                            E2E=y SKIP_BUILD=y SKIP_IMAGE_BUILD=y SKIP_UP=y SKIP_TEST=y SKIP_DOWN=y ./hack/e2e.sh
                            echo "info: change ownerships for jenkins"
                            # we run as root in our pods, this is required
                            # otherwise jenkins agent will fail because of the lack of permission
                            chown -R 1000:1000 .
                            """
                        }
                        stash excludes: "vendor/**,deploy/**,tests/**", name: "tidb-operator"
                    }
                }
            }
        }
        }

        def GLOBALS = "KIND_DATA_HOSTPATH=/kind-data KIND_ETCD_DATADIR=/mnt/tmpfs/etcd E2E=y SKIP_BUILD=y SKIP_IMAGE_BUILD=y DOCKER_REPO=hub.pingcap.net/tidb-operator-e2e IMAGE_TAG=${IMAGE_TAG} DELETE_NAMESPACE_ON_FAILURE=${params.DELETE_NAMESPACE_ON_FAILURE} GINKGO_NO_COLOR=y"
        build("tidb-operator", "${GLOBALS} GINKGO_NODES=${params.GINKGO_NODES} ./hack/e2e.sh -- ${params.E2E_ARGS}")

        if (GIT_REF ==~ /^(master|)$/ || GIT_REF ==~ /^(release-.*)$/
            || GIT_REF ==~ /^v(0|[1-9]\d*)\.(0|[1-9]\d*)\.(0|[1-9]\d*)(?:-((?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*)(?:\.(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*))*))?(?:\+([0-9a-zA-Z-]+(?:\.[0-9a-zA-Z-]+)*))?$/) {
            // Upload assets if the git ref is the master branch or version tag
            podTemplate(yaml: buildPodYAML(resources: [requests: [cpu: "1", memory: "2Gi"]]),
              cloud: "kubernetes-ng", namespace: "jenkins-tidb-operator",
            ) {
                node(POD_LABEL) {
                    container("main") {
                        dir("${PROJECT_DIR}") {
                            unstash 'tidb-operator'
                            stage('upload tidb-operator binaries and charts'){
                              sh """
                              export BUILD_BRANCH=${GIT_REF}
                              export GITHASH=${GITHASH}
                              ./ci/upload-binaries-charts.sh
                              """
                            }
                        }
                    }
                }
            }
        }
    }
    currentBuild.result = "SUCCESS"
} catch (err) {
    println("fatal: " + err)
    currentBuild.result = 'FAILURE'
}

// vim: et sw=4 ts=4
