@integration @AUTORETRY
Feature: Integration tests


@private
Scenario: General smoke test to validate private agent image
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink with default configuration type already exists
    When a new agent is created with 1 orb tag(s)
        And the agent container is started on an available port
        And the agent status is online
        And referred agent is subscribed to 1 group
        And 2 simple policies are applied to the group
    Then backends route must be enabled
        And handlers route must be enabled
        And taps route must be enabled
        And inputs route must be enabled
        And pktvisor state is running
        And this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 2 policies are applied and all has status running
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 120 seconds
        And 2 dataset(s) have validity valid and 0 have validity invalid in 30 seconds


@smoke
  Scenario: Full happy path using new orb account
    Given that there is an unregistered valid email with 12345678 password
    When the Orb user request this account registration with email used for testing process as company and Test process as fullname
    Then the status code must be 201
      And account is registered with email, with password, email used for testing process company and Test process full name
    Given that a sink with default configuration type already exists
    When a new agent is created with 1 orb tag(s)
        And the agent container is started on an available port
        And the agent status is online
        And referred agent is subscribed to 1 group
        And 1 simple policies are applied to the group
    Then backends route must be enabled
        And handlers route must be enabled
        And taps route must be enabled
        And inputs route must be enabled
        And pktvisor state is running
        And this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 1 policies are applied and all has status running
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 120 seconds
        And 2 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
    When this agent is removed
        And one of applied policies is removed
        And 1 group(s) to which the agent is linked is removed
        And remove 1 of the linked sinks from orb
        And the dataset is removed



@smoke
Scenario: Test agents backend routes
    Given the Orb user has a registered account
        And the Orb user logs in
    When a new agent is created with 1 orb tag(s)
        And the agent container is started on an available port
    Then backends route must be enabled
        And handlers route must be enabled
        And taps route must be enabled
        And inputs route must be enabled


@smoke @sanity
Scenario: Create dataset with name conflict
    Given the Orb user has a registered account
        And the Orb user logs in
        And 1 Agent Group(s) is created with 2 orb tag(s)
        And that a sink with default configuration type already exists
        And 1 simple policies are applied to the group
    When a new dataset is requested to be created with the same name as an existent one
    Then the error message on response is failed to create dataset


@MUTE
Scenario: Edit dataset using an already existent name (conflict)
    Given the Orb user has a registered account
        And the Orb user logs in
        And 1 Agent Group(s) is created with 2 orb tag(s)
        And that a sink with default configuration type already exists
        And 1 simple policies are applied to the group
        And 1 new dataset is created using the policy, last group and 1 sink
    When editing a dataset using a name already in use
    Then the error message on response is entity already exists


@smoke
Scenario: Apply multiple advanced policies to an agent
    Given the Orb user has a registered account
        And the Orb user logs in
        And that an agent with 1 orb tag(s) already exists and is online
        And pktvisor state is running
        And referred agent is subscribed to 1 group
        And this agent's heartbeat shows that 1 groups are matching the agent
        And that a sink with default configuration type already exists
    When 14 mixed policies are applied to the group
    Then this agent's heartbeat shows that 14 policies are applied and all has status running
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And 14 dataset(s) have validity valid and 0 have validity invalid in 30 seconds


@smoke
Scenario: Apply two simple policies to an agent
    Given the Orb user has a registered account
        And the Orb user logs in
        And that an agent with 1 orb tag(s) already exists and is online
        And pktvisor state is running
        And referred agent is subscribed to 1 group
        And this agent's heartbeat shows that 1 groups are matching the agent
        And that a sink with default configuration type already exists
    When 2 simple policies are applied to the group
    Then this agent's heartbeat shows that 2 policies are applied and all has status running
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 120 seconds
        And 2 dataset(s) have validity valid and 0 have validity invalid in 30 seconds


@sanity
Scenario: apply one policy using multiple datasets to the same group
    Given the Orb user has a registered account
        And the Orb user logs in
        And that an agent with 2 orb tag(s) already exists and is online
        And pktvisor state is running
        And referred agent is subscribed to 1 group
        And this agent's heartbeat shows that 1 groups are matching the agent
        And that a sink with default configuration type already exists
    When 2 simple policies are applied to the group by 3 datasets each
    Then this agent's heartbeat shows that 2 policies are applied and all has status running
        And 3 datasets are linked with each policy on agent's heartbeat within 180 seconds
        And the container logs contain the message "policy applied successfully" referred to each policy within 180 seconds
        And referred sink must have active state on response within 180 seconds
        And 6 dataset(s) have validity valid and 0 have validity invalid in 30 seconds


@smoke
Scenario: Remove group to which agent is linked
    Given the Orb user has a registered account
        And the Orb user logs in
        And that an agent with 1 orb tag(s) already exists and is online
        And pktvisor state is running
        And referred agent is subscribed to 1 group
        And this agent's heartbeat shows that 1 groups are matching the agent
        And that a sink with default configuration type already exists
        And 2 simple policies are applied to the group
        And this agent's heartbeat shows that 2 policies are applied and all has status running
    When 1 group(s) to which the agent is linked is removed
    Then the container logs should contain the message "completed RPC unsubscription to group" within 30 seconds
        And the container logs contain the message "policy no longer used by any group, removing" referred to each policy within 30 seconds
        And this agent's heartbeat shows that 0 policies are applied to the agent
        And this agent's heartbeat shows that 0 groups are matching the agent
        And no dataset should be linked to the removed group anymore
        And 0 dataset(s) have validity valid and 2 have validity invalid in 30 seconds


@smoke
Scenario: Remove policy from agent
    Given the Orb user has a registered account
        And the Orb user logs in
        And that an agent with 3 orb tag(s) already exists and is online
        And pktvisor state is running
        And referred agent is subscribed to 1 group
        And this agent's heartbeat shows that 1 groups are matching the agent
        And that a sink with default configuration type already exists
        And 2 simple policies are applied to the group
        And this agent's heartbeat shows that 2 policies are applied and all has status running
    When one of applied policies is removed
    Then referred policy must not be listed on the orb policies list
        And no dataset should be linked to the removed policy anymore
        And 1 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And this agent's heartbeat shows that 1 policies are applied and all has status running
        And container logs should inform that removed policy was stopped and removed within 30 seconds
        And the container logs that were output after the policy have been removed contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And the container logs that were output after the policy have been removed does not contain the message "scraped metrics for policy" referred to deleted policy anymore


@smoke
Scenario: Remove dataset from agent with just one dataset linked
    Given the Orb user has a registered account
        And the Orb user logs in
        And that an agent with 3 orb tag(s) already exists and is online
        And pktvisor state is running
        And referred agent is subscribed to 1 group
        And this agent's heartbeat shows that 1 groups are matching the agent
        And that a sink with default configuration type already exists
        And 1 simple policies are applied to the group
        And this agent's heartbeat shows that 1 policies are applied and all has status running
    When a dataset linked to this agent is removed
    Then referred dataset must not be listed on the orb datasets list
        And this agent's heartbeat shows that 0 policies are applied and all has status running
        And container logs should inform that removed policy was stopped and removed within 30 seconds
        And the container logs that were output after removing dataset contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And the container logs that were output after removing dataset does not contain the message "scraped metrics for policy" referred to deleted policy anymore


@smoke
Scenario: Remove dataset from agent with more than one dataset linked
    Given the Orb user has a registered account
        And the Orb user logs in
        And that an agent with 4 orb tag(s) already exists and is online
        And pktvisor state is running
        And referred agent is subscribed to 1 group
        And this agent's heartbeat shows that 1 groups are matching the agent
        And that a sink with default configuration type already exists
        And 3 simple policies are applied to the group
        And this agent's heartbeat shows that 3 policies are applied and all has status running
    When a dataset linked to this agent is removed
    Then referred dataset must not be listed on the orb datasets list
        And this agent's heartbeat shows that 2 policies are applied and all has status running
        And container logs should inform that removed policy was stopped and removed within 30 seconds
        And the container logs that were output after removing dataset contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And the container logs that were output after removing dataset does not contain the message "scraped metrics for policy" referred to deleted policy anymore


@smoke
Scenario: Provision agent with tags matching an existent group
    Given the Orb user has a registered account
        And the Orb user logs in
        And 1 Agent Group(s) is created with 2 orb tag(s)
    When a new agent is created with orb tags matching 1 existing group
        And the agent container is started on an available port
        And the agent status is online
    Then the agent status in Orb should be online within 30 seconds
        And the container logs should contain the message "completed RPC subscription to group" within 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent


@smoke
Scenario: Provision agent with tag matching existing group linked to a valid dataset
    Given the Orb user has a registered account
        And the Orb user logs in
        And 1 Agent Group(s) is created with 3 orb tag(s)
        And that a sink with default configuration type already exists
        And 2 simple policies are applied to the group
    When a new agent is created with orb tags matching 1 existing group
        And the agent container is started on an available port
        And the agent status is online
    Then this agent's heartbeat shows that 2 policies are applied and all has status running
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 120 seconds
        And 2 dataset(s) have validity valid and 0 have validity invalid in 30 seconds


@smoke
Scenario: Provision agent with tag matching existing group with multiple policies linked to a valid dataset
    Given the Orb user has a registered account
        And the Orb user logs in
        And 1 Agent Group(s) is created with 2 orb tag(s)
        And that a sink with default configuration type already exists
        And 14 mixed policies are applied to the group
    When a new agent is created with orb tags matching 1 existing group
        And the agent container is started on an available port
        And the agent status is online
    Then this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 14 policies are applied and all has status running
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 120 seconds
        And 14 dataset(s) have validity valid and 0 have validity invalid in 30 seconds


@smoke
Scenario: Provision agent with tag matching existing edited group with multiple policies linked to a valid dataset
    Given the Orb user has a registered account
        And the Orb user logs in
        And 1 Agent Group(s) is created with 3 orb tag(s)
        And the name, tags, description of last Agent Group is edited using: name=edited_before_policy/ tags=2 orb tag(s)/ description=edited
        And that a sink with default configuration type already exists
        And 14 mixed policies are applied to the group
    When a new agent is created with orb tags matching 1 existing group
        And the agent container is started on an available port
        And the agent status is online
    Then this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 14 policies are applied and all has status running
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 120 seconds
        And 14 dataset(s) have validity valid and 0 have validity invalid in 30 seconds


@smoke
Scenario: Provision agent with tag matching existing group with multiple policies
    Given the Orb user has a registered account
        And the Orb user logs in
        And 1 Agent Group(s) is created with 3 orb tag(s)
        And that a sink with default configuration type already exists
        And 20 mixed policies are applied to the group
        And the name, tags, description of last Agent Group is edited using: name=edited_after_policy/ tags=2 orb tag(s)/ description=edited
    When a new agent is created with orb tags matching 1 existing group
        And the agent container is started on an available port
        And the agent status is online
    Then this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 20 policies are applied and all has status running
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 120 seconds
        And 20 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
@MUTE
#@sanity @sink_status_idle
Scenario: Sink idle after 5 minutes without metrics flow
    Given the Orb user has a registered account
        And the Orb user logs in
        And that an agent with 1 orb tag(s) already exists and is online
        And pktvisor state is running
        And referred agent is subscribed to 1 group
        And this agent's heartbeat shows that 1 groups are matching the agent
        And that a sink with default configuration type already exists
        And 2 simple policies are applied to the group
        And this agent's heartbeat shows that 2 policies are applied and all has status running
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 120 seconds
        And 2 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
    When stop the orb-agent container
    Then referred sink must have idle state on response after 660 seconds

@MUTE
#@sanity @sink_status_error
Scenario: Sink with invalid endpoint
    Given the Orb user has a registered account
        And the Orb user logs in
        And that an agent with 1 orb tag(s) already exists and is online
        And pktvisor state is running
        And referred agent is subscribed to 1 group
        And this agent's heartbeat shows that 1 groups are matching the agent
        And that a sink with invalid endpoint already exists
        And 3 simple policies are applied to the group
        And that a policy using: handler=dns, description='policy_dns', host_specification=10.0.1.0/24,10.0.2.1/32,2001:db8::/64, bpf_filter_expression=udp port 53, pcap_source=libpcap, only_qname_suffix=[.foo.com/ .example.com], only_rcode=2 already exists
    When 1 new dataset is created using the policy, last group and 1 sink
    Then this agent's heartbeat shows that 4 policies are applied and all has status running
        And the container logs should contain the message "managing agent policy from core" within 30 seconds
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And the container logs should contain the message "scraped metrics for policy" within 180 seconds
        And referred sink must have error state on response within 120 seconds
        And 4 dataset(s) have validity valid and 0 have validity invalid in 30 seconds


@smoke
Scenario: Unapplying policies that failed by editing agent orb tags to unsubscribe from group
    Given the Orb user has a registered account
        And the Orb user logs in
        And that an agent with 1 orb tag(s) already exists and is online
        And pktvisor state is running
        And referred agent is subscribed to 1 group
        And this agent's heartbeat shows that 1 groups are matching the agent
        And that a sink with default configuration type already exists
        And 3 simple policies are applied to the group
        And that a policy using: handler=dns, description='policy_dns', bpf_filter_expression=ufp pot 53 already exists
        And 1 new dataset is created using the policy, last group and 1 sink
        And this agent's heartbeat shows that 4 policies are applied and 3 has status running
        And this agent's heartbeat shows that 4 policies are applied and 1 has status failed_to_apply
    When edit the orb tags on agent and use 2 orb tag(s)
    Then the container logs should contain the message "completed RPC unsubscription to group" within 30 seconds
        And this agent's heartbeat shows that 0 policies are applied to the agent



@smoke
Scenario: Unapplying policies that failed by editing group tags to unsubscribe agent from group
    Given the Orb user has a registered account
        And the Orb user logs in
        And that an agent with 1 orb tag(s) already exists and is online
        And pktvisor state is running
        And referred agent is subscribed to 1 group
        And this agent's heartbeat shows that 1 groups are matching the agent
        And that a sink with default configuration type already exists
        And 3 simple policies are applied to the group
        And that a policy using: handler=dns, description='policy_dns', bpf_filter_expression=ufp pot 53 already exists
        And 1 new dataset is created using the policy, last group and 1 sink
        And this agent's heartbeat shows that 4 policies are applied and 3 has status running
        And this agent's heartbeat shows that 4 policies are applied and 1 has status failed_to_apply
    When the tags of last Agent Group is edited using: tags=2 orb tag(s)
    Then the container logs should contain the message "completed RPC unsubscription to group" within 30 seconds
        And this agent's heartbeat shows that 0 policies are applied to the agent



@smoke
Scenario: Unapplying policies that failed by removing group
    Given the Orb user has a registered account
        And the Orb user logs in
        And that an agent with 1 orb tag(s) already exists and is online
        And pktvisor state is running
        And referred agent is subscribed to 1 group
        And this agent's heartbeat shows that 1 groups are matching the agent
        And that a sink with default configuration type already exists
        And 3 simple policies are applied to the group
        And that a policy using: handler=dns, description='policy_dns', bpf_filter_expression=ufp pot 53 already exists
        And 1 new dataset is created using the policy, last group and 1 sink
        And this agent's heartbeat shows that 4 policies are applied and 3 has status running
        And this agent's heartbeat shows that 4 policies are applied and 1 has status failed_to_apply
    When 1 group(s) to which the agent is linked is removed
    Then the container logs should contain the message "completed RPC unsubscription to group" within 30 seconds
        And this agent's heartbeat shows that 0 policies are applied to the agent
        And no dataset should be linked to the removed group anymore
        And 0 dataset(s) have validity valid and 4 have validity invalid in 30 seconds

@MUTE
#@smoke @sink_status_error
Scenario: Sink with invalid username
    Given the Orb user has a registered account
        And the Orb user logs in
        And that an agent with 1 orb tag(s) already exists and is online
        And pktvisor state is running
        And referred agent is subscribed to 1 group
        And this agent's heartbeat shows that 1 groups are matching the agent
        And that a sink with invalid username already exists
        And 3 simple policies are applied to the group
        And that a policy using: handler=dns, description='policy_dns', host_specification=10.0.1.0/24,10.0.2.1/32,2001:db8::/64, bpf_filter_expression=udp port 53, pcap_source=libpcap, only_qname_suffix=[.foo.com/ .example.com], only_rcode=3 already exists
    When 1 new dataset is created using the policy, last group and 1 sink
    Then the container logs should contain the message "managing agent policy from core" within 30 seconds
        And this agent's heartbeat shows that 4 policies are applied and all has status running
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And the container logs should contain the message "scraped metrics for policy" within 180 seconds
        And referred sink must have error state on response within 120 seconds
        And 4 dataset(s) have validity valid and 0 have validity invalid in 30 seconds

@MUTE
#@sanity @sink_status_error
Scenario: Sink with invalid password
    Given the Orb user has a registered account
        And the Orb user logs in
        And that an agent with 1 orb tag(s) already exists and is online
        And pktvisor state is running
        And referred agent is subscribed to 1 group
        And this agent's heartbeat shows that 1 groups are matching the agent
        And that a sink with invalid password already exists
        And 3 simple policies are applied to the group
        And that a policy using: handler=dns, description='policy_dns', host_specification=10.0.1.0/24,10.0.2.1/32,2001:db8::/64, bpf_filter_expression=udp port 53, pcap_source=libpcap, only_qname_suffix=[.foo.com/ .example.com], only_rcode=5 already exists
    When 1 new dataset is created using the policy, last group and 1 sink
    Then this agent's heartbeat shows that 4 policies are applied and all has status running
        And the container logs should contain the message "managing agent policy from core" within 30 seconds
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And the container logs should contain the message "scraped metrics for policy" within 180 seconds
        And referred sink must have error state on response within 120 seconds
        And 4 dataset(s) have validity valid and 0 have validity invalid in 30 seconds


@smoke
Scenario: Agent subscription to multiple groups created after provisioning agent
    Given the Orb user has a registered account
        And the Orb user logs in
        And a new agent is created with 3 orb tag(s)
        And the agent container is started on an available port
        And the agent status is online
    When 1 Agent Group(s) is created with 2 tags contained in the agent
        And 1 Agent Group(s) is created with 1 tags contained in the agent
    Then the container logs contain the message "completed RPC subscription to group" referred to each matching group within 30 seconds
        And this agent's heartbeat shows that 2 groups are matching the agent


@smoke
Scenario: Agent subscription to multiple groups created before provisioning agent
    Given the Orb user has a registered account
        And the Orb user logs in
        And 1 Agent Group(s) is created with 2 orb tag(s)
        And 1 Agent Group(s) is created with 1 orb tag(s)
        And 1 Agent Group(s) is created with 1 orb tag(s)
    When a new agent is created with orb tags matching all existing group
        And the agent container is started on an available port
        And the agent status is online
    Then the container logs contain the message "completed RPC subscription to group" referred to each matching group within 30 seconds
        And this agent's heartbeat shows that 3 groups are matching the agent


@smoke
Scenario: Agent subscription to group after editing orb agent's tags (agent provisioned before editing and group created after)
    Given the Orb user has a registered account
        And the Orb user logs in
        And a new agent is created with 1 orb tag(s)
        And the agent container is started on an available port
        And the agent status is online
    When edit the orb tags on agent and use 2 orb tag(s)
        And 1 Agent Group(s) is created with all tags contained in the agent
    Then the container logs contain the message "completed RPC subscription to group" referred to each matching group within 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent


@smoke
Scenario: Agent subscription to group after editing orb agent's tags (editing tags after agent provision)
    Given the Orb user has a registered account
        And the Orb user logs in
        And 1 Agent Group(s) is created with 2 orb tag(s)
        And 1 Agent Group(s) is created with 1 orb tag(s)
        And a new agent is created with orb tags matching 1 existing group
        And the agent container is started on an available port
        And the agent status is online
    When edit the orb tags on agent and use orb tags matching 2 existing group
    Then the container logs contain the message "completed RPC subscription to group" referred to each matching group within 30 seconds
        And this agent's heartbeat shows that 2 groups are matching the agent


@smoke
Scenario: Agent subscription to group after editing orb agent's tags (editing tags before agent provision)
    Given the Orb user has a registered account
        And the Orb user logs in
        And a new agent is created with 2 orb tag(s)
        And edit the orb tags on agent and use 1 orb tag(s)
        And the agent container is started on an available port
        And the agent status is online
    When 1 Agent Group(s) is created with 1 tags contained in the agent
    Then the container logs contain the message "completed RPC subscription to group" referred to each matching group within 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent


@smoke
Scenario: Agent subscription to multiple group with policies after editing orb agent's tags (editing tags after agent provision)
    Given the Orb user has a registered account
        And the Orb user logs in
        And that an agent with 1 orb tag(s) already exists and is online
        And pktvisor state is running
        And referred agent is subscribed to 1 group
        And this agent's heartbeat shows that 1 groups are matching the agent
        And that a sink with default configuration type already exists
        And 2 simple policies are applied to the group
        And this agent's heartbeat shows that 2 policies are applied and all has status running
        And 1 Agent Group(s) is created with 1 orb tag(s)
        And 1 simple policies are applied to the group
    When edit the orb tags on agent and use orb tags matching all existing group
    Then the container logs contain the message "completed RPC subscription to group" referred to each matching group within 30 seconds
        And this agent's heartbeat shows that 3 policies are applied and all has status running
        And this agent's heartbeat shows that 2 groups are matching the agent
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds


@smoke
Scenario: Agent subscription to group with policies after editing orb agent's tags (editing tags after agent provision)
    Given the Orb user has a registered account
        And the Orb user logs in
        And that an agent with 1 orb tag(s) already exists and is online
        And pktvisor state is running
        And referred agent is subscribed to 1 group
        And this agent's heartbeat shows that 1 groups are matching the agent
        And that a sink with default configuration type already exists
        And 2 simple policies are applied to the group
        And this agent's heartbeat shows that 2 policies are applied and all has status running
        And 1 Agent Group(s) is created with 1 orb tag(s)
        And 1 simple policies are applied to the group
    When edit the orb tags on agent and use orb tags matching last existing group
    Then the container logs contain the message "completed RPC subscription to group" referred to each matching group within 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 1 policies are applied and all has status running
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds


#@smoke
@MUTE
Scenario: Remove one of the groups that applies the same policy on the agent
    Given the Orb user has a registered account
        And the Orb user logs in
        And that an agent with 1 orb tag(s) already exists and is online
        And pktvisor state is running
        And referred agent is subscribed to 2 groups
        And this agent's heartbeat shows that 2 groups are matching the agent
        And that a sink with default configuration type already exists
        And a new policy is created using: handler=dns, description='policy_dns_2_groups'
        And 2 new dataset is created using the policy, an existing group and 1 sink
        And this agent's heartbeat shows that 1 policies are applied and 1 has status running
        And 2 datasets are linked with each policy on agent's heartbeat within 30 seconds
    When 1 group(s) to which the agent is linked is removed
    Then this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 1 policies are applied and 1 has status running
        And 1 datasets are linked with each policy on agent's heartbeat within 60 seconds
        And no dataset should be linked to the removed group anymore
        And 1 dataset(s) have validity valid and 1 have validity invalid in 30 seconds


@smoke
Scenario: Remove one of the datasets that applies the same policy on the agent
    Given the Orb user has a registered account
        And the Orb user logs in
        And that an agent with 1 orb tag(s) already exists and is online
        And pktvisor state is running
        And referred agent is subscribed to 2 groups
        And this agent's heartbeat shows that 2 groups are matching the agent
        And that a sink with default configuration type already exists
        And a new policy is created using: handler=dns, description='policy_dns_2_groups'
        And 2 new dataset is created using the policy, an existing group and 1 sink
        And this agent's heartbeat shows that 1 policies are applied and 1 has status running
    When a dataset linked to this agent is removed
    Then this agent's heartbeat shows that 2 groups are matching the agent
        And this agent's heartbeat shows that 1 policies are applied and 1 has status running
        And 1 datasets are linked with each policy on agent's heartbeat within 60 seconds


@smoke
Scenario: Insert tags in agents created without tags and apply policies to group matching new tags
    Given the Orb user has a registered account
        And the Orb user logs in
        And a new agent is created with 0 orb tag(s)
        And the agent container is started on an available port
        And the agent status is online
        And that a sink with default configuration type already exists
    When edit the orb tags on agent and use 2 orb tag(s)
        And 1 Agent Group(s) is created with same tag as the agent and without description
        And 1 simple policies are applied to the group
    Then the container logs contain the message "completed RPC subscription to group" referred to each matching group within 30 seconds
        And this agent's heartbeat shows that 1 policies are applied and all has status running
        And this agent's heartbeat shows that 1 groups are matching the agent
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds


@smoke
Scenario: Edit agent name and apply policies to then
    Given the Orb user has a registered account
        And the Orb user logs in
        And that an agent with 5 orb tag(s) already exists and is online
        And pktvisor state is running
        And 1 Agent Group(s) is created with all tags contained in the agent
        And 1 agent must be matching on response field matching_agents of the last group created
        And that a sink with default configuration type already exists
    When edit the agent name
        And 1 simple policies are applied to the group
    Then this agent's heartbeat shows that 1 policies are applied and all has status running
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds


@smoke
Scenario: Editing tags of an Agent Group with policies (unsubscription - provision agent before editing)
    Given the Orb user has a registered account
        And the Orb user logs in
        And that an agent with 1 orb tag(s) already exists and is online
        And 1 Agent Group(s) is created with same tag as the agent and without description
        And this agent's heartbeat shows that 1 groups are matching the agent
        And that a sink with default configuration type already exists
        And 2 simple policies are applied to the group
    When the name, tags, description of last Agent Group is edited using: name=new_name/ tags=2 orb tag(s)/ description=None
    Then 0 agent must be matching on response field matching_agents of the last group created
        And the container logs should contain the message "completed RPC unsubscription to group" within 30 seconds
        And the agent status in Orb should be online within 30 seconds


@smoke
Scenario: Editing tags of an Agent Group with policies (subscription - provision agent before editing)
    Given the Orb user has a registered account
        And the Orb user logs in
        And that an agent with 2 orb tag(s) already exists and is online
        And pktvisor state is running
        And 1 Agent Group(s) is created with 1 orb tag(s) and without description
        And that a sink with default configuration type already exists
        And 2 simple policies are applied to the group
    When the name, tags, description of last Agent Group is edited using: name=new_name/ tags=matching the agent/ description=None
    Then 1 agent must be matching on response field matching_agents of the last group created
        And the container logs should contain the message "completed RPC subscription to group" within 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 2 policies are applied and all has status running


@smoke
Scenario: Editing tags of an Agent Group with policies (provision agent after editing)
    Given the Orb user has a registered account
        And the Orb user logs in
        And 1 Agent Group(s) is created with 1 orb tag(s) and without description
        And that a sink with default configuration type already exists
        And 2 simple policies are applied to the group
        And a new agent is created with orb tags matching 1 existing group
        And 1 agent must be matching on response field matching_agents of the last group created
    When the name, tags, description of last Agent Group is edited using: name=new_name/ tags=2 orb tag/ description=None
        And the agent container is started on an available port
        And the agent status is online
    Then 0 agent must be matching on response field matching_agents of the last group created
        And the agent status in Orb should be online within 30 seconds


@smoke
Scenario: Editing tags of an Agent Group with policies (subscription - provision agent after editing)
    Given the Orb user has a registered account
        And the Orb user logs in
        And 1 Agent Group(s) is created with 1 orb tag(s) and without description
        And that a sink with default configuration type already exists
        And 2 simple policies are applied to the group
    When the name, tags, description of last Agent Group is edited using: name=new_name/ tags=2 orb tag(s)/ description=None
        And a new agent is created with orb tags matching 1 existing group
        And the agent container is started on an available port
        And the agent status is online
    Then 1 agent must be matching on response field matching_agents of the last group created
        And the container logs should contain the message "completed RPC subscription to group" within 30 seconds
        And the agent status in Orb should be online within 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 2 policies are applied and all has status running


@smoke
Scenario: Editing tags of an Agent and Agent Group with policies (unsubscription - provision agent before editing)
    Given the Orb user has a registered account
        And the Orb user logs in
        And that an agent with 1 orb tag(s) already exists and is online
        And pktvisor state is running
        And 1 Agent Group(s) is created with same tag as the agent and without description
        And that a sink with default configuration type already exists
        And 2 simple policies are applied to the group
    When the name, tags, description of last Agent Group is edited using: name=new_name/ tags=2 orb tag(s)/ description=None
        And edit the orb tags on agent and use 1 orb tag(s)
    Then 0 agent must be matching on response field matching_agents of the last group created
        And the container logs should contain the message "completed RPC unsubscription to group" within 30 seconds
        And this agent's heartbeat shows that 0 groups are matching the agent
        And this agent's heartbeat shows that 0 policies are applied to the agent
        And the agent status in Orb should be online within 30 seconds


@smoke
Scenario: Editing tags of an Agent and Agent Group with policies (subscription - provision agent before editing)
    Given the Orb user has a registered account
        And the Orb user logs in
        And that an agent with 3 orb tag(s) already exists and is online
        And pktvisor state is running
        And 1 Agent Group(s) is created with 1 orb tag(s) and without description
        And that a sink with default configuration type already exists
        And 2 simple policies are applied to the group
    When edit the orb tags on agent and use 2 orb tag(s)
        And the name, tags, description of last Agent Group is edited using: name=new_name/ tags=matching the agent/ description=None
    Then 1 agent must be matching on response field matching_agents of the last group created
        And the container logs should contain the message "completed RPC subscription to group" within 30 seconds
        And the agent status in Orb should be online within 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 2 policies are applied and all has status running


@smoke
Scenario: Editing tags of an Agent and Agent Group with policies (provision agent after editing)
    Given the Orb user has a registered account
        And the Orb user logs in
        And 1 Agent Group(s) is created with 1 orb tag(s) and without description
        And that a sink with default configuration type already exists
        And 2 simple policies are applied to the group
    When the name, tags, description of last Agent Group is edited using: name=new_name/ tags=2 orb tag(s)/ description=None
        And a new agent is created with 1 orb tag(s)
        And edit the orb tags on agent and use 1 orb tag(s)
        And the agent container is started on an available port
        And the agent status is online
    Then 0 agent must be matching on response field matching_agents of the last group created
        And the agent status in Orb should be online within 30 seconds


@smoke
Scenario: Editing tags of an Agent and Agent Group with policies (subscription - provision agent after editing)
    Given the Orb user has a registered account
        And the Orb user logs in
        And 1 Agent Group(s) is created with 1 orb tag(s) and without description
        And that a sink with default configuration type already exists
        And 2 simple policies are applied to the group
    When the name, tags, description of last Agent Group is edited using: name=new_name/ tags=1 orb tag/ description=None
        And a new agent is created with 2 orb tag(s)
        And edit the orb tags on agent and use orb tags matching 1 existing group
        And the agent container is started on an available port
        And the agent status is online
    Then 1 agent must be matching on response field matching_agents of the last group created
        And the container logs should contain the message "completed RPC subscription to group" within 30 seconds
        And the agent status in Orb should be online within 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 2 policies are applied and all has status running



@smoke
Scenario: Edit an advanced policy with handler dns changing the handler to net
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink with default configuration type already exists
        And that an agent with 1 orb tag(s) already exists and is online
        And pktvisor state is running
        And 1 Agent Group(s) is created with all tags contained in the agent
        And a new policy is created using: handler=dns, description='policy_dns', bpf_filter_expression=udp port 53, pcap_source=libpcap, only_qname_suffix=[.foo.com/ .example.com], only_rcode=0
        And 1 new dataset is created using the policy, last group and 1 sink
    When editing a policy using name=my_policy, handler=net, only_qname_suffix=None, only_rcode=None
    Then policy version must be 1
        And policy name must be my_policy
        And policy handler must be net
        And policy only_qname_suffix must be None
        And policy only_rcode must be None
        And the version of policy in agent must be 1
        And this agent's heartbeat shows that 1 policies are applied and all has status running
        And the container logs that were output after editing policies contain the message "policy applied successfully" referred to each applied policy within 10 seconds



@smoke
Scenario: Edit an advanced policy with handler dns changing the handler to dhcp
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink with default configuration type already exists
        And that an agent with 1 orb tag(s) already exists and is online
        And pktvisor state is running
        And 1 Agent Group(s) is created with all tags contained in the agent
        And this agent's heartbeat shows that 1 groups are matching the agent
        And a new policy is created using: handler=dns, host_specification=10.0.1.0/24,10.0.2.1/32,2001:db8::/64, bpf_filter_expression=udp port 53, pcap_source=libpcap, only_qname_suffix=[.foo.com/ .example.com], only_rcode=2
        And 1 new dataset is created using the policy, last group and 1 sink
    When editing a policy using name=second_policy, handler=dhcp, only_qname_suffix=None, only_rcode=None
    Then policy version must be 1
        And policy name must be second_policy
        And policy handler must be dhcp
        And the version of policy in agent must be 1
        And this agent's heartbeat shows that 1 policies are applied and all has status running
        And the container logs that were output after editing policies contain the message "policy applied successfully" referred to each applied policy within 10 seconds


@smoke
Scenario: Edit a simple policy with handler dhcp changing the handler to net
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink with default configuration type already exists
        And that an agent with 1 orb tag(s) already exists and is online
        And pktvisor state is running
        And 1 Agent Group(s) is created with all tags contained in the agent
        And this agent's heartbeat shows that 1 groups are matching the agent
        And a new policy is created using: handler=dhcp
        And 1 new dataset is created using the policy, last group and 1 sink
    When editing a policy using handler=net, description="policy_net"
    Then policy version must be 1
        And policy handler must be net
        And the version of policy in agent must be 1
        And this agent's heartbeat shows that 1 policies are applied and all has status running
        And the container logs that were output after editing policies contain the message "policy applied successfully" referred to each applied policy within 10 seconds


@smoke
Scenario: Edit a simple policy with handler net changing the handler to dns and inserting advanced parameters
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink with default configuration type already exists
        And that an agent with 1 orb tag(s) already exists and is online
        And pktvisor state is running
        And 1 Agent Group(s) is created with all tags contained in the agent
        And a new policy is created using: handler=net
        And 1 new dataset is created using the policy, last group and 1 sink
    When editing a policy using handler=dns, host_specification=10.0.1.0/24,10.0.2.1/32,2001:db8::/64, bpf_filter_expression=udp port 53, pcap_source=libpcap, only_qname_suffix=[.foo.com/ .example.com], only_rcode=2
    Then policy version must be 1
        And policy handler must be dns
        And policy host_specification must be 10.0.1.0/24,10.0.2.1/32,2001:db8::/64
        And policy bpf_filter_expression must be udp port 53
        And policy only_qname_suffix must be ['.foo.com', '.example.com']
        And policy only_rcode must be 2
        And the version of policy in agent must be 1
        And this agent's heartbeat shows that 1 policies are applied and all has status running
        And the container logs that were output after editing policies contain the message "policy applied successfully" referred to each applied policy within 10 seconds


@smoke
Scenario: remove 1 sink from a dataset with 2 sinks
    Given the Orb user has a registered account
        And the Orb user logs in
        And that 2 sinks already exists
        And that an agent with 1 orb tag(s) already exists and is online
        And pktvisor state is running
        And 1 Agent Group(s) is created with all tags contained in the agent
        And this agent's heartbeat shows that 1 groups are matching the agent
        And a new policy is created using: handler=dhcp
        And 1 new dataset is created using the policy, last group and 2 sinks
        And 1 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
    When remove 1 of the linked sinks from orb
    Then 1 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And this agent's heartbeat shows that 1 policies are applied and all has status running
        And no dataset should be linked to the removed sink anymore


@smoke
Scenario: remove 1 sink from a dataset with 1 sinks
    Given the Orb user has a registered account
        And the Orb user logs in
        And that 2 sinks already exists
        And that an agent with 1 orb tag(s) already exists and is online
        And pktvisor state is running
        And 1 Agent Group(s) is created with all tags contained in the agent
        And a new policy is created using: handler=dhcp
        And 1 new dataset is created using the policy, last group and 1 sinks
        And 1 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
    When remove 1 of the linked sinks from orb
    Then 0 dataset(s) have validity valid and 1 have validity invalid in 30 seconds
        And no dataset should be linked to the removed sink anymore
        And this agent's heartbeat shows that 0 policies are applied to the agent
        And the container logs should contain the message "completed RPC subscription to group" within 30 seconds


@smoke
Scenario: remove one sink from a dataset with 1 sinks, edit the dataset and insert another sink
    Given the Orb user has a registered account
        And the Orb user logs in
        And that 2 sinks already exists
        And that an agent with 1 orb tag(s) already exists and is online
        And pktvisor state is running
        And 1 Agent Group(s) is created with all tags contained in the agent
        And a new policy is created using: handler=dns
        And 1 new dataset is created using the policy, last group and 1 sinks
        And this agent's heartbeat shows that 1 policies are applied and all has status running
        And remove 1 of the linked sinks from orb
        And this agent's heartbeat shows that 1 groups are matching the agent
        And 0 dataset(s) have validity valid and 1 have validity invalid in 30 seconds
        And this agent's heartbeat shows that 0 policies are applied to the agent
    When the dataset is edited and 1 sinks are linked
    Then 1 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And no dataset should be linked to the removed sink anymore
        And this agent's heartbeat shows that 1 policies are applied and all has status running


@MUTE
#@smoke
Scenario: Remotely restart agents with policies applied
    Given the Orb user has a registered account
        And the Orb user logs in
        And that an agent with 1 orb tag(s) already exists and is online
        And pktvisor state is running
        And referred agent is subscribed to 1 group
        And this agent's heartbeat shows that 1 groups are matching the agent
        And that a sink with default configuration type already exists
        And 2 simple policies are applied to the group
        And this agent's heartbeat shows that 2 policies are applied and all has status running
    When remotely restart the agent
    Then the container logs that were output after reset the agent contain the message "pktvisor process stopped" within 30 seconds
        And the container logs should contain the message "all backends and comms were restarted" within 30 seconds
        And the container logs that were output after reset the agent contain the message "removing policies" within 30 seconds
        And the container logs that were output after reset the agent contain the message "resetting backend" within 30 seconds
        And the container logs that were output after reset the agent contain the message "all backends and comms were restarted" within 30 seconds
        And the container logs that were output after reset the agent contain the message "policy applied successfully" referred to each applied policy within 30 seconds
        And the container logs that were output after reset the agent contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds

@MUTE
#@smoke
Scenario: Remotely restart agents without policies applied
    Given the Orb user has a registered account
        And the Orb user logs in
        And that an agent with 1 orb tag(s) already exists and is online
        And pktvisor state is running
        And referred agent is subscribed to 1 group
        And this agent's heartbeat shows that 1 groups are matching the agent
        And that a sink with default configuration type already exists
    When remotely restart the agent
        And the container logs that were output after reset the agent contain the message "resetting backend" within 30 seconds
        And the container logs that were output after reset the agent contain the message "pktvisor process stopped" within 30 seconds
        And the container logs that were output after reset the agent contain the message "all backends and comms were restarted" within 30 seconds
        And 2 simple policies are applied to the group
    Then the container logs should contain the message "all backends and comms were restarted" within 30 seconds
        And this agent's heartbeat shows that 2 policies are applied and all has status running
        And the container logs that were output after reset the agent contain the message "policy applied successfully" referred to each applied policy within 20 seconds
        And the container logs that were output after reset the agent contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds


@smoke
Scenario: Create duplicated policy
    Given the Orb user has a registered account
        And the Orb user logs in
        And that an agent with 1 orb tag(s) already exists and is online
        And pktvisor state is running
        And referred agent is subscribed to 1 group
        And this agent's heartbeat shows that 1 groups are matching the agent
        And that a sink with default configuration type already exists
    When 1 simple policies are applied to the group
        And 1 duplicated policies is applied to the group
    Then this agent's heartbeat shows that 2 policies are applied and all has status running
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 120 seconds
        And 2 dataset(s) have validity valid and 0 have validity invalid in 30 seconds


@smoke
Scenario: Remove agent (check dataset)
    Given the Orb user has a registered account
        And the Orb user logs in
        And a new agent is created with 1 orb tag(s)
        And the agent container is started on an available port
        And the agent status is online
        And referred agent is subscribed to 1 group
        And that a sink with default configuration type already exists
        And 2 simple policies are applied to the group
    When this agent is removed
    Then 0 agent must be matching on response field matching_agents of the last group created
        And the container logs should contain the message "ERROR mqtt log" within 120 seconds
        And the container logs should contain the message "error reconnecting with MQTT, stopping agent" within 120 seconds
        And last container created is exited within 70 seconds
        And the container logs should not contain any panic message
        And last container created is exited after 120 seconds
        And 2 dataset(s) have validity valid and 0 have validity invalid in 30 seconds

@MUTE
#@sanity @sink_status_error
Scenario: Edit sink active and use invalid remote host
    Given the Orb user has a registered account
        And the Orb user logs in
        And 1 Agent Group(s) is created with 3 orb tag(s)
        And the name, tags, description of last Agent Group is edited using: name=edited_before_policy/ tags=2 orb tag(s)/ description=edited
        And that a sink with default configuration type already exists
        And 10 simple policies are applied to the group
    When a new agent is created with orb tags matching 1 existing group
        And the agent container is started on an available port
        And the agent status is online
        And this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 10 policies are applied and all has status running
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 120 seconds
        And the sink remote host is edited and an invalid one is used
    Then referred sink must have error state on response within 120 seconds
        And 10 dataset(s) have validity valid and 0 have validity invalid in 30 seconds

@MUTE
#@sanity @sink_status_error
Scenario: Edit sink active and use invalid username
    Given the Orb user has a registered account
        And the Orb user logs in
        And 1 Agent Group(s) is created with 3 orb tag(s)
        And the name, tags, description of last Agent Group is edited using: name=edited_before_policy/ tags=2 orb tag(s)/ description=edited
        And that a sink with default configuration type already exists
        And 10 simple policies are applied to the group
    When a new agent is created with orb tags matching 1 existing group
        And the agent container is started on an available port
        And the agent status is online
        And this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 10 policies are applied and all has status running
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 120 seconds
        And the sink username is edited and an invalid one is used
    Then referred sink must have error state on response within 120 seconds
        And 10 dataset(s) have validity valid and 0 have validity invalid in 30 seconds

@MUTE
#@sanity @sink_status_error
Scenario: Edit sink active and use invalid password
    Given the Orb user has a registered account
        And the Orb user logs in
        And 1 Agent Group(s) is created with 3 orb tag(s)
        And the name, tags, description of last Agent Group is edited using: name=edited_before_policy/ tags=2 orb tag(s)/ description=edited
        And that a sink with default configuration type already exists
        And 10 simple policies are applied to the group
    When a new agent is created with orb tags matching 1 existing group
        And the agent container is started on an available port
        And the agent status is online
        And this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 10 policies are applied and all has status running
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 120 seconds
        And the sink password is edited and an invalid one is used
    Then referred sink must have error state on response within 120 seconds
        And 10 dataset(s) have validity valid and 0 have validity invalid in 30 seconds

@MUTE
#@sanity @sink_status_error
Scenario: Edit sink with invalid username and use valid one
    Given the Orb user has a registered account
        And the Orb user logs in
        And that an agent with 1 orb tag(s) already exists and is online
        And pktvisor state is running
        And referred agent is subscribed to 1 group
        And this agent's heartbeat shows that 1 groups are matching the agent
        And that a sink with invalid username already exists
        And 3 simple policies are applied to the group
        And that a policy using: handler=dns, description='policy_dns', host_specification=10.0.1.0/24,10.0.2.1/32,2001:db8::/64, bpf_filter_expression=udp port 53, pcap_source=libpcap, only_qname_suffix=[.foo.com/ .example.com], only_rcode=3 already exists
    When 1 new dataset is created using the policy, last group and 1 sink
        And this agent's heartbeat shows that 4 policies are applied and all has status running
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have error state on response within 120 seconds
        And 4 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And the sink username is edited and an valid one is used
    Then referred sink must have active state on response within 120 seconds
        And 4 dataset(s) have validity valid and 0 have validity invalid in 30 seconds

@MUTE
#@sanity @sink_status_error
Scenario: Edit sink with password and use valid one
    Given the Orb user has a registered account
        And the Orb user logs in
        And that an agent with 1 orb tag(s) already exists and is online
        And pktvisor state is running
        And referred agent is subscribed to 1 group
        And this agent's heartbeat shows that 1 groups are matching the agent
        And that a sink with invalid password already exists
        And 3 simple policies are applied to the group
        And that a policy using: handler=dns, description='policy_dns', host_specification=10.0.1.0/24,10.0.2.1/32,2001:db8::/64, bpf_filter_expression=udp port 53, pcap_source=libpcap, only_qname_suffix=[.foo.com/ .example.com], only_rcode=5 already exists
        And 1 new dataset is created using the policy, last group and 1 sink
    When this agent's heartbeat shows that 4 policies are applied and all has status running
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have error state on response within 120 seconds
        And 4 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And the sink password is edited and an valid one is used
    Then referred sink must have active state on response within 120 seconds
        And 4 dataset(s) have validity valid and 0 have validity invalid in 30 seconds


@sanity
Scenario: Check policies status when agent backend stop running
    Given the Orb user has a registered account
        And the Orb user logs in
        And that an agent with 1 orb tag(s) already exists and is online
        And pktvisor state is running
        And referred agent is subscribed to 1 group
        And this agent's heartbeat shows that 1 groups are matching the agent
        And that a sink with default configuration type already exists
        And 5 mixed policies are applied to the group
        And this agent's heartbeat shows that 5 policies are applied and all has status running
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And 5 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
    When agent backend (pktvisor) stops running
        And pktvisor state is backend_error
    Then this agent's heartbeat shows that 5 policies are applied and all has status unknown


@sanity
Scenario: Check backend status when agent backend stop running
    Given the Orb user has a registered account
        And the Orb user logs in
        And that an agent with 1 orb tag(s) already exists and is online
        And pktvisor state is running
    When agent backend (pktvisor) stops running
    Then pktvisor state is backend_error
        And pktvisor error is failed to retrieve backend status: signal: killed


@sanity
Scenario: Check auto reset after pktvisor stop running
    Given the Orb user has a registered account
        And the Orb user logs in
        And that an agent with 1 orb tag(s) already exists and is online
        And pktvisor state is running
        And referred agent is subscribed to 1 group
        And this agent's heartbeat shows that 1 groups are matching the agent
        And that a sink with default configuration type already exists
        And 10 mixed policies are applied to the group
        Then this agent's heartbeat shows that 10 policies are applied and all has status running
    When agent backend (pktvisor) stops running
        And pktvisor state is backend_error
        And pktvisor error is failed to retrieve backend status: signal: killed
        And this agent's heartbeat shows that 10 policies are applied and all has status unknown
    Then agent backend pktvisor restart_count is 1
        And pktvisor state is running
        And this agent's heartbeat shows that 10 policies are applied and all has status running


@sanity
Scenario: Check new policies applied after pktvisor stop running
    Given the Orb user has a registered account
        And the Orb user logs in
        And that an agent with 1 orb tag(s) already exists and is online
        And pktvisor state is running
        And referred agent is subscribed to 1 group
        And this agent's heartbeat shows that 1 groups are matching the agent
        And that a sink with default configuration type already exists
        And 10 mixed policies are applied to the group
        And this agent's heartbeat shows that 10 policies are applied and all has status running
        And agent backend (pktvisor) stops running
        And pktvisor state is backend_error
        And pktvisor error is failed to retrieve backend status: signal: killed
        And this agent's heartbeat shows that 10 policies are applied and all has status unknown
        And agent backend pktvisor restart_count is 1
        And pktvisor state is running
        And this agent's heartbeat shows that 10 policies are applied and all has status running
    When 2 mixed policies are applied to the group
    Then this agent's heartbeat shows that 12 policies are applied and all has status running


@sanity @sink_partial_update
Scenario: Partial Update: sink status after updating only sink name
    Given the Orb user has a registered account
        And the Orb user logs in
        And that an agent with 1 orb tag(s) already exists and is online
        And pktvisor state is running
        And referred agent is subscribed to 1 group
        And this agent's heartbeat shows that 1 groups are matching the agent
        And that a sink with default configuration type already exists
        And 2 simple policies are applied to the group
        And this agent's heartbeat shows that 2 policies are applied and all has status running
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 360 seconds
    When the name of this sink is updated
    Then the name updates to the new value and other fields remains the same
        And referred sink must have active state on response after 360 seconds


@sanity @sink_partial_update
Scenario: Partial Update: sink status after updating only sink description
    Given the Orb user has a registered account
        And the Orb user logs in
        And that an agent with 1 orb tag(s) already exists and is online
        And pktvisor state is running
        And referred agent is subscribed to 1 group
        And this agent's heartbeat shows that 1 groups are matching the agent
        And that a sink with default configuration type already exists
        And 2 simple policies are applied to the group
        And this agent's heartbeat shows that 2 policies are applied and all has status running
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 360 seconds
    When the description of this sink is updated
    Then the description updates to the new value and other fields remains the same
      And referred sink must have active state on response after 360 seconds


@sanity @sink_partial_update
Scenario: Partial Update: sink status after updating only sink tags
    Given the Orb user has a registered account
        And the Orb user logs in
        And that an agent with 1 orb tag(s) already exists and is online
        And pktvisor state is running
        And referred agent is subscribed to 1 group
        And this agent's heartbeat shows that 1 groups are matching the agent
        And that a sink with default configuration type already exists
        And 2 simple policies are applied to the group
        And this agent's heartbeat shows that 2 policies are applied and all has status running
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 360 seconds
    When the tags of this sink is updated
    Then the tags updates to the new value and other fields remains the same
        And referred sink must have active state on response after 360 seconds


@sanity @sink_partial_update
Scenario: Partial Update: sink status after updating only sink configs
    Given the Orb user has a registered account
        And the Orb user logs in
        And that an agent with 1 orb tag(s) already exists and is online
        And pktvisor state is running
        And referred agent is subscribed to 1 group
        And this agent's heartbeat shows that 1 groups are matching the agent
        And that a sink with invalid password already exists
        And 3 simple policies are applied to the group
        And this agent's heartbeat shows that 3 policies are applied and all has status running
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And referred sink must have error state on response within 360 seconds
    When the config of this sink is updated
    Then the config updates to the new value and other fields remains the same
        And referred sink must have active state on response within 360 seconds


@sanity @sink_partial_update
Scenario: Partial Update: sink status after updating only sink name and description
    Given the Orb user has a registered account
        And the Orb user logs in
        And that an agent with 1 orb tag(s) already exists and is online
        And pktvisor state is running
        And referred agent is subscribed to 1 group
        And this agent's heartbeat shows that 1 groups are matching the agent
        And that a sink with default configuration type already exists
        And 2 simple policies are applied to the group
        And this agent's heartbeat shows that 2 policies are applied and all has status running
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 360 seconds
    When the name and description of this sink is updated
    Then the name and description updates to the new value and other fields remains the same
        And referred sink must have active state on response after 360 seconds


@sanity @sink_partial_update
Scenario: Partial Update: sink status after updating only sink name and configs
    Given the Orb user has a registered account
        And the Orb user logs in
        And that an agent with 1 orb tag(s) already exists and is online
        And pktvisor state is running
        And referred agent is subscribed to 1 group
        And this agent's heartbeat shows that 1 groups are matching the agent
        And that a sink with invalid password already exists
        And 3 simple policies are applied to the group
        And this agent's heartbeat shows that 3 policies are applied and all has status running
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And referred sink must have error state on response within 360 seconds
    When the name and config of this sink is updated
    Then the name and config updates to the new value and other fields remains the same
        And referred sink must have active state on response within 360 seconds


@sanity @sink_partial_update
Scenario: Partial Update: sink status after updating only sink name and tags
    Given the Orb user has a registered account
        And the Orb user logs in
        And that an agent with 1 orb tag(s) already exists and is online
        And pktvisor state is running
        And referred agent is subscribed to 1 group
        And this agent's heartbeat shows that 1 groups are matching the agent
        And that a sink with default configuration type already exists
        And 2 simple policies are applied to the group
        And this agent's heartbeat shows that 2 policies are applied and all has status running
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 360 seconds
    When the name and tags of this sink is updated
    Then the name and tags updates to the new value and other fields remains the same
        And referred sink must have active state on response after 360 seconds


@sanity @sink_partial_update
Scenario: Partial Update: sink status after updating only sink description and tags
    Given the Orb user has a registered account
        And the Orb user logs in
        And that an agent with 1 orb tag(s) already exists and is online
        And pktvisor state is running
        And referred agent is subscribed to 1 group
        And this agent's heartbeat shows that 1 groups are matching the agent
        And that a sink with default configuration type already exists
        And 2 simple policies are applied to the group
        And this agent's heartbeat shows that 2 policies are applied and all has status running
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 360 seconds
    When the description and tags of this sink is updated
    Then the description and tags updates to the new value and other fields remains the same
        And referred sink must have active state on response after 360 seconds


@sanity @sink_partial_update
Scenario: Partial Update: sink status after updating only sink description and configs
    Given the Orb user has a registered account
        And the Orb user logs in
        And that an agent with 1 orb tag(s) already exists and is online
        And pktvisor state is running
        And referred agent is subscribed to 1 group
        And this agent's heartbeat shows that 1 groups are matching the agent
        And that a sink with invalid password already exists
        And 3 simple policies are applied to the group
        And this agent's heartbeat shows that 3 policies are applied and all has status running
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And referred sink must have error state on response within 360 seconds
    When the description and config of this sink is updated
    Then the description and config updates to the new value and other fields remains the same
        And referred sink must have active state on response within 360 seconds


@sanity @sink_partial_update
Scenario: Partial Update: sink status after updating only sink tags and configs
    Given the Orb user has a registered account
        And the Orb user logs in
        And that an agent with 1 orb tag(s) already exists and is online
        And pktvisor state is running
        And referred agent is subscribed to 1 group
        And this agent's heartbeat shows that 1 groups are matching the agent
        And that a sink with invalid password already exists
        And 3 simple policies are applied to the group
        And this agent's heartbeat shows that 3 policies are applied and all has status running
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And referred sink must have error state on response within 360 seconds
    When the tags and config of this sink is updated
    Then the tags and config updates to the new value and other fields remains the same
        And referred sink must have active state on response within 360 seconds


@sanity @sink_partial_update
Scenario: Partial Update: sink status after updating only sink name, description and tags
    Given the Orb user has a registered account
        And the Orb user logs in
        And that an agent with 1 orb tag(s) already exists and is online
        And pktvisor state is running
        And referred agent is subscribed to 1 group
        And this agent's heartbeat shows that 1 groups are matching the agent
        And that a sink with default configuration type already exists
        And 2 simple policies are applied to the group
        And this agent's heartbeat shows that 2 policies are applied and all has status running
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 360 seconds
    When the name, description and tags of this sink is updated
    Then the name, description and tags updates to the new value and other fields remains the same
        And referred sink must have active state on response after 360 seconds



@sanity @sink_partial_update
Scenario: Partial Update: sink status after updating only sink name, description and configs
    Given the Orb user has a registered account
        And the Orb user logs in
        And that an agent with 1 orb tag(s) already exists and is online
        And pktvisor state is running
        And referred agent is subscribed to 1 group
        And this agent's heartbeat shows that 1 groups are matching the agent
        And that a sink with invalid password already exists
        And 3 simple policies are applied to the group
        And this agent's heartbeat shows that 3 policies are applied and all has status running
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And referred sink must have error state on response within 360 seconds
    When the name, description and config of this sink is updated
    Then the name, description and config updates to the new value and other fields remains the same
        And referred sink must have active state on response within 360 seconds


@sanity @sink_partial_update
Scenario: Partial Update: sink status after updating only sink name, tags and configs
    Given the Orb user has a registered account
        And the Orb user logs in
        And that an agent with 1 orb tag(s) already exists and is online
        And pktvisor state is running
        And referred agent is subscribed to 1 group
        And this agent's heartbeat shows that 1 groups are matching the agent
        And that a sink with invalid password already exists
        And 3 simple policies are applied to the group
        And this agent's heartbeat shows that 3 policies are applied and all has status running
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And referred sink must have error state on response within 360 seconds
    When the name, tags and config of this sink is updated
    Then the name, tags and config updates to the new value and other fields remains the same
        And referred sink must have active state on response within 360 seconds


@sanity @sink_partial_update
Scenario: Partial Update: sink status after updating only sink description, tags and configs
    Given the Orb user has a registered account
        And the Orb user logs in
        And that an agent with 1 orb tag(s) already exists and is online
        And pktvisor state is running
        And referred agent is subscribed to 1 group
        And this agent's heartbeat shows that 1 groups are matching the agent
        And that a sink with invalid password already exists
        And 3 simple policies are applied to the group
        And this agent's heartbeat shows that 3 policies are applied and all has status running
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And referred sink must have error state on response within 360 seconds
    When the description, tags and config of this sink is updated
    Then the description, tags and config updates to the new value and other fields remains the same
        And referred sink must have active state on response within 360 seconds

@smoke @sink_yaml
Scenario: Using sink with yaml configuration
    Given the Orb user has a registered account
        And the Orb user logs in
        And that an agent with 1 orb tag(s) already exists and is online
        And pktvisor state is running
        And referred agent is subscribed to 1 group
        And this agent's heartbeat shows that 1 groups are matching the agent
        And that a sink with yaml configuration type already exists
    When 10 mixed policies are applied to the group
    Then this agent's heartbeat shows that 10 policies are applied and all has status running
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And 10 dataset(s) have validity valid and 0 have validity invalid in 30 seconds

@smoke
Scenario: Using sink with json configuration specified
    Given the Orb user has a registered account
        And the Orb user logs in
        And that an agent with 1 orb tag(s) already exists and is online
        And pktvisor state is running
        And referred agent is subscribed to 1 group
        And this agent's heartbeat shows that 1 groups are matching the agent
        And that a sink with json configuration type already exists
    When 10 mixed policies are applied to the group
    Then this agent's heartbeat shows that 10 policies are applied and all has status running
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And 10 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
