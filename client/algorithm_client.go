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
	"log"
)

type AlgorithmClient struct {
	url  string
	auth string
}

type Algorithms struct {
	Algorithms []*AlgorithmInfo
}

type Algorithm struct {
	Algorithm *AlgorithmInfo
}

type AlgorithmInfo struct {
	Description   string
	Interface     string
	MaxCloudCover int `json:"max_cloud_cover"`
	Name          string
	ServiceId     string `json:"service_id"`
	Version       string
}

func (a *Algorithms) String() string {
	s := ""
	for _, v := range a.Algorithms {
		s += v.String() + "\n"
	}
	return s
}

func (a *Algorithm) String() string {
	return a.Algorithm.String()
}

func (a *AlgorithmInfo) String() string {
	return fmt.Sprintf("[algorithm %s]", a.ServiceId)
}

//---------------------------------------------------------------------

func NewAlgorithmClient() (*AlgorithmClient, error) {

	fields, err := ReadBeachfrontrcFields([]string{"domain", "auth"})
	if err != nil {
		return nil, err
	}

	url := "https://" + apiServer + "." + fields["domain"]

	return &AlgorithmClient{
		url:  url,
		auth: fields["auth"],
	}, nil
}

//---------------------------------------------------------------------

func (c *AlgorithmClient) GetInfoForAll() (string, error) {

	log.Print("Algorithm.GetInfoForAll")
	path := "/v0/algorithm"
	url := fmt.Sprintf("%s%s", c.url, path)

	jsn, err := doHttpGetJSONWithAuth(url, c.auth, 200)
	if err != nil {
		return "", err
	}

	obj := &Algorithms{}
	err = json.Unmarshal([]byte(jsn), obj)
	if err != nil {
		return "", err
	}

	return obj.String(), nil
}

func (c *AlgorithmClient) GetInfoForOne(id string) (string, error) {

	log.Print("Algorithm.GetInfoForOne")

	path := "/v0/algorithm"
	url := fmt.Sprintf("%s%s/%s", c.url, path, id)

	jsn, err := doHttpGetJSONWithAuth(url, c.auth, 200)
	if err != nil {
		return "", err
	}

	obj := &Algorithm{}
	err = json.Unmarshal([]byte(jsn), obj)
	if err != nil {
		return "", err
	}

	return obj.String(), nil
}
