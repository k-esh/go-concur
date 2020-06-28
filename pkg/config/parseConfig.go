package config

import (
	"github.com/BurntSushi/toml"
	"log"
)

type tomlConfig struct {
	Title    string
	CommandUnits map[string]CommandUnit
}

type CommandUnit struct {
	Nickname string
	Command  string
	Args     []string
	//TODO support a max run elapsed time
}


//TODO add a test case
func InitJobSteps(configFile string) []CommandUnit {
	var config tomlConfig
	if _, err := toml.DecodeFile(configFile, &config); err != nil {
		log.Fatal(err) //TODO cleaner way to fail or log it
	}
	log.Printf("config title: %s from file: %s\n", config.Title, configFile)

	jobsteps := make([]CommandUnit, 0, len(config.CommandUnits))

	nicknames := make([]string, 0, len(jobsteps))
	for _, js := range config.CommandUnits { //no need for name (the map's key)
		jobsteps = append(jobsteps, js)
		nicknames = append(nicknames, js.Nickname)
		log.Printf("configured jobstep: %s command: %s with args: [%s]\n", js.Nickname, js.Command, js.Args)
	}
	log.Printf("configured names to run: %s\n", nicknames)
	return jobsteps
}
