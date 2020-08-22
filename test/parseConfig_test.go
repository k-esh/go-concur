package test

import (
	"example.com/runner/pkg/config"
	"example.com/runner/pkg/runner"
	"testing"
)

const checkMark = "\u2713"
const ballotBox = "\u2717"


func TestFileNotFoundYieldsError(t *testing.T) {
	bogusFile := "/yyyyyoo/bar.yyyyy"
	t.Log("when instantiation from bogus file")

	_, err := config.InitJobSteps(bogusFile)
	if err == nil {
		t.Errorf("\tNew() should give error on bad file %s %s", bogusFile, ballotBox)
	} else {
		t.Logf("\tNew() should give error %s", checkMark)

	}
}

func TestConstruct(t *testing.T) {
	file := "sampleConfig.toml"
	t.Logf("when instantiation from file: %s", file)

	r , err := runner.New(file)
	if err != nil {
		t.Errorf("\tshould be able to instantiate a runner %s, %s", err, ballotBox)
	} else {
		t.Logf("\tshould instantiate a: %T %s", r, checkCheck)
		//here call to a.Start() is possible but wont do in this test
	}
}



	func TestSampleFileParses(t *testing.T) {
	configFile := "sampleConfig.toml" //the default config file  TODO put in folder test
	testConfig, err := config.InitJobSteps(configFile)
	if err != nil {
		t.Errorf("unexpected error parsing file; sample file: %s error: %v", configFile, err)
	}
	//expect three command elements
	if len(testConfig.CommandUnits) != 3 {
		t.Errorf("test config parsed incorrectly size wanted 3, got: %d", len(testConfig.CommandUnits))
	}
	//expect maxRuntime is here as expected
	if testConfig.MaxRuntime != 120 {
		t.Errorf("test config maxRuntime wanted 120, got: %d", testConfig.MaxRuntime)
	}


	//for each CommandUnit we have put in the sample test file, assert is as expected
	///- firstCommandUnit value {/bin/sh [date]}
	first := testConfig.CommandUnits["firstCommandUnit"]
	if first.Command != "date" {
		t.Errorf("firstCommandUnit name wanted date, got: %s", first.Command)
	}
	if first.Args != nil {
		t.Errorf("firstCommandUnit argument wanted nil, got: %v", first)
	}

	///- favoriteCommandUnit value {echo [hi gopher]}
	fave := testConfig.CommandUnits["favoriteCommandUnit"]
	if fave.Command != "echo" {
		t.Errorf("favoriteCommandUnit name wanted echo, got: %s", fave.Command)
	}
	if len(fave.Args) != 2 {
		t.Errorf("favoriteCommandUnit argument count wanted 2 , got: %d", len(fave.Args))
	}
	if fave.Args[0] != "hi" {
		t.Errorf("favoriteCommandUnit solo argument wanted date, got: %s", fave.Args[0])
	}
	if fave.Args[1] != "gopher" {
		t.Errorf("favoriteCommandUnit solo argument wanted gopher, got: %s", fave.Args[1])
	}



}
