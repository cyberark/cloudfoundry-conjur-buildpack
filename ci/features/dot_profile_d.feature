Feature: .profile.d script
  Script to populate the environment variables with secrets from Conjur

  @BUILD_DIR
  Scenario: Populates environment with secrets from Conjur
    Given the 'compile' script is run
    And VCAP_SERVICES contains cybark-conjur credentials
    And the build directory has a secrets.yml file
    When the .profile.d script is sourced
    Then the environment contains the secret values as per secrets.yml
