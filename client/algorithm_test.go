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
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAlgorithmInfoForAll(t *testing.T) {
	assert := assert.New(t)

	c, err := NewAlgorithmClient()
	assert.NoError(err)

	s, err := c.GetInfoForAll()
	assert.NoError(err)

	obj := map[string]interface{}{}
	err = json.Unmarshal([]byte(s), &obj)
	assert.NoError(err)

	obj2 := obj["algorithms"].([]interface{})

	found := false
	for _, v := range obj2 {
		obj3 := v.(map[string]interface{})["name"]
		if obj3 == "NDWI_PY" {
			found = true
			break
		}
	}
	assert.True(found)
}

func TestAlgorithmInfoForOne(t *testing.T) {
	assert := assert.New(t)

	c, err := NewAlgorithmClient()
	assert.NoError(err)

	// TODO: get id from a getall call(), not hard-coded here
	s, err := c.GetInfoForOne("f64d4845-0b9d-4bf1-8d49-45bafd639875")
	assert.NoError(err)

	obj := map[string]interface{}{}
	err = json.Unmarshal([]byte(s), &obj)
	assert.NoError(err)
	log.Print(s)

	obj2 := obj["algorithms"].([]interface{})

	found := false
	for _, v := range obj2 {
		obj3 := v.(map[string]interface{})["name"]
		if obj3 == "NDWI_PY" {
			found = true
			break
		}
	}
	assert.True(found)
}
