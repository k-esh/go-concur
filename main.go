package main

import (
	"example.com/runner/pkg/config"
	"example.com/runner/pkg/runner"
	"log"
	"os"
)

type CommandUnitResult struct {
	Nickname string
	Command string
	Args []string
	OutputStdAndError string
	ExitCode          int
	Error             error
}

//TODO logging in json format maybe zip library
//TODO enforce have at least one command ?
func main() {
	//take input arg config file and if there is no arg use default
	const defaultConfigFile = "pkg/config/jobStep.toml"
	var configFile string
	if len(os.Args) > 1 {
		configFile = os.Args[1]
	} else {
		configFile = defaultConfigFile
	}
	log.Printf("config file set to: %s \n", configFile)
	myCommandUnits := config.InitJobSteps(configFile)

	cuResults := runner.NewHampshire(myCommandUnits)
	log.Printf("all %d of the jobsteps are done", len(myCommandUnits))
//TODO use below instead	log.Printf("results: %+v\n", cuResults)

	var failCount int
	for _, result := range cuResults {
		log.Printf("result: %+v\n", result)
		if result.ExitCode != 0 {
			failCount++
		}
	}
	log.Printf("count of failures: %+d\n", failCount)
	//fail when any of the commands have failed?
	os.Exit(failCount)

}
