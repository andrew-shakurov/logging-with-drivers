package logdriverjson

import (
	"encoding/json"
	"os"

	"example.com/log"
)

type JSONLogDriver struct {
	config Config
}

func (d *JSONLogDriver) IsSelected(keyFromConfig string) bool {
	return DriverKey == keyFromConfig
}

func (d *JSONLogDriver) Configure(rawConfig []byte) error {
	err := json.Unmarshal(rawConfig, &d.config)
	if err != nil {
		return err
	}

	// @todo validate configuration
	return nil
}

func (d *JSONLogDriver) NewLog() log.Log {
	copyOfConfig := d.config

	// @todo re-think how to handle an error
	out, _ := os.OpenFile(d.config.OutputFile, os.O_RDWR|os.O_CREATE, 0644)

	defaultLog := JSONLog{
		MessageLogLevelOfLog: log.L_INFO,
		config:               copyOfConfig,
		out:                  out,
	}

	return &defaultLog
}

func init() {
	log.GlobalLogFactory.AddDriver(DriverKey, &JSONLogDriver{})
}
