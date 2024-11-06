package log

import (
	"encoding/json"
	"os"
)

type LogDriver interface {
	IsSelected(keyFromConfig string) bool
	Configure(rawConfig []byte) error
	NewLog() Log
}

type LogFactory struct {
	drivers   []LogDriver
	config    DefaultLogFactoryConfig
	configRaw []byte
}

func (f *LogFactory) AddDriver(name string, driver LogDriver) {
	// not concurrency safe
	f.drivers = append(f.drivers, driver)
}

func (f *LogFactory) ConfigureFromFile(filePath string) error {
	var err error
	f.configRaw, err = os.ReadFile(filePath)
	if err != nil {
		return err
	}
	err = json.Unmarshal(f.configRaw, &f.config)
	if err != nil {
		return err
	}

	drive := f.getSelectedDriver()
	err = drive.Configure(f.configRaw)
	if err != nil {
		return err
	}

	return nil
}

func (r *LogFactory) NewLog() Log {
	return r.getSelectedDriver().NewLog()
}

func (f *LogFactory) getSelectedDriver() LogDriver {
	for _, driver := range f.drivers {
		if driver.IsSelected(f.config.DriverKey) {
			return driver
		}
	}

	return &DefaultLogDriver{}
}

var GlobalLogFactory = &LogFactory{}
