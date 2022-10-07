package fs

import (
	"errors"
	"os"
	"testing"

	"github.com/madappgang/identifo/v2/model"
)

func TestNewFSConnectionTesterWithOneExpectedFile(t *testing.T) {
	settings := model.FileStorageLocal{
		Path: "../../test/artifacts/spa",
	}
	tester := NewFSConnectionTester(settings, []string{"index.html"})

	err := tester.Connect()
	if err != nil {
		t.Errorf("unable to find expected files with error: %+v", err)
	}
}

func TestNewFSConnectionTesterWithNoFiles(t *testing.T) {
	settings := model.FileStorageLocal{
		Path: "../../test/artifacts",
	}
	tester := NewFSConnectionTester(settings, nil)

	err := tester.Connect()
	if err != nil {
		t.Errorf("unable to find expected files with error: %+v", err)
	}
}

func TestNewFSConnectionTesterWithFile(t *testing.T) {
	settings := model.FileStorageLocal{
		Path: "../../test/artifacts/templates",
	}
	tester := NewFSConnectionTester(settings, []string{"mail1.template", "mail2.template"})

	err := tester.Connect()
	if err != nil {
		t.Errorf("unable to find expected files with error: %+v", err)
	}
}

func TestNewFSConnectionTesterWithNoFilesFail(t *testing.T) {
	settings := model.FileStorageLocal{
		Path: "../../test/artifacts/failpath",
	}
	tester := NewFSConnectionTester(settings, nil)

	err := tester.Connect()
	if !errors.Is(err, os.ErrNotExist) {
		t.Errorf("expected to get fs.ErrNotExists, got: %+v", err)
	}
}

func TestNewFSConnectionTesterWithFileWhichIsAbsent(t *testing.T) {
	settings := model.FileStorageLocal{
		Path: "../../test/artifacts",
	}
	tester := NewFSConnectionTester(settings, []string{"index.html"})

	err := tester.Connect()
	if err == nil {
		t.Error("expected to get error, but did not")
	}
}

func TestNewFSConnectionTesterWithFilesAndOneIsAbsent(t *testing.T) {
	settings := model.FileStorageLocal{
		Path: "../../test/artifacts/templates",
	}
	tester := NewFSConnectionTester(settings, []string{"mail1.template", "mail2.template", "mail3.template"})

	err := tester.Connect()
	if err == nil {
		t.Error("expected to get error, but did not")
	}
}
