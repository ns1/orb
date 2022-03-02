@agents_ui
Feature: Create agents using orb ui

    @smoke
    Scenario: Create agent
        Given the Orb user logs in through the UI
            And that fleet Management is clickable on ORB Menu
            And that Agents is clickable on ORB Menu
        When a new agent is created through the UI with region:br, demo:true, ns1:true orb tag(s)
        Then the agents list and the agents view should display agent's status as New within 10 seconds

    @smoke
    Scenario: Provision agent
        Given the Orb user logs in through the UI
            And that fleet Management is clickable on ORB Menu
            And that Agents is clickable on ORB Menu
        When a new agent is created through the UI with 2 orb tag(s)
            And the agent container is started using the command provided by the UI on port default
        Then the agents list and the agents view should display agent's status as Online within 10 seconds
            And the agent status in Orb should be online
            And the container logs should contain the message "sending capabilities" within 10 seconds

    @smoke
    Scenario: Run two orb agents on the same port
        Given the Orb user logs in through the UI
            And that the user is on the orb Agent page
            And a new agent is created through the UI with 1 orb tag(s)
            And the agent container is started using the command provided by the UI on port default
            And that the user is on the orb Agent page
        When a new agent is created through the UI with 1 orb tag(s)
            And the agent container is started using the command provided by the UI on port default
        Then last container created is exited after 2 seconds
            And the container logs should contain the message "agent startup error" within 2 seconds
            And container on port default is running after 2 seconds

    @smoke
    Scenario: Run two orb agents on the same port
        Given the Orb user logs in through the UI
            And that the user is on the orb Agent page
            And a new agent is created through the UI with 1 orb tag(s)
            And the agent container is started using the command provided by the UI on port default
            And that the user is on the orb Agent page
        When a new agent is created through the UI with 1 orb tag(s)
            And the agent container is started using the command provided by the UI on port 10854
        Then last container created is running after 2 seconds
            And container on port default is running after 2 seconds