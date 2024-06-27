package storage

import (
	"io/fs"
	"log/slog"
	"strings"
	"sync"
	"time"

	"github.com/madappgang/identifo/v2/logging"
	"golang.org/x/exp/slices"
)

// FSWatcher watch for files changes in FS and notifies on file change
type FSWatcher struct {
	logger     *slog.Logger
	f          fs.FS
	poll       time.Duration
	change     chan []string
	err        chan error
	done       chan bool
	isWatching bool

	keys          []string
	watchingSince map[string]time.Time
	keysLock      sync.RWMutex
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

func NewFSWatcher(
	logger *slog.Logger,
	f fs.FS,
	keys []string,
	poll time.Duration,
) *FSWatcher {
	// let's remove trailing
	return &FSWatcher{
		logger:        logger,
		f:             f,
		poll:          poll,
		change:        make(chan []string),
		err:           make(chan error),
		done:          make(chan bool),
		isWatching:    false,
		keys:          KeysWithFixedSlashed(keys),
		watchingSince: make(map[string]time.Time),
		keysLock:      sync.RWMutex{},
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
			w.logger.Debug("fs watcher checking the files ...")
			go w.checkUpdatedFiles()
		case <-w.done:
			w.isWatching = false
			w.logger.Debug("fs watcher has received done signal and stopping itself ...")
			return
		}
	}
}

func (w *FSWatcher) checkUpdatedFiles() {
	var modifiedKeys []string

	w.keysLock.RLock()

	for _, key := range w.keys {
		stat, err := fs.Stat(w.f, key)
		if err != nil {
			w.logger.Error("fs watcher getting error", logging.FieldError, err)
			w.err <- err
			continue
		}
		if stat.ModTime().After(w.watchingSince[key]) {
			modifiedKeys = append(modifiedKeys, key)
		}
	}

	w.keysLock.RUnlock()

	if len(modifiedKeys) == 0 {
		return
	}

	w.keysLock.Lock()

	for _, key := range modifiedKeys {
		w.watchingSince[key] = time.Now()
	}

	w.keysLock.Unlock()

	w.logger.Info("fs files has changed response", "keys", modifiedKeys)
	// report file change
	w.change <- modifiedKeys
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
	w.keysLock.Lock()
	defer w.keysLock.Unlock()

	if !slices.Contains(w.keys, path) {
		w.keys = append(w.keys, path)
		w.watchingSince[path] = time.Now()
	}
}
