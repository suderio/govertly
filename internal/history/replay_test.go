package history_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/suderio/govertly/internal/history"
)

// Global to count replays
var replayedCommands []string

type ReplayableCommand struct {
	Name string `json:"name"`
}

func (r *ReplayableCommand) Do() error {
	replayedCommands = append(replayedCommands, r.Name)
	return nil
}
func (r *ReplayableCommand) Type() string { return "replayable" }
func (r *ReplayableCommand) Payload() any { return r }

func TestReplayCommands(t *testing.T) {
	// Reset tracking
	replayedCommands = []string{}

	// Register
	history.RegisterCommandType("replayable", func(p json.RawMessage) (history.Command, error) {
		var cmd ReplayableCommand
		err := json.Unmarshal(p, &cmd)
		return &cmd, err
	})

	// Setup temp dir
	tmp := t.TempDir()
	hist := &history.HistoryManager{Dir: tmp}

	// Create and save multiple commands
	cmd1 := &ReplayableCommand{Name: "first"}
	cmd2 := &ReplayableCommand{Name: "second"}

	// Fake commitId
	commit := "abc123"

	// Save with 100ms gap to preserve timestamp order
	hist.SaveCommand(cmd1, commit)
	time.Sleep(100 * time.Millisecond)
	hist.SaveCommand(cmd2, commit)

	// Replay
	err := hist.ReplayCommands()
	assert.NoError(t, err)

	assert.Equal(t, 2, len(replayedCommands))
	assert.Equal(t, "first", replayedCommands[0])
	assert.Equal(t, "second", replayedCommands[1])
}
