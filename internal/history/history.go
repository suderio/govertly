package history

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Interface all commands must implement
type Command interface {
	Do() error
	Type() string
	Payload() any
}

// JSON structure written to file
type CommandEnvelope struct {
	Type      string          `json:"type"`
	CommitId  string          `json:"commitId"`
	Timestamp time.Time       `json:"timestamp"`
	Payload   json.RawMessage `json:"payload"`
	Hash      string          `json:"-"`
}

// HistoryManager handles storage and retrieval
type HistoryManager struct {
	Dir string // full path to $XDG_STATE_HOME/govertly
}

func NewHistoryManager() (*HistoryManager, error) {
	stateDir := os.Getenv("XDG_STATE_HOME")
	if stateDir == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, err
		}
		stateDir = filepath.Join(home, ".local", "state")
	}
	full := filepath.Join(stateDir, "govertly", "commands")
	if err := os.MkdirAll(full, 0700); err != nil {
		return nil, err
	}
	return &HistoryManager{Dir: full}, nil
}

// SaveCommand serializes and writes the command to a file
func (h *HistoryManager) SaveCommand(cmd Command, commitId string) (string, error) {
	payload, err := json.Marshal(cmd.Payload())
	if err != nil {
		return "", err
	}

	env := CommandEnvelope{
		Type:      cmd.Type(),
		CommitId:  commitId,
		Timestamp: time.Now().UTC(),
		Payload:   payload,
	}

	hash := computeHash(env)
	env.Hash = hash

	filePath := filepath.Join(h.Dir, hash+".json")
	data, err := json.MarshalIndent(env, "", "  ")
	if err != nil {
		return "", err
	}

	if err := os.WriteFile(filePath, data, 0600); err != nil {
		return "", err
	}

	return hash, nil
}

// LoadAll returns all envelopes, sorted by timestamp
func (h *HistoryManager) LoadAll() ([]CommandEnvelope, error) {
	files, err := os.ReadDir(h.Dir)
	if err != nil {
		return nil, err
	}

	var result []CommandEnvelope
	for _, file := range files {
		if file.IsDir() || filepath.Ext(file.Name()) != ".json" {
			continue
		}
		data, err := os.ReadFile(filepath.Join(h.Dir, file.Name()))
		if err != nil {
			return nil, err
		}

		var env CommandEnvelope
		if err := json.Unmarshal(data, &env); err != nil {
			return nil, err
		}
		env.Hash = file.Name()[:len(file.Name())-5]
		result = append(result, env)
	}
	return result, nil
}

// Helper to compute SHA-256 hash of envelope contents
func computeHash(env CommandEnvelope) string {
	blob := fmt.Sprintf("%s|%s|%s|%s",
		env.Type,
		env.CommitId,
		env.Timestamp.Format(time.RFC3339Nano),
		string(env.Payload),
	)
	sum := sha256.Sum256([]byte(blob))
	return fmt.Sprintf("%x", sum)
}

