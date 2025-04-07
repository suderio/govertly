package history

import (
	"fmt"
	"sort"
)

// ReplayCommands replays all stored commands in order
func (h *HistoryManager) ReplayCommands() error {
	cmds, err := h.LoadAll()
	if err != nil {
		return err
	}

	// Sort by timestamp
	sort.Slice(cmds, func(i, j int) bool {
		return cmds[i].Timestamp.Before(cmds[j].Timestamp)
	})

	for _, env := range cmds {
		cmd, err := CreateCommand(env)
		if err != nil {
			return fmt.Errorf("could not create command %s: %w", env.Hash, err)
		}

		fmt.Printf("Replaying: %s (%s)\n", env.Type, env.Hash)
		if err := cmd.Do(); err != nil {
			return fmt.Errorf("command %s failed: %w", env.Hash, err)
		}
	}

	return nil
}
