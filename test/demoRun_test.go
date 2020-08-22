package test

import (
	"example.com/runner/pkg/runner"
	"testing"
)

const checkCheck = checkMark
const ballotBox = "\u2717"

func TestConstruct(t *testing.T) {
	file := "/Users/mac/craft/innov/goBatch/test/sampleConfig.toml"

	t.Logf("when instantiation from file: %s", file)

	r , err := runner.New(file)
	if err != nil {
		t.Errorf("\tshould be able to instantiate a runner %s, %s", err, ballotBox)
	} else {
		t.Logf("\tshould instantiate a: %T %s", r, checkCheck)
		//here call to a.Start() is possible but wont do in this test
	}

}
