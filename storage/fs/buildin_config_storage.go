package fs

import (
	"fmt"

	"github.com/madappgang/identifo/v2/model"
)

type BuildinConfigStorage struct {
	UpdateChan       chan interface{}
	updateChanClosed bool
	cached           *model.ServerSettings
}

func NewBuildingConfigurationStorage() *BuildinConfigStorage {
	return &BuildinConfigStorage{
		UpdateChan:       make(chan interface{}, 1),
		updateChanClosed: false,
	}
}

func (cs *BuildinConfigStorage) WriteConfig(settings model.ServerSettings) error {
	return fmt.Errorf("Building configuration storage is static and does not supports mutation")
}

func (cs *BuildinConfigStorage) LoadServerSettings(validate bool) (model.ServerSettings, []error) {
	cs.cached = &model.DefaultServerSettings
	return model.DefaultServerSettings, nil
}

func (cs *BuildinConfigStorage) GetUpdateChan() chan interface{} {
	return cs.UpdateChan
}

// CloseUpdateChan closes update channel.
func (cs *BuildinConfigStorage) CloseUpdateChan() {
	close(cs.UpdateChan)
	cs.updateChanClosed = true
}

func (cs *BuildinConfigStorage) ForceReloadOnWriteConfig() bool {
	return false
}

func (cs *BuildinConfigStorage) Errors() []error {
	return nil
}

func (cs *BuildinConfigStorage) LoadedSettings() *model.ServerSettings {
	return cs.cached
}
