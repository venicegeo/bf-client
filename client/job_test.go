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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJobInfoForJobs(t *testing.T) {
	assert := assert.New(t)

	c, err := NewJobClient()
	assert.NoError(err)

	s, err := c.GetJobInfoForJobs()
	assert.NoError(err)

	obj := map[string]interface{}{}
	err = json.Unmarshal([]byte(s), &obj)
	assert.NoError(err)

	obj2 := obj["jobs"].(map[string]interface{})
	obj3 := obj2["features"].([]interface{})
	obj4 := obj3[0].(map[string]interface{})

	assert.Equal("Feature", obj4["type"])
	assert.NotEqual("", obj4["id"])
}

func TestJobInfoForJob(t *testing.T) {
	assert := assert.New(t)

	const id = "6f475d16-d8f1-4f15-9dda-cf9b2d502241"

	c, err := NewJobClient()
	assert.NoError(err)

	s, err := c.GetJobInfoForJob(id)
	assert.NoError(err)

	obj := map[string]interface{}{}
	err = json.Unmarshal([]byte(s), &obj)
	assert.NoError(err)

	obj2 := obj["job"].(map[string]interface{})

	assert.Equal(id, obj2["id"])
}

func TestJobSubmit(t *testing.T) {
	assert := assert.New(t)

	c, err := NewJobClient()
	assert.NoError(err)

	err = c.DoJobSubmit()
	assert.Error(err)
}

func TestJobDelete(t *testing.T) {
	assert := assert.New(t)

	c, err := NewJobClient()
	assert.NoError(err)

	err = c.DoJobDelete("")
	assert.Error(err)
}
