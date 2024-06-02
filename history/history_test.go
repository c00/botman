package history

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"testing"
	"time"

	"github.com/c00/botman/config"
	"github.com/c00/botman/models"
)

//todo make this save to another folder and use known starting points.

func mustBeTestBuild() {
	if config.APP_FOLDER != ".botman-tests" {
		panic("not a test build")
	}
}

func mustParse(input string) time.Time {
	date, err := time.Parse(time.RFC3339, input)
	if err != nil {
		panic("cannot parse date")
	}
	return date
}

func saveTestEntry(date string) models.HistoryEntry {
	entry := models.HistoryEntry{
		Name: fmt.Sprintf("%v.yaml", date),
		Date: mustParse(date),
		// Date: mustParse("2024-01-20T13:00:00+07:00"),
		Messages: []models.ChatMessage{
			{Role: "user", Content: "test message 1"},
			{Role: "assistant", Content: "test message 2"},
		},
	}

	saveFile(entry)
	return entry
}

func cleanupTestFolder() {
	mustBeTestBuild()
	u, _ := user.Current()
	path := filepath.Join(u.HomeDir, config.APP_FOLDER, "history")
	os.RemoveAll(path)
}

func TestSaveChat(t *testing.T) {
	cleanupTestFolder()

	messages := []models.ChatMessage{
		{Role: "test role", Content: "test message"},
	}

	entry, err := SaveChat(messages)

	if err != nil {
		t.Fatal("Nothing was saved", err)
	}

	//Check file now exists.
	u, _ := user.Current()
	path := filepath.Join(u.HomeDir, config.APP_FOLDER, "history", entry.Name)
	_, err = os.Stat(path)
	if err != nil {
		t.Fatal("Stat error", err)
	}
}

func TestListChats(t *testing.T) {
	cleanupTestFolder()

	saveTestEntry("2024-01-20T13:00:00+07:00")
	saveTestEntry("2024-01-21T13:00:00+07:00")
	saveTestEntry("2024-01-22T13:00:00+07:00")

	entries, err := List()
	if err != nil {
		t.Fatal("could not list chats", err)
	}

	if len(entries) != 3 {
		t.Fatal("there should be 3 items")
	}
}

func TestLoadChat(t *testing.T) {
	cleanupTestFolder()

	entry1 := saveTestEntry("2024-01-20T13:00:00+07:00")
	entry2 := saveTestEntry("2024-01-21T13:00:00+07:00")

	entry, err := LoadChat(0)
	if err != nil {
		t.Fatal("could not load chat", err)
	}
	if entry.Name != entry2.Name {
		t.Fatal("wrong entry retrieved")
	}

	entry, err = LoadChat(1)
	if err != nil {
		t.Fatal("could not load chat", err)
	}
	if entry.Name != entry1.Name {
		t.Fatal("wrong entry retrieved")
	}
}
