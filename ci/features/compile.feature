Feature: Compile script
  Compile script installs summon, summon-conjur and .profile.d script

  @BUILD_DIR
  Scenario: Successfully installs summon, summon-conjur and .profile.d scripts
    When the 'compile' script is run
    Then the result should have a 0 exit status
    And summon is installed
    And summon-conjur is installed
    And the retrieve secrets .profile.d script is installed
