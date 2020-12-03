Feature: Get Handler
  In order to ensure that the Get handler is functioning correctly
  As a developer
  I need to be able to test the code from the outside in

  Scenario: Get returns 7 branches
    Given the application is running
    When I call the Get endpoint
    Then there should be 7 branches returned
  
  Scenario: Insert adds a new branch
    Given the application is running
    When I call the Insert endpoint with the folowing JSON
    ```
    {
      "name":"test", 
      "street":"test_street",
      "city":"test_city", 
      "zip": "12345"
    }
    ```
    And I call the Get endpoint
    Then there should be 8 branches returned
