package main

import (
	"testing"

	"gotest.tools/v3/assert"
)

func Test_newBadgesCommand(t *testing.T) {
	c := newBadgesCommand("1.0.0")
	assert.Assert(t, c != nil)
}
