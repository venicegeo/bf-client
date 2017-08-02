package client

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClient(t *testing.T) {
	assert := assert.New(t)

	c, err := NewClient()
	assert.NoError(err)
	assert.NotNil(c)

	assert.NotNil(c.Catalog)
	assert.NotNil(c.Coastline)
	assert.NotNil(c.Job)
	assert.NotNil(c.Algorithm)
}
