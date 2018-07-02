Feature: RAML test

  Scenario: local RAML
    Given ramltestee is running
    And  contract is run with following parameres
    """
    --verbose --no-color test /opt/bbtest/raml/api.raml
    """
    Then logs contains following
    """
    PASS POST http://ramltestee:8080/v1/person
    PASS GET http://ramltestee:8080/ping
    PASS GET http://ramltestee:8080/v1/person/
    PASS GET http://ramltestee:8080/v1/person
    PASS DELETE http://ramltestee:8080/v1/person/
    """
    And ramltestee is not running

  Scenario: RAML v0.8
    Given ramltestee is running
    And  contract is run with following parameres
    """
    --verbose --no-color test /opt/spec/raml/v08/api.raml
    """
    And ramltestee is not running

  Scenario: RAML v1.0
    Given ramltestee is running
    And  contract is run with following parameres
    """
    --verbose --no-color test /opt/spec/raml/v10/api.raml
    """
    And ramltestee is not running
