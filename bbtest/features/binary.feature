Feature: Binary test

  Scenario: local RAML test
    Given ramltestee is running
    Then  contract is run with following parameres
    """
    --verbose test /opt/bbtest/raml/api.raml
    """
    And   ramltestee is not running
