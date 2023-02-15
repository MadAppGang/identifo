package storage

import (
	"io/fs"
	"log"
	"strings"
	"time"

	"golang.org/x/exp/slices"
)

// FSWatcher watch for files changes in FS and notifies on file change
type FSWatcher struct {
	f             fs.FS
	keys          []string
	poll          time.Duration
	change        chan []string
	err           chan error
	done          chan bool
	isWatching    bool
	watchingSince map[string]time.Time
}

// KeysWithFixedSlashed remove prefixed or postfixed slashed in a path
// because it is not valid path and will fail fs.ValidPath validation
func KeysWithFixedSlashed(keys []string) []string {
	result := []string{}
	for _, k := range keys {
		k = strings.TrimPrefix(k, "/")
		k = strings.TrimSuffix(k, "/")
		result = append(result, k)
	}
	return result
}

func NewFSWatcher(f fs.FS, keys []string, poll time.Duration) *FSWatcher {
	// let's remove trailing
	return &FSWatcher{
		f:             f,
		keys:          KeysWithFixedSlashed(keys),
		poll:          poll,
		change:        make(chan []string),
		err:           make(chan error),
		done:          make(chan bool),
		watchingSince: make(map[string]time.Time),
		isWatching:    false,
	}
}

func (w *FSWatcher) Watch() {
	if w.isWatching {
		return
	}

	w.isWatching = true
	for _, k := range w.keys {
		w.watchingSince[k] = time.Now()
	}
	// non blocking run of watch function
	go w.runWatch()
}

// watch runloop using go channels
func (w *FSWatcher) runWatch() {
	defer func() {
		w.isWatching = false
		w.watchingSince = make(map[string]time.Time)
	}()

	for {
		select {
		case <-time.After(w.poll):
			log.Println("fs watcher checking the files ...")
			go w.checkUpdatedFiles()
		case <-w.done:
			w.isWatching = false
			log.Println("fs watcher has received done signal and stopping itself ...")
			return
		}
	}
}

func (w *FSWatcher) checkUpdatedFiles() {
	var modifiedKeys []string
	for _, key := range w.keys {
		stat, err := fs.Stat(w.f, key)
		if err != nil {
			log.Printf("getting error: %+v\n", err)
			w.err <- err
			continue
		}
		if stat.ModTime().After(w.watchingSince[key]) {
			w.watchingSince[key] = time.Now()
			modifiedKeys = append(modifiedKeys, key)
		}
	}
	if len(modifiedKeys) > 0 {
		log.Printf("fs files has changed response: %+v", modifiedKeys)
		// report file change
		w.change <- modifiedKeys
	}
}

func (w *FSWatcher) IsWatching() bool {
	return w.isWatching
}

func (w *FSWatcher) WatchChan() <-chan []string {
	return w.change
}

func (w *FSWatcher) ErrorChan() <-chan error {
	return w.err
}

func (w *FSWatcher) Stop() {
	if !w.isWatching {
		return
	}
	// non blocking stop
	go func() {
		w.done <- true
	}()
}

func (w *FSWatcher) AppendForWatching(path string) {
	if !slices.Contains(w.keys, path) {
		w.keys = append(w.keys, path)
		w.watchingSince[path] = time.Now()
	}
}
