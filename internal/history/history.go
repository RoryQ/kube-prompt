package history

import (
	"bufio"
	"fmt"
	"github.com/c-bata/go-prompt"
	"github.com/c-bata/kube-prompt/internal/debug"
	"os"
	"path"
	"strings"
)

const (
	envEnablePersistence = "KUBE_PROMPT_PERSIST_HISTORY"
	historyFileName  = ".kube-prompt.history"
)

var (
	historyFile *os.File
)

func Close() {
	if historyFile != nil {
		_ = historyFile.Close()
	}
}

func historyFilePath() string {
	home, _ := os.UserHomeDir()
	return path.Join(home, historyFileName)
}

func readPersisted() []string {
	var lines []string
	var err error
	history, err := os.Open(historyFilePath())
	if err != nil {
		debug.Log(err.Error())
		return lines
	}
	defer history.Close()

	scan := bufio.NewScanner(history)
	for scan.Scan() {
		lines = append(lines, scan.Text())
	}

	return lines
}

func openForAppend() {
	var err error
	historyFile, err = os.OpenFile(historyFilePath(), os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		debug.Log(err.Error())
	}
}

func LoadPersisted() prompt.Option {
	enabled := os.Getenv(envEnablePersistence)
	if enabled == "true" || enabled == "1" {
		persisted := readPersisted()
		openForAppend()
		return prompt.OptionHistory(persisted)
	}

	return func(prompt *prompt.Prompt) error {
		return nil
	}
}

func LogHistory(s string) {
	if historyFile != nil {
		// bash ignorespace behaviour
		if !strings.HasPrefix(s, " "){
			_, _ = historyFile.WriteString(fmt.Sprintf("%s\n", s))
			_ = historyFile.Sync()
		}
	}
}