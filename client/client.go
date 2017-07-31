/* Copyright 2017, RadiantBlue Technologies, Inc.

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

package client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"time"

	"gopkg.in/urfave/cli.v1"
)

type BFClient struct {
	Timeout time.Duration

	bfDomain  string // "int.geointservices.io"
	bfAuth    string
	planetKey string
}

func NewBFClient() (*BFClient, error) {

	fields, err := ReadBeachfrontrcFields([]string{"domain", "auth", "planet_key"})
	if err != nil {
		return nil, err
	}

	return &BFClient{
		Timeout:   2 * time.Minute,
		bfDomain:  fields["domain"],
		bfAuth:    fields["auth"],
		planetKey: fields["planet_key"],
	}, nil
}

// TODO: remove the need for this
type serverType string

const (
	apiServer    = "bf-api"
	brokerServer = "bf-ia-broker"
)

//---------------------------------------------------------------------

func (c *BFClient) GetCatalogInfoForCatalogs() error {

	//log.Printf("GetCatalogInfoForCatalogs")

	return cli.NewExitError("catalog: --info for catalogs not yet supported", 2)
}

func (c *BFClient) GetCatalogInfoForScene(id string) (string, error) {

	//log.Printf("GetCatalogInfoForScene")

	ids := strings.Split(id, ":")
	if len(ids) != 2 || ids[0] == "" || ids[1] == "" {
		return "", fmt.Errorf("scene id format must be '<catalogname>:<sceneid>'")
	}

	path := "/planet/" + ids[0] + "/" + ids[1]

	params := "PL_API_KEY=" + c.planetKey

	url := fmt.Sprintf("https://%s.%s%s?%s", brokerServer, c.bfDomain, path, params)

	status, responseBody, err := doHttpGetJSON(url, c.Timeout)
	if err != nil {
		return "", err
	}
	if status != 200 {
		return "", fmt.Errorf("HTTP request failed with status %d", status)
	}

	log.Print(responseBody)

	return string(responseBody), nil
}

func (c *BFClient) GetCatalogInfoForCatalog(id string) error {

	//log.Printf("GetCatalogInfoForCatalog")

	path := "/planet/discover/" + id

	params := "PL_API_KEY=" + c.planetKey

	url := fmt.Sprintf("https://%s.%s%s?%s", brokerServer, c.bfDomain, path, params)

	status, responseBody, err := doHttpGetJSON(url, c.Timeout)
	if err != nil {
		return err
	}
	if status != 200 {
		return fmt.Errorf("HTTP request failed with status %d", status)
	}

	log.Print(responseBody)

	return nil
}

func (c *BFClient) DoCatalogSceneDownload(id string) error {

	log.Printf("DoCatalogSceneDownload")

	body, err := c.GetCatalogInfoForScene(id)

	obj := map[string]interface{}{}
	err = json.Unmarshal([]byte(body), &obj)
	if err != nil {
		return err
	}

	propsX, ok := obj["properties"]
	if !ok {
		return fmt.Errorf("failed to parse scene properties")
	}
	props, ok := propsX.(map[string]interface{})
	if !ok {
		return fmt.Errorf("failed to parse scene properties")
	}
	bandsX, ok := props["bands"]
	if !ok {
		return fmt.Errorf("failed to parse scene bands1")
	}
	bands, ok := bandsX.(map[string]interface{})
	if !ok {
		return fmt.Errorf("failed to parse scene bands")
	}

	for _, valueX := range bands {
		value, ok := valueX.(string)
		if !ok {
			return fmt.Errorf("failed to parse band value")
		}

		status, byts, err := doHttpGetBytes(value, c.Timeout)
		if err != nil {
			return err
		}
		if status != 200 {
			return fmt.Errorf("HTTP download failed with status %d", status)
		}

		idx := strings.LastIndex(value, "/")
		if idx == -1 {
			return fmt.Errorf("unable to parse URL path: %s", value)
		}
		filename := value[idx+1 : len(value)]
		log.Print(filename)

		// TODO: fix path root
		err = ioutil.WriteFile("./"+filename, byts, 0600)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *BFClient) GetJobInfoForJobs() error {

	path := "/v0/job"
	url := fmt.Sprintf("https://%s.%s%s", apiServer, c.bfDomain, path)

	status, responseBody, err := doHttpGetJSONWithAuth(url, c.Timeout, c.bfAuth)
	if err != nil {
		return err
	}
	if status != 200 {
		return fmt.Errorf("HTTP request failed with status %d", status)
	}

	log.Print(responseBody)

	return nil
}

func (c *BFClient) GetJobInfoForJob(id string) error {

	path := "/v0/job"
	url := fmt.Sprintf("https://%s.%s%s/%s", apiServer, c.bfDomain, path, id)

	status, responseBody, err := doHttpGetJSONWithAuth(url, c.Timeout, c.bfAuth)
	if err != nil {
		return err
	}
	if status != 200 {
		return fmt.Errorf("HTTP request failed with status %d", status)
	}

	log.Print(responseBody)

	return nil
}

func (c *BFClient) DoJobSubmit() error {
	return cli.NewExitError("job: --submit not yet supported", 2)
}

func (c *BFClient) DoJobDelete(id string) error {
	return cli.NewExitError("job: --delete not yet supported", 2)
}

func (c *BFClient) DoCoastlineDownload(id string) error {

	path := "/v0/job"
	url := fmt.Sprintf("https://%s.%s%s/%s.geojson", apiServer, c.bfDomain, path, id)

	status, responseBody, err := doHttpGetJSONWithAuth(url, c.Timeout, c.bfAuth)
	if err != nil {
		return err
	}
	if status != 200 {
		return fmt.Errorf("HTTP request failed with status %d", status)
	}

	// TODO: fix root path
	err = ioutil.WriteFile("./"+id+".geojson", []byte(responseBody), 0600)
	if err != nil {
		return err
	}

	return nil
}

func (c *BFClient) GetAlgorithmInfoForAll() error {

	path := "/v0/algorithm"
	url := fmt.Sprintf("https://%s.%s%s", apiServer, c.bfDomain, path)

	status, responseBody, err := doHttpGetJSONWithAuth(url, c.Timeout, c.bfAuth)
	if err != nil {
		return err
	}
	if status != 200 {
		return fmt.Errorf("HTTP request failed with status %d", status)
	}

	log.Print(responseBody)

	return nil
}

func (c *BFClient) GetAlgorithmInfoForOne(id string) error {

	path := "/v0/algorithm"
	url := fmt.Sprintf("https://%s.%s%s/%s", apiServer, c.bfDomain, path, id)

	status, responseBody, err := doHttpGetJSONWithAuth(url, c.Timeout, c.bfAuth)
	if err != nil {
		return err
	}
	if status != 200 {
		return fmt.Errorf("HTTP request failed with status %d", status)
	}

	log.Print(responseBody)

	return nil
}
