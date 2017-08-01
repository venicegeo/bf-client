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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCoastlineDownload(t *testing.T) {
	assert := assert.New(t)

	c, err := NewCoastlineClient()
	assert.NoError(err)

	const id = "6f475d16-d8f1-4f15-9dda-cf9b2d502241"

	//	t.Skip("downloads take too long and may time out")
	m, err := c.DoCoastlineDownload(id)
	assert.NoError(err)

	assert.Contains(m, "3530022")
}
