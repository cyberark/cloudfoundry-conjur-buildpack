@integration
Feature: Integrations Tests for remote TAS foundation

  These tests verify how the buildpack interacts with other
  final buildpacks. For the given language-specific final
  buildpack, the `conjur-env` should run and populate
  the application environment with the values from
  `secrets.yml`.

  With these tests, we currently do not connect to a Conjur
  instance, but only test the buildpack interactions.

    Background:
      Given I create an org and space
      And I install the buildpack

    Scenario: Python offline buildpack integration
      When I push a "python" app with the "offline" buildpack
      Then the secrets.yml values are available in the app

    Scenario: Ruby offline buildpack integration
      When I push a "ruby" app with the "offline" buildpack
      Then the secrets.yml values are available in the app

    Scenario: Java offline buildpack integration
      When I push a "java" app with the "offline" buildpack
      Then the secrets.yml values are available in the app

#    # The online buildpack tests are only valid if the latest commits
#    # are push to the Github remote branch.
    Scenario: Python online buildpack integration
      When I push a "python" app with the "online" buildpack
      Then the secrets.yml values are available in the app

    Scenario: Ruby online buildpack integration
      When I push a "ruby" app with the "online" buildpack
      Then the secrets.yml values are available in the app

    Scenario: Java online buildpack integration
      When I push a "java" app with the "online" buildpack
      Then the secrets.yml values are available in the app
