package runner

import (
	"errors"
	"example.com/runner/pkg/config"
	"log"
	"math"
	"os/exec"
)

type CommandUnitResult struct {
	Name              string
	Command           string
	Args              []string
	OutputStdAndError string
	ExitCode          int
	Error             error
}

const defaultFailureCode = math.MaxInt16

//TODO logging in json format maybe zip library
func NewJersey(commandDocket config.CommandConfig) ([]CommandUnitResult, error) {
	//get out error if the config in has zero to do
	if len(commandDocket.CommandUnits) == 0 {
		return nil, errors.New("invalid config has zero elements to run")
	}

	resultChannel := make(chan CommandUnitResult) //channel for each execution to write its results
	for name, commandElement := range commandDocket.CommandUnits {
		go runStep(name, commandElement, resultChannel)
	}
	//TODO am I doing make corerctly with capacity ?
	cuResults := make([]CommandUnitResult, 0, len(commandDocket.CommandUnits))
	for i := 0; i < len(commandDocket.CommandUnits); i++ {
		result := <-resultChannel
		cuResults = append(cuResults, result)
		log.Printf("jobstep: %s has finished with status: %d\n", result.Name, result.ExitCode)
	}

	log.Printf("all %d of the jobsteps are done", len(commandDocket.CommandUnits))
	log.Printf("results: %+v\n", cuResults)
	return cuResults, nil
}

func runStep(name string, step config.CommandUnit, resultChannel chan CommandUnitResult) {
	log.Printf("executing step: %s\n", name)

	stepCommand := exec.Command(step.Command, step.Args[0:]...)
	stdoutStdErr, err := stepCommand.CombinedOutput()

	//capture result into the struct; do not set the exit code as we do that below after check err
	var finishedCommandResult = CommandUnitResult{
		Name:              name,
		Command:           step.Command,
		Args:              step.Args,
		OutputStdAndError: string(stdoutStdErr),
		Error:             err,
	}

	if err != nil {
		log.Printf("execution result fail on step named: %s with errorType: %T error: %s\n", name, err, err)
		// If the command starts but does not complete successfully, the error is of
		// type *ExitError. Other error types may be returned for other situations.
		/// is if it is ExitError and if so rely on its ExitCode otherwise do what we can w/ default
		_, ok := err.(*exec.ExitError)
		if true == ok {
			finishedCommandResult.ExitCode = err.(*exec.ExitError).ExitCode()
		} else {
			finishedCommandResult.ExitCode = defaultFailureCode
		}
	} else {
		log.Printf("execution result success on step: %s\n", name)

		finishedCommandResult.ExitCode = 0
	}
	resultChannel <- finishedCommandResult
}
