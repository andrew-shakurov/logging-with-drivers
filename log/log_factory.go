package log

import (
	"encoding/json"
	"os"
)

type logFactoryConfig struct {
	DriverKey string `json:"driver"`
}

type logFactory struct {
	drivers   []LogDriver
	config    logFactoryConfig
	configRaw []byte
}

func (f *logFactory) addDriver(name string, driver LogDriver) {
	// not concurrency safe
	// @todo reconsider the driver selection approach
	f.drivers = append(f.drivers, driver)
}

func (f *logFactory) configureFromFile(filePath string) error {
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

func (r *logFactory) newLog() Log {
	return r.getSelectedDriver().NewLog()
}

func (f *logFactory) getSelectedDriver() LogDriver {
	for _, driver := range f.drivers {
		if driver.IsSelected(f.config.DriverKey) {
			return driver
		}
	}

	return f.drivers[0]
}

var logFactoryInstance = &logFactory{}

func NewLog() Log {
	return logFactoryInstance.newLog()
}

func ConfigureFromFile(pathToConfigFile string) error {
	return logFactoryInstance.configureFromFile(pathToConfigFile)
}

func AddDriver(name string, driver LogDriver) {
	logFactoryInstance.addDriver(name, driver)
}

func Close() {
	logFactoryInstance.getSelectedDriver().Close()
}

func init() {
	AddDriver(DriverKey, newDefaultLogDriver())
}
