/*
 * Copyright 2024 Ardan Labs
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * This file is part of the [Service] project by Ardan Labs.
 * Repository URL: https://github.com/ardanlabs/service
 *
 * Changes Made:
 * - Stephen O'Dwyer - originally dbtest.go. Transformed it to work with vault
 * For more information, see the repository's changelog or commit history.
 */
package secretsengine

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"testing"
	"time"

	vault "github.com/hashicorp/vault/api"
)

const token = "root"
const SUCCESS string = "😃"
const FAILURE string = "😅"

func BuildPlugin(name string) error {
	fmt.Printf("******************** LOGS (%s) ********************\n", "build output")

	app := "go"
	arg0 := "build"
	arg1 := "-o"
	arg2 := fmt.Sprintf("vault/plugins/%s/pwmanager", name)
	arg3 := "cmd/vault-plugin-secrets-pwmanager/main.go"
	fmt.Println("running build command: ", app, arg0, arg1, arg2, arg3)
	cmd := exec.Command(app, arg0, arg1, arg2, arg3)
	stdout, err := cmd.Output()

	fmt.Println(string(stdout))

	return err

}

func TailContainerLogs(c Container) {
	app := "docker"
	arg0 := "logs"
	arg1 := "-f"
	arg2 := c.Name
	cmd := exec.Command(app, arg0, arg1, arg2)

	stderr, _ := cmd.StderrPipe()
	cmd.Start()

	scanner := bufio.NewScanner(stderr)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		m := scanner.Text()
		fmt.Println(m)
	}
	cmd.Wait()
}

// StartDB starts a database instance.
func StartDB(name string) (Container, error) {
	const address = "0.0.0.0:8200"
	const token = token
	const image = "hashicorp/vault:1.18.3"
	const port = "8200"

	dockerArgs := []string{"-e", "VAULT_DEV_ROOT_TOKEN_ID=" + token, "-e", "VAULT_DEV_LISTEN_ADDRESS=" + address, "-v", fmt.Sprintf("./vault/plugins/%s:/plugins", name)}
	appArgs := []string{"server", "-dev", "-dev-root-token-id=root", "-dev-plugin-dir=/plugins", "-log-level=debug"}

	c, err := StartContainer(image, name, port, dockerArgs, appArgs)
	if err != nil {
		return Container{}, fmt.Errorf("starting container: %w", err)
	}

	fmt.Printf("Image:       %s\n", image)
	fmt.Printf("ContainerID: %s\n", c.Name)
	fmt.Printf("Host:        %s\n", c.HostPort)

	return c, nil
}

// StopDB stops a running database instance.
func StopDB(c *Container) {
	_ = StopContainer(c.Name)
	fmt.Println("Stopped:", c.Name)
}

// WaitForDB waits for vault to be ready then returns
func WaitForDB(ctx context.Context, api *vault.Client, t *testing.T) {
	for attempts := 1; ; attempts++ {
		var err error
		if _, err = api.Sys().Health(); err == nil {
			t.Log("Connected To Vault")
			break
		}

		t.Log("Waiting For Vault")

		if ctx.Err() != nil {
			t.Fatalf("context error aborting test: %s", err)
		}

		time.Sleep(time.Duration(attempts) * 100 * time.Millisecond)

		if ctx.Err() != nil {
			t.Fatalf("context error aborting: %s", err)
		}
	}
}

// Test owns state for running and shutting down tests.
type Test struct {
	Log      *Logger
	Buff     *bytes.Buffer
	Teardown func()
	t        *testing.T
	client   *pwmanagerClient
	c        *Container
}

// NewTestHarness creates a test Vault Server inside a Docker container. It returns
// the Vault client to use as well as a function to call at the end of the test.
func NewTestHarness(t *testing.T, name string) (*Test, error) {
	if err := BuildPlugin(name); err != nil {
		t.Fatalf("failed to build plugin: %s", err)
	}

	c, err := StartDB(name)
	if err != nil {
		return nil, err
	}

	go TailContainerLogs(c)

	// teardown is the function that should be invoked when the caller is done
	// with the database.
	teardown := func() {
		t.Helper()
		StopDB(&c)
		fmt.Printf("******************** LOGS (%s) ********************\n", name)
	}

	config := vault.DefaultConfig()
	config.Address = "http://" + c.HostPort

	v, err := vault.NewClient(config)
	if err != nil {
		t.Fatalf("unable to initialize Vault client: %v", err)
	}

	// Authenticate
	// WARNING: This is just for testing.
	// Don't do this in production!
	v.SetToken(token)
	if err != nil {
		t.Fatalf("error connecting to vault: %s", err)
	}

	var buf bytes.Buffer
	log := NewLogger(&buf, LevelInfo, name, func(context.Context) string { return "00000000-0000-0000-0000-000000000000" })

	client := NewPwmanagerClient(v)

	test := Test{
		Log:      log,
		Buff:     &buf,
		Teardown: teardown,
		t:        t,
		client:   client,
		c:        &c,
	}
	ctx, _ := context.WithTimeout(context.Background(), time.Second*10)
	WaitForDB(ctx, v, t)

	return &test, nil
}

// StringPointer is a helper to get a *string from a string. It is in the tests
// package because we normally don't want to deal with pointers to basic types
// but it's useful in some tests.
func StringPointer(s string) *string {
	return &s
}

// IntPointer is a helper to get a *int from a int. It is in the tests package
// because we normally don't want to deal with pointers to basic types but it's
// useful in some tests.
func IntPointer(i int) *int {
	return &i
}

// FloatPointer is a helper to get a *float64 from a float64. It is in the tests
// package because we normally don't want to deal with pointers to basic types
// but it's useful in some tests.
func FloatPointer(f float64) *float64 {
	return &f
}
