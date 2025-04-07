package sync

import (
	"bufio"
	"fmt"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
	"os"

	"github.com/suderio/govertly/internal/gitutils"
	"github.com/suderio/govertly/internal/history"
)

func confirm(prompt string) bool {
	fmt.Print(prompt + " [y/N]: ")
	sc := bufio.NewScanner(os.Stdin)
	sc.Scan()
	return sc.Text() == "y" || sc.Text() == "Y"
}

func SyncAndReplay(storePath string, hm *history.HistoryManager) error {
	remote := viper.GetString("git.remote")
	branch := viper.GetString("git.branch")

	local, err := gitutils.GetCurrentCommit(storePath)
	if err != nil {
		return fmt.Errorf("get local commit: %w", err)
	}

	remoteHash, err := gitutils.GetRemoteCommit(storePath, remote, branch)
	if err != nil {
		return fmt.Errorf("get remote commit: %w", err)
	}

	if local == remoteHash {
		fmt.Println("Local and remote are in sync.")
		return nil
	}

	if !confirm(fmt.Sprintf("Local is behind %s/%s. Reset, replay commands, and push?", remote, branch)) {
		fmt.Println("Aborted.")
		return nil
	}

	if err := gitutils.ResetToRemote(storePath, remote, branch); err != nil {
		return fmt.Errorf("reset failed: %w", err)
	}

	fmt.Println("Replaying command history...")
	if err := hm.ReplayCommands(); err != nil {
		return fmt.Errorf("replay failed: %w", err)
	}

	fmt.Println("Committing and pushing...")
	if err := gitutils.CommitAndPush(storePath, "Replayed local commands after sync"); err != nil {
		return fmt.Errorf("push failed: %w", err)
	}

	return nil
}
