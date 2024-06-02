package history

import (
	"errors"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"time"

	"github.com/c00/botman/config"
	"github.com/c00/botman/models"
	"gopkg.in/yaml.v3"
)

func SaveChat(messages []models.ChatMessage) (models.HistoryEntry, error) {
	//Convert the things to a history entry
	entry := models.HistoryEntry{
		Name:     fmt.Sprintf("%v.yaml", time.Now().Format(time.RFC3339)),
		Date:     time.Now(),
		Messages: messages,
	}

	//Save that to disk as yaml
	err := saveFile(entry)
	if err != nil {
		return models.HistoryEntry{}, err
	}

	return entry, nil
}

func saveFile(entry models.HistoryEntry) error {
	if len(entry.Messages) == 0 {
		return nil
	}

	if entry.Name == "" {
		return errors.New("entry does not have a name")
	}

	u, err := user.Current()
	if err != nil {
		return err
	}

	savePath := filepath.Join(u.HomeDir, config.APP_FOLDER, "history", entry.Name)

	dir := filepath.Dir(savePath)
	info, err := os.Stat(dir)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.MkdirAll(dir, 0700)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	} else if !info.IsDir() {
		return fmt.Errorf("%v should be a directory", dir)
	}

	bytes, err := yaml.Marshal(entry)
	if err != nil {
		return err
	}
	err = os.WriteFile(savePath, bytes, 0700)
	if err != nil {
		return err
	}

	return nil
}

func LoadChat(lookback int) (models.HistoryEntry, error) {
	files, err := List()
	if err != nil {
		return models.HistoryEntry{}, err
	}

	length := len(files)
	if length == 0 {
		return models.HistoryEntry{}, errors.New("no history")
	}

	if lookback >= length {
		return models.HistoryEntry{}, errors.New("lookback is further than available entries")
	}

	//Get the looback index
	index := length - lookback - 1
	return loadFile(files[index])

}

func loadFile(path string) (models.HistoryEntry, error) {
	u, err := user.Current()
	if err != nil {
		return models.HistoryEntry{}, err
	}
	//Read the history dir
	filePath := filepath.Join(u.HomeDir, config.APP_FOLDER, "history", path)

	bytes, err := os.ReadFile(filePath)
	if err != nil {
		return models.HistoryEntry{}, err
	}

	entry := models.HistoryEntry{}
	err = yaml.Unmarshal(bytes, &entry)
	if err != nil {
		return models.HistoryEntry{}, err
	}

	return entry, nil
}

func List() ([]string, error) {
	u, err := user.Current()
	if err != nil {
		return []string{}, err
	}
	//Read the history dir
	historyPath := filepath.Join(u.HomeDir, config.APP_FOLDER, "history")
	files := []string{}

	err = filepath.Walk(historyPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		//Check if the filename ends in .yaml, if not continue
		if !strings.HasSuffix(info.Name(), ".yaml") {
			return nil
		}
		files = append(files, info.Name())
		return nil
	})
	if err != nil {
		return []string{}, err
	}

	return files, nil
}
