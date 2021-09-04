package s3

import (
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/madappgang/identifo/model"
)

// S3 SDK is missing this error code
const ErrCodeNotModified = "NotModified"

type PollWatcher struct {
	client        *s3.S3
	settings      model.S3StorageSettings
	poll          time.Duration
	change        chan bool
	err           chan error
	done          chan bool
	isWatching    bool
	watchingSince *time.Time
}

func NewPollWatcher(settings model.S3StorageSettings, poll time.Duration) (model.ConfigurationWatcher, error) {
	s3client, err := NewS3Client(settings.Region)
	if err != nil {
		return nil, fmt.Errorf("Cannot initialize S3 client for s3 poll config watcher: %s.", err)
	}

	return &PollWatcher{
		client:     s3client,
		settings:   settings,
		poll:       poll,
		change:     make(chan bool),
		err:        make(chan error),
		done:       make(chan bool),
		isWatching: false,
	}, nil
}

func (w *PollWatcher) Watch() {
	// non blocking run of watch function
	go w.runWatch()
}

// blocking version of Watch
func (w *PollWatcher) runWatch() {
	w.isWatching = true
	t := time.Now()
	w.watchingSince = &t

	defer func() {
		w.isWatching = false
		w.watchingSince = nil
	}()

	for {
		select {
		case <-time.After(w.poll):
			log.Println("s3 config poll watcher checking the config file ...")
			go w.requestFileVersion()
		case <-w.done:
			log.Println("s3 config poll watcher has received done signal and stopping itself ...")
			return
		}
	}
}

func (w *PollWatcher) requestFileVersion() {
	input := &s3.HeadObjectInput{
		Bucket:          aws.String(w.settings.Bucket),
		Key:             aws.String(w.settings.Key),
		IfModifiedSince: w.watchingSince, // Return the object only if it has been modified since the specified time, otherwise return a 304 (not modified).
	}

	_, err := w.client.HeadObject(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok && aerr.Code() == ErrCodeNotModified {
			// file has not been changed
			// just silently returning
			return
		} else {
			log.Printf("gettings error: %+v\n", err)
			// report error
			w.err <- err
		}
	} else {
		// no error received, it means file changed
		// report file change
		w.change <- true
	}
}

func (w *PollWatcher) IsWatching() bool {
	return w.isWatching
}

func (w *PollWatcher) WatchChan() <-chan bool {
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
