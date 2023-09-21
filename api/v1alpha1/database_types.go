/*
Copyright 2022-2023 Nutanix, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

/*
GENERATED by operator-sdk
Changes added
*/

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// DatabaseSpec defines the desired state of Database
type DatabaseSpec struct {
	NDB      NDB      `json:"ndb"`
	Instance Instance `json:"databaseInstance"`
}

// DatabaseStatus defines the observed state of Database
type DatabaseStatus struct {
	IPAddress               string `json:"ipAddress"`
	Id                      string `json:"id"`
	Status                  string `json:"status"`
	DatabaseServerId        string `json:"dbServerId"`
	Type                    string `json:"type"`
	ProvisioningOperationId string `json:"provisioningOperationId"`
}

// Database is the Schema for the databases API
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:shortName={"db","dbs"}
// +kubebuilder:printcolumn:name="IP Address",type=string,JSONPath=`.status.ipAddress`
// +kubebuilder:printcolumn:name="Status",type=string,JSONPath=`.status.status`
// +kubebuilder:printcolumn:name="Type",type=string,JSONPath=`.status.type`
type Database struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DatabaseSpec   `json:"spec,omitempty"`
	Status DatabaseStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true
// DatabaseList contains a list of Database
type DatabaseList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Database `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Database{}, &DatabaseList{})
}

// These are required to have a deep copy, object interface implementation
// These are the structs for the Spec and Status

// Details of the NDB installation
type NDB struct {
	ClusterId string `json:"clusterId"`
	// Name of the secret holding the credentials for NDB (username and password)
	CredentialSecret string `json:"credentialSecret"`
	// NDB Server URL
	Server string `json:"server"`
	// Skip server's certificate and hostname verification
	SkipCertificateVerification bool `json:"skipCertificateVerification"`
}

// Database instance specific details
type Instance struct {
	// Name of the database instance
	DatabaseInstanceName string `json:"databaseInstanceName"`
	// Description of the database instance
	// +optional
	Description string `json:"description"`
	// Name(s) of the database(s) to be provisiond inside the database instance
	// default [ "database_one", "database_two", "database_three" ]
	// +optional
	DatabaseNames []string `json:"databaseNames"`
	// Name of the secret holding the credentials for the database instance (password and ssh key)
	CredentialSecret string `json:"credentialSecret"`
	// Size of the database instance, minimum 10 (GBs)
	Size int `json:"size"`
	// default UTC
	// +optional
	TimeZone string `json:"timezone"`
	Type     string `json:"type"`
	// +optional
	Profiles *Profiles `json:"profiles"`
	// +optional
	// Information related to time machine that is to be associated with this database
	TMInfo *DBTimeMachineInfo `json:"timeMachine"`
	// +optional
	// Additional database engine specific action arguments
	TypeDetails map[string]string `json:"typeDetails"`
}

// Time Machine details
type DBTimeMachineInfo struct {
	// +optional
	Name string `json:"name"`
	// +optional
	Description string `json:"description"`
	// +optional
	// Name of the SLA to be used, default NONE
	SLAName string `json:"sla"`
	// +optional
	// Daily snapshot time in HH:MM:SS (24 hour format)
	DailySnapshotTime string `json:"dailySnapshotTime"`
	// +optional
	// Number of snapshots per day
	SnapshotsPerDay int `json:"snapshotsPerDay"`
	// +optional
	// Log catch up frequency in minutes
	LogCatchUpFrequency int `json:"logCatchUpFrequency"`
	// +optional
	// Day of the week for weekly snapshot
	WeeklySnapshotDay string `json:"weeklySnapshotDay"`
	// +optional
	// Day of the month for monthly snapshot
	MonthlySnapshotDay int `json:"monthlySnapshotDay"`
	// +optional
	// Start month for the quarterly snapshot
	// Jan => Jan, Apr, Jul, Oct.
	// Feb => Feb, May, Aug, Nov.
	// Mar => Mar, Jun, Sep, Dec.
	QuarterlySnapshotMonth string `json:"quarterlySnapshotMonth"`
}

type Profiles struct {
	// +optional
	Software Profile `json:"software"`
	// +optional
	Compute Profile `json:"compute"`
	// +optional
	Network Profile `json:"network"`
	// +optional
	DbParam Profile `json:"dbParam"`
	// +optional
	DbParamInstance Profile `json:"dbParamInstance"`
}

type Profile struct {
	// +optional
	Id string `json:"id"`
	// +optional
	Name string `json:"name"`
}
