package test

import (
	"example.com/runner/pkg/config"
	"testing"
)

const checkMark = "\u2713"

func TestEmptyFileFails(t *testing.T) {
	bogusFile := "/yyyyyoo/bar.yyyyy"
	_, err := config.InitJobSteps(bogusFile)
	if err == nil {
		t.Error("expected error but did not get it; passed file that does not exist")
	} else {
		t.Log("should be error on invalid file as config", checkMark)

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
	if first.Command != "/bin/sh" {
		t.Errorf("firstCommandUnit name wanted /bin/sh, got: %s", first.Command)
	}
	if first.Args[0] != "date" {
		t.Errorf("firstCommandUnit solo argument wanted date, got: %s", first.Args[0])
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
