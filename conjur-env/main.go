package main

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/cyberark/conjur-api-go/conjurapi"
	"github.com/cyberark/summon/secretsyml"
)

const (
	// Environment variable names must consist solely of letters (a-z, A-Z),
	// digits (0-9), and underscores, and must not begin with a digit in
	// order to work as part of a Bash export command.
	envVarNameRegExp = "^[a-zA-Z_][a-zA-Z0-9_]*$"
	serviceLabel     = "cyberark-conjur"
)

// Provider is an interface for retrieving Conjur secrets
type Provider interface {
	RetrieveSecret(string) ([]byte, error)
}

// VcapServices implements the json.Unmarshaler interface for VcapServices,
// which allows us to only unmarshal the `cyberark-conjur` service object
// using the ConjurInfo struct.
type VcapServices struct {
	ConjurInfo ConjurInfo
}

// ConjurInfo contains information required for connecting to a Conjur server
type ConjurInfo struct {
	Credentials ConjurCredentials `json:"credentials"`
}

// ConjurCredentials store credentials that are required to connect with
// a Conjur server
type ConjurCredentials struct {
	ApplianceURL   string `json:"appliance_url"`
	APIKey         string `json:"authn_api_key"`
	Login          string `json:"authn_login"`
	Account        string `json:"account"`
	SSLCertificate string `json:"ssl_certificate"`
	Version        int    `json:"version"`
}

func (ci ConjurInfo) setEnv() {
	ci.Credentials.setEnv()
}

func (c ConjurCredentials) setEnv() {
	os.Setenv("CONJUR_APPLIANCE_URL", c.ApplianceURL)
	os.Setenv("CONJUR_AUTHN_LOGIN", c.Login)
	os.Setenv("CONJUR_AUTHN_API_KEY", c.APIKey)
	os.Setenv("CONJUR_ACCOUNT", c.Account)
	os.Setenv("CONJUR_SSL_CERTIFICATE", c.SSLCertificate)
	os.Setenv("CONJUR_VERSION", strconv.Itoa(c.Version))
}

func setConjurCredentialsEnv() error {
	// Get the Conjur connection information from the VCAP_SERVICES
	vcapServices := os.Getenv("VCAP_SERVICES")

	if vcapServices == "" {
		return fmt.Errorf("VCAP_SERVICES environment variable is empty or doesn't exist")
	}

	services := VcapServices{}
	err := json.Unmarshal([]byte(vcapServices), &services)
	if err != nil {
		return fmt.Errorf("error parsing Conjur connection information: %v",
			err.Error())
	}

	conjurInfo := services.ConjurInfo
	conjurInfo.setEnv()

	return nil
}

// newProvider is a function that returns a Provider for retrieving Conjur
// secrets
type newProvider func() (Provider, error)

// NewAPIProvider returns a Conjur API client based on Conjur credentials
// environment variable settings
func NewAPIProvider() (Provider, error) {
	err := setConjurCredentialsEnv()
	if err != nil {
		return nil, err
	}

	config, err := conjurapi.LoadConfig()
	if err != nil {
		return nil, err
	}

	return conjurapi.NewClientFromEnvironment(config)
}

func main() {
	// Get the path of the secrets YAML file.
	secretsYamlPath, exists := os.LookupEnv("SECRETS_YAML_PATH")
	if !exists {
		secretsYamlPath = "secrets.yml"
	}
	secretsEnv, exists := os.LookupEnv("SECRETS_YAML_ENVIRONMENT")
	if !exists {
		secretsEnv = ""
	}

	// Parse the secrets YAML file.
	secrets, err := parseSecretsYamlFile(secretsYamlPath, secretsEnv)
	printAndExitIfError(err)

	// Confirm that environment variable names parsed from the secrets YAML
	// are valid Bash environment variable names.
	err = validateEnvVarNames(secrets)
	printAndExitIfError(err)

	// Create a temporary file factory. No need to defer cleanup because
	// we're injecting values to the environment.
	tempFactory := newTempFactory()

	// Retrieve secrets and generate a concatenation of export statements.
	settings, err := retrieveSecrets(secrets, NewAPIProvider, &tempFactory)
	printAndExitIfError(err)

	// Return the export strings in stdout
	fmt.Print(strings.Join(settings, "\n"))
}

// parseSecretsYamlFile parses a secrets YAML file at a specified path
// and returns an error if either the file doesn't exist or it contains
// invalid secrets YAML syntax.
func parseSecretsYamlFile(secretsYamlPath string, env string) (secretsyml.SecretsMap, error) {
	secrets, err := secretsyml.ParseFromFile(secretsYamlPath, env, nil)
	if os.IsNotExist(err) {
		err = fmt.Errorf("error: %s not found", secretsYamlPath)
	}
	return secrets, err
}

// validateEnvVarNames validates environment variable names that are
// contained in a secretsyml.SecretsMap. It returns an error if any
// environment variable name contains characters other than a-z, A-Z,
// 0-9, or underscores.
func validateEnvVarNames(secrets secretsyml.SecretsMap) error {
	r := regexp.MustCompile(envVarNameRegExp)
	for name := range secrets {
		if !r.MatchString(name) {
			return fmt.Errorf("invalid env variable name: %s", name)
		}
	}
	return nil
}

// retrieveSecrets retrieves secrets and generates export command strings
// by doing the following:
// - Create/load a specified secrets provider
// - Retrieve secrets using that provider based on a given secrets map
// - Create temp files for any secrets that require their values to
//   be stored in a file
// - Generate a concatenation of export command strings that can be used
//   to inject the secrets into a shell environment.
func retrieveSecrets(
	secrets secretsyml.SecretsMap,
	newProvider newProvider,
	tempFactory *tempFactory) ([]string, error) {

	type result struct {
		key   string
		bytes []byte
		error
	}

	// Lazy loading provider
	var provider Provider
	var err error
	for _, spec := range secrets {
		if provider == nil && spec.IsVar() {
			provider, err = newProvider()
			if err != nil {
				return nil, err
			}
		}
	}

	// Run provider calls concurrently
	results := make(chan result, len(secrets))
	var wg sync.WaitGroup

	for key, spec := range secrets {
		wg.Add(1)
		go func(key string, spec secretsyml.SecretSpec) {
			var (
				secretBytes []byte
				err         error
			)

			if spec.IsVar() {
				secretBytes, err = provider.RetrieveSecret(spec.Path)

				if spec.IsFile() {
					fname := tempFactory.push(secretBytes)
					secretBytes = []byte(fname)
				}
			} else {
				// If the spec isn't a variable, use its value as-is
				secretBytes = []byte(spec.Path)
			}

			results <- result{key, secretBytes, err}
			wg.Done()
			return
		}(key, spec)
	}
	wg.Wait()
	close(results)

	// Inline function to generate an export string
	makeSetting := func(result result) string {
		// Base64 encode the value
		encodedValue := base64.StdEncoding.EncodeToString(result.bytes)
		// Create a setting of the form "<variable>: <base64-encoded-value>"
		return fmt.Sprintf("%s: %s", result.key, encodedValue)
	}

	// Generate settings
	var settings []string
	for result := range results {
		if result.error != nil {
			return nil, fmt.Errorf("%s fetching variable: %s",
				result.error,
				result.key)
		}
		settings = append(settings, makeSetting(result))
	}
	return settings, nil
}

func printAndExitIfError(err error) {
	if err == nil {
		return
	}
	os.Stderr.Write([]byte(err.Error()))
	os.Exit(1)
}

// UnmarshalJSON implements the json.Unmarshaler interface for VcapServices,
// which allows us to only unmarshal the `cyberark-conjur` service object
// using the ConjurInfo struct.
func (vcapServices *VcapServices) UnmarshalJSON(b []byte) error {
	services := make(map[string][]interface{})
	err := json.Unmarshal(b, &services)
	if err != nil {
		return err
	}

	conjurServices, ok := services[serviceLabel]
	if !ok || len(conjurServices) == 0 {
		return errors.New("no Conjur services are bound to this application")
	}

	infoBytes, err := json.Marshal(conjurServices[0])
	if err != nil {
		return err
	}

	info := ConjurInfo{}
	err = json.Unmarshal(infoBytes, &info)
	if err != nil {
		return err
	}

	vcapServices.ConjurInfo = info
	return nil
}
