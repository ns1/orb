@integration
Feature: Integration tests

Scenario: Apply two policies to an agent
    Given the Orb user has a registered account
        And the Orb user logs in
        And that an agent already exists and is online
        And referred agent is subscribed to a group
        And that a sink already exists
    When 2 policies are applied to the agent
    Then this agent's heartbeat shows that 2 policies are successfully applied
        And the container logs contain the message "policy applied successfully" referred to each policy within 10 seconds
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 10 seconds
        And datasets related to all existing policies have validity valid


Scenario: Remove policy from agent
    Given the Orb user has a registered account
        And the Orb user logs in
        And that an agent already exists and is online
        And referred agent is subscribed to a group
        And that a sink already exists
        And 2 policies are applied to the agent
        And this agent's heartbeat shows that 2 policies are successfully applied
    When one of applied policies is removed
    Then referred policy must not be listed on the orb policies list
        And datasets related to removed policy has validity invalid
        And datasets related to all existing policies have validity valid
        And this agent's heartbeat shows that 1 policies are successfully applied
        And container logs should inform that removed policy was stopped and removed within 10 seconds
        And the container logs that were output after the policy have been removed contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And the container logs that were output after the policy have been removed does not contain the message "scraped metrics for policy" referred to deleted policy anymore


Scenario: Remove dataset from agent with just one dataset linked
    Given the Orb user has a registered account
        And the Orb user logs in
        And that an agent already exists and is online
        And referred agent is subscribed to a group
        And that a sink already exists
        And 1 policies are applied to the agent
        And this agent's heartbeat shows that 1 policies are successfully applied
    When a dataset linked to this agent is removed
    Then referred dataset must not be listed on the orb datasets list
        And this agent's heartbeat shows that 0 policies are successfully applied
        And container logs should inform that removed policy was stopped and removed within 10 seconds
        And the container logs that were output after removing dataset contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And the container logs that were output after removing dataset does not contain the message "scraped metrics for policy" referred to deleted policy anymore


Scenario: Remove dataset from agent with more than one dataset linked
    Given the Orb user has a registered account
        And the Orb user logs in
        And that an agent already exists and is online
        And referred agent is subscribed to a group
        And that a sink already exists
        And 3 policies are applied to the agent
        And this agent's heartbeat shows that 3 policies are successfully applied
    When a dataset linked to this agent is removed
    Then referred dataset must not be listed on the orb datasets list
        And this agent's heartbeat shows that 2 policies are successfully applied
        And container logs should inform that removed policy was stopped and removed within 10 seconds
        And the container logs that were output after removing dataset contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And the container logs that were output after removing dataset does not contain the message "scraped metrics for policy" referred to deleted policy anymore