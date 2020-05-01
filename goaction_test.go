package goaction

import (
	"testing"

	"github.com/stretchr/testify/assert"
)
func TestVars(t *testing.T) {
	if !CI {
		t.Skip("Only runs in CI mode")
	}

	assert.Equal(t, "/home/runner", Home)
	assert.Equal(t, ".github/workflows/testgo.yml", Workflow)
	assert.Equal(t, "posener/goaction", Repository)
	assert.NotEmpty(t, RunID)
	assert.NotEmpty(t, RunNum)
	assert.NotEmpty(t, ActionID)
	assert.NotEmpty(t, Actor)
	assert.NotEmpty(t, Workspace)
	assert.NotEmpty(t, SHA)
	assert.NotEmpty(t, Ref)

	assert.Equal(t, "posener", Owner())
	assert.Equal(t, "goaction", Project())
	switch Event{
	case EventPush:
		assert.Equal(t, "master", Branch())
	case EventPullRequest:
		assert.Less(t, 0, PrNum())
	}
}