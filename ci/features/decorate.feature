Feature: Decorate script
  Decorate script should return a 0 when secrets.yml exist and the "cyberark-conjur" key is present in VCAP_SERVICES, otherwise non-zero exit status

  @BUILD_DIR
  Scenario: When conditions are met returns a zero exit status
    Given VCAP_SERVICES has a cyberark-conjur key
    And the build directory has a secrets.yml file
    When the "decorate" script is run
    Then the result should have a 0 exit status

  @BUILD_DIR
  Scenario: When conditions are not met returns a non-zero exit status
    Given VCAP_SERVICES does not have a cyberark-conjur key
    When the "decorate" script is run
    Then the result should have a 1 exit status
