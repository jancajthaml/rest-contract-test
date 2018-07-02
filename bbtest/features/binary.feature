Feature: Binary test

  Scenario: run binary
    Given mock is running
    Then  contract is run with following parameres
    """
    --verbose --no-color test /opt/bbtest/raml/api.raml
    """
    And mock is not running
