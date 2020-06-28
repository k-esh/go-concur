package main

import (
	"example.com/runner/pkg/config"
	"log"
	"math"
	"os"
	"os/exec"
)



//TODO rename to better name, and do not export it if possible
type Jobstep struct {
	Nickname string
	Command  string
	Args     []string
	//TODO support a max run elapsed time
}

//TODO do not embed here to keep mine private; instead have command and args only
type JobStepResult struct {
	wrappedJopstep    config.CommandUnit
	outputStdAndError string
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
	myCommandUnits := config.InitJobSteps(configFile) //print each thing we are going to run on one line
	//TODO see if there is a concise functional way to do this
	theNicknamesToRun := make([]string, 0, len(myCommandUnits))
	for _, js := range myCommandUnits {
		theNicknamesToRun = append(theNicknamesToRun, js.Nickname)
	}
	log.Printf("running in parallel:: %s\n", theNicknamesToRun)

	resultChannel := make(chan JobStepResult) //channel for each execution to write its results
	for _, js := range myCommandUnits {
		go runStep(js, resultChannel)
	}


//TODO return structure that summarizes the steps and exit codes
	jsResults := make([]JobStepResult, 0, len(myCommandUnits))

	for i := 0; i < len(myCommandUnits); i++ {
		result := <-resultChannel
		jsResults = append(jsResults, result)
		log.Printf("jobstep: %s has finished with status: %d\n", result.wrappedJopstep.Nickname, result.ExitCode)
	}

	log.Printf("all %d of the jobsteps are done", len(myCommandUnits))
	log.Printf("results: %+v\n", jsResults)
}

const defaultFailureCode = math.MaxInt16

func runStep(step config.CommandUnit, resultChannel chan JobStepResult) {
	log.Printf("executing step: %s\n", step.Nickname)

	stepCommand := exec.Command(step.Command, step.Args[0:]...)
	stdoutStdErr, err := stepCommand.CombinedOutput()

	var finished = JobStepResult{
		wrappedJopstep:    step,
		outputStdAndError: string(stdoutStdErr),
		Error:             err,
	}

	if err != nil {
		log.Printf("result failed on step named: %s with errorType: %T error: %s\n", step.Nickname, err, err)
		// If the command starts but does not complete successfully, the error is of
		// type *ExitError. Other error types may be returned for other situations.
		/// is if it is ExitError and if so rely on its ExitCode otherwise do what we can w/ default
		_, ok := err.(*exec.ExitError)
		if true == ok {
			finished.ExitCode = err.(*exec.ExitError).ExitCode()
		} else {
			finished.ExitCode = defaultFailureCode
		}

		resultChannel <- finished
	} else {
		log.Printf("result success on step: %s\n", step.Nickname)

		finished.ExitCode = 0
		resultChannel <- finished
	}
}

//TODO add a test case

