Feature: Detect script
  Detect script returns a non-zero exit status

Scenario: Returns a non-zero exit status
  Given the 'detect' script is run
  Then the result should have a non-zero exit status
