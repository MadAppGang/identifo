package storage

import (
	"io/fs"
	"log"
	"time"
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

func NewFSWatcher(f fs.FS, keys []string, poll time.Duration) *FSWatcher {
	return &FSWatcher{
		f:          f,
		keys:       keys,
		poll:       poll,
		change:     make(chan []string),
		err:        make(chan error),
		done:       make(chan bool),
		isWatching: false,
	}
}

func (w *FSWatcher) Watch() {
	// non blocking run of watch function
	go w.runWatch()
}

// watch runloop using go channels
func (w *FSWatcher) runWatch() {
	w.isWatching = true
	w.watchingSince = make(map[string]time.Time)
	for _, k := range w.keys {
		w.watchingSince[k] = time.Now()
	}

	defer func() {
		w.isWatching = false
		w.watchingSince = nil
	}()

	for {
		select {
		case <-time.After(w.poll):
			log.Println("fs watcher checking the files ...")
			go w.checkUpdatedFiles()
		case <-w.done:
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
	// non blocking stop
	go func() {
		w.done <- true
	}()
}
