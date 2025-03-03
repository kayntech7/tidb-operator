---
title: TiDB Operator API Document
summary: Reference of TiDB Operator API
category: how-to
---
<h1>API Document</h1>
<h2 id="federation.pingcap.com/v1alpha1">federation.pingcap.com/v1alpha1</h2>
Resource Types:
<ul><li>
<a href="#volumebackup">VolumeBackup</a>
</li><li>
<a href="#volumebackupschedule">VolumeBackupSchedule</a>
</li><li>
<a href="#volumerestore">VolumeRestore</a>
</li></ul>
<h3 id="volumebackup">VolumeBackup</h3>
<p>
<p>VolumeBackup is the control script&rsquo;s spec</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>apiVersion</code></br>
string</td>
<td>
<code>
federation.pingcap.com/v1alpha1
</code>
</td>
</tr>
<tr>
<td>
<code>kind</code></br>
string
</td>
<td><code>VolumeBackup</code></td>
</tr>
<tr>
<td>
<code>metadata</code></br>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.19/#objectmeta-v1-meta">
Kubernetes meta/v1.ObjectMeta
</a>
</em>
</td>
<td>
Refer to the Kubernetes API documentation for the fields of the
<code>metadata</code> field.
</td>
</tr>
<tr>
<td>
<code>spec</code></br>
<em>
<a href="#volumebackupspec">
VolumeBackupSpec
</a>
</em>
</td>
<td>
<br/>
<br/>
<table>
<tr>
<td>
<code>clusters</code></br>
<em>
<a href="#volumebackupmembercluster">
[]VolumeBackupMemberCluster
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>template</code></br>
<em>
<a href="#volumebackupmemberspec">
VolumeBackupMemberSpec
</a>
</em>
</td>
<td>
</td>
</tr>
</table>
</td>
</tr>
<tr>
<td>
<code>status</code></br>
<em>
<a href="#volumebackupstatus">
VolumeBackupStatus
</a>
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="volumebackupschedule">VolumeBackupSchedule</h3>
<p>
<p>VolumeBackupSchedule is the control script&rsquo;s spec</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>apiVersion</code></br>
string</td>
<td>
<code>
federation.pingcap.com/v1alpha1
</code>
</td>
</tr>
<tr>
<td>
<code>kind</code></br>
string
</td>
<td><code>VolumeBackupSchedule</code></td>
</tr>
<tr>
<td>
<code>metadata</code></br>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.19/#objectmeta-v1-meta">
Kubernetes meta/v1.ObjectMeta
</a>
</em>
</td>
<td>
Refer to the Kubernetes API documentation for the fields of the
<code>metadata</code> field.
</td>
</tr>
<tr>
<td>
<code>spec</code></br>
<em>
<a href="#volumebackupschedulespec">
VolumeBackupScheduleSpec
</a>
</em>
</td>
<td>
<br/>
<br/>
<table>
</table>
</td>
</tr>
<tr>
<td>
<code>status</code></br>
<em>
<a href="#volumebackupschedulestatus">
VolumeBackupScheduleStatus
</a>
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="volumerestore">VolumeRestore</h3>
<p>
<p>VolumeRestore is the control script&rsquo;s spec</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>apiVersion</code></br>
string</td>
<td>
<code>
federation.pingcap.com/v1alpha1
</code>
</td>
</tr>
<tr>
<td>
<code>kind</code></br>
string
</td>
<td><code>VolumeRestore</code></td>
</tr>
<tr>
<td>
<code>metadata</code></br>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.19/#objectmeta-v1-meta">
Kubernetes meta/v1.ObjectMeta
</a>
</em>
</td>
<td>
Refer to the Kubernetes API documentation for the fields of the
<code>metadata</code> field.
</td>
</tr>
<tr>
<td>
<code>spec</code></br>
<em>
<a href="#volumerestorespec">
VolumeRestoreSpec
</a>
</em>
</td>
<td>
<br/>
<br/>
<table>
</table>
</td>
</tr>
<tr>
<td>
<code>status</code></br>
<em>
<a href="#volumerestorestatus">
VolumeRestoreStatus
</a>
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="brconfig">BRConfig</h3>
<p>
(<em>Appears on:</em>
<a href="#volumebackupmemberspec">VolumeBackupMemberSpec</a>)
</p>
<p>
<p>BRConfig contains config for BR</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>concurrency</code></br>
<em>
uint32
</em>
</td>
<td>
<p>Concurrency is the size of thread pool on each node that execute the backup task</p>
</td>
</tr>
<tr>
<td>
<code>checkRequirements</code></br>
<em>
bool
</em>
</td>
<td>
<p>CheckRequirements specifies whether to check requirements</p>
</td>
</tr>
<tr>
<td>
<code>sendCredToTikv</code></br>
<em>
bool
</em>
</td>
<td>
<p>SendCredToTikv specifies whether to send credentials to TiKV</p>
</td>
</tr>
<tr>
<td>
<code>options</code></br>
<em>
[]string
</em>
</td>
<td>
<p>Options means options for backup data to remote storage with BR. These options has highest priority.</p>
</td>
</tr>
</tbody>
</table>
<h3 id="volumebackupcondition">VolumeBackupCondition</h3>
<p>
(<em>Appears on:</em>
<a href="#volumebackupstatus">VolumeBackupStatus</a>)
</p>
<p>
<p>VolumeBackupCondition describes the observed state of a VolumeBackup at a certain point.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>status</code></br>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.19/#conditionstatus-v1-core">
Kubernetes core/v1.ConditionStatus
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>type</code></br>
<em>
<a href="#volumebackupconditiontype">
VolumeBackupConditionType
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>lastTransitionTime</code></br>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.19/#time-v1-meta">
Kubernetes meta/v1.Time
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>reason</code></br>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>message</code></br>
<em>
string
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="volumebackupconditiontype">VolumeBackupConditionType</h3>
<p>
(<em>Appears on:</em>
<a href="#volumebackupcondition">VolumeBackupCondition</a>, 
<a href="#volumebackupstatus">VolumeBackupStatus</a>)
</p>
<p>
<p>VolumeBackupConditionType represents a valid condition of a VolumeBackup.</p>
</p>
<h3 id="volumebackupmembercluster">VolumeBackupMemberCluster</h3>
<p>
(<em>Appears on:</em>
<a href="#volumebackupspec">VolumeBackupSpec</a>)
</p>
<p>
<p>VolumeBackupMemberCluster contains the TiDB cluster which need to execute volume backup</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>k8sClusterName</code></br>
<em>
string
</em>
</td>
<td>
<p>K8sClusterName is the name of the k8s cluster where the tc locates</p>
</td>
</tr>
<tr>
<td>
<code>tcName</code></br>
<em>
string
</em>
</td>
<td>
<p>TCName is the name of the TiDBCluster CR which need to execute volume backup</p>
</td>
</tr>
<tr>
<td>
<code>tcNamespace</code></br>
<em>
string
</em>
</td>
<td>
<p>TCNamespace is the namespace of the TiDBCluster CR</p>
</td>
</tr>
</tbody>
</table>
<h3 id="volumebackupmemberspec">VolumeBackupMemberSpec</h3>
<p>
(<em>Appears on:</em>
<a href="#volumebackupspec">VolumeBackupSpec</a>)
</p>
<p>
<p>VolumeBackupMemberSpec contains the backup specification for one tidb cluster</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>resources</code></br>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.19/#resourcerequirements-v1-core">
Kubernetes core/v1.ResourceRequirements
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>env</code></br>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.19/#envvar-v1-core">
[]Kubernetes core/v1.EnvVar
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>List of environment variables to set in the container, like v1.Container.Env.
Note that the following builtin env vars will be overwritten by values set here
- S3_PROVIDER
- S3_ENDPOINT
- AWS_REGION
- AWS_ACL
- AWS_STORAGE_CLASS
- AWS_DEFAULT_REGION
- AWS_ACCESS_KEY_ID
- AWS_SECRET_ACCESS_KEY
- GCS_PROJECT_ID
- GCS_OBJECT_ACL
- GCS_BUCKET_ACL
- GCS_LOCATION
- GCS_STORAGE_CLASS
- GCS_SERVICE_ACCOUNT_JSON_KEY
- BR_LOG_TO_TERM</p>
</td>
</tr>
<tr>
<td>
<code>br</code></br>
<em>
<a href="#brconfig">
BRConfig
</a>
</em>
</td>
<td>
<p>BRConfig is the configs for BR</p>
</td>
</tr>
<tr>
<td>
<code>StorageProvider</code></br>
<em>
github.com/pingcap/tidb-operator/pkg/apis/pingcap/v1alpha1.StorageProvider
</em>
</td>
<td>
<p>
(Members of <code>StorageProvider</code> are embedded into this type.)
</p>
</td>
</tr>
<tr>
<td>
<code>tolerations</code></br>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.19/#toleration-v1-core">
[]Kubernetes core/v1.Toleration
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>toolImage</code></br>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>ToolImage specifies the tool image used in <code>Backup</code>, which supports BR and Dumpling images.
For examples <code>spec.toolImage: pingcap/br:v4.0.8</code> or <code>spec.toolImage: pingcap/dumpling:v4.0.8</code>
For BR image, if it does not contain tag, Pod will use image &lsquo;ToolImage:${TiKV_Version}&rsquo;.</p>
</td>
</tr>
<tr>
<td>
<code>imagePullSecrets</code></br>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.19/#localobjectreference-v1-core">
[]Kubernetes core/v1.LocalObjectReference
</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>ImagePullSecrets is an optional list of references to secrets in the same namespace to use for pulling any of the images.</p>
</td>
</tr>
<tr>
<td>
<code>serviceAccount</code></br>
<em>
string
</em>
</td>
<td>
<p>Specify service account of backup</p>
</td>
</tr>
<tr>
<td>
<code>cleanPolicy</code></br>
<em>
github.com/pingcap/tidb-operator/pkg/apis/pingcap/v1alpha1.CleanPolicyType
</em>
</td>
<td>
<p>CleanPolicy denotes whether to clean backup data when the object is deleted from the cluster, if not set, the backup data will be retained</p>
</td>
</tr>
<tr>
<td>
<code>priorityClassName</code></br>
<em>
string
</em>
</td>
<td>
<p>PriorityClassName of Backup Job Pods</p>
</td>
</tr>
</tbody>
</table>
<h3 id="volumebackupmemberstatus">VolumeBackupMemberStatus</h3>
<p>
(<em>Appears on:</em>
<a href="#volumebackupstatus">VolumeBackupStatus</a>)
</p>
<p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>k8sClusterName</code></br>
<em>
string
</em>
</td>
<td>
<p>K8sClusterName is the name of the k8s cluster where the tc locates</p>
</td>
</tr>
<tr>
<td>
<code>tcName</code></br>
<em>
string
</em>
</td>
<td>
<p>TCName is the name of the TiDBCluster CR which need to execute volume backup</p>
</td>
</tr>
<tr>
<td>
<code>tcNamespace</code></br>
<em>
string
</em>
</td>
<td>
<p>TCNamespace is the namespace of the TiDBCluster CR</p>
</td>
</tr>
<tr>
<td>
<code>backupName</code></br>
<em>
string
</em>
</td>
<td>
<p>BackupName is the name of Backup CR</p>
</td>
</tr>
<tr>
<td>
<code>backupPath</code></br>
<em>
string
</em>
</td>
<td>
<p>BackupPath is the location of the backup</p>
</td>
</tr>
<tr>
<td>
<code>backupSize</code></br>
<em>
int64
</em>
</td>
<td>
<p>BackupSize is the data size of the backup</p>
</td>
</tr>
<tr>
<td>
<code>commitTs</code></br>
<em>
string
</em>
</td>
<td>
<p>CommitTs is the commit ts of the backup</p>
</td>
</tr>
</tbody>
</table>
<h3 id="volumebackupschedulespec">VolumeBackupScheduleSpec</h3>
<p>
(<em>Appears on:</em>
<a href="#volumebackupschedule">VolumeBackupSchedule</a>)
</p>
<p>
<p>VolumeBackupScheduleSpec describes the attributes that a user creates on a volume backup schedule.</p>
</p>
<h3 id="volumebackupschedulestatus">VolumeBackupScheduleStatus</h3>
<p>
(<em>Appears on:</em>
<a href="#volumebackupschedule">VolumeBackupSchedule</a>)
</p>
<p>
<p>VolumeBackupScheduleStatus represents the current status of a volume backup schedule.</p>
</p>
<h3 id="volumebackupspec">VolumeBackupSpec</h3>
<p>
(<em>Appears on:</em>
<a href="#volumebackup">VolumeBackup</a>)
</p>
<p>
<p>VolumeBackupSpec describes the attributes that a user creates on a volume backup.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>clusters</code></br>
<em>
<a href="#volumebackupmembercluster">
[]VolumeBackupMemberCluster
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>template</code></br>
<em>
<a href="#volumebackupmemberspec">
VolumeBackupMemberSpec
</a>
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="volumebackupstatus">VolumeBackupStatus</h3>
<p>
(<em>Appears on:</em>
<a href="#volumebackup">VolumeBackup</a>)
</p>
<p>
<p>VolumeBackupStatus represents the current status of a volume backup.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>backups</code></br>
<em>
<a href="#volumebackupmemberstatus">
[]VolumeBackupMemberStatus
</a>
</em>
</td>
<td>
<p>Backups are volume backups&rsquo; information in data plane</p>
</td>
</tr>
<tr>
<td>
<code>timeStarted</code></br>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.19/#time-v1-meta">
Kubernetes meta/v1.Time
</a>
</em>
</td>
<td>
<p>TimeStarted is the time at which the backup was started.</p>
</td>
</tr>
<tr>
<td>
<code>timeCompleted</code></br>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.19/#time-v1-meta">
Kubernetes meta/v1.Time
</a>
</em>
</td>
<td>
<p>TimeCompleted is the time at which the backup was completed.</p>
</td>
</tr>
<tr>
<td>
<code>timeTaken</code></br>
<em>
string
</em>
</td>
<td>
<p>TimeTaken is the time that volume backup federation takes, it is TimeCompleted - TimeStarted</p>
</td>
</tr>
<tr>
<td>
<code>backupSizeReadable</code></br>
<em>
string
</em>
</td>
<td>
<p>BackupSizeReadable is the data size of the backup.
the difference with BackupSize is that its format is human readable</p>
</td>
</tr>
<tr>
<td>
<code>backupSize</code></br>
<em>
int64
</em>
</td>
<td>
<p>BackupSize is the data size of the backup.</p>
</td>
</tr>
<tr>
<td>
<code>commitTs</code></br>
<em>
string
</em>
</td>
<td>
<p>CommitTs is the commit ts of the backup, snapshot ts for full backup or start ts for log backup.</p>
</td>
</tr>
<tr>
<td>
<code>phase</code></br>
<em>
<a href="#volumebackupconditiontype">
VolumeBackupConditionType
</a>
</em>
</td>
<td>
<p>Phase is a user readable state inferred from the underlying Backup conditions</p>
</td>
</tr>
<tr>
<td>
<code>conditions</code></br>
<em>
<a href="#volumebackupcondition">
[]VolumeBackupCondition
</a>
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="volumerestorecondition">VolumeRestoreCondition</h3>
<p>
(<em>Appears on:</em>
<a href="#volumerestorestatus">VolumeRestoreStatus</a>)
</p>
<p>
<p>VolumeRestoreCondition describes the observed state of a VolumeRestore at a certain point.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>status</code></br>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.19/#conditionstatus-v1-core">
Kubernetes core/v1.ConditionStatus
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>lastTransitionTime</code></br>
<em>
<a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.19/#time-v1-meta">
Kubernetes meta/v1.Time
</a>
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>reason</code></br>
<em>
string
</em>
</td>
<td>
</td>
</tr>
<tr>
<td>
<code>message</code></br>
<em>
string
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<h3 id="volumerestorespec">VolumeRestoreSpec</h3>
<p>
(<em>Appears on:</em>
<a href="#volumerestore">VolumeRestore</a>)
</p>
<p>
<p>VolumeRestoreSpec describes the attributes that a user creates on a volume restore.</p>
</p>
<h3 id="volumerestorestatus">VolumeRestoreStatus</h3>
<p>
(<em>Appears on:</em>
<a href="#volumerestore">VolumeRestore</a>)
</p>
<p>
<p>VolumeRestoreStatus represents the current status of a volume restore.</p>
</p>
<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>
<tr>
<td>
<code>conditions</code></br>
<em>
<a href="#volumerestorecondition">
[]VolumeRestoreCondition
</a>
</em>
</td>
<td>
</td>
</tr>
</tbody>
</table>
<hr/>
<p><em>
Generated with <code>gen-crd-api-reference-docs</code>
</em></p>
