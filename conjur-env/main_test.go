package main

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/cyberark/summon/secretsyml"
	"github.com/stretchr/testify/assert"
)

const (
	validSecretsYaml = `
SOME_VAR: !var prod/user/robot/api_key
SOME_VAR_FILE: !var:file prod/user/robot/private_key
SOME_LOCAL_FILE: !file my-local-file`

	validSecretsYamlWithEnv = `
common:
  COMMON_VAR: common_secret

wanted:
  SOME_VAR: !var prod/user/robot/api_key

unwanted:
  SOME_VAR: !var prod/user/robot/my_secret`

	invalidSecretsYaml = `
SOME_VAR: !var prod/user/robot/api_key
OOPS_NO_COLON !var prod/user/robot/my_secret`
)

type exportCmd struct {
	format          string
	tempFileContent string
}

var (
	validSecretsMap = secretsyml.SecretsMap{
		"SOME_VAR": {
			Tags: []secretsyml.YamlTag{secretsyml.Var},
			Path: "prod/user/robot/api_key",
		},
		"SOME_VAR_FILE": {
			Tags: []secretsyml.YamlTag{secretsyml.Var, secretsyml.File},
			Path: "prod/user/robot/private_key",
		},
		"SOME_LOCAL_FILE": {
			Tags: []secretsyml.YamlTag{secretsyml.File},
			Path: "my-local-file",
		},
	}

	validSecretsFromWantedEnv = secretsyml.SecretsMap{
		"COMMON_VAR": {
			Tags: []secretsyml.YamlTag{secretsyml.Literal},
			Path: "common_secret",
		},
		"SOME_VAR": {
			Tags: []secretsyml.YamlTag{secretsyml.Var},
			Path: "prod/user/robot/api_key",
		},
	}

	validSecretsValues = map[string]string{
		"prod/user/robot/api_key":     "Secret-API-Key",
		"prod/user/robot/private_key": "Secret-Private-Key",
	}

	validSecretsExports = []exportCmd{
		{"SOME_VAR: " + base64.StdEncoding.EncodeToString([]byte("Secret-API-Key")), ""},
		{"SOME_LOCAL_FILE: " + base64.StdEncoding.EncodeToString([]byte("my-local-file")), ""},
		{"SOME_VAR_FILE: %s", "Secret-Private-Key"},
	}

	secretsMapWithUnknownSecret = secretsyml.SecretsMap{
		"UNKNOWN_VAR": {
			Tags: []secretsyml.YamlTag{secretsyml.Var},
			Path: "completely/nonexistent/path",
		},
	}
)

// mockProvider implements the Provider interface. It is initialized with a
// map that allows secrets values to be retrieved based upon secrets paths.
type mockProvider struct {
	valuesMap map[string]string
}

// newMockProvider instantiates a mockProvider.
func newMockProvider() (Provider, error) {
	return mockProvider{valuesMap: validSecretsValues}, nil
}

// RetrieveSecret returns a secret value based upon a secret path.
func (prov mockProvider) RetrieveSecret(path string) ([]byte, error) {
	val, ok := prov.valuesMap[path]
	if !ok {
		return []byte{}, fmt.Errorf("secret %s not found", path)
	}
	return []byte(val), nil
}

// newMockProviderError is a newProvider function that simply returns an
// error. It is intended to simulate an error in creating or loading a
// secrets provider.
func newMockProviderError() (Provider, error) {
	return nil, fmt.Errorf("could not create a provider")
}

// createSecretsYamlFile creates a temp file containing secrets YAML data.
func createSecretsYamlFile(yamlData string) string {
	f, _ := ioutil.TempFile("", "secrets_*.yaml")
	defer f.Close()

	f.Write([]byte(yamlData))
	return f.Name()
}

// TestParseSecretsYamlFile test the parseSecretsYamlFile function. It
// creates a temp file with specified secrets YAML data, and then calls
// parseSecretsYamlFile to parse the data and confirms expected results.
func TestParseSecretsYamlFile(t *testing.T) {
	testCases := []struct {
		description string
		env         string
		yamlData    string
		expectedMap secretsyml.SecretsMap
		expectError bool
	}{
		{
			description: "Processes a valid secrets YAML file",
			env:         "",
			yamlData:    validSecretsYaml,
			expectedMap: validSecretsMap,
			expectError: false,
		}, {
			description: "Errors on an invalid secrets YAML file",
			env:         "",
			yamlData:    invalidSecretsYaml,
			expectError: true,
		}, {
			description: "Uses the specified secrets environment",
			env:         "wanted",
			yamlData:    validSecretsYamlWithEnv,
			expectedMap: validSecretsFromWantedEnv,
			expectError: false,
		}, {
			description: "Throws an error when env is invalid",
			env:         "doesnt.exist",
			yamlData:    validSecretsYamlWithEnv,
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			// Create a secrets.yml file
			file := createSecretsYamlFile(tc.yamlData)
			defer os.Remove(file)

			// Parse the secrets.yml file
			secrets, err := parseSecretsYamlFile(file, tc.env)

			// Check for error if expected
			if tc.expectError {
				assert.NotNil(t, err)
				return
			}

			// No error expected
			assert.Nil(t, err)
			assert.Equal(t, tc.expectedMap, secrets)
		})
	}
}

// TestValidateEnvVarNames tests the validateEnvVarNames function.
func TestValidateEnvVarNames(t *testing.T) {
	testCases := []struct {
		description    string
		names          []string
		expectedErrStr string
	}{
		{
			"Valid Names",
			[]string{"ALLCAPS", "lowercase", "CamelCase", "Snake_Name_1",
				"_starting_with_underscore", "A", "z", "_"},
			"",
		}, {
			"Invalid name beginning with a digit",
			[]string{"7_UP_THE_UNCOLA"},
			"7_UP_THE_UNCOLA",
		}, {
			"Invalid name with dash",
			[]string{"With-Dash", "Valid_Name_1", "Valid_Name_2"},
			"With-Dash",
		}, {
			"Invalid name with dots",
			[]string{"Valid_Name_1", "With.Dots", "Valid_Name_2"},
			"With.Dots",
		}, {}, {
			"Invalid name with '$' characters",
			[]string{"Valid_Name_1", "Valid_Name_2", "With$Dollar$Signs"},
			"With$Dollar$Signs",
		}, {
			"Invalid name with '=' characters",
			[]string{"With=Equal=Signs"},
			"With=Equal=Signs",
		}, {
			"Invalid name with '@' characters",
			[]string{"With@At@Signs"},
			"With@At@Signs",
		}, {
			"Invalid name with '&' characters",
			[]string{"With&Ampersand"},
			"With&Ampersand",
		}, {
			"Invalid name with '*' characters",
			[]string{"With*Asterisk"},
			"With*Asterisk",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			// Create a SecretsMap with an entry for each name for this test
			secrets := secretsyml.SecretsMap{}
			for _, name := range tc.names {
				secrets[name] = secretsyml.SecretSpec{}
			}

			// Validate names for this test case
			err := validateEnvVarNames(secrets)

			// Check for error if expected
			if tc.expectedErrStr != "" {
				assert.Contains(t, err.Error(), tc.expectedErrStr)
				return
			}

			// No error expected
			assert.Nil(t, err)
		})
	}
}

// TestRetrieveSecrets tests the testRetrieveSecrets function.
func TestRetrieveSecrets(t *testing.T) {
	testCases := []struct {
		description     string
		secretsMap      secretsyml.SecretsMap
		newProvider     newProvider
		expectedExports []exportCmd
		expectedErrStr  string
	}{
		{
			description:     "Successful secrets retrieval",
			secretsMap:      validSecretsMap,
			newProvider:     newMockProvider,
			expectedExports: validSecretsExports,
		}, {
			description:    "Handles error on creating secrets provider",
			secretsMap:     validSecretsMap,
			newProvider:    newMockProviderError,
			expectedErrStr: "could not create a provider",
		}, {
			description:    "Handles error retrieving unknown secret",
			secretsMap:     secretsMapWithUnknownSecret,
			newProvider:    newMockProvider,
			expectedErrStr: "not found",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			tempFactory := newTempFactory()
			defer tempFactory.cleanup()

			exportStrings, err := retrieveSecrets(
				tc.secretsMap,
				tc.newProvider,
				&tempFactory)

			// Check for error if expected
			if tc.expectedErrStr != "" {
				assert.Contains(t, err.Error(), tc.expectedErrStr)
				return
			}

			// No errors are expected. Check that export strings and
			// temporary files are generated as expected.
			assert.Nil(t, err)

			expNum := len(tc.expectedExports)
			assert.Equal(t, len(exportStrings), expNum,
				"%d export strings are expected to be generated", expNum)

			alreadyChecked := map[string]interface{}{}
			for _, str := range exportStrings {
				assert.NotEqual(t, str, "", "Null export string generated")

				_, checked := alreadyChecked[str]
				assert.False(t, checked, "Duplicate export string: %s", str)

				// Match this string with expected export commands.
				matchCmd, tempFile, err := findMatchingExport(str,
					tc.expectedExports, tempFactory.files)
				assert.Nil(t, err)
				assert.NotNil(t, matchCmd, "Unexpected export cmd: %s", str)

				// If the export command contains a temp file name, verify
				// that the content of that temp file are as expected.
				if matchCmd.tempFileContent != "" {
					decodedFileName, err := base64.StdEncoding.DecodeString(tempFile)
					err = checkFileContent(string(decodedFileName), matchCmd.tempFileContent)
					assert.Nil(t, err)
				}
			}
		})
	}
}

// findMatchingExport compares a generated export command string with entries
// in a list of expected export commands. For each entry in the list:
// - If the expected command does not refer to a temp file, then the
//   generated command string is compared directly with the expected command.
// - If the expected command refers to a temp file, then:
//      (a) A temp filename is extracted from the generated command
//      (b) The extracted temp filename is compared with a list of
//          available temp files.
// It optionally returns:
// - A pointer to the matching export command, if there is a match
// - A temp filename, if the filename in the generated export string matches
//   an available temp file.
// - Any errors that occur
func findMatchingExport(str string, exportCmds []exportCmd,
	tempFiles []string) (*exportCmd, string, error) {
	for _, cmd := range exportCmds {
		if cmd.tempFileContent == "" {
			if str == cmd.format {
				return &cmd, "", nil
			}
		} else {
			var file string
			n, err := fmt.Sscanf(str, cmd.format, &file)
			if err != nil {
				return nil, "", err
			}
			if n != 1 {
				continue
			}
			// Format of scanned file name should be '<filename>';
			// Trim trailing semicolon and single quotes.
			for _, tempFile := range tempFiles {
				if file == base64.StdEncoding.EncodeToString([]byte(tempFile)) {
					return &cmd, file, nil
				}
			}
			// Short cut: The string matches one of the expected export
			// formats, but the filename that it includes does not match
			// any existing temp files, so no need to continue.
			break
		}
	}
	return nil, "", nil
}

// checkFileContent reads a file and compares its content with expected
// content. It returns error for either a file read error or for
// mismatching content.
func checkFileContent(filename, expectedContent string) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	if string(data) != expectedContent {
		err = fmt.Errorf("content of %s does not match. Expected: '%s', Read: %s",
			filename, expectedContent, string(data))
		return err
	}
	return nil
}
