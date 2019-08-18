package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestZoneValidate(t *testing.T) {
	z := &Zone{
		Name:    "example-com",
		Records: []string{},
	}

	err := z.validate()
	assert.Error(t, err)
}
