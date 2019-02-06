Feature: Supply script
  Supply script installs conjur-env and .profile.d script

  @BUILD_DIR
  Scenario: Successfully installs conjur-env and .profile.d scripts
    When the 'supply' script is run
    Then the result should have a 0 exit status
    And conjur-env is installed
    And the retrieve secrets .profile.d script is installed
