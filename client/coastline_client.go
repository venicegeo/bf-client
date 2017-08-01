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
	"fmt"
	"io/ioutil"
	"log"
)

type CoastlineClient struct {
	url  string
	auth string
}

func NewCoastlineClient() (*CoastlineClient, error) {

	fields, err := ReadBeachfrontrcFields([]string{"domain", "auth"})
	if err != nil {
		return nil, err
	}

	url := "https://" + apiServer + "." + fields["domain"]

	return &CoastlineClient{
		url:  url,
		auth: fields["auth"],
	}, nil
}

//---------------------------------------------------------------------

func (c *CoastlineClient) DoCoastlineDownload(id string) (string, error) {

	log.Printf("DoCoastlineDownload")

	path := "/v0/job"
	url := fmt.Sprintf("%s%s/%s.geojson", c.url, path, id)

	responseBody, err := doHttpGetJSONWithAuth(url, c.auth, 200)
	if err != nil {
		return "", err
	}

	// TODO: fix root path
	err = ioutil.WriteFile("./"+id+".geojson", []byte(responseBody), 0600)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("Wrote %d bytes of geojson\n", len(responseBody)), nil
}
