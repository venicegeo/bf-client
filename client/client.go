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

import "time"

const apiServer = "bf-api"

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

//---------------------------------------------------------------------
