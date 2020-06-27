package main

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"log"
	"os/exec"
)

type tomlConfig struct {
	Title    string
	Jobsteps map[string]jobstep
}
type jobstep struct {
	Nickname string
	Shell    string
	Command  string
	Args     []string
}
type completedJobstep struct {
	Nickname string
	StdoutStderr []byte
	ExitCode int
	Error error
	//TODO the output and or errors
}


//TODO pass argument with the config file
func main() {
	const defaultConfigFile = "jobStep.toml"
	jobsteps := initJobSteps(defaultConfigFile)
	fmt.Println("vardump: \n", len(jobsteps))

	resultChannel := make(chan completedJobstep) //channel for each execution to write its results

	for _, js := range jobsteps {
		fmt.Printf("starting go routine for step: %s\n", js.Nickname)//step.Nickname)
		go runStep(js, resultChannel)
}

	for i := 0; i < len(jobsteps); i++ {
		result := <-resultChannel
		fmt.Printf("jobstep: %s is finished with status: %d\n", result.Nickname,result.ExitCode)
	}
	fmt.Printf("all %d jobsteps are all done", len(jobsteps))
}

func runStep(step jobstep, resultChannel chan completedJobstep) {
	fmt.Printf("executing step: %s Shell: %s Command: %s Args: %s\n", step.Nickname, step.Shell, step.Command, step.Args)
	stepCommand := exec.Command(step.Shell, step.Command)
	stdoutStdErr, err := stepCommand.CombinedOutput()
	var finished = completedJobstep{
		Nickname: step.Nickname,
		StdoutStderr: stdoutStdErr,
		Error: err,
	}
	//TODO capture output from the process too
	const variableFailCode = 86 //TODO get actual exit code
	const successCode = 0
	if err != nil {
		fmt.Printf("failed result on step named: %s\n", step.Nickname)
		finished.ExitCode = variableFailCode
		resultChannel <- finished
	} else {
		fmt.Printf("success result on step named: %s\n", step.Nickname)
		finished.ExitCode = successCode
		resultChannel <- finished
	}
}

//TODO add a test case

func initJobSteps(configFile string) []jobstep {
	var config tomlConfig
	if _, err := toml.DecodeFile(configFile, &config); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("config created from file: %s title of config: %s\n", configFile, config.Title)

	jobsteps := make([]jobstep, 0, len(config.Jobsteps))

	for _, js := range config.Jobsteps { //no need for name (the map's key)
		jobsteps = append(jobsteps, js)
		fmt.Printf("configured jobstep: %s (%s, %s) with args: %s\n", js.Nickname, js.Shell, js.Command, js.Args)
	}
	return jobsteps
}
