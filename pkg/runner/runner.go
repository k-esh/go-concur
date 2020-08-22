package runner

import (
	"errors"
	"example.com/runner/pkg/config"
	"log"
	"math"
	"os/exec"
)

type Runner struct {
	config config.CommandConfig
}
//TODO can return a parent child / header childs...nicer than slice
type CommandUnitResult struct {
	Name              string
	Command           string
	Args              []string
	OutputStdAndError string
	ExitCode          int
	Error             error
}

const defaultFailureCode = math.MaxInt16

//parses input file returns ready to use runner, or parse file error
func New(configFile string) (*Runner , error) {
	//if any inputfile/parsing errors return error here
	configuration, err := config.InitJobSteps(configFile)
	if err != nil {
		return nil, err
	}
	if len(configuration.CommandUnits) == 0 {
		return nil, errors.New("invalid config has zero elements to be run")
	}


	return &Runner{
		config: configuration,
	}, nil
}

func (r *Runner) Start () ([]CommandUnitResult, error) {
	docket := r.config
	resultChannel := make(chan CommandUnitResult) //channel for each execution to write its results
	for name, commandElement := range docket.CommandUnits {
		go runStep(name, commandElement, resultChannel)
	}

	cuResults := make([]CommandUnitResult, 0, len(docket.CommandUnits))
	for i := 0; i < len(docket.CommandUnits); i++ {
		result := <-resultChannel
		cuResults = append(cuResults, result)
		log.Printf("jobstep: %s has finished with status: %d\n", result.Name, result.ExitCode)
	}

	log.Printf("all %d of the jobsteps are done", len(docket.CommandUnits))

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
