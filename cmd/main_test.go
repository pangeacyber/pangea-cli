package main_test

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

const pangeaCLICommand = "pangea"

func TestMain(m *testing.M) {
	// Run setup code here
	setUp()

	// Run tests
	exitCode := m.Run()

	// Run teardown code here

	// Exit with the test exit code
	os.Exit(exitCode)
}

func printArgs(args ...string) string {
	r := ""
	for _, a := range args {
		r = fmt.Sprintf("%s%s ", r, a)
	}
	return r
}

func printCaller() {
	// Get the caller's program counter
	_, file, line, ok := runtime.Caller(2)
	if ok {
		fmt.Printf("Called from file %s line %d\n", file, line)
	}
}

func setUp() {
	// Build CLI
	cmd := exec.Command("go", "build", "-o", pangeaCLICommand, "./main.go")
	output, err := cmd.Output()
	if err != nil {
		log.Fatalf("Error executing 'go build' command: %v.\n", err)
	}

	if len(output) != 0 {
		log.Fatalf("Error executing 'go build' command. Output: %s\n", string(output))
	}
}

func runRaw(args ...string) string {
	cmd := exec.Command("./"+pangeaCLICommand, args...)
	output, err := cmd.Output()
	if err != nil {
		log.Fatalf("Error executing raw CLI command: %v", err)
	}
	return string(output)
}

func run(args ...string) map[string]any {
	cmd := exec.Command("./"+pangeaCLICommand, args...)
	output, err := cmd.Output()
	if err != nil {
		log.Fatalf("Error executing CLI command: '%s'. Error: %v.\n", printArgs(args...), err)
	}

	if strings.Contains(string(output), "APIerror") {
		fmt.Println(printArgs(args...))

		// Get the caller's program counter
		_, file, line, ok := runtime.Caller(1)
		if ok {
			fmt.Printf("Called from file %s line %d\n", file, line)
		}
		log.Fatalf("API error in CLI command:\n%s\n", string(output))
	}

	var result map[string]interface{}
	if err := json.Unmarshal(output, &result); err != nil {
		fmt.Println(printArgs(args...))
		printCaller()
		log.Fatalf("Error parsing result:\n%s\n", string(output))
	}
	return result
}

func TestRunCLI(t *testing.T) {
	output := runRaw()
	assert.NotEmpty(t, output)
}
