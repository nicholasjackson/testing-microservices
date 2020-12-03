package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"time"

	"github.com/cucumber/godog"
	"github.com/cucumber/messages-go/v10"
	"github.com/nicholasjackson/testing-microservices/handlers"
)

func main() {
	opts := godog.Options{
		Format:    "pretty",
		Paths:     []string{"features"},
		Randomize: time.Now().UTC().UnixNano(), // randomize scenario execution order
	}

	status := godog.TestSuite{
		Name:                "godogs",
		ScenarioInitializer: InitializeScenario,
		Options:             &opts,
	}.Run()

	os.Exit(status)
}

func startApplication() *exec.Cmd {
	cmd := exec.Command(
		"go",
		"build",
		"-o",
		"./bin/testbuild",
		"main.go",
	)

	cmd.Dir = "../"
	cmd.Env = os.Environ()

	err := cmd.Start()
	if err != nil {
		panic(err)
	}

	err = cmd.Wait()
	if err != nil {
		panic(err)
	}

	cmd = exec.Command(
		"./bin/testbuild",
	)

	cmd.Dir = "../"
	cmd.Env = append(os.Environ(), "DB_PORT=5433")
	cmd.Stdout = logWriter
	cmd.Stderr = logWriter

	err = cmd.Start()
	if err != nil {
		panic(err)
	}

	return cmd
}

func startDB() *exec.Cmd {
	cmd := exec.Command(
		"docker",
		"run",
		"-d",
		"--name", "testing",
		"-e", "POSTGRES_DB=root",
		"-e", "POSTGRES_USER=root",
		"-e", "POSTGRES_PASSWORD=password",
		"-p", "5433:5432",
		"shipyardrun/postgres:9.6",
	)
	cmd.Stdout = logWriter
	cmd.Stderr = logWriter

	err := cmd.Start()
	if err != nil {
		panic(err)
	}

	err = cmd.Wait()
	if err != nil {
		panic(err)
	}

	return cmd
}

func stopDB() {
	cmd := exec.Command(
		"docker",
		"rm",
		"-f",
		"testing",
	)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Start()
	if err != nil {
		panic(err)
	}

	err = cmd.Wait()
	if err != nil {
		panic(err)
	}
}

func stopCmd(c *exec.Cmd) {
	c.Process.Kill()
}

var responseData []byte

// redirect stdout and stderr to a buffer
var logWriter *bytes.Buffer

func InitializeScenario(ctx *godog.ScenarioContext) {
	var cmdApp *exec.Cmd

	ctx.BeforeScenario(func(*godog.Scenario) {
		responseData = nil
		logWriter = bytes.NewBufferString("")

		cmdApp = startApplication()
		startDB()
	})

	ctx.AfterScenario(func(sc *messages.Pickle, err error) {
		stopCmd(cmdApp)
		stopDB()

		// only write the output log on an error or when debug is enabled
		if err != nil || os.Getenv("DEBUG") == "true" {
			fmt.Printf(logWriter.String())
		}
	})

	ctx.Step(`^I call the Get endpoint$`, iCallTheGetEndpoint)
	ctx.Step(`^I call the Insert endpoint with the folowing JSON$`, iCallTheInsertEndpointWithTheFolowingJSON)
	ctx.Step(`^the application is running$`, theApplicationIsRunning)
	ctx.Step(`^there should be (\d+) branches returned$`, thereShouldBeBranchesReturned)
}

func iCallTheGetEndpoint() error {
	resp, err := http.Get("http://localhost:9090/branches")
	if err != nil {
		return err
	}

	responseData, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func theApplicationIsRunning() error {
	// test the app is running, call the health endpoint
	for n := 0; n < 30; n++ {
		resp, _ := http.Get("http://localhost:9090/health")

		if resp != nil && resp.StatusCode == http.StatusOK {
			return nil
		}

		time.Sleep(1 * time.Second)
	}

	return fmt.Errorf("Error waiting for service to start")
}

func thereShouldBeBranchesReturned(arg1 int) error {
	branches := []handlers.Branch{}
	json.Unmarshal(responseData, &branches)

	if len(branches) != arg1 {
		return fmt.Errorf("Expected %d branches, got %d, data: %s", arg1, len(branches), string(responseData))
	}

	return nil
}

func iCallTheInsertEndpointWithTheFolowingJSON(arg1 *messages.PickleStepArgument_PickleDocString) error {
	resp, err := http.Post("http://localhost:9090/branches", "application/json", bytes.NewBufferString(arg1.Content))
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	return nil
}
