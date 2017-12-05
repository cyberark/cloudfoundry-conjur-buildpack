Feature: Detect script
  Detect script should return a non-zero exit status

Scenario:
  Given the "detect" script is run
  Then the result should have a non-zero exit status
