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

package util

import (
	"errors"
	"strings"

	"github.com/pingcap/tidb-operator/pkg/apis/pingcap/v1alpha1"
	crdutils "github.com/yisaer/crd-validation/pkg"
	extensionsobj "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
)

var (
	tidbClusteradditionalPrinterColumns []extensionsobj.CustomResourceColumnDefinition
	tidbClusterReadyColumn              = extensionsobj.CustomResourceColumnDefinition{
		Name:     "Ready",
		Type:     "string",
		JSONPath: `.status.conditions[?(@.type=="Ready")].status`,
	}
	tidbClusterStatusMessageColumn = extensionsobj.CustomResourceColumnDefinition{
		Name:     "Status",
		Type:     "string",
		JSONPath: `.status.conditions[?(@.type=="Ready")].message`,
		Priority: 1,
	}
	tidbClusterPDColumn = extensionsobj.CustomResourceColumnDefinition{
		Name:        "PD",
		Type:        "string",
		Description: "The image for PD cluster",
		JSONPath:    ".status.pd.image",
	}
	tidbClusterPDStorageColumn = extensionsobj.CustomResourceColumnDefinition{
		Name:        "Storage",
		Type:        "string",
		Description: "The storage size specified for PD node",
		JSONPath:    ".spec.pd.requests.storage",
	}
	tidbClusterPDReadyColumn = extensionsobj.CustomResourceColumnDefinition{
		Name:        "Ready",
		Type:        "integer",
		Description: "The desired replicas number of PD cluster",
		JSONPath:    ".status.pd.statefulSet.readyReplicas",
	}
	tidbClusterPDDesireColumn = extensionsobj.CustomResourceColumnDefinition{
		Name:        "Desire",
		Type:        "integer",
		Description: "The desired replicas number of PD cluster",
		JSONPath:    ".spec.pd.replicas",
	}
	tidbClusterTiKVColumn = extensionsobj.CustomResourceColumnDefinition{
		Name:        "TiKV",
		Type:        "string",
		Description: "The image for TiKV cluster",
		JSONPath:    ".status.tikv.image",
	}
	tidbClusterTiKVStorageColumn = extensionsobj.CustomResourceColumnDefinition{
		Name:        "Storage",
		Type:        "string",
		Description: "The storage size specified for TiKV node",
		JSONPath:    ".spec.tikv.requests.storage",
	}
	tidbClusterTiKVReadyColumn = extensionsobj.CustomResourceColumnDefinition{
		Name:        "Ready",
		Type:        "integer",
		Description: "The ready replicas number of TiKV cluster",
		JSONPath:    ".status.tikv.statefulSet.readyReplicas",
	}
	tidbClusterTiKVDesireColumn = extensionsobj.CustomResourceColumnDefinition{
		Name:        "Desire",
		Type:        "integer",
		Description: "The desired replicas number of TiKV cluster",
		JSONPath:    ".spec.tikv.replicas",
	}
	tidbClusterTiDBColumn = extensionsobj.CustomResourceColumnDefinition{
		Name:        "TiDB",
		Type:        "string",
		Description: "The image for TiDB cluster",
		JSONPath:    ".status.tidb.image",
	}
	tidbClusterTiDBReadyColumn = extensionsobj.CustomResourceColumnDefinition{
		Name:        "Ready",
		Type:        "integer",
		Description: "The ready replicas number of TiDB cluster",
		JSONPath:    ".status.tidb.statefulSet.readyReplicas",
	}
	tidbClusterTiDBDesireColumn = extensionsobj.CustomResourceColumnDefinition{
		Name:        "Desire",
		Type:        "integer",
		Description: "The desired replicas number of TiDB cluster",
		JSONPath:    ".spec.tidb.replicas",
	}
	dmClusteradditionalPrinterColumns []extensionsobj.CustomResourceColumnDefinition
	dmClusterReadyColumn              = extensionsobj.CustomResourceColumnDefinition{
		Name:     "Ready",
		Type:     "string",
		JSONPath: `.status.conditions[?(@.type=="Ready")].status`,
	}
	dmClusterStatusMessageColumn = extensionsobj.CustomResourceColumnDefinition{
		Name:     "Status",
		Type:     "string",
		JSONPath: `.status.conditions[?(@.type=="Ready")].message`,
		Priority: 1,
	}
	dmClusterMasterColumn = extensionsobj.CustomResourceColumnDefinition{
		Name:        "Master",
		Type:        "string",
		Description: "The image for dm-master cluster",
		JSONPath:    ".status.master.image",
	}
	dmClusterMasterStorageColumn = extensionsobj.CustomResourceColumnDefinition{
		Name:        "Storage",
		Type:        "string",
		Description: "The storage size specified for dm-master node",
		JSONPath:    ".spec.master.storageSize",
	}
	dmClusterMasterReadyColumn = extensionsobj.CustomResourceColumnDefinition{
		Name:        "Ready",
		Type:        "integer",
		Description: "The desired replicas number of dm-master cluster",
		JSONPath:    ".status.master.statefulSet.readyReplicas",
	}
	dmClusterMasterDesireColumn = extensionsobj.CustomResourceColumnDefinition{
		Name:        "Desire",
		Type:        "integer",
		Description: "The desired replicas number of dm-master cluster",
		JSONPath:    ".spec.master.replicas",
	}
	dmClusterWorkerColumn = extensionsobj.CustomResourceColumnDefinition{
		Name:        "Worker",
		Type:        "string",
		Description: "The image for dm-worker cluster",
		JSONPath:    ".status.worker.image",
	}
	dmClusterWorkerStorageColumn = extensionsobj.CustomResourceColumnDefinition{
		Name:        "Storage",
		Type:        "string",
		Description: "The storage size specified for dm-worker node",
		JSONPath:    ".spec.worker.storageSize",
	}
	dmClusterWorkerReadyColumn = extensionsobj.CustomResourceColumnDefinition{
		Name:        "Ready",
		Type:        "integer",
		Description: "The ready replicas number of dm-worker cluster",
		JSONPath:    ".status.worker.statefulSet.readyReplicas",
	}
	dmClusterWorkerDesireColumn = extensionsobj.CustomResourceColumnDefinition{
		Name:        "Desire",
		Type:        "integer",
		Description: "The desired replicas number of dm-worker cluster",
		JSONPath:    ".spec.worker.replicas",
	}
	backupAdditionalPrinterColumns []extensionsobj.CustomResourceColumnDefinition
	backupStatusColumn             = extensionsobj.CustomResourceColumnDefinition{
		Name:        "Status",
		Type:        "string",
		Description: "The current status of the backup",
		JSONPath:    ".status.phase",
	}
	backupPathColumn = extensionsobj.CustomResourceColumnDefinition{
		Name:        "BackupPath",
		Type:        "string",
		Description: "The full path of backup data",
		JSONPath:    ".status.backupPath",
	}
	backupBackupSizeColumn = extensionsobj.CustomResourceColumnDefinition{
		Name:        "BackupSize",
		Type:        "string",
		Description: "The data size of the backup",
		JSONPath:    ".status.backupSizeReadable",
	}
	backupCommitTSColumn = extensionsobj.CustomResourceColumnDefinition{
		Name:        "CommitTS",
		Type:        "string",
		Description: "The commit ts of tidb cluster dump",
		JSONPath:    ".status.commitTs",
	}
	backupStartedColumn = extensionsobj.CustomResourceColumnDefinition{
		Name:        "Started",
		Type:        "string",
		Format:      "date-time",
		Description: "The time at which the backup was started",
		Priority:    1,
		JSONPath:    ".status.timeStarted",
	}
	backupCompletedColumn = extensionsobj.CustomResourceColumnDefinition{
		Name:        "Completed",
		Type:        "string",
		Format:      "date-time",
		Description: "The time at which the backup was completed",
		Priority:    1,
		JSONPath:    ".status.timeCompleted",
	}
	restoreAdditionalPrinterColumns []extensionsobj.CustomResourceColumnDefinition
	restoreStatusColumn             = extensionsobj.CustomResourceColumnDefinition{
		Name:        "Status",
		Type:        "string",
		Description: "The current status of the restore",
		JSONPath:    ".status.phase",
	}
	restoreStartedColumn = extensionsobj.CustomResourceColumnDefinition{
		Name:        "Started",
		Type:        "string",
		Format:      "date-time",
		Description: "The time at which the backup was started",
		JSONPath:    ".status.timeStarted",
	}
	restoreCompletedColumn = extensionsobj.CustomResourceColumnDefinition{
		Name:        "Completed",
		Type:        "string",
		Format:      "date-time",
		Description: "The time at which the restore was completed",
		JSONPath:    ".status.timeCompleted",
	}
	restoreCommitTSColumn = extensionsobj.CustomResourceColumnDefinition{
		Name:        "CommitTS",
		Type:        "string",
		Description: "The commit ts of tidb cluster restore",
		JSONPath:    ".status.commitTs",
	}
	bksAdditionalPrinterColumns []extensionsobj.CustomResourceColumnDefinition
	bksScheduleColumn           = extensionsobj.CustomResourceColumnDefinition{
		Name:        "Schedule",
		Type:        "string",
		Description: "The cron format string used for backup scheduling.",
		JSONPath:    ".spec.schedule",
	}
	bksMaxBackups = extensionsobj.CustomResourceColumnDefinition{
		Name:        "MaxBackups",
		Type:        "integer",
		Description: "The max number of backups we want to keep.",
		JSONPath:    ".spec.maxBackups",
	}
	bksLastBackup = extensionsobj.CustomResourceColumnDefinition{
		Name:        "LastBackup",
		Type:        "string",
		Description: "The last backup CR name",
		Priority:    1,
		JSONPath:    ".status.lastBackup",
	}
	bksLastBackupTime = extensionsobj.CustomResourceColumnDefinition{
		Name:        "LastBackupTime",
		Type:        "date",
		Description: "The last time the backup was successfully created",
		Priority:    1,
		JSONPath:    ".status.lastBackupTime",
	}
	tidbInitializerPrinterColumns []extensionsobj.CustomResourceColumnDefinition
	tidbInitializerPhase          = extensionsobj.CustomResourceColumnDefinition{
		Name:        "Phase",
		Type:        "string",
		Description: "The current phase of initialization",
		Priority:    1,
		JSONPath:    ".status.phase",
	}
	autoScalerPrinterColumns []extensionsobj.CustomResourceColumnDefinition
	// TODO add The current replicas number of TiKV cluster
	autoScalerTiKVMaxReplicasColumn = extensionsobj.CustomResourceColumnDefinition{
		Name:        "TiKV-MaxReplicas",
		Type:        "integer",
		Description: "The maximal replicas of TiKV",
		JSONPath:    ".spec.tikv.maxReplicas",
	}
	autoScalerTiKVMinReplicasColumn = extensionsobj.CustomResourceColumnDefinition{
		Name:        "TiKV-MinReplicas",
		Type:        "integer",
		Description: "The minimal replicas of TiKV",
		JSONPath:    ".spec.tikv.minReplicas",
	}
	// TODO add The current replicas number of TiDB cluster
	autoScalerTiDBMaxReplicasColumn = extensionsobj.CustomResourceColumnDefinition{
		Name:        "TiDB-MaxReplicas",
		Type:        "integer",
		Description: "The maximal replicas of TiDB",
		JSONPath:    ".spec.tidb.maxReplicas",
	}
	autoScalerTiDBMinReplicasColumn = extensionsobj.CustomResourceColumnDefinition{
		Name:        "TiDB-MinReplicas",
		Type:        "integer",
		Description: "The minimal replicas of TiDB",
		JSONPath:    ".spec.tidb.minReplicas",
	}
	ageColumn = extensionsobj.CustomResourceColumnDefinition{
		Name:     "Age",
		Type:     "date",
		JSONPath: ".metadata.creationTimestamp",
	}
	tidbMonitorAdditionalPrinterColumns []extensionsobj.CustomResourceColumnDefinition

	tidbMonitorDesiredColumn = extensionsobj.CustomResourceColumnDefinition{
		Name:        "DESIRED",
		Type:        "integer",
		Description: "The desired replicas number of tidbMonitor cluster",
		JSONPath:    ".spec.replicas",
	}
	tidbMonitorReadyColumn = extensionsobj.CustomResourceColumnDefinition{
		Name:     "READY",
		Type:     "integer",
		JSONPath: ".status.statefulSet.readyReplicas",
	}
	tidbMonitorUpdatedColumn = extensionsobj.CustomResourceColumnDefinition{
		Name:     "UP-TO-DATE",
		Type:     "integer",
		JSONPath: ".status.statefulSet.updatedReplicas",
	}
)

func init() {
	tidbClusteradditionalPrinterColumns = append(tidbClusteradditionalPrinterColumns,
		tidbClusterReadyColumn,
		tidbClusterPDColumn, tidbClusterPDStorageColumn, tidbClusterPDReadyColumn, tidbClusterPDDesireColumn,
		tidbClusterTiKVColumn, tidbClusterTiKVStorageColumn, tidbClusterTiKVReadyColumn, tidbClusterTiKVDesireColumn,
		tidbClusterTiDBColumn, tidbClusterTiDBReadyColumn, tidbClusterTiDBDesireColumn, tidbClusterStatusMessageColumn, ageColumn)
	dmClusteradditionalPrinterColumns = append(dmClusteradditionalPrinterColumns,
		dmClusterReadyColumn,
		dmClusterMasterColumn, dmClusterMasterStorageColumn, dmClusterMasterReadyColumn, dmClusterMasterDesireColumn,
		dmClusterWorkerColumn, dmClusterWorkerStorageColumn, dmClusterWorkerReadyColumn, dmClusterWorkerDesireColumn,
		dmClusterStatusMessageColumn, ageColumn)
	backupAdditionalPrinterColumns = append(backupAdditionalPrinterColumns, backupStatusColumn, backupPathColumn, backupBackupSizeColumn, backupCommitTSColumn, backupStartedColumn, backupCompletedColumn, ageColumn)
	restoreAdditionalPrinterColumns = append(restoreAdditionalPrinterColumns, restoreStatusColumn, restoreStartedColumn, restoreCompletedColumn, restoreCommitTSColumn, ageColumn)
	bksAdditionalPrinterColumns = append(bksAdditionalPrinterColumns, bksScheduleColumn, bksMaxBackups, bksLastBackup, bksLastBackupTime, ageColumn)
	tidbInitializerPrinterColumns = append(tidbInitializerPrinterColumns, tidbInitializerPhase, ageColumn)
	autoScalerPrinterColumns = append(autoScalerPrinterColumns, autoScalerTiDBMaxReplicasColumn, autoScalerTiDBMinReplicasColumn,
		autoScalerTiKVMaxReplicasColumn, autoScalerTiKVMinReplicasColumn, ageColumn)
	tidbMonitorAdditionalPrinterColumns = append(tidbMonitorAdditionalPrinterColumns, tidbMonitorDesiredColumn, tidbMonitorReadyColumn, tidbMonitorUpdatedColumn, ageColumn)
}

func NewCustomResourceDefinition(crdKind v1alpha1.CrdKind, group string, labels map[string]string, validation bool) *extensionsobj.CustomResourceDefinition {
	crd := crdutils.NewCustomResourceDefinition(crdutils.Config{
		SpecDefinitionName:    crdKind.SpecName,
		EnableValidation:      validation,
		Labels:                crdutils.Labels{LabelsMap: labels},
		ResourceScope:         string(extensionsobj.NamespaceScoped),
		Group:                 group,
		Kind:                  crdKind.Kind,
		Version:               v1alpha1.Version,
		Plural:                crdKind.Plural,
		ShortNames:            crdKind.ShortNames,
		GetOpenAPIDefinitions: v1alpha1.GetOpenAPIDefinitions,
	})
	addAdditionalPrinterColumnsForCRD(crd, crdKind)
	return crd
}

func GetCrdKindFromKindName(kindName string) (v1alpha1.CrdKind, error) {
	switch strings.ToLower(kindName) {
	case v1alpha1.TiDBClusterKindKey:
		return v1alpha1.DefaultCrdKinds.TiDBCluster, nil
	case v1alpha1.DMClusterKindKey:
		return v1alpha1.DefaultCrdKinds.DMCluster, nil
	case v1alpha1.BackupKindKey:
		return v1alpha1.DefaultCrdKinds.Backup, nil
	case v1alpha1.RestoreKindKey:
		return v1alpha1.DefaultCrdKinds.Restore, nil
	case v1alpha1.BackupScheduleKindKey:
		return v1alpha1.DefaultCrdKinds.BackupSchedule, nil
	case v1alpha1.TiDBMonitorKindKey:
		return v1alpha1.DefaultCrdKinds.TiDBMonitor, nil
	case v1alpha1.TiDBInitializerKindKey:
		return v1alpha1.DefaultCrdKinds.TiDBInitializer, nil
	case v1alpha1.TidbClusterAutoScalerKindKey:
		return v1alpha1.DefaultCrdKinds.TidbClusterAutoScaler, nil
	default:
		return v1alpha1.CrdKind{}, errors.New("unknown CrdKind Name")
	}
}

func addAdditionalPrinterColumnsForCRD(crd *extensionsobj.CustomResourceDefinition, crdKind v1alpha1.CrdKind) {
	switch crdKind.Kind {
	case v1alpha1.DefaultCrdKinds.TiDBCluster.Kind:
		crd.Spec.AdditionalPrinterColumns = tidbClusteradditionalPrinterColumns
	case v1alpha1.DefaultCrdKinds.DMCluster.Kind:
		crd.Spec.AdditionalPrinterColumns = dmClusteradditionalPrinterColumns
	case v1alpha1.DefaultCrdKinds.Backup.Kind:
		crd.Spec.AdditionalPrinterColumns = backupAdditionalPrinterColumns
	case v1alpha1.DefaultCrdKinds.Restore.Kind:
		crd.Spec.AdditionalPrinterColumns = restoreAdditionalPrinterColumns
	case v1alpha1.DefaultCrdKinds.BackupSchedule.Kind:
		crd.Spec.AdditionalPrinterColumns = bksAdditionalPrinterColumns
	case v1alpha1.DefaultCrdKinds.TiDBMonitor.Kind:
		crd.Spec.AdditionalPrinterColumns = tidbMonitorAdditionalPrinterColumns
	case v1alpha1.DefaultCrdKinds.TiDBInitializer.Kind:
		crd.Spec.AdditionalPrinterColumns = tidbInitializerPrinterColumns
	case v1alpha1.DefaultCrdKinds.TidbClusterAutoScaler.Kind:
		crd.Spec.AdditionalPrinterColumns = autoScalerPrinterColumns
	default:
	}
}
