package storage_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/madappgang/identifo/v2/model"
	"github.com/madappgang/identifo/v2/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewFSWatcher(t *testing.T) {
	settings := model.FileStorageSettings{
		Type: model.FileStorageTypeLocal,
		Local: model.FileStorageLocal{
			FolderPath: "../fs_test",
		},
	}
	fss, err := storage.NewFS(settings)
	require.NoError(t, err)
	makeFiles()
	defer deleteFiles()

	fileChanged := false
	updatedCount := 0
	changedFiles := []string{}

	time.Sleep(time.Second * 1)

	watcher := storage.NewFSWatcher(fss, []string{"test1.txt", "test2.txt"}, time.Second*2)
	watcher.Watch()

	go func() {
		time.Sleep(time.Second * 1)
		go updateFile("../fs_test/test1.txt")
		go updateFile("../fs_test/test2.txt")
	}()

outFor:
	for {
		select {
		case err := <-watcher.ErrorChan():
			t.Error(err)
			return
		case files := <-watcher.WatchChan():
			updatedCount++
			fileChanged = true
			changedFiles = append(changedFiles, files...)
			t.Log("getting file changed")
		case <-time.After(time.Second * 5):
			// wait for all updates here for 5 secs
			break outFor
		}
	}

	assert.True(t, fileChanged)
	// should fire update only once!
	assert.LessOrEqual(t, updatedCount, 2) // could update in one butch or two bunches
	assert.GreaterOrEqual(t, updatedCount, 1)
	assert.Contains(t, changedFiles, "test1.txt")
	assert.Contains(t, changedFiles, "test2.txt")
	assert.NotContains(t, changedFiles, "file2.txt")
}

func makeFiles() {
	os.Mkdir("../fs_test", os.ModePerm)
	f, _ := os.Create("../fs_test/test1.txt")
	f.WriteString("Hello")
	f.Close()
	f, _ = os.Create("../fs_test/test2.txt")
	f.WriteString("Hello2")
	f.Close()
	f, _ = os.Create("../fs_test/test3.txt")
	f.WriteString("Hello3")
	f.Close()
}

func updateFile(name string) {
	data := fmt.Sprintf("This is file data has been created at: %v", time.Now())
	_ = ioutil.WriteFile(name, []byte(data), 0o644)
}

func deleteFiles() {
	os.RemoveAll("../fs_test")
}
