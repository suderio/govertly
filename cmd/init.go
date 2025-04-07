package cmd

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	keyFile   string
	subPath   string
	enableGit bool
	remoteGit string
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize new password storage and assign key-file for encryption",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Resolve base store path
		baseDir := viper.GetString("store_dir")
		if baseDir == "" {
			baseDir = os.Getenv("XDG_DATA_HOME")
			if baseDir == "" {
				home, err := os.UserHomeDir()
				if err != nil {
					return fmt.Errorf("could not resolve user home dir: %w", err)
				}
				baseDir = filepath.Join(home, ".local", "share", "govertly")
			} else {
				baseDir = filepath.Join(baseDir, "govertly")
			}
		}

		// Resolve path to subdirectory if --path was given
		targetDir := baseDir
		if subPath != "" {
			targetDir = filepath.Join(baseDir, subPath)
		}

		// Ensure target directory exists
		if err := os.MkdirAll(targetDir, 0700); err != nil {
			return fmt.Errorf("could not create target directory: %w", err)
		}

		// Handle key assignment
		keyFilePath := filepath.Join(targetDir, ".key")
		if keyFile == "" {
			// Remove existing key file
			if err := os.Remove(keyFilePath); err != nil && !errors.Is(err, os.ErrNotExist) {
				return fmt.Errorf("failed to remove key file: %w", err)
			}
			fmt.Println("Removed key file at:", keyFilePath)
		} else {
			// Write keys to file (support multi-key, newline separated)
			keys := strings.Split(keyFile, ",")
			f, err := os.Create(keyFilePath)
			if err != nil {
				return fmt.Errorf("could not create key file: %w", err)
			}
			defer f.Close()
			for _, k := range keys {
				_, err := f.WriteString(strings.TrimSpace(k) + "\n")
				if err != nil {
					return fmt.Errorf("failed to write key: %w", err)
				}
			}
			fmt.Println("Key file written at:", keyFilePath)

			// TODO: re-encrypt existing files in targetDir if keys changed
		}

		// Initialize Git if specified
		if enableGit || remoteGit != "" {
			if err := initGitRepo(targetDir, remoteGit); err != nil {
				return fmt.Errorf("git init failed: %w", err)
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	initCmd.Flags().StringVarP(&keyFile, "key-file", "k", "", "Comma-separated list of age keys or path to key file")
	initCmd.Flags().StringVarP(&subPath, "path", "p", "", "Sub-path inside the password store to assign keys")
	initCmd.Flags().BoolVarP(&enableGit, "git", "g", false, "Initialize a local git repository")
	initCmd.Flags().StringVarP(&remoteGit, "remote", "r", "", "Initialize git and add remote repository URL")

	// Default path and config bindings
	viper.AutomaticEnv()
	viper.BindEnv("store_dir", "GOVERTLY_DIR")
	viper.BindEnv("default_keys", "GOVERTLY_KEY")
	viper.BindEnv("age_opts", "GOVERTLY_AGE_OPTS")
}

// -- Helpers --

func initGitRepo(path string, remote string) error {
	if _, err := os.Stat(filepath.Join(path, ".git")); err == nil {
		fmt.Println("Git repo already initialized.")
		return nil
	}
	// TODO use git-go
	cmd := exec.Command("git", "init")
	cmd.Dir = path
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to run git init: %w", err)
	}
	fmt.Println("Initialized git repo at:", path)

	if remote != "" {
		cmd = exec.Command("git", "remote", "add", "origin", remote)
		cmd.Dir = path
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to add remote: %w", err)
		}
		fmt.Println("Added remote:", remote)
	}

	return nil
}
