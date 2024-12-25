package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestX(t *testing.T) {
	d, _ := time.ParseDuration("7s")
	assert.Equal(t, d, time.Second*7)
}
