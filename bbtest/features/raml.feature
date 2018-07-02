Feature: RAML test

  Scenario: local RAML
    Given mock is running
    And contract is run with following parameres
    """
    --verbose --no-color test /opt/bbtest/raml/api.raml
    """
    Then logs contains following
    """
    PASS POST http://mock:8080/v1/person
    PASS GET http://mock:8080/ping
    PASS GET http://mock:8080/v1/person/
    PASS GET http://mock:8080/v1/person
    PASS DELETE http://mock:8080/v1/person/
    """
    And mock is not running

  Scenario: RAML v0.8
    Given mock is running
    And contract is run with following parameres
    """
    --verbose --no-color test /opt/spec/raml/v08/api.raml
    """
    And mock is not running

  Scenario: RAML v1.0
    Given mock is running
    And contract is run with following parameres
    """
    --verbose --no-color test /opt/spec/raml/v10/api.raml
    """
    And mock is not running
