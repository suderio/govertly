package history_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/suderio/govertly/internal/history"
)

// Sample command reused
type AnotherAddPasswordCommand struct {
	Path      string `json:"path"`
	Encrypted string `json:"encrypted"`
}

func (a *AnotherAddPasswordCommand) Do() error    { return nil }
func (a *AnotherAddPasswordCommand) Type() string { return "add-password" }
func (a *AnotherAddPasswordCommand) Payload() any { return a }

func TestRegisterAndCreateCommand(t *testing.T) {
	// Register
	history.RegisterCommandType("add-password", func(p json.RawMessage) (history.Command, error) {
		var cmd AnotherAddPasswordCommand
		err := json.Unmarshal(p, &cmd)
		return &cmd, err
	})

	// Create payload
	payload := AnotherAddPasswordCommand{
		Path:      "secure/site",
		Encrypted: "age-encrypted-payload",
	}
	raw, err := json.Marshal(payload)
	assert.NoError(t, err)

	env := history.CommandEnvelope{
		Type:    "add-password",
		Payload: raw,
	}

	// Reconstruct command
	cmd, err := history.CreateCommand(env)
	assert.NoError(t, err)

	casted, ok := cmd.(*AnotherAddPasswordCommand)
	assert.True(t, ok)
	assert.Equal(t, "secure/site", casted.Path)
	assert.Equal(t, "age-encrypted-payload", casted.Encrypted)
}
