Feature: Compile script
  Compile script should install summon, summon-conjur and .profile.d script

#  @BUILD_DIR
#  Scenario:
#    When the "compile" script is run
#    Then the result should have a 0 exit status

  @BUILD_DIR
  Scenario:
    When the "compile" script is run
    Then summon is installed
    And summon-conjur is installed
    And the retrieve secrets .profile.d script is installed
