Feature: RAML test

  Scenario: local RAML test
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
    And   ramltestee is not running
