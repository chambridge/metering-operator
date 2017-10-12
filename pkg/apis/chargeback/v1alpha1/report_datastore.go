package v1alpha1

import (
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type ReportDataStoreList struct {
	meta.TypeMeta `json:",inline"`
	meta.ListMeta `json:"metadata,omitempty"`
	Items         []*ReportDataStore `json:"items"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type ReportDataStore struct {
	meta.TypeMeta   `json:",inline"`
	meta.ObjectMeta `json:"metadata,omitempty"`

	Spec      ReportDataStoreSpec `json:"spec"`
	TableName string              `json:"table_name"`
}

type ReportDataStoreSpec struct {
	Storage ReportDataStoreStorage `json:"storage"`
	Queries []string               `json:"queries"`
}

type ReportDataStoreStorage struct {
	Type   string `json:"type"`
	Format string `json:"format"`
	Bucket string `json:"bucket"`
	Prefix string `json:"prefix"`
}