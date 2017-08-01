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

func TestCatalogInfoForCatalogs(t *testing.T) {
	assert := assert.New(t)

	c, err := NewCatalogClient()
	assert.NoError(err)

	_, err = c.GetCatalogInfoForCatalogs()
	assert.Error(err)
}

func TestCatalogInfoForScene(t *testing.T) {
	assert := assert.New(t)

	c, err := NewCatalogClient()
	assert.NoError(err)

	s, err := c.GetCatalogInfoForScene("landsat:LC81260322017212LGN00")
	assert.NoError(err)

	obj := map[string]interface{}{}
	err = json.Unmarshal([]byte(s), &obj)
	assert.NoError(err)

	assert.Equal("Feature", obj["type"])
	assert.Equal("LC81260322017212LGN00", obj["id"])
}

func TestCatalogInfoForCatalog(t *testing.T) {
	assert := assert.New(t)

	c, err := NewCatalogClient()
	assert.NoError(err)

	s, err := c.GetCatalogInfoForCatalog("landsat")
	assert.NoError(err)

	obj := map[string]interface{}{}
	err = json.Unmarshal([]byte(s), &obj)
	assert.NoError(err)

	obj2 := obj["features"].([]interface{})[0]
	assert.IsType(map[string]interface{}{}, obj2)
	obj3 := obj2.(map[string]interface{})

	//fmt.Println(obj3)
	assert.Equal("Feature", obj3["type"])
}

func TestCatalogSceneDownload(t *testing.T) {
	assert := assert.New(t)

	c, err := NewCatalogClient()
	assert.NoError(err)

	t.Skip("downloads take too long and may time out")
	m, err := c.DoCatalogSceneDownload("landsat:LC81260322017212LGN00")
	assert.NoError(err)

	assert.Equal(11, m)
}
