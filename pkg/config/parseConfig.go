package config

import (
	"github.com/BurntSushi/toml"
	"log"
)

//Commands and their args in this structure
type CommandConfig struct {
	Title        string
	MaxRuntime	 int
	CommandUnits map[string]CommandUnit
}
//commands and their optional args
type CommandUnit struct {
	Command string
	Args    []string
}

//TODO add a test case
func InitJobSteps(configFile string) (CommandConfig, error) {
	var config CommandConfig

	//retern error here if any problem with file
	if _, err := toml.DecodeFile(configFile, &config); err != nil {
		log.Printf("error: %s parsing config file: %s will return error\n", err, configFile)
		return config, err
	}

	log.Printf("Config from file: %s has title: %s maxRuntime: %s commands: +%v\n", configFile, config.Title, config.MaxRuntime, config.CommandUnits)
	return config, nil
}
