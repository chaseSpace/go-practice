package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestX(t *testing.T) {
	d, _ := time.ParseDuration("7m")
	assert.Equal(t, d, time.Minute*7)

	d, _ = time.ParseDuration("7h7s")
	assert.Equal(t, d, time.Hour*7+time.Second*7)

	d, err := time.ParseDuration("7d")
	t.Log(err) /// time: unknown unit "d" in duration "7d"
}
