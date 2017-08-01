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

	"gopkg.in/urfave/cli.v1"
)

type JobClient struct {
	url  string
	auth string
}

func NewJobClient() (*JobClient, error) {

	fields, err := ReadBeachfrontrcFields([]string{"domain", "auth"})
	if err != nil {
		return nil, err
	}

	url := "https://" + apiServer + "." + fields["domain"]

	return &JobClient{
		url:  url,
		auth: fields["auth"],
	}, nil
}

//---------------------------------------------------------------------

func (c *JobClient) GetJobInfoForJobs() (string, error) {

	path := "/v0/job"
	url := fmt.Sprintf("%s%s", c.url, path)

	responseBody, err := doHttpGetJSONWithAuth(url, c.auth, 200)
	if err != nil {
		return "", err
	}

	//log.Print(responseBody)

	return responseBody, nil
}

func (c *JobClient) GetJobInfoForJob(id string) (string, error) {

	path := "/v0/job"
	url := fmt.Sprintf("%s%s/%s", c.url, path, id)

	responseBody, err := doHttpGetJSONWithAuth(url, c.auth, 200)
	if err != nil {
		return "", err
	}

	//log.Print(responseBody)

	return responseBody, nil
}

func (c *JobClient) DoJobSubmit() error {
	return cli.NewExitError("job: --submit not yet supported", 2)
}

func (c *JobClient) DoJobDelete(id string) error {
	return cli.NewExitError("job: --delete not yet supported", 2)
}
