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

package ndb_api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/nutanix-cloud-native/ndb-operator/ndb_client"
	ctrllog "sigs.k8s.io/controller-runtime/pkg/log"
)

// Fetches all the databases on the NDB instance and retutns a slice of the databases
func GetAllDatabases(ctx context.Context, ndbClient *ndb_client.NDBClient) (databases []DatabaseResponse, err error) {
	log := ctrllog.FromContext(ctx)
	log.Info("Entered ndb_api.GetAllDatabases")
	if ndbClient == nil {
		err = errors.New("nil reference: received nil reference for ndbClient")
		log.Error(err, "Received nil ndbClient reference")
		return
	}
	res, err := ndbClient.Get("databases?detailed=true")
	if err != nil || res == nil || res.StatusCode != http.StatusOK {
		if err == nil {
			if res != nil {
				err = fmt.Errorf("GET /databases responded with %d", res.StatusCode)
			} else {
				err = fmt.Errorf("GET /databases responded with a nil response")
			}
		}
		log.Error(err, "Error occurred fetching all databases")
		return
	}
	log.Info("GET /databases", "HTTP status code", res.StatusCode)
	body, err := io.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		log.Error(err, "Error occurred reading response.Body in GetAllDatabases")
		return
	}
	err = json.Unmarshal(body, &databases)
	if err != nil {
		log.Error(err, "Error occurred trying to unmarshal.")
		return
	}
	log.Info("Returning from ndb_api.GetAllDatabases")
	return
}

// Fetches and returns a database by an Id
func GetDatabaseById(ctx context.Context, ndbClient *ndb_client.NDBClient, id string) (database DatabaseResponse, err error) {
	log := ctrllog.FromContext(ctx)
	log.Info("Entered ndb_api.GetDatabaseById", "databaseId", id)
	if ndbClient == nil {
		err = errors.New("nil reference")
		log.Error(err, "Received nil ndbClient reference")
		return
	}
	// Checking if id is empty, this is necessary otherwise the request becomes a call to get all databases (/databases)
	if id == "" {
		err = fmt.Errorf("database id is empty")
		log.Error(err, "no database id provided")
		return
	}
	getDbDetailedPath := fmt.Sprintf("databases/%s?detailed=true", id)
	res, err := ndbClient.Get(getDbDetailedPath)
	if err != nil || res == nil || res.StatusCode != http.StatusOK {
		if err == nil {
			if res != nil {
				err = fmt.Errorf("GET %s responded with %d", getDbDetailedPath, res.StatusCode)
			} else {
				err = fmt.Errorf("GET %s responded with a nil response", getDbDetailedPath)
			}
		}
		log.Error(err, "Error occurred fetching database")
		return
	}
	log.Info(getDbDetailedPath, "HTTP status code", res.StatusCode)
	body, err := io.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		log.Error(err, "Error occurred reading response.Body in Get Database by ID")
		return
	}
	err = json.Unmarshal(body, &database)
	if err != nil {
		log.Error(err, "Error occurred trying to unmarshal.")
		return
	}
	log.Info("Returning from ndb_api.GetDatabaseById")
	return
}

// Fetches and returns a database by name
func GetDatabaseByName(ctx context.Context, ndbClient *ndb_client.NDBClient, name string) (database *DatabaseResponse, err error) {
	log := ctrllog.FromContext(ctx)
	log.Info("Entered ndb_api.GetDatabaseByName", "name", name)
	if ndbClient == nil {
		err = errors.New("nil reference")
		log.Error(err, "Received nil ndbClient reference")
		return
	}
	// Checking if id is empty, this is necessary otherwise the request becomes a call to get all databases (/databases)
	if name == "" {
		err = fmt.Errorf("database name is empty")
		log.Error(err, "no database name provided")
		return
	}
	getDbDetailedPath := fmt.Sprintf("databases/%s?value-type=name&detailed=true", name)
	res, err := ndbClient.Get(getDbDetailedPath)
	if err != nil || res == nil || res.StatusCode != http.StatusOK {
		if err == nil {
			if res != nil {
				err = fmt.Errorf("GET %s responded with %d", getDbDetailedPath, res.StatusCode)
			} else {
				err = fmt.Errorf("GET %s responded with a nil response", getDbDetailedPath)
			}
		}
		log.Error(err, "Error occurred fetching database")
		return
	}
	log.Info(getDbDetailedPath, "HTTP status code", res.StatusCode)
	body, err := io.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		log.Error(err, "Error occurred reading response.Body in Get Database by name")
		return
	}
	err = json.Unmarshal(body, &database)
	if err != nil {
		log.Error(err, "Error occurred trying to unmarshal.")
		return
	}
	log.Info("Returning from ndb_api.GetDatabaseByName")
	return
}

// Provisions a database instance based on the database provisioning request
// Returns the task info summary response for the operation
func ProvisionDatabase(ctx context.Context, ndbClient *ndb_client.NDBClient, req *DatabaseProvisionRequest) (task TaskInfoSummaryResponse, err error) {
	log := ctrllog.FromContext(ctx)
	log.Info("Entered ndb_api.ProvisionDatabase")
	if ndbClient == nil {
		err = errors.New("nil reference")
		log.Error(err, "Received nil ndbClient reference")
		return
	}
	res, err := ndbClient.Post("databases/provision", req)
	if err != nil || res == nil || res.StatusCode != http.StatusOK {
		if err == nil {
			if res != nil {
				err = fmt.Errorf("POST databases/provision responded with %d", res.StatusCode)
			} else {
				err = fmt.Errorf("POST databases/provision responded with nil response")
			}
		}
		body, readErr := io.ReadAll(res.Body)
		log.Error(err, fmt.Sprintf("Error occurred provisioning database: %s", string(body))) // added this temporarily as currently there is no error message as to why the request failed. This body generally contains an explanation
		if readErr != nil {
			log.Info(readErr.Error())
		}
		return
	}
	log.Info("POST databases/provision", "HTTP status code", res.StatusCode)
	body, err := io.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		log.Error(err, "Error occurred reading response in ProvisionDatabase")
		return
	}
	err = json.Unmarshal(body, &task)
	if err != nil {
		log.Error(err, "Error occurred trying to unmarshal.")
		return
	}
	log.Info("Returning from ndb_api.ProvisionDatabase")
	return
}

// Deprovisions a database instance given a database id
// Returns the task info summary response for the operation
func DeprovisionDatabase(ctx context.Context, ndbClient *ndb_client.NDBClient, id string, req DatabaseDeprovisionRequest) (task TaskInfoSummaryResponse, err error) {
	log := ctrllog.FromContext(ctx)
	log.Info("Entered ndb_api.DeprovisionDatabase", "databaseId", id)
	if ndbClient == nil {
		err = errors.New("nil reference")
		log.Error(err, "Received nil ndbClient reference")
		return
	}
	if id == "" {
		err = fmt.Errorf("id is empty")
		log.Error(err, "no database id provided")
		return
	}
	res, err := ndbClient.Delete("databases/"+id, req)
	if err != nil || res == nil || res.StatusCode != http.StatusOK {
		if err == nil {
			if res != nil {
				err = fmt.Errorf("DELETE /databases/%s responded with %d", id, res.StatusCode)
			} else {
				err = fmt.Errorf("DELETE /databases/%s responded with nil response", id)
			}
		}
		log.Error(err, "Error occurred deprovisioning database")
		return
	}
	log.Info("DELETE /databases/"+id, "HTTP status code", res.StatusCode)
	body, err := io.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		log.Error(err, "Error occurred reading response.Body")
		return
	}
	err = json.Unmarshal(body, &task)
	if err != nil {
		log.Error(err, "Error occurred trying to unmarshal.")
		return
	}
	log.Info("Returning from ndb_api.DeprovisionDatabase")
	return
}
