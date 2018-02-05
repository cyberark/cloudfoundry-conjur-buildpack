Feature: profile d scripts without the buildpack
  Scripts to simultaneously install summon and summon-conjur, and populate the environment variables with secrets from Conjur

  @BUILD_DIR
  Scenario: Populates environment with secrets from Conjur
    Given the compile script is run against the app's root folder
    And conjur-env is installed
    And a root policy:
    """
    - !variable conjur_single_line_secret_id
    - !variable conjur_multi_line_secret_id
    """
    And the 'conjur_single_line_secret_id' variable has a secret value
    """
    single line
    """
    And the 'conjur_multi_line_secret_id' variable has a secret value
    """
    first line
    second line
    """
    And VCAP_SERVICES contains cyberark-conjur credentials
    And the build directory has this secrets.yml file
    """
    CONJUR_SINGLE_LINE_SECRET: !var conjur_single_line_secret_id
    CONJUR_MULTI_LINE_SECRET: !var conjur_multi_line_secret_id
    LITERAL_SECRET: some literal secret
    """
    When the .profile.d scripts are sourced
    And the 'env' command is run
    Then the environment contains
    """
    LITERAL_SECRET=some literal secret

    """
    And the environment contains
    """
    CONJUR_SINGLE_LINE_SECRET=single line

    """
    And the environment contains
    """
    CONJUR_MULTI_LINE_SECRET=first line
    second line

    """
