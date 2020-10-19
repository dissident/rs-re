package main

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestMain(t *testing.T) {
	foo := 1
	assert.Equal(t, foo, 1)
}
