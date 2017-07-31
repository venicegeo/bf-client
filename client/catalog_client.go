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

const catalogServer = "bf-ia-broker"

const catalogTimeout = time.Duration(2 * time.Minute)

type CatalogClient struct {
	url      string // "https://<bf-ia-broker>.<int.geointservices.io>"
	keyParam string // "PL_API_KEY=<123abc>"
}

func NewCatalogClient() (*CatalogClient, error) {

	fields, err := ReadBeachfrontrcFields([]string{"domain", "planet_key"})
	if err != nil {
		return nil, err
	}

	keyParam := "PL_API_KEY=" + fields["planet_key"]

	return &CatalogClient{
		url:      "https://" + catalogServer + "." + fields["domain"],
		keyParam: keyParam,
	}, nil
}

//---------------------------------------------------------------------

type sceneInfo struct {
	Properties propertiesInfo
}

type propertiesInfo struct {
	Bands map[string]string
}

//---------------------------------------------------------------------

// landsat:LC80480102017209LGN00 -> (landsat, LC80480102017209LGN00)
func splitId(id string) (string, string, error) {

	ids := strings.Split(id, ":")
	if len(ids) != 2 || ids[0] == "" || ids[1] == "" {
		return "", "", fmt.Errorf("scene id format must be '<catalogname>:<sceneid>'")
	}

	return ids[0], ids[1], nil
}

//---------------------------------------------------------------------

func (c *CatalogClient) GetCatalogInfoForCatalogs() (string, error) {

	log.Printf("GetCatalogInfoForCatalogs")

	return "", cli.NewExitError("catalog: --info for catalogs not yet supported", 2)
}

func (c *CatalogClient) GetCatalogInfoForScene(id string) (string, error) {

	log.Printf("GetCatalogInfoForScene")

	sensor, scene, err := splitId(id)
	if err != nil {
		return "", err
	}

	path := "/planet/" + sensor + "/" + scene

	params := c.keyParam

	url := fmt.Sprintf("%s%s?%s", c.url, path, params)

	return doHttpGetJSON(url, catalogTimeout, 200)
}

func (c *CatalogClient) GetCatalogInfoForCatalog(id string) (string, error) {

	log.Printf("GetCatalogInfoForCatalog")

	path := "/planet/discover/" + id

	params := c.keyParam

	url := fmt.Sprintf("%s%s?%s", c.url, path, params)

	return doHttpGetJSON(url, catalogTimeout, 200)
}

// returns map: file name -> file size
func (c *CatalogClient) DoCatalogSceneDownload(id string) (map[string]int, error) {

	log.Printf("DoCatalogSceneDownload")

	body, err := c.GetCatalogInfoForScene(id)

	info := sceneInfo{}
	err = json.Unmarshal([]byte(body), &info)
	if err != nil {
		return nil, err
	}

	fileInfo := map[string]int{}

	i := 1
	for bandName, value := range info.Properties.Bands {

		status, byts, err := doHttpGetBytes(value, catalogTimeout)
		if err != nil {
			return nil, err
		}
		if status != 200 {
			return nil, fmt.Errorf("HTTP download failed with status %d", status)
		}

		idx := strings.LastIndex(value, "/")
		if idx == -1 {
			return nil, fmt.Errorf("unable to parse URL path: %s", value)
		}
		filename := value[idx+1 : len(value)]

		// TODO: fix path root
		err = ioutil.WriteFile("./"+filename, byts, 0600)
		if err != nil {
			return nil, err
		}

		fileInfo[filename] = len(byts)

		log.Printf("%d/%d: %s\n", i, len(info.Properties.Bands), bandName)

		i++
	}

	return fileInfo, nil
}
