package fs_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/madappgang/identifo/v2/storage/fs"
)

// TestWatcher test file watcher
func TestWatcher(t *testing.T) {
	// let's create temp file, and use current time to create random file name
	filename := fmt.Sprintf("/tmp/file_ratcher_test_%v.txt", time.Now().Unix())
	if err := ioutil.WriteFile(filename, []byte("Hello"), 0o755); err != nil {
		t.Fatalf("error creating test file: %v", err)
	}

	// create and start watcher
	watcher := fs.NewWatcher(filename)
	fileChanged := false
	watcher.Watch()

	// this go-routine is adding some data to file after 1 second, letting us to start select statement
	go func() {
		time.Sleep(time.Second * 1)
		f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0o644)
		if err != nil {
			t.Errorf("error opening file for changes: %v", err)
			return
		}
		_, err = f.Write([]byte("/nFile changed"))
		if err != nil {
			t.Errorf("error writing to file: %v", err)
			return
		}
		f.Close()
	}()

	// select statement which listens to error, file change and timeout
	select {
	case err := <-watcher.ErrorChan():
		t.Errorf("error getting file update %v", err)
		return
	case <-watcher.WatchChan():
		fileChanged = true
		t.Log("getting file changed")
	case <-time.After(time.Second * 5):
		t.Error("timeout waiting file update")
	}

	// let's check if file has been changed after select finished
	if !fileChanged {
		t.Error("no file change handled")
	}
}
