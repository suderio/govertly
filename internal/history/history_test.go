package history_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/suderio/govertly/internal/history"
)

// -- Sample Command --

type AddPasswordCommand struct {
	Path      string `json:"path"`
	Encrypted string `json:"encrypted"`
}

func (a *AddPasswordCommand) Do() error {
	return nil
}

func (a *AddPasswordCommand) Type() string {
	return "add-password"
}

func (a *AddPasswordCommand) Payload() any {
	return a
}

// -- Test --

func TestCommandSaveAndLoad(t *testing.T) {
	// Setup test-specific dir
	tmpDir := t.TempDir()
	hist := &history.HistoryManager{Dir: tmpDir}

	// Create test command
	cmd := &AddPasswordCommand{
		Path:      "test/site",
		Encrypted: "AGE-ENCRYPTED-123",
	}
	commitId := "deadbeef123"

	// Save command
	hash, err := hist.SaveCommand(cmd, commitId)
	assert.NoError(t, err)
	assert.True(t, len(hash) == 64, "Hash should be SHA256")

	// Load command back
	entries, err := hist.LoadAll()
	assert.NoError(t, err)
	assert.Equal(t, 1, len(entries))

	entry := entries[0]
	assert.Equal(t, "add-password", entry.Type)
	assert.Equal(t, commitId, entry.CommitId)
	assert.True(t, time.Since(entry.Timestamp) < time.Minute)
	assert.Equal(t, hash, entry.Hash)

	// Check payload content
	var decoded AddPasswordCommand
	err = json.Unmarshal(entry.Payload, &decoded)
	assert.NoError(t, err)
	assert.Equal(t, "test/site", decoded.Path)
	assert.Equal(t, "AGE-ENCRYPTED-123", decoded.Encrypted)

	// Check if the file exists with correct name
	files, err := os.ReadDir(tmpDir)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(files))
	assert.True(t, strings.HasPrefix(files[0].Name(), hash))
	assert.Equal(t, filepath.Ext(files[0].Name()), ".json")
}
