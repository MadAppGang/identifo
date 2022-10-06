package s3

import (
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/madappgang/identifo/v2/model"
)

// S3 SDK is missing this error code
const ErrCodeNotModified = "NotModified"

type PollWatcher struct {
	client        *s3.S3
	keys          []string
	settings      model.S3StorageSettings
	poll          time.Duration
	change        chan []string
	err           chan error
	done          chan bool
	isWatching    bool
	watchingSince map[string]time.Time
}

func NewPollWatcher(settings model.S3StorageSettings, poll time.Duration) (model.ConfigurationWatcher, error) {
	s3client, err := NewS3Client(settings.Region, settings.Endpoint)
	if err != nil {
		return nil, fmt.Errorf("Cannot initialize S3 client for s3 poll config watcher: %s.", err)
	}

	return &PollWatcher{
		client:     s3client,
		settings:   settings,
		keys:       []string{settings.Key},
		poll:       poll,
		change:     make(chan []string),
		err:        make(chan error),
		done:       make(chan bool),
		isWatching: false,
	}, nil
}

func NewPollWatcherForKeyList(settings model.S3StorageSettings, keys []string, poll time.Duration) (model.ConfigurationWatcher, error) {
	s3client, err := NewS3Client(settings.Region, settings.Endpoint)
	if err != nil {
		return nil, fmt.Errorf("Cannot initialize S3 client for s3 poll config watcher: %s.", err)
	}

	return &PollWatcher{
		client:     s3client,
		settings:   settings,
		poll:       poll,
		change:     make(chan []string),
		err:        make(chan error),
		done:       make(chan bool),
		isWatching: false,
		keys:       keys,
	}, nil
}

func (w *PollWatcher) Watch() {
	// non blocking run of watch function
	go w.runWatch()
}

// blocking version of Watch
func (w *PollWatcher) runWatch() {
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
			log.Println("s3 config poll watcher checking the files ...")
			go w.checkUpdatedFiles()
		case <-w.done:
			log.Println("s3 config poll watcher has received done signal and stopping itself ...")
			return
		}
	}
}

func (w *PollWatcher) checkUpdatedFiles() {
	var modifiedKeys []string
	for _, key := range w.keys {
		input := &s3.HeadObjectInput{
			Bucket: aws.String(w.settings.Bucket),
			Key:    aws.String(key),
			// Return the object only if it has been modified since the specified time, otherwise return a 304 (not modified).
			IfModifiedSince: aws.Time(w.watchingSince[key]),
		}

		_, err := w.client.HeadObject(input)
		if err != nil {
			if aerr, ok := err.(awserr.Error); ok && aerr.Code() == ErrCodeNotModified {
				// file has not been changed
				// just silently ignore
				// log.Printf("The file has not been modified: %s", *input.Key)
			} else {
				log.Printf("getting error: %+v\n", err)
				// report error
				w.err <- err
			}
		} else {
			w.watchingSince[key] = time.Now()
			modifiedKeys = append(modifiedKeys, key)
		}
	}
	if len(modifiedKeys) > 0 {
		log.Printf("s3 files has changed response: %+v", modifiedKeys)
		// report file change
		w.change <- modifiedKeys
	}
}

func (w *PollWatcher) IsWatching() bool {
	return w.isWatching
}

func (w *PollWatcher) WatchChan() <-chan []string {
	return w.change
}

func (w *PollWatcher) ErrorChan() <-chan error {
	return w.err
}

func (w *PollWatcher) Stop() {
	// non blocking stop
	go func() {
		w.done <- true
	}()
}
