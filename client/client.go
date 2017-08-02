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

const apiServer = "bf-api"

type Client struct {
	Catalog   *CatalogClient
	Job       *JobClient
	Coastline *CoastlineClient
	Algorithm *AlgorithmClient
}

func NewClient() (*Client, error) {

	var err error
	c := &Client{}

	c.Catalog, err = NewCatalogClient()
	if err != nil {
		return nil, err
	}

	c.Job, err = NewJobClient()
	if err != nil {
		return nil, err
	}

	c.Coastline, err = NewCoastlineClient()
	if err != nil {
		return nil, err
	}

	c.Algorithm, err = NewAlgorithmClient()
	if err != nil {
		return nil, err
	}

	return c, nil
}

//---------------------------------------------------------------------
