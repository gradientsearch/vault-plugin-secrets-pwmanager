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
	"crypto/rand"
	"fmt"
	"os/exec"
	"testing"
	"time"

	"github.com/hashicorp/vault/api"
)

const rootToken = "root"
const SUCCESS string = "😃"
const FAILURE string = "💣"

func BuildPlugin(name string) error {

	app := "go"
	arg0 := "build"
	arg1 := "-o"
	arg2 := fmt.Sprintf("vault/plugins/%s/pwmanager", name)
	arg3 := "cmd/vault-plugin-secrets-pwmanager/main.go"
	cmd := exec.Command(app, arg0, arg1, arg2, arg3)
	stdout, err := cmd.Output()

	fmt.Println(string(stdout))

	return err

}

// StartDB starts a database instance.
func StartDB(name string, t *testing.T) (Container, error) {
	const address = "0.0.0.0:8200"
	const token = rootToken
	const image = "hashicorp/vault:1.18.3"
	const port = "8200"

	dockerArgs := []string{"-e", "VAULT_DEV_ROOT_TOKEN_ID=" + token, "-e", "VAULT_DEV_LISTEN_ADDRESS=" + address, "-v", fmt.Sprintf("./vault/plugins/%s:/plugins", name)}
	appArgs := []string{"server", "-dev", "-dev-root-token-id=root", "-dev-plugin-dir=/plugins", "-log-level=debug"}

	c, err := StartContainer(image, name, port, dockerArgs, appArgs)
	if err != nil {
		return Container{}, fmt.Errorf("starting container: %w", err)
	}

	t.Logf("Image:       %s\n", image)
	t.Logf("ContainerID: %s\n", c.Name)
	t.Logf("Host:        %s\n\n", c.HostPort)

	return c, nil
}

// StopDB stops a running database instance.
func StopDB(c *Container, t *testing.T) {
	_ = StopContainer(c.Name)
	t.Log("Stopped:", c.Name)
}

// Tail container logs will stream docker logs -f <container name>
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

// TestHarness owns state for running and shutting down tests.
type TestHarness struct {
	Log       *Logger
	Buff      *bytes.Buffer
	Teardown  func()
	Testing   *testing.T
	Client    *pwmanagerClient
	Container *Container
}

// NewTestHarness creates a test Vault Server inside a Docker container. It returns
// the Vault root client to use as well as a teardown function to call.
func NewTestHarness(t *testing.T, name string, tailContainer bool) (*TestHarness, error) {
	t.Logf("******************** LOGS (%s) ********************\n", name)

	if err := BuildPlugin(name); err != nil {
		t.Fatalf("failed to build plugin: %s", err)
	}

	c, err := StartDB(name, t)
	if err != nil {
		return nil, err
	}

	if tailContainer {
		go TailContainerLogs(c)
	}

	// teardown is the function that should be invoked when the caller is done
	// with the database.
	teardown := func() {
		t.Helper()
		StopDB(&c, t)
	}

	var buf bytes.Buffer
	log := NewLogger(&buf, LevelInfo, name, func(context.Context) string { return "00000000-0000-0000-0000-000000000000" })

	test := TestHarness{
		Log:      log,
		Buff:     &buf,
		Teardown: teardown,
		Testing:  t,

		Container: &c,
	}

	test.WithClient(rootToken)

	ctx, _ := context.WithTimeout(context.Background(), time.Second*10)
	test.WaitForDB(ctx)

	return &test, nil
}

// WaitForDB waits for vault to be ready then returns
func (t *TestHarness) WaitForDB(ctx context.Context) {
	for attempts := 1; ; attempts++ {
		var err error
		if _, err = t.Client.c.Sys().Health(); err == nil {
			t.Testing.Log("Connected To Vault")
			break
		}

		t.Testing.Log("Waiting For Vault")

		if ctx.Err() != nil {
			t.Testing.Fatalf("context error aborting test: %s", err)
		}

		time.Sleep(time.Duration(attempts) * 100 * time.Millisecond)

		if ctx.Err() != nil {
			t.Testing.Fatalf("context error aborting: %s", err)
		}
	}
}

// Add pwmanagerClient to this testHarness
func (t *TestHarness) WithClient(token string) {
	if client, err := NewClient(token, t.Container.HostPort); err != nil {
		t.Testing.Fatal(err)
	} else {
		t.Client = client
	}

}

// Add a pwmanager secrets mount to vault server
func (t *TestHarness) WithPwManagerMount() {
	mi := api.MountInput{
		Type:        "pwmanager",
		Description: "password manager for users",
	}

	if err := t.Client.c.Sys().Mount("/pwmanager", &mi); err != nil {
		t.Testing.Fatalf("failed to create pwmanager mount")
	}
	t.Testing.Log("successfully created pwmanager mount")
}

// Add policies
func (t *TestHarness) WithPolicies(policies map[string]string) {
	for k, v := range policies {
		if err := t.Client.c.Sys().PutPolicy(k, v); err != nil {
			t.Testing.Fatalf("error creating policy %s: %s", k, err)
		}
		t.Testing.Logf("successfully created %s policy\n", k)
	}
}

// add userpass auth mount to vault server with users
func (t *TestHarness) WithUserpassAuth(mount string, users []string, adminUser string) map[string]TestUser {
	if err := t.Client.c.Sys().EnableAuth("/userpass", "userpass", "userpass used for pwmanager users"); err != nil {
		t.Testing.Fatalf("failed to create userpass mount")
	}
	t.Testing.Log("successfully created /userpass mount ")

	lrs := map[string]TestUser{}
	for _, u := range users {

		userInfo := UserInfo{
			Password:      "gophers",
			TokenPolicies: []string{"default", fmt.Sprintf("%s/user/default", mount), fmt.Sprintf("%s/entity/%s", mount, u)},
		}

		if u == adminUser {
			userInfo.TokenPolicies = append(userInfo.TokenPolicies, fmt.Sprintf("%s/admin/default", mount))
		}

		if err := t.Client.Userpass().User("userpass", u, userInfo); err != nil {
			t.Testing.Fatalf("failed to create user %s", err)
		}
		t.Testing.Logf("successfully created user %s in /userpass\n", u)

		if lr, err := t.Client.Userpass().Login("userpass", u, userInfo); err != nil {
			t.Testing.Fatalf("failed to create user %s", err)
		} else {
			t.Testing.Logf("successfully logged in user %s to /userpass\n", u)
			tu := TestUser{LoginResponse: lr, PwManagerMount: mount}
			tu.WithClient(t)
			lrs[u] = tu
		}
	}

	return lrs
}

// func policy() {
// 	// user has full access to all their safes. Safes are pathed using entity-id
// 			// so if my id is 1 i have access to all safes under <mount>/1/*. a safe would be
// 			// <mount>/1/uuid (we want to keep the safe name secret). The name of the safe will
// 			// be stored in the safe metadata the user will be able to decrypt and the client will
// 			// display the actual safe name
// 			userSafesPaths := fmt.Sprintf("%s/%s", "pwmanager", lr.Auth.EntityID)
// 			userAccess := Access{lr.Auth.EntityID, map[SafePath]Capabilities{SafePath(userSafesPaths): Capabilities{"create", "read", "update", "patch", "delete", "list"}}}
// 			defaultUserPolicyPath := "policies/pwmanager_user_default.tmpl"
// 			fs, err := os.OpenFile(defaultUserPolicyPath, os.O_RDONLY, 0444)
// 			if err != nil {
// 				t.Testing.Fatalf("error opening %s policy: %s", defaultUserPolicyPath, err)
// 			}

// 			tmplFile, err := io.ReadAll(fs)
// 			if err != nil {
// 				t.Testing.Fatalf("error reading %s file: %s", defaultUserPolicyPath, err)
// 			}

// 			tmpl, err := template.New("test").Parse(string(tmplFile))
// 			if err != nil {
// 				t.Testing.Fatalf("error parsing template: %s", err)
// 			}
// 			var b bytes.Buffer
// 			err = tmpl.Execute(&b, userAccess)
// 			if err != nil {
// 				t.Testing.Fatalf("error executing template: %s", err)
// 			}

// }
type TestUser struct {
	LoginResponse  LoginResponse
	PwManagerMount string
	Client         *pwmanagerClient
	UUK            UUK
	SecretData     TestUsersSecrets
}

type TestUsersSecrets struct {
	Password  string
	SecretKey []byte
}

// add pwmangerClient to this testuser
func (t *TestUser) WithClient(th *TestHarness) {
	if client, err := NewClient(t.LoginResponse.Auth.ClientToken, th.Container.HostPort); err != nil {
		th.Testing.Fatal(err)
	} else {
		t.Client = client
	}
}

func (t *TestUser) WithUUK(th *TestHarness) {
	uuk := UUK{}
	userSecrets := TestUsersSecrets{}
	userSecrets.Password = "gophers"
	userSecrets.SecretKey = make([]byte, 32)
	if _, err := rand.Read(userSecrets.SecretKey); err != nil {
		th.Testing.Fatalf("error creating secret key: %s; ", err)
	}

	uuk.Build([]byte(userSecrets.Password), []byte(t.PwManagerMount), userSecrets.SecretKey, []byte(t.LoginResponse.Auth.EntityID))
	t.UUK = uuk
}

type SafePath string
type Capabilities []string
type Access struct {
	EntityID string
	Safes    map[SafePath]Capabilities
}
