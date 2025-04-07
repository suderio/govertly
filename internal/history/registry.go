package history

import (
	"encoding/json"
	"fmt"
)

// CommandFactory creates a Command from raw JSON payload
type CommandFactory func(json.RawMessage) (Command, error)

var registry = make(map[string]CommandFactory)

// RegisterCommandType registers a command type by name
func RegisterCommandType(name string, factory CommandFactory) {
	registry[name] = factory
}

// CreateCommand attempts to reconstruct a Command from envelope
func CreateCommand(env CommandEnvelope) (Command, error) {
	factory, ok := registry[env.Type]
	if !ok {
		return nil, fmt.Errorf("unknown command type: %s", env.Type)
	}
	return factory(env.Payload)
}
