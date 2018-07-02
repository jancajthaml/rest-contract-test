Feature: Binary test

  Scenario: run binary
    Given ramltestee is running
    Then  contract is run with following parameres
    """
    --verbose --no-color test /opt/bbtest/raml/api.raml
    """
    And   ramltestee is not running
