package main

import (
	"fmt"
	"sync"
	"github.com/cyberark/summon/secretsyml"
	"github.com/cyberark/conjur-api-go/conjurapi"
	"os"
	"strings"
	"io/ioutil"
	"encoding/json"
)

type Provider interface {
	RetrieveSecret(string) ([]byte, error)
}

type CatProvider struct {
}
func (CatProvider) RetrieveSecret(path string) ([]byte, error) {
	return ioutil.ReadFile(path)
}

type ConjurInfo struct {
	Credentials ConjurCredentials `json:"credentials"`
}

type ConjurCredentials struct {
	ApplianceURL     string `json:"appliance_url"`
	APIKey     		 string `json:"authn_api_key"`
	Login     		 string `json:"authn_login"`
	Account 		 string `json:"account"`
}

func (c ConjurCredentials) setEnv() {
	os.Setenv("CONJUR_APPLIANCE_URL", c.ApplianceURL)
	os.Setenv("CONJUR_AUTHN_LOGIN", c.Login)
	os.Setenv("CONJUR_AUTHN_API_KEY", c.APIKey)
	os.Setenv("CONJUR_ACCOUNT", c.Account)
}

const SERVICE_LABEL="cyberark-conjur"
func setConjurCredentialsEnv() error {
	// Get the Conjur connection information from the VCAP_SERVICES
	s := os.Getenv("VCAP_SERVICES")

	if s == "" {
		return fmt.Errorf("VCAP_SERVICES environment variable not found\n")
	}

	services := make(map[string][]ConjurInfo)
	err := json.Unmarshal([]byte(s), &services)
	if err != nil {
		return fmt.Errorf("Error parsing Conjur connection information: %v\n", err.Error())
	}

	info := services[SERVICE_LABEL]
	if len(info) == 0 {
		return fmt.Errorf("No Conjur services are bound to this application.\n")
	}

	creds := info[0].Credentials
	creds.setEnv()

	return nil
}

func NewProvider() (Provider, error) {
	//return CatProvider{}, nil
	err := setConjurCredentialsEnv()
	if err != nil {
		return nil, err
	}

	return conjurapi.NewClientFromEnvironment(conjurapi.LoadConfig())
}

func main() {
	var (
		provider Provider
		err error
	)

	var (
		secrets secretsyml.SecretsMap
	)

	secrets, err = secretsyml.ParseFromFile("secrets.yml", "", nil)
	if os.IsNotExist(err) {
		printAndExitIfError(fmt.Errorf("secrets.yml not found\n"))
	}
	printAndExitIfError(err)

	tempFactory := NewTempFactory("")
	// defer tempFactory.Cleanup()
	// no need to cleanup because we're injecting values to the environment

	type Result struct {
		key string
		bytes []byte
		error
	}

	// Run provider calls concurrently
	results := make(chan Result, len(secrets))
	var wg sync.WaitGroup

	for _, spec := range secrets {
		if provider == nil && spec.IsVar() {
			provider, err = NewProvider()
			printAndExitIfError(err)
		}
	}

	for key, spec := range secrets {
		wg.Add(1)
		go func(key string, spec secretsyml.SecretSpec) {
			var (
				secretBytes []byte
				err error
			)

			if spec.IsVar() {
				secretBytes, err = provider.RetrieveSecret(spec.Path)

				if spec.IsFile() {
					fname := tempFactory.Push(secretBytes)
					secretBytes = []byte(fname)
				}
			} else {
				// If the spec isn't a variable, use its value as-is
				secretBytes = []byte(spec.Path)
			}

			results <- Result{key, secretBytes, err}
			wg.Done()
			return
		}(key, spec)
	}
	wg.Wait()
	close(results)

	var exportStrings []string

	for result := range results {
		if result.error == nil {
			exportString := fmt.Sprintf("export %s='%s';", result.key, strings.Replace(string(result.bytes), "'", "'\"'\"'", -1))
			exportStrings = append(exportStrings, exportString)
		} else {
			printAndExitIfError(fmt.Errorf("error fetching variable - %s", result.error))
		}
	}

	fmt.Print(strings.Join(exportStrings, "\n"))
}

func printAndExitIfError(err error) {
	if err == nil {
		return
	}
	os.Stderr.Write([]byte(err.Error()))
	os.Exit(1)
}
