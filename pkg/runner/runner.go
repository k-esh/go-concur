package runner

import (
	"example.com/runner/pkg/config"
	"log"
	"math"
	"os/exec"
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
func NewHampshire(myCommandUnits []config.CommandUnit) []CommandUnitResult{

	resultChannel := make(chan CommandUnitResult) //channel for each execution to write its results
	for _, js := range myCommandUnits {
		go runStep(js, resultChannel)
	}
	//TODO am I doing make corerctly with capacity ?
	cuResults := make([]CommandUnitResult, 0, len(myCommandUnits))
	for i := 0; i < len(myCommandUnits); i++ {
		result := <-resultChannel
		cuResults = append(cuResults, result)
		log.Printf("jobstep: %s has finished with status: %d\n", result.Nickname, result.ExitCode)
	}

	log.Printf("all %d of the jobsteps are done", len(myCommandUnits))
	log.Printf("results: %+v\n", cuResults)
	//TODO should we fail when any of the commands have failed?
	return cuResults
}

const defaultFailureCode = math.MaxInt16

func runStep(step config.CommandUnit, resultChannel chan CommandUnitResult) {
	log.Printf("executing step: %s\n", step.Nickname)

	stepCommand := exec.Command(step.Command, step.Args[0:]...)
	stdoutStdErr, err := stepCommand.CombinedOutput()

	var finished = CommandUnitResult{
		Nickname: step.Nickname,
		Command: step.Command,
		Args: step.Args,
		OutputStdAndError: string(stdoutStdErr),
		Error:             err,
	}

	if err != nil {
		log.Printf("execution result fail on step named: %s with errorType: %T error: %s\n", step.Nickname, err, err)
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
		log.Printf("execution result success on step: %s\n", step.Nickname)

		finished.ExitCode = 0
		resultChannel <- finished
	}
}
