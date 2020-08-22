package main

import
(
	"example.com/runner/pkg/runner"
	"fmt"
	"os"
)


func main() {
	if len(os.Args)!=2 {
		panic("missing argument for configuration file location")
	}

	file := os.Args[1] //"/Users/mac/craft/innov/goBatch/test/sampleConfig.toml"
	runner, err := runner.New(file)
	if err != nil {
		panic(err)
	}

	result, err := runner.Start()
	if err != nil {
		panic(err)
	}
	fmt.Printf("result is %v", result)

	
}
