@metrics @AUTORETRY
Feature: Integration tests validating metric groups

#### netprobe

@sanity @metric_groups @metrics_netprobe @auto_provision
Scenario: netprobe handler with default metric groups configuration
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink with default configuration type already exists
    When an agent(backend_type:pktvisor, settings: {"input_type":"netprobe","test_type":"ping", "targets": {"www.google.com": {"target": "www.google.com"}, "orb_community": {"target": "orb.community"}}}) is self-provisioned via a configuration file on port available with 1 agent tags and has status online. [Overwrite default: False. Paste only file: True]
        And pktvisor state is running
        And 1 Agent Group(s) is created with all tags contained in the agent
        And a netprobe policy netprobe with tap_selector matching all tag(s) of the tap from an agent, default metric_groups enabled, default metric_groups disabled and settings: default is applied to the group
        And 1 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 1 policies are applied and all has status running
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 240 seconds
    Then metrics must be correctly generated for netprobe handler
        And remove the agent .yaml generated on each scenario


@sanity @metric_groups @metrics_netprobe @auto_provision
Scenario: netprobe handler with all metric groups enabled
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink with default configuration type already exists
    When an agent(backend_type:pktvisor, settings: {"input_type":"netprobe","test_type":"ping", "targets": {"www.google.com": {"target": "www.google.com"}, "orb_community": {"target": "orb.community"}}}) is self-provisioned via a configuration file on port available with 1 agent tags and has status online. [Overwrite default: False. Paste only file: True]
        And pktvisor state is running
        And 1 Agent Group(s) is created with all tags contained in the agent
        And a netprobe policy netprobe with tap_selector matching all tag(s) of the tap from an agent, all metric_groups enabled, none metric_groups disabled and settings: default is applied to the group
        And 1 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 1 policies are applied and all has status running
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 240 seconds
    Then metrics must be correctly generated for netprobe handler
        And remove the agent .yaml generated on each scenario


@sanity @metric_groups @metrics_netprobe @auto_provision
Scenario: netprobe handler with all metric groups disabled
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink with default configuration type already exists
    When an agent(backend_type:pktvisor, settings: {"input_type":"netprobe","test_type":"ping", "targets": {"www.google.com": {"target": "www.google.com"}, "orb_community": {"target": "orb.community"}}}) is self-provisioned via a configuration file on port available with 1 agent tags and has status online. [Overwrite default: False. Paste only file: True]
        And pktvisor state is running
        And 1 Agent Group(s) is created with all tags contained in the agent
        And a netprobe policy netprobe with tap_selector matching all tag(s) of the tap from an agent, none metric_groups enabled, all metric_groups disabled and settings: default is applied to the group
        And 1 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 1 policies are applied and all has status running
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
    Then metrics must be correctly generated for netprobe handler
        And remove the agent .yaml generated on each scenario


@sanity @metric_groups @metrics_netprobe @auto_provision
Scenario: netprobe handler with only counters metric groups enabled
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink with default configuration type already exists
    When an agent(backend_type:pktvisor, settings: {"input_type":"netprobe","test_type":"ping", "targets": {"www.google.com": {"target": "www.google.com"}, "orb_community": {"target": "orb.community"}}}) is self-provisioned via a configuration file on port available with 1 agent tags and has status online. [Overwrite default: False. Paste only file: True]
        And pktvisor state is running
        And 1 Agent Group(s) is created with all tags contained in the agent
        And a netprobe policy netprobe with tap_selector matching all tag(s) of the tap from an agent, counters metric_groups enabled, quantiles, histograms metric_groups disabled and settings: default is applied to the group
        And 1 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 1 policies are applied and all has status running
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 240 seconds
    Then metrics must be correctly generated for netprobe handler
        And remove the agent .yaml generated on each scenario


@sanity @metric_groups @metrics_netprobe @auto_provision
Scenario: netprobe handler with only quantiles metric groups enabled
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink with default configuration type already exists
    When an agent(backend_type:pktvisor, settings: {"input_type":"netprobe","test_type":"ping", "targets": {"www.google.com": {"target": "www.google.com"}, "orb_community": {"target": "orb.community"}}}) is self-provisioned via a configuration file on port available with 1 agent tags and has status online. [Overwrite default: False. Paste only file: True]
        And pktvisor state is running
        And 1 Agent Group(s) is created with all tags contained in the agent
        And a netprobe policy netprobe with tap_selector matching all tag(s) of the tap from an agent, quantiles metric_groups enabled, counters, histograms metric_groups disabled and settings: default is applied to the group
        And 1 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 1 policies are applied and all has status running
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 240 seconds
    Then metrics must be correctly generated for netprobe handler
        And remove the agent .yaml generated on each scenario


@sanity @metric_groups @metrics_netprobe @auto_provision
Scenario: netprobe handler with only histograms metric groups enabled
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink with default configuration type already exists
    When an agent(backend_type:pktvisor, settings: {"input_type":"netprobe","test_type":"ping", "targets": {"www.google.com": {"target": "www.google.com"}, "orb_community": {"target": "orb.community"}}}) is self-provisioned via a configuration file on port available with 1 agent tags and has status online. [Overwrite default: False. Paste only file: True]
        And pktvisor state is running
        And 1 Agent Group(s) is created with all tags contained in the agent
        And a netprobe policy netprobe with tap_selector matching all tag(s) of the tap from an agent, histograms metric_groups enabled, quantiles, counters metric_groups disabled and settings: default is applied to the group
        And 1 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 1 policies are applied and all has status running
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 240 seconds
    Then metrics must be correctly generated for netprobe handler
        And remove the agent .yaml generated on each scenario


@sanity @metric_groups @metrics_netprobe @auto_provision
Scenario: netprobe handler with counters and histograms metric groups enabled and quantiles disabled
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink with default configuration type already exists
    When an agent(backend_type:pktvisor, settings: {"input_type":"netprobe","test_type":"ping", "targets": {"www.google.com": {"target": "www.google.com"}, "orb_community": {"target": "orb.community"}}}) is self-provisioned via a configuration file on port available with 1 agent tags and has status online. [Overwrite default: False. Paste only file: True]
        And pktvisor state is running
        And 1 Agent Group(s) is created with all tags contained in the agent
        And a netprobe policy netprobe with tap_selector matching all tag(s) of the tap from an agent, counters, histograms metric_groups enabled, quantiles metric_groups disabled and settings: default is applied to the group
        And 1 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 1 policies are applied and all has status running
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 240 seconds
    Then metrics must be correctly generated for netprobe handler
        And remove the agent .yaml generated on each scenario


@sanity @metric_groups @metrics_netprobe @auto_provision
Scenario: netprobe handler with counters and quantiles metric groups enabled and histograms disabled
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink with default configuration type already exists
    When an agent(backend_type:pktvisor, settings: {"input_type":"netprobe","test_type":"ping", "targets": {"www.google.com": {"target": "www.google.com"}, "orb_community": {"target": "orb.community"}}}) is self-provisioned via a configuration file on port available with 1 agent tags and has status online. [Overwrite default: False. Paste only file: True]
        And pktvisor state is running
        And 1 Agent Group(s) is created with all tags contained in the agent
        And a netprobe policy netprobe with tap_selector matching all tag(s) of the tap from an agent, counters, quantiles metric_groups enabled, histograms metric_groups disabled and settings: default is applied to the group
        And 1 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 1 policies are applied and all has status running
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 240 seconds
    Then metrics must be correctly generated for netprobe handler
        And remove the agent .yaml generated on each scenario


@sanity @metric_groups @metrics_netprobe @auto_provision
Scenario: netprobe handler with histograms and quantiles metric groups enabled and counters disabled
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink with default configuration type already exists
    When an agent(backend_type:pktvisor, settings: {"input_type":"netprobe","test_type":"ping", "targets": {"www.google.com": {"target": "www.google.com"}, "orb_community": {"target": "orb.community"}}}) is self-provisioned via a configuration file on port available with 1 agent tags and has status online. [Overwrite default: False. Paste only file: True]
        And pktvisor state is running
        And 1 Agent Group(s) is created with all tags contained in the agent
        And a netprobe policy netprobe with tap_selector matching all tag(s) of the tap from an agent, histograms, quantiles metric_groups enabled, counters metric_groups disabled and settings: default is applied to the group
        And 1 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 1 policies are applied and all has status running
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 240 seconds
    Then metrics must be correctly generated for netprobe handler
        And remove the agent .yaml generated on each scenario

#### flow netflow

@sanity @metric_groups @metrics_flow @root @mocked_interface @auto_provision
Scenario: flow handler type netflow with default metric groups configuration
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink with default configuration type already exists
        And a mocked interface is configured with mtu: 65000 and ip: 192.168.100.2/24
        And a virtual switch is configured and is up with type netflow and target \"192.168.100.2:available\"
    When an agent(backend_type:pktvisor, settings: {"input_type":"flow","bind":"192.168.100.2", "port":"switch", "flow_type":"netflow"}) is self-provisioned via a configuration file on port available with 1 agent tags and has status online. [Overwrite default: False. Paste only file: True]
        And pktvisor state is running
        And 1 Agent Group(s) is created with all tags contained in the agent
        And a flow policy flow with tap_selector matching all tag(s) of the tap from an agent, default metric_groups enabled, default metric_groups disabled and settings: default is applied to the group
        And 1 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 1 policies are applied and all has status running
        And run mocked data dns_ipv4_tcp.pcap, dns_ipv4_udp.pcap, dns_ipv6_tcp.pcap, dns_ipv6_udp.pcap, dns_udp_mixed_rcode.pcap, dnssec.pcap, dhcpv6.pcap, dhcp-flow.pcap on the created virtual switch
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 240 seconds
    Then metrics must be correctly generated for flow handler
        And remove the agent .yaml generated on each scenario


@sanity @metric_groups @metrics_flow @root @mocked_interface @auto_provision
Scenario: flow handler type netflow with all metric groups enabled
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink with default configuration type already exists
        And a mocked interface is configured with mtu: 65000 and ip: 192.168.100.2/24
        And a virtual switch is configured and is up with type netflow and target \"192.168.100.2:available\"
    When an agent(backend_type:pktvisor, settings: {"input_type":"flow","bind":"192.168.100.2", "port":"switch", "flow_type":"netflow"}) is self-provisioned via a configuration file on port available with 1 agent tags and has status online. [Overwrite default: False. Paste only file: True]
        And pktvisor state is running
        And 1 Agent Group(s) is created with all tags contained in the agent
        And a flow policy flow with tap_selector matching all tag(s) of the tap from an agent, all metric_groups enabled, none metric_groups disabled and settings: default is applied to the group
        And 1 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 1 policies are applied and all has status running
        And run mocked data dns_ipv4_tcp.pcap, dns_ipv4_udp.pcap, dns_ipv6_tcp.pcap, dns_ipv6_udp.pcap, dns_udp_mixed_rcode.pcap, dnssec.pcap, dhcpv6.pcap, dhcp-flow.pcap on the created virtual switch
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 240 seconds
    Then metrics must be correctly generated for flow handler
        And remove the agent .yaml generated on each scenario


@sanity @metric_groups @metrics_flow @root @mocked_interface @auto_provision
Scenario: flow handler type netflow with all metric groups disabled
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink with default configuration type already exists
        And a mocked interface is configured with mtu: 65000 and ip: 192.168.100.2/24
        And a virtual switch is configured and is up with type netflow and target \"192.168.100.2:available\"
    When an agent(backend_type:pktvisor, settings: {"input_type":"flow","bind":"192.168.100.2", "port":"switch", "flow_type":"netflow"}) is self-provisioned via a configuration file on port available with 1 agent tags and has status online. [Overwrite default: False. Paste only file: True]
        And pktvisor state is running
        And 1 Agent Group(s) is created with all tags contained in the agent
        And a flow policy flow with tap_selector matching all tag(s) of the tap from an agent, none metric_groups enabled, all metric_groups disabled and settings: default is applied to the group
        And 1 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 1 policies are applied and all has status running
        And run mocked data dns_ipv4_tcp.pcap, dns_ipv4_udp.pcap, dns_ipv6_tcp.pcap, dns_ipv6_udp.pcap, dns_udp_mixed_rcode.pcap, dnssec.pcap, dhcpv6.pcap, dhcp-flow.pcap on the created virtual switch
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
    Then metrics must be correctly generated for flow handler
        And remove the agent .yaml generated on each scenario


@sanity @metric_groups @metrics_flow @root @mocked_interface @auto_provision
Scenario: flow handler type netflow with only cardinality metric groups enabled
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink with default configuration type already exists
        And a mocked interface is configured with mtu: 65000 and ip: 192.168.100.2/24
        And a virtual switch is configured and is up with type netflow and target \"192.168.100.2:available\"
    When an agent(backend_type:pktvisor, settings: {"input_type":"flow","bind":"192.168.100.2", "port":"switch", "flow_type":"netflow"}) is self-provisioned via a configuration file on port available with 1 agent tags and has status online. [Overwrite default: False. Paste only file: True]
        And pktvisor state is running
        And 1 Agent Group(s) is created with all tags contained in the agent
        And a flow policy flow with tap_selector matching all tag(s) of the tap from an agent, cardinality metric_groups enabled, top_tos, counters, by_packets, by_bytes, top_geo, conversations, top_ports, top_ips_ports, top_interfaces, top_ips metric_groups disabled and settings: default is applied to the group
        And 1 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 1 policies are applied and all has status running
        And run mocked data dns_ipv4_tcp.pcap, dns_ipv4_udp.pcap, dns_ipv6_tcp.pcap, dns_ipv6_udp.pcap, dns_udp_mixed_rcode.pcap, dnssec.pcap, dhcpv6.pcap, dhcp-flow.pcap on the created virtual switch
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 240 seconds
    Then metrics must be correctly generated for flow handler
        And remove the agent .yaml generated on each scenario


@sanity @metric_groups @metrics_flow @root @mocked_interface @auto_provision
Scenario: flow handler type netflow with only counters metric groups enabled
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink with default configuration type already exists
        And a mocked interface is configured with mtu: 65000 and ip: 192.168.100.2/24
        And a virtual switch is configured and is up with type netflow and target \"192.168.100.2:available\"
    When an agent(backend_type:pktvisor, settings: {"input_type":"flow","bind":"192.168.100.2", "port":"switch", "flow_type":"netflow"}) is self-provisioned via a configuration file on port available with 1 agent tags and has status online. [Overwrite default: False. Paste only file: True]
        And pktvisor state is running
        And 1 Agent Group(s) is created with all tags contained in the agent
        And a flow policy flow with tap_selector matching all tag(s) of the tap from an agent, counters metric_groups enabled, top_tos, cardinality, by_packets, by_bytes, top_geo, conversations, top_ports, top_ips_ports, top_interfaces, top_ips metric_groups disabled and settings: default is applied to the group
        And 1 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 1 policies are applied and all has status running
        And run mocked data dns_ipv4_tcp.pcap, dns_ipv4_udp.pcap, dns_ipv6_tcp.pcap, dns_ipv6_udp.pcap, dns_udp_mixed_rcode.pcap, dnssec.pcap, dhcpv6.pcap, dhcp-flow.pcap on the created virtual switch
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 240 seconds
    Then metrics must be correctly generated for flow handler
        And remove the agent .yaml generated on each scenario


@sanity @metric_groups @metrics_flow @root @mocked_interface @auto_provision
Scenario: flow handler type netflow with only counters and by_bytes metric groups enabled
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink with default configuration type already exists
        And a mocked interface is configured with mtu: 65000 and ip: 192.168.100.2/24
        And a virtual switch is configured and is up with type netflow and target \"192.168.100.2:available\"
    When an agent(backend_type:pktvisor, settings: {"input_type":"flow","bind":"192.168.100.2", "port":"switch", "flow_type":"netflow"}) is self-provisioned via a configuration file on port available with 1 agent tags and has status online. [Overwrite default: False. Paste only file: True]
        And pktvisor state is running
        And 1 Agent Group(s) is created with all tags contained in the agent
        And a flow policy flow with tap_selector matching all tag(s) of the tap from an agent, counters, by_bytes metric_groups enabled, top_tos, cardinality, by_packets, top_geo, conversations, top_ports, top_ips_ports, top_interfaces, top_ips metric_groups disabled and settings: default is applied to the group
        And 1 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 1 policies are applied and all has status running
        And run mocked data dns_ipv4_tcp.pcap, dns_ipv4_udp.pcap, dns_ipv6_tcp.pcap, dns_ipv6_udp.pcap, dns_udp_mixed_rcode.pcap, dnssec.pcap, dhcpv6.pcap, dhcp-flow.pcap on the created virtual switch
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 240 seconds
    Then metrics must be correctly generated for flow handler
        And remove the agent .yaml generated on each scenario


@sanity @metric_groups @metrics_flow @root @mocked_interface @auto_provision
Scenario: flow handler type netflow with only counters and by_packets metric groups enabled
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink with default configuration type already exists
        And a mocked interface is configured with mtu: 65000 and ip: 192.168.100.2/24
        And a virtual switch is configured and is up with type netflow and target \"192.168.100.2:available\"
    When an agent(backend_type:pktvisor, settings: {"input_type":"flow","bind":"192.168.100.2", "port":"switch", "flow_type":"netflow"}) is self-provisioned via a configuration file on port available with 1 agent tags and has status online. [Overwrite default: False. Paste only file: True]
        And pktvisor state is running
        And 1 Agent Group(s) is created with all tags contained in the agent
        And a flow policy flow with tap_selector matching all tag(s) of the tap from an agent, counters, by_packets metric_groups enabled, top_tos, cardinality, by_bytes, top_geo, conversations, top_ports, top_ips_ports, top_interfaces, top_ips metric_groups disabled and settings: default is applied to the group
        And 1 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 1 policies are applied and all has status running
        And run mocked data dns_ipv4_tcp.pcap, dns_ipv4_udp.pcap, dns_ipv6_tcp.pcap, dns_ipv6_udp.pcap, dns_udp_mixed_rcode.pcap, dnssec.pcap, dhcpv6.pcap, dhcp-flow.pcap on the created virtual switch
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 240 seconds
    Then metrics must be correctly generated for flow handler
        And remove the agent .yaml generated on each scenario


@sanity @metric_groups @metrics_flow @root @mocked_interface @auto_provision
Scenario: flow handler type netflow with only by_packets metric groups enabled
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink with default configuration type already exists
        And a mocked interface is configured with mtu: 65000 and ip: 192.168.100.2/24
        And a virtual switch is configured and is up with type netflow and target \"192.168.100.2:available\"
    When an agent(backend_type:pktvisor, settings: {"input_type":"flow","bind":"192.168.100.2", "port":"switch", "flow_type":"netflow"}) is self-provisioned via a configuration file on port available with 1 agent tags and has status online. [Overwrite default: False. Paste only file: True]
        And pktvisor state is running
        And 1 Agent Group(s) is created with all tags contained in the agent
        And a flow policy flow with tap_selector matching all tag(s) of the tap from an agent, by_packets metric_groups enabled, top_tos, counters, cardinality, by_bytes, top_geo, conversations, top_ports, top_ips_ports, top_interfaces, top_ips metric_groups disabled and settings: default is applied to the group
        And 1 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 1 policies are applied and all has status running
        And run mocked data dns_ipv4_tcp.pcap, dns_ipv4_udp.pcap, dns_ipv6_tcp.pcap, dns_ipv6_udp.pcap, dns_udp_mixed_rcode.pcap, dnssec.pcap, dhcpv6.pcap, dhcp-flow.pcap on the created virtual switch
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
    Then metrics must be correctly generated for flow handler
        And remove the agent .yaml generated on each scenario


@sanity @metric_groups @metrics_flow @root @mocked_interface @auto_provision
Scenario: flow handler type netflow with only by_bytes metric groups enabled
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink with default configuration type already exists
        And a mocked interface is configured with mtu: 65000 and ip: 192.168.100.2/24
        And a virtual switch is configured and is up with type netflow and target \"192.168.100.2:available\"
    When an agent(backend_type:pktvisor, settings: {"input_type":"flow","bind":"192.168.100.2", "port":"switch", "flow_type":"netflow"}) is self-provisioned via a configuration file on port available with 1 agent tags and has status online. [Overwrite default: False. Paste only file: True]
        And pktvisor state is running
        And 1 Agent Group(s) is created with all tags contained in the agent
        And a flow policy flow with tap_selector matching all tag(s) of the tap from an agent, by_bytes metric_groups enabled, top_tos, counters, by_packets, cardinality, top_geo, conversations, top_ports, top_ips_ports, top_interfaces, top_ips metric_groups disabled and settings: default is applied to the group
        And 1 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 1 policies are applied and all has status running
        And run mocked data dns_ipv4_tcp.pcap, dns_ipv4_udp.pcap, dns_ipv6_tcp.pcap, dns_ipv6_udp.pcap, dns_udp_mixed_rcode.pcap, dnssec.pcap, dhcpv6.pcap, dhcp-flow.pcap on the created virtual switch
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
    Then metrics must be correctly generated for flow handler
        And remove the agent .yaml generated on each scenario


@sanity @metric_groups @metrics_flow @root @mocked_interface @auto_provision
Scenario: flow handler type netflow with only top_geo metric groups enabled
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink with default configuration type already exists
        And a mocked interface is configured with mtu: 65000 and ip: 192.168.100.2/24
        And a virtual switch is configured and is up with type netflow and target \"192.168.100.2:available\"
    When an agent(backend_type:pktvisor, settings: {"input_type":"flow","bind":"192.168.100.2", "port":"switch", "flow_type":"netflow"}) is self-provisioned via a configuration file on port available with 1 agent tags and has status online. [Overwrite default: False. Paste only file: True]
        And pktvisor state is running
        And 1 Agent Group(s) is created with all tags contained in the agent
        And a flow policy flow with tap_selector matching all tag(s) of the tap from an agent, top_geo metric_groups enabled, top_tos, counters, by_packets, by_bytes, cardinality, conversations, top_ports, top_ips_ports, top_interfaces, top_ips metric_groups disabled and settings: default is applied to the group
        And 1 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 1 policies are applied and all has status running
        And run mocked data dns_ipv4_tcp.pcap, dns_ipv4_udp.pcap, dns_ipv6_tcp.pcap, dns_ipv6_udp.pcap, dns_udp_mixed_rcode.pcap, dnssec.pcap, dhcpv6.pcap, dhcp-flow.pcap on the created virtual switch
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
    Then metrics must be correctly generated for flow handler
        And remove the agent .yaml generated on each scenario


@sanity @metric_groups @metrics_flow @root @mocked_interface @auto_provision
Scenario: flow handler type netflow with only top_geo and by_bytes metric groups enabled
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink with default configuration type already exists
        And a mocked interface is configured with mtu: 65000 and ip: 192.168.100.2/24
        And a virtual switch is configured and is up with type netflow and target \"192.168.100.2:available\"
    When an agent(backend_type:pktvisor, settings: {"input_type":"flow","bind":"192.168.100.2", "port":"switch", "flow_type":"netflow"}) is self-provisioned via a configuration file on port available with 1 agent tags and has status online. [Overwrite default: False. Paste only file: True]
        And pktvisor state is running
        And 1 Agent Group(s) is created with all tags contained in the agent
        And a flow policy flow with tap_selector matching all tag(s) of the tap from an agent, top_geo, by_bytes metric_groups enabled, top_tos, counters, by_packets, cardinality, conversations, top_ports, top_ips_ports, top_interfaces, top_ips metric_groups disabled and settings: default is applied to the group
        And 1 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 1 policies are applied and all has status running
        And run mocked data dns_ipv4_tcp.pcap, dns_ipv4_udp.pcap, dns_ipv6_tcp.pcap, dns_ipv6_udp.pcap, dns_udp_mixed_rcode.pcap, dnssec.pcap, dhcpv6.pcap, dhcp-flow.pcap on the created virtual switch
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 240 seconds
    Then metrics must be correctly generated for flow handler
        And remove the agent .yaml generated on each scenario

@sanity @metric_groups @metrics_flow @root @mocked_interface @auto_provision
Scenario: flow handler type netflow with only top_geo and by_packets metric groups enabled
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink with default configuration type already exists
        And a mocked interface is configured with mtu: 65000 and ip: 192.168.100.2/24
        And a virtual switch is configured and is up with type netflow and target \"192.168.100.2:available\"
    When an agent(backend_type:pktvisor, settings: {"input_type":"flow","bind":"192.168.100.2", "port":"switch", "flow_type":"netflow"}) is self-provisioned via a configuration file on port available with 1 agent tags and has status online. [Overwrite default: False. Paste only file: True]
        And pktvisor state is running
        And 1 Agent Group(s) is created with all tags contained in the agent
        And a flow policy flow with tap_selector matching all tag(s) of the tap from an agent, top_geo, by_packets metric_groups enabled, top_tos, counters, by_bytes, cardinality, conversations, top_ports, top_ips_ports, top_interfaces, top_ips metric_groups disabled and settings: default is applied to the group
        And 1 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 1 policies are applied and all has status running
        And run mocked data dns_ipv4_tcp.pcap, dns_ipv4_udp.pcap, dns_ipv6_tcp.pcap, dns_ipv6_udp.pcap, dns_udp_mixed_rcode.pcap, dnssec.pcap, dhcpv6.pcap, dhcp-flow.pcap on the created virtual switch
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 240 seconds
    Then metrics must be correctly generated for flow handler
        And remove the agent .yaml generated on each scenario


@sanity @metric_groups @metrics_flow @root @mocked_interface @auto_provision
Scenario: flow handler type netflow with only conversations metric groups enabled
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink with default configuration type already exists
        And a mocked interface is configured with mtu: 65000 and ip: 192.168.100.2/24
        And a virtual switch is configured and is up with type netflow and target \"192.168.100.2:available\"
    When an agent(backend_type:pktvisor, settings: {"input_type":"flow","bind":"192.168.100.2", "port":"switch", "flow_type":"netflow"}) is self-provisioned via a configuration file on port available with 1 agent tags and has status online. [Overwrite default: False. Paste only file: True]
        And pktvisor state is running
        And 1 Agent Group(s) is created with all tags contained in the agent
        And a flow policy flow with tap_selector matching all tag(s) of the tap from an agent, conversations metric_groups enabled, top_tos, counters, by_packets, by_bytes, top_geo, cardinality, top_ports, top_ips_ports, top_interfaces, top_ips metric_groups disabled and settings: default is applied to the group
        And 1 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 1 policies are applied and all has status running
        And run mocked data dns_ipv4_tcp.pcap, dns_ipv4_udp.pcap, dns_ipv6_tcp.pcap, dns_ipv6_udp.pcap, dns_udp_mixed_rcode.pcap, dnssec.pcap, dhcpv6.pcap, dhcp-flow.pcap on the created virtual switch
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
    Then metrics must be correctly generated for flow handler
        And remove the agent .yaml generated on each scenario


@sanity @metric_groups @metrics_flow @root @mocked_interface @auto_provision
Scenario: flow handler type netflow with only conversations and cardinality metric groups enabled
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink with default configuration type already exists
        And a mocked interface is configured with mtu: 65000 and ip: 192.168.100.2/24
        And a virtual switch is configured and is up with type netflow and target \"192.168.100.2:available\"
    When an agent(backend_type:pktvisor, settings: {"input_type":"flow","bind":"192.168.100.2", "port":"switch", "flow_type":"netflow"}) is self-provisioned via a configuration file on port available with 1 agent tags and has status online. [Overwrite default: False. Paste only file: True]
        And pktvisor state is running
        And 1 Agent Group(s) is created with all tags contained in the agent
        And a flow policy flow with tap_selector matching all tag(s) of the tap from an agent, conversations, cardinality metric_groups enabled, top_tos, counters, by_packets, by_bytes, top_geo, top_ports, top_ips_ports, top_interfaces, top_ips metric_groups disabled and settings: default is applied to the group
        And 1 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 1 policies are applied and all has status running
        And run mocked data dns_ipv4_tcp.pcap, dns_ipv4_udp.pcap, dns_ipv6_tcp.pcap, dns_ipv6_udp.pcap, dns_udp_mixed_rcode.pcap, dnssec.pcap, dhcpv6.pcap, dhcp-flow.pcap on the created virtual switch
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 240 seconds
    Then metrics must be correctly generated for flow handler
        And remove the agent .yaml generated on each scenario

@sanity @metric_groups @metrics_flow @root @mocked_interface @auto_provision
Scenario: flow handler type netflow with only conversations and by_bytes metric groups enabled
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink with default configuration type already exists
        And a mocked interface is configured with mtu: 65000 and ip: 192.168.100.2/24
        And a virtual switch is configured and is up with type netflow and target \"192.168.100.2:available\"
    When an agent(backend_type:pktvisor, settings: {"input_type":"flow","bind":"192.168.100.2", "port":"switch", "flow_type":"netflow"}) is self-provisioned via a configuration file on port available with 1 agent tags and has status online. [Overwrite default: False. Paste only file: True]
        And pktvisor state is running
        And 1 Agent Group(s) is created with all tags contained in the agent
        And a flow policy flow with tap_selector matching all tag(s) of the tap from an agent, conversations, by_bytes metric_groups enabled, top_tos, counters, by_packets, top_geo, cardinality, top_ports, top_ips_ports, top_interfaces, top_ips metric_groups disabled and settings: default is applied to the group
        And 1 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 1 policies are applied and all has status running
        And run mocked data dns_ipv4_tcp.pcap, dns_ipv4_udp.pcap, dns_ipv6_tcp.pcap, dns_ipv6_udp.pcap, dns_udp_mixed_rcode.pcap, dnssec.pcap, dhcpv6.pcap, dhcp-flow.pcap on the created virtual switch
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 240 seconds
    Then metrics must be correctly generated for flow handler
        And remove the agent .yaml generated on each scenario


@sanity @metric_groups @metrics_flow @root @mocked_interface @auto_provision
Scenario: flow handler type netflow with only conversations and by_packets metric groups enabled
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink with default configuration type already exists
        And a mocked interface is configured with mtu: 65000 and ip: 192.168.100.2/24
        And a virtual switch is configured and is up with type netflow and target \"192.168.100.2:available\"
    When an agent(backend_type:pktvisor, settings: {"input_type":"flow","bind":"192.168.100.2", "port":"switch", "flow_type":"netflow"}) is self-provisioned via a configuration file on port available with 1 agent tags and has status online. [Overwrite default: False. Paste only file: True]
        And pktvisor state is running
        And 1 Agent Group(s) is created with all tags contained in the agent
        And a flow policy flow with tap_selector matching all tag(s) of the tap from an agent, conversations, by_packets metric_groups enabled, top_tos, counters, by_bytes, top_geo, cardinality, top_ports, top_ips_ports, top_interfaces, top_ips metric_groups disabled and settings: default is applied to the group
        And 1 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 1 policies are applied and all has status running
        And run mocked data dns_ipv4_tcp.pcap, dns_ipv4_udp.pcap, dns_ipv6_tcp.pcap, dns_ipv6_udp.pcap, dns_udp_mixed_rcode.pcap, dnssec.pcap, dhcpv6.pcap, dhcp-flow.pcap on the created virtual switch
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 240 seconds
    Then metrics must be correctly generated for flow handler
        And remove the agent .yaml generated on each scenario


@sanity @metric_groups @metrics_flow @root @mocked_interface @auto_provision
Scenario: flow handler type netflow with only top_ports metric groups enabled
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink with default configuration type already exists
        And a mocked interface is configured with mtu: 65000 and ip: 192.168.100.2/24
        And a virtual switch is configured and is up with type netflow and target \"192.168.100.2:available\"
    When an agent(backend_type:pktvisor, settings: {"input_type":"flow","bind":"192.168.100.2", "port":"switch", "flow_type":"netflow"}) is self-provisioned via a configuration file on port available with 1 agent tags and has status online. [Overwrite default: False. Paste only file: True]
        And pktvisor state is running
        And 1 Agent Group(s) is created with all tags contained in the agent
        And a flow policy flow with tap_selector matching all tag(s) of the tap from an agent, top_ports metric_groups enabled, top_tos, conversations, counters, by_packets, by_bytes, top_geo, cardinality, top_ips_ports, top_interfaces, top_ips metric_groups disabled and settings: default is applied to the group
        And 1 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 1 policies are applied and all has status running
        And run mocked data dns_ipv4_tcp.pcap, dns_ipv4_udp.pcap, dns_ipv6_tcp.pcap, dns_ipv6_udp.pcap, dns_udp_mixed_rcode.pcap, dnssec.pcap, dhcpv6.pcap, dhcp-flow.pcap on the created virtual switch
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
    Then metrics must be correctly generated for flow handler
        And remove the agent .yaml generated on each scenario


@sanity @metric_groups @metrics_flow @root @mocked_interface @auto_provision
Scenario: flow handler type netflow with only top_ports and by_packets metric groups enabled
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink with default configuration type already exists
        And a mocked interface is configured with mtu: 65000 and ip: 192.168.100.2/24
        And a virtual switch is configured and is up with type netflow and target \"192.168.100.2:available\"
    When an agent(backend_type:pktvisor, settings: {"input_type":"flow","bind":"192.168.100.2", "port":"switch", "flow_type":"netflow"}) is self-provisioned via a configuration file on port available with 1 agent tags and has status online. [Overwrite default: False. Paste only file: True]
        And pktvisor state is running
        And 1 Agent Group(s) is created with all tags contained in the agent
        And a flow policy flow with tap_selector matching all tag(s) of the tap from an agent, top_ports, by_packets metric_groups enabled, top_tos, conversations, counters, by_bytes, top_geo, cardinality, top_ips_ports, top_interfaces, top_ips metric_groups disabled and settings: default is applied to the group
        And 1 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 1 policies are applied and all has status running
        And run mocked data dns_ipv4_tcp.pcap, dns_ipv4_udp.pcap, dns_ipv6_tcp.pcap, dns_ipv6_udp.pcap, dns_udp_mixed_rcode.pcap, dnssec.pcap, dhcpv6.pcap, dhcp-flow.pcap on the created virtual switch
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 240 seconds
    Then metrics must be correctly generated for flow handler
        And remove the agent .yaml generated on each scenario


@sanity @metric_groups @metrics_flow @root @mocked_interface @auto_provision
Scenario: flow handler type netflow with only top_ports and by_bytes metric groups enabled
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink with default configuration type already exists
        And a mocked interface is configured with mtu: 65000 and ip: 192.168.100.2/24
        And a virtual switch is configured and is up with type netflow and target \"192.168.100.2:available\"
    When an agent(backend_type:pktvisor, settings: {"input_type":"flow","bind":"192.168.100.2", "port":"switch", "flow_type":"netflow"}) is self-provisioned via a configuration file on port available with 1 agent tags and has status online. [Overwrite default: False. Paste only file: True]
        And pktvisor state is running
        And 1 Agent Group(s) is created with all tags contained in the agent
        And a flow policy flow with tap_selector matching all tag(s) of the tap from an agent, top_ports, by_bytes metric_groups enabled, top_tos, conversations, counters, by_packets, top_geo, cardinality, top_ips_ports, top_interfaces, top_ips metric_groups disabled and settings: default is applied to the group
        And 1 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 1 policies are applied and all has status running
        And run mocked data dns_ipv4_tcp.pcap, dns_ipv4_udp.pcap, dns_ipv6_tcp.pcap, dns_ipv6_udp.pcap, dns_udp_mixed_rcode.pcap, dnssec.pcap, dhcpv6.pcap, dhcp-flow.pcap on the created virtual switch
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 240 seconds
    Then metrics must be correctly generated for flow handler
        And remove the agent .yaml generated on each scenario


@sanity @metric_groups @metrics_flow @root @mocked_interface @auto_provision
Scenario: flow handler type netflow with only top_ips_ports metric groups enabled
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink with default configuration type already exists
        And a mocked interface is configured with mtu: 65000 and ip: 192.168.100.2/24
        And a virtual switch is configured and is up with type netflow and target \"192.168.100.2:available\"
    When an agent(backend_type:pktvisor, settings: {"input_type":"flow","bind":"192.168.100.2", "port":"switch", "flow_type":"netflow"}) is self-provisioned via a configuration file on port available with 1 agent tags and has status online. [Overwrite default: False. Paste only file: True]
        And pktvisor state is running
        And 1 Agent Group(s) is created with all tags contained in the agent
        And a flow policy flow with tap_selector matching all tag(s) of the tap from an agent, top_ips_ports metric_groups enabled, top_tos, conversations, counters, by_packets, by_bytes, top_geo, cardinality, top_ports, top_interfaces, top_ips metric_groups disabled and settings: default is applied to the group
        And 1 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 1 policies are applied and all has status running
        And run mocked data dns_ipv4_tcp.pcap, dns_ipv4_udp.pcap, dns_ipv6_tcp.pcap, dns_ipv6_udp.pcap, dns_udp_mixed_rcode.pcap, dnssec.pcap, dhcpv6.pcap, dhcp-flow.pcap on the created virtual switch
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
    Then metrics must be correctly generated for flow handler
        And remove the agent .yaml generated on each scenario


    @sanity @metric_groups @metrics_flow @root @mocked_interface @auto_provision
Scenario: flow handler type netflow with only top_ips_ports and by_bytes metric groups enabled
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink with default configuration type already exists
        And a mocked interface is configured with mtu: 65000 and ip: 192.168.100.2/24
        And a virtual switch is configured and is up with type netflow and target \"192.168.100.2:available\"
    When an agent(backend_type:pktvisor, settings: {"input_type":"flow","bind":"192.168.100.2", "port":"switch", "flow_type":"netflow"}) is self-provisioned via a configuration file on port available with 1 agent tags and has status online. [Overwrite default: False. Paste only file: True]
        And pktvisor state is running
        And 1 Agent Group(s) is created with all tags contained in the agent
        And a flow policy flow with tap_selector matching all tag(s) of the tap from an agent, top_ips_ports, by_bytes metric_groups enabled, top_tos, conversations, counters, by_packets, top_geo, cardinality, top_ports, top_interfaces, top_ips metric_groups disabled and settings: default is applied to the group
        And 1 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 1 policies are applied and all has status running
        And run mocked data dns_ipv4_tcp.pcap, dns_ipv4_udp.pcap, dns_ipv6_tcp.pcap, dns_ipv6_udp.pcap, dns_udp_mixed_rcode.pcap, dnssec.pcap, dhcpv6.pcap, dhcp-flow.pcap on the created virtual switch
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 240 seconds
    Then metrics must be correctly generated for flow handler
        And remove the agent .yaml generated on each scenario


@sanity @metric_groups @metrics_flow @root @mocked_interface @auto_provision
Scenario: flow handler type netflow with only top_ips_ports and by_packets metric groups enabled
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink with default configuration type already exists
        And a mocked interface is configured with mtu: 65000 and ip: 192.168.100.2/24
        And a virtual switch is configured and is up with type netflow and target \"192.168.100.2:available\"
    When an agent(backend_type:pktvisor, settings: {"input_type":"flow","bind":"192.168.100.2", "port":"switch", "flow_type":"netflow"}) is self-provisioned via a configuration file on port available with 1 agent tags and has status online. [Overwrite default: False. Paste only file: True]
        And pktvisor state is running
        And 1 Agent Group(s) is created with all tags contained in the agent
        And a flow policy flow with tap_selector matching all tag(s) of the tap from an agent, top_ips_ports, by_packets metric_groups enabled, top_tos, conversations, counters, by_bytes, top_geo, cardinality, top_ports, top_interfaces, top_ips metric_groups disabled and settings: default is applied to the group
        And 1 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 1 policies are applied and all has status running
        And run mocked data dns_ipv4_tcp.pcap, dns_ipv4_udp.pcap, dns_ipv6_tcp.pcap, dns_ipv6_udp.pcap, dns_udp_mixed_rcode.pcap, dnssec.pcap, dhcpv6.pcap, dhcp-flow.pcap on the created virtual switch
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 240 seconds
    Then metrics must be correctly generated for flow handler
        And remove the agent .yaml generated on each scenario


@sanity @metric_groups @metrics_flow @root @mocked_interface @auto_provision
Scenario: flow handler type netflow with only top_interfaces metric groups enabled
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink with default configuration type already exists
        And a mocked interface is configured with mtu: 65000 and ip: 192.168.100.2/24
        And a virtual switch is configured and is up with type netflow and target \"192.168.100.2:available\"
    When an agent(backend_type:pktvisor, settings: {"input_type":"flow","bind":"192.168.100.2", "port":"switch", "flow_type":"netflow"}) is self-provisioned via a configuration file on port available with 1 agent tags and has status online. [Overwrite default: False. Paste only file: True]
        And pktvisor state is running
        And 1 Agent Group(s) is created with all tags contained in the agent
        And a flow policy flow with tap_selector matching all tag(s) of the tap from an agent, top_interfaces metric_groups enabled, top_tos, conversations, counters, by_packets, by_bytes, top_geo, cardinality, top_ports, top_ips_ports, top_ips metric_groups disabled and settings: default is applied to the group
        And 1 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 1 policies are applied and all has status running
        And run mocked data dns_ipv4_tcp.pcap, dns_ipv4_udp.pcap, dns_ipv6_tcp.pcap, dns_ipv6_udp.pcap, dns_udp_mixed_rcode.pcap, dnssec.pcap, dhcpv6.pcap, dhcp-flow.pcap on the created virtual switch
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
    Then metrics must be correctly generated for flow handler
        And remove the agent .yaml generated on each scenario


@sanity @metric_groups @metrics_flow @root @mocked_interface @auto_provision
Scenario: flow handler type netflow with only top_interfaces and by_bytes metric groups enabled
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink with default configuration type already exists
        And a mocked interface is configured with mtu: 65000 and ip: 192.168.100.2/24
        And a virtual switch is configured and is up with type netflow and target \"192.168.100.2:available\"
    When an agent(backend_type:pktvisor, settings: {"input_type":"flow","bind":"192.168.100.2", "port":"switch", "flow_type":"netflow"}) is self-provisioned via a configuration file on port available with 1 agent tags and has status online. [Overwrite default: False. Paste only file: True]
        And pktvisor state is running
        And 1 Agent Group(s) is created with all tags contained in the agent
        And a flow policy flow with tap_selector matching all tag(s) of the tap from an agent, top_interfaces, by_bytes metric_groups enabled, top_tos, conversations, counters, by_packets, top_geo, cardinality, top_ports, top_ips_ports, top_ips metric_groups disabled and settings: default is applied to the group
        And 1 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 1 policies are applied and all has status running
        And run mocked data dns_ipv4_tcp.pcap, dns_ipv4_udp.pcap, dns_ipv6_tcp.pcap, dns_ipv6_udp.pcap, dns_udp_mixed_rcode.pcap, dnssec.pcap, dhcpv6.pcap, dhcp-flow.pcap on the created virtual switch
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 240 seconds
    Then metrics must be correctly generated for flow handler
        And remove the agent .yaml generated on each scenario


@sanity @metric_groups @metrics_flow @root @mocked_interface @auto_provision
Scenario: flow handler type netflow with only top_interfaces and by_packets metric groups enabled
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink with default configuration type already exists
        And a mocked interface is configured with mtu: 65000 and ip: 192.168.100.2/24
        And a virtual switch is configured and is up with type netflow and target \"192.168.100.2:available\"
    When an agent(backend_type:pktvisor, settings: {"input_type":"flow","bind":"192.168.100.2", "port":"switch", "flow_type":"netflow"}) is self-provisioned via a configuration file on port available with 1 agent tags and has status online. [Overwrite default: False. Paste only file: True]
        And pktvisor state is running
        And 1 Agent Group(s) is created with all tags contained in the agent
        And a flow policy flow with tap_selector matching all tag(s) of the tap from an agent, top_interfaces, by_packets metric_groups enabled, top_tos, conversations, counters, by_bytes, top_geo, cardinality, top_ports, top_ips_ports, top_ips metric_groups disabled and settings: default is applied to the group
        And 1 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 1 policies are applied and all has status running
        And run mocked data dns_ipv4_tcp.pcap, dns_ipv4_udp.pcap, dns_ipv6_tcp.pcap, dns_ipv6_udp.pcap, dns_udp_mixed_rcode.pcap, dnssec.pcap, dhcpv6.pcap, dhcp-flow.pcap on the created virtual switch
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 240 seconds
    Then metrics must be correctly generated for flow handler
        And remove the agent .yaml generated on each scenario

@sanity @metric_groups @metrics_flow @root @mocked_interface @auto_provision
Scenario: flow handler type netflow with only top_ips metric groups enabled
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink with default configuration type already exists
        And a mocked interface is configured with mtu: 65000 and ip: 192.168.100.2/24
        And a virtual switch is configured and is up with type netflow and target \"192.168.100.2:available\"
    When an agent(backend_type:pktvisor, settings: {"input_type":"flow","bind":"192.168.100.2", "port":"switch", "flow_type":"netflow"}) is self-provisioned via a configuration file on port available with 1 agent tags and has status online. [Overwrite default: False. Paste only file: True]
        And pktvisor state is running
        And 1 Agent Group(s) is created with all tags contained in the agent
        And a flow policy flow with tap_selector matching all tag(s) of the tap from an agent, top_ips metric_groups enabled, top_tos, conversations, counters, by_packets, by_bytes, top_geo, cardinality, top_ports, top_ips_ports, top_interfaces metric_groups disabled and settings: default is applied to the group
        And 1 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 1 policies are applied and all has status running
        And run mocked data dns_ipv4_tcp.pcap, dns_ipv4_udp.pcap, dns_ipv6_tcp.pcap, dns_ipv6_udp.pcap, dns_udp_mixed_rcode.pcap, dnssec.pcap, dhcpv6.pcap, dhcp-flow.pcap on the created virtual switch
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
    Then metrics must be correctly generated for flow handler
        And remove the agent .yaml generated on each scenario


@sanity @metric_groups @metrics_flow @root @mocked_interface @auto_provision
Scenario: flow handler type netflow with only top_ips and by_packets metric groups enabled
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink with default configuration type already exists
        And a mocked interface is configured with mtu: 65000 and ip: 192.168.100.2/24
        And a virtual switch is configured and is up with type netflow and target \"192.168.100.2:available\"
    When an agent(backend_type:pktvisor, settings: {"input_type":"flow","bind":"192.168.100.2", "port":"switch", "flow_type":"netflow"}) is self-provisioned via a configuration file on port available with 1 agent tags and has status online. [Overwrite default: False. Paste only file: True]
        And pktvisor state is running
        And 1 Agent Group(s) is created with all tags contained in the agent
        And a flow policy flow with tap_selector matching all tag(s) of the tap from an agent, top_ips, by_packets metric_groups enabled, top_tos, conversations, counters, by_bytes, top_geo, cardinality, top_ports, top_ips_ports, top_interfaces metric_groups disabled and settings: default is applied to the group
        And 1 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 1 policies are applied and all has status running
        And run mocked data dns_ipv4_tcp.pcap, dns_ipv4_udp.pcap, dns_ipv6_tcp.pcap, dns_ipv6_udp.pcap, dns_udp_mixed_rcode.pcap, dnssec.pcap, dhcpv6.pcap, dhcp-flow.pcap on the created virtual switch
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 240 seconds
    Then metrics must be correctly generated for flow handler
        And remove the agent .yaml generated on each scenario


@sanity @metric_groups @metrics_flow @root @mocked_interface @auto_provision
Scenario: flow handler type netflow with only top_ips and by_bytes metric groups enabled
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink with default configuration type already exists
        And a mocked interface is configured with mtu: 65000 and ip: 192.168.100.2/24
        And a virtual switch is configured and is up with type netflow and target \"192.168.100.2:available\"
    When an agent(backend_type:pktvisor, settings: {"input_type":"flow","bind":"192.168.100.2", "port":"switch", "flow_type":"netflow"}) is self-provisioned via a configuration file on port available with 1 agent tags and has status online. [Overwrite default: False. Paste only file: True]
        And pktvisor state is running
        And 1 Agent Group(s) is created with all tags contained in the agent
        And a flow policy flow with tap_selector matching all tag(s) of the tap from an agent, top_ips, by_bytes metric_groups enabled, top_tos, conversations, counters, by_packets, top_geo, cardinality, top_ports, top_ips_ports, top_interfaces metric_groups disabled and settings: default is applied to the group
        And 1 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 1 policies are applied and all has status running
        And run mocked data dns_ipv4_tcp.pcap, dns_ipv4_udp.pcap, dns_ipv6_tcp.pcap, dns_ipv6_udp.pcap, dns_udp_mixed_rcode.pcap, dnssec.pcap, dhcpv6.pcap, dhcp-flow.pcap on the created virtual switch
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 240 seconds
    Then metrics must be correctly generated for flow handler
        And remove the agent .yaml generated on each scenario


@sanity @metric_groups @metrics_flow @root @mocked_interface @auto_provision
Scenario: flow handler type netflow with only top_tos metric groups enabled
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink with default configuration type already exists
        And a mocked interface is configured with mtu: 65000 and ip: 192.168.100.2/24
        And a virtual switch is configured and is up with type netflow and target \"192.168.100.2:available\"
    When an agent(backend_type:pktvisor, settings: {"input_type":"flow","bind":"192.168.100.2", "port":"switch", "flow_type":"netflow"}) is self-provisioned via a configuration file on port available with 1 agent tags and has status online. [Overwrite default: False. Paste only file: True]
        And pktvisor state is running
        And 1 Agent Group(s) is created with all tags contained in the agent
        And a flow policy flow with tap_selector matching all tag(s) of the tap from an agent, top_tos metric_groups enabled, top_ips, by_bytes, conversations, counters, by_packets, top_geo, cardinality, top_ports, top_ips_ports, top_interfaces metric_groups disabled and settings: default is applied to the group
        And 1 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 1 policies are applied and all has status running
        And run mocked data dns_ipv4_tcp.pcap, dns_ipv4_udp.pcap, dns_ipv6_tcp.pcap, dns_ipv6_udp.pcap, dns_udp_mixed_rcode.pcap, dnssec.pcap, dhcpv6.pcap, dhcp-flow.pcap on the created virtual switch
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
    Then metrics must be correctly generated for flow handler
        And remove the agent .yaml generated on each scenario


@sanity @metric_groups @metrics_flow @root @mocked_interface @auto_provision
Scenario: flow handler type netflow with only top_tos and by_bytes metric groups enabled
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink with default configuration type already exists
        And a mocked interface is configured with mtu: 65000 and ip: 192.168.100.2/24
        And a virtual switch is configured and is up with type netflow and target \"192.168.100.2:available\"
    When an agent(backend_type:pktvisor, settings: {"input_type":"flow","bind":"192.168.100.2", "port":"switch", "flow_type":"netflow"}) is self-provisioned via a configuration file on port available with 1 agent tags and has status online. [Overwrite default: False. Paste only file: True]
        And pktvisor state is running
        And 1 Agent Group(s) is created with all tags contained in the agent
        And a flow policy flow with tap_selector matching all tag(s) of the tap from an agent, top_tos, by_bytes metric_groups enabled, top_ips, conversations, counters, by_packets, top_geo, cardinality, top_ports, top_ips_ports, top_interfaces metric_groups disabled and settings: default is applied to the group
        And 1 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 1 policies are applied and all has status running
        And run mocked data dns_ipv4_tcp.pcap, dns_ipv4_udp.pcap, dns_ipv6_tcp.pcap, dns_ipv6_udp.pcap, dns_udp_mixed_rcode.pcap, dnssec.pcap, dhcpv6.pcap, dhcp-flow.pcap on the created virtual switch
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 240 seconds
    Then metrics must be correctly generated for flow handler
        And remove the agent .yaml generated on each scenario


@sanity @metric_groups @metrics_flow @root @mocked_interface @auto_provision
Scenario: flow handler type netflow with only top_tos and by_packets metric groups enabled
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink with default configuration type already exists
        And a mocked interface is configured with mtu: 65000 and ip: 192.168.100.2/24
        And a virtual switch is configured and is up with type netflow and target \"192.168.100.2:available\"
    When an agent(backend_type:pktvisor, settings: {"input_type":"flow","bind":"192.168.100.2", "port":"switch", "flow_type":"netflow"}) is self-provisioned via a configuration file on port available with 1 agent tags and has status online. [Overwrite default: False. Paste only file: True]
        And pktvisor state is running
        And 1 Agent Group(s) is created with all tags contained in the agent
        And a flow policy flow with tap_selector matching all tag(s) of the tap from an agent, top_tos, by_packets metric_groups enabled, top_ips, by_bytes, conversations, counters, top_geo, cardinality, top_ports, top_ips_ports, top_interfaces metric_groups disabled and settings: default is applied to the group
        And 1 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 1 policies are applied and all has status running
        And run mocked data dns_ipv4_tcp.pcap, dns_ipv4_udp.pcap, dns_ipv6_tcp.pcap, dns_ipv6_udp.pcap, dns_udp_mixed_rcode.pcap, dnssec.pcap, dhcpv6.pcap, dhcp-flow.pcap on the created virtual switch
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 240 seconds
    Then metrics must be correctly generated for flow handler
        And remove the agent .yaml generated on each scenario


#### flow sflow

@sanity @metric_groups @metrics_flow @root @mocked_interface @auto_provision
Scenario: flow handler type sflow with default metric groups configuration
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink with default configuration type already exists
        And a mocked interface is configured with mtu: 65000 and ip: 192.168.100.2/24
        And a virtual switch is configured and is up with type sflow and target \"192.168.100.2:available\"
    When an agent(backend_type:pktvisor, settings: {"input_type":"flow","bind":"192.168.100.2", "port":"switch", "flow_type":"sflow"}) is self-provisioned via a configuration file on port available with 1 agent tags and has status online. [Overwrite default: False. Paste only file: True]
        And pktvisor state is running
        And 1 Agent Group(s) is created with all tags contained in the agent
        And a flow policy flow with tap_selector matching all tag(s) of the tap from an agent, default metric_groups enabled, default metric_groups disabled and settings: default is applied to the group
        And 1 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 1 policies are applied and all has status running
        And run mocked data dns_ipv4_tcp.pcap, dns_ipv4_udp.pcap, dns_ipv6_tcp.pcap, dns_ipv6_udp.pcap, dns_udp_mixed_rcode.pcap, dnssec.pcap, dhcpv6.pcap, dhcp-flow.pcap on the created virtual switch
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 240 seconds
    Then metrics must be correctly generated for flow handler
        And remove the agent .yaml generated on each scenario


@sanity @metric_groups @metrics_flow @root @mocked_interface @auto_provision
Scenario: flow handler type sflow with all metric groups enabled
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink with default configuration type already exists
        And a mocked interface is configured with mtu: 65000 and ip: 192.168.100.2/24
        And a virtual switch is configured and is up with type sflow and target \"192.168.100.2:available\"
    When an agent(backend_type:pktvisor, settings: {"input_type":"flow","bind":"192.168.100.2", "port":"switch", "flow_type":"sflow"}) is self-provisioned via a configuration file on port available with 1 agent tags and has status online. [Overwrite default: False. Paste only file: True]
        And pktvisor state is running
        And 1 Agent Group(s) is created with all tags contained in the agent
        And a flow policy flow with tap_selector matching all tag(s) of the tap from an agent, all metric_groups enabled, none metric_groups disabled and settings: default is applied to the group
        And 1 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 1 policies are applied and all has status running
        And run mocked data dns_ipv4_tcp.pcap, dns_ipv4_udp.pcap, dns_ipv6_tcp.pcap, dns_ipv6_udp.pcap, dns_udp_mixed_rcode.pcap, dnssec.pcap, dhcpv6.pcap, dhcp-flow.pcap on the created virtual switch
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 240 seconds
    Then metrics must be correctly generated for flow handler
        And remove the agent .yaml generated on each scenario


@sanity @metric_groups @metrics_flow @root @mocked_interface @auto_provision
Scenario: flow handler type sflow with all metric groups disabled
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink with default configuration type already exists
        And a mocked interface is configured with mtu: 65000 and ip: 192.168.100.2/24
        And a virtual switch is configured and is up with type sflow and target \"192.168.100.2:available\"
    When an agent(backend_type:pktvisor, settings: {"input_type":"flow","bind":"192.168.100.2", "port":"switch", "flow_type":"sflow"}) is self-provisioned via a configuration file on port available with 1 agent tags and has status online. [Overwrite default: False. Paste only file: True]
        And pktvisor state is running
        And 1 Agent Group(s) is created with all tags contained in the agent
        And a flow policy flow with tap_selector matching all tag(s) of the tap from an agent, none metric_groups enabled, all metric_groups disabled and settings: default is applied to the group
        And 1 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 1 policies are applied and all has status running
        And run mocked data dns_ipv4_tcp.pcap, dns_ipv4_udp.pcap, dns_ipv6_tcp.pcap, dns_ipv6_udp.pcap, dns_udp_mixed_rcode.pcap, dnssec.pcap, dhcpv6.pcap, dhcp-flow.pcap on the created virtual switch
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
    Then metrics must be correctly generated for flow handler
        And remove the agent .yaml generated on each scenario


@sanity @metric_groups @metrics_flow @root @mocked_interface @auto_provision
Scenario: flow handler type sflow with only cardinality metric groups enabled
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink with default configuration type already exists
        And a mocked interface is configured with mtu: 65000 and ip: 192.168.100.2/24
        And a virtual switch is configured and is up with type sflow and target \"192.168.100.2:available\"
    When an agent(backend_type:pktvisor, settings: {"input_type":"flow","bind":"192.168.100.2", "port":"switch", "flow_type":"sflow"}) is self-provisioned via a configuration file on port available with 1 agent tags and has status online. [Overwrite default: False. Paste only file: True]
        And pktvisor state is running
        And 1 Agent Group(s) is created with all tags contained in the agent
        And a flow policy flow with tap_selector matching all tag(s) of the tap from an agent, cardinality metric_groups enabled, top_tos, counters, by_packets, by_bytes, top_geo, conversations, top_ports, top_ips_ports, top_interfaces, top_ips metric_groups disabled and settings: default is applied to the group
        And 1 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 1 policies are applied and all has status running
        And run mocked data dns_ipv4_tcp.pcap, dns_ipv4_udp.pcap, dns_ipv6_tcp.pcap, dns_ipv6_udp.pcap, dns_udp_mixed_rcode.pcap, dnssec.pcap, dhcpv6.pcap, dhcp-flow.pcap on the created virtual switch
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 240 seconds
    Then metrics must be correctly generated for flow handler
        And remove the agent .yaml generated on each scenario


@sanity @metric_groups @metrics_flow @root @mocked_interface @auto_provision
Scenario: flow handler type sflow with only counters metric groups enabled
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink with default configuration type already exists
        And a mocked interface is configured with mtu: 65000 and ip: 192.168.100.2/24
        And a virtual switch is configured and is up with type sflow and target \"192.168.100.2:available\"
    When an agent(backend_type:pktvisor, settings: {"input_type":"flow","bind":"192.168.100.2", "port":"switch", "flow_type":"sflow"}) is self-provisioned via a configuration file on port available with 1 agent tags and has status online. [Overwrite default: False. Paste only file: True]
        And pktvisor state is running
        And 1 Agent Group(s) is created with all tags contained in the agent
        And a flow policy flow with tap_selector matching all tag(s) of the tap from an agent, counters metric_groups enabled, top_tos, cardinality, by_packets, by_bytes, top_geo, conversations, top_ports, top_ips_ports, top_interfaces, top_ips metric_groups disabled and settings: default is applied to the group
        And 1 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 1 policies are applied and all has status running
        And run mocked data dns_ipv4_tcp.pcap, dns_ipv4_udp.pcap, dns_ipv6_tcp.pcap, dns_ipv6_udp.pcap, dns_udp_mixed_rcode.pcap, dnssec.pcap, dhcpv6.pcap, dhcp-flow.pcap on the created virtual switch
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 240 seconds
    Then metrics must be correctly generated for flow handler
        And remove the agent .yaml generated on each scenario


@sanity @metric_groups @metrics_flow @root @mocked_interface @auto_provision
Scenario: flow handler type sflow with only counters and by_bytes metric groups enabled
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink with default configuration type already exists
        And a mocked interface is configured with mtu: 65000 and ip: 192.168.100.2/24
        And a virtual switch is configured and is up with type sflow and target \"192.168.100.2:available\"
    When an agent(backend_type:pktvisor, settings: {"input_type":"flow","bind":"192.168.100.2", "port":"switch", "flow_type":"sflow"}) is self-provisioned via a configuration file on port available with 1 agent tags and has status online. [Overwrite default: False. Paste only file: True]
        And pktvisor state is running
        And 1 Agent Group(s) is created with all tags contained in the agent
        And a flow policy flow with tap_selector matching all tag(s) of the tap from an agent, counters, by_bytes metric_groups enabled, top_tos, cardinality, by_packets, top_geo, conversations, top_ports, top_ips_ports, top_interfaces, top_ips metric_groups disabled and settings: default is applied to the group
        And 1 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 1 policies are applied and all has status running
        And run mocked data dns_ipv4_tcp.pcap, dns_ipv4_udp.pcap, dns_ipv6_tcp.pcap, dns_ipv6_udp.pcap, dns_udp_mixed_rcode.pcap, dnssec.pcap, dhcpv6.pcap, dhcp-flow.pcap on the created virtual switch
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 240 seconds
    Then metrics must be correctly generated for flow handler
        And remove the agent .yaml generated on each scenario


@sanity @metric_groups @metrics_flow @root @mocked_interface @auto_provision
Scenario: flow handler type sflow with only counters and by_packets metric groups enabled
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink with default configuration type already exists
        And a mocked interface is configured with mtu: 65000 and ip: 192.168.100.2/24
        And a virtual switch is configured and is up with type sflow and target \"192.168.100.2:available\"
    When an agent(backend_type:pktvisor, settings: {"input_type":"flow","bind":"192.168.100.2", "port":"switch", "flow_type":"sflow"}) is self-provisioned via a configuration file on port available with 1 agent tags and has status online. [Overwrite default: False. Paste only file: True]
        And pktvisor state is running
        And 1 Agent Group(s) is created with all tags contained in the agent
        And a flow policy flow with tap_selector matching all tag(s) of the tap from an agent, counters, by_packets metric_groups enabled, top_tos, cardinality, by_bytes, top_geo, conversations, top_ports, top_ips_ports, top_interfaces, top_ips metric_groups disabled and settings: default is applied to the group
        And 1 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 1 policies are applied and all has status running
        And run mocked data dns_ipv4_tcp.pcap, dns_ipv4_udp.pcap, dns_ipv6_tcp.pcap, dns_ipv6_udp.pcap, dns_udp_mixed_rcode.pcap, dnssec.pcap, dhcpv6.pcap, dhcp-flow.pcap on the created virtual switch
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 240 seconds
    Then metrics must be correctly generated for flow handler
        And remove the agent .yaml generated on each scenario


@sanity @metric_groups @metrics_flow @root @mocked_interface @auto_provision
Scenario: flow handler type sflow with only by_packets metric groups enabled
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink with default configuration type already exists
        And a mocked interface is configured with mtu: 65000 and ip: 192.168.100.2/24
        And a virtual switch is configured and is up with type sflow and target \"192.168.100.2:available\"
    When an agent(backend_type:pktvisor, settings: {"input_type":"flow","bind":"192.168.100.2", "port":"switch", "flow_type":"sflow"}) is self-provisioned via a configuration file on port available with 1 agent tags and has status online. [Overwrite default: False. Paste only file: True]
        And pktvisor state is running
        And 1 Agent Group(s) is created with all tags contained in the agent
        And a flow policy flow with tap_selector matching all tag(s) of the tap from an agent, by_packets metric_groups enabled, top_tos, counters, cardinality, by_bytes, top_geo, conversations, top_ports, top_ips_ports, top_interfaces, top_ips metric_groups disabled and settings: default is applied to the group
        And 1 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 1 policies are applied and all has status running
        And run mocked data dns_ipv4_tcp.pcap, dns_ipv4_udp.pcap, dns_ipv6_tcp.pcap, dns_ipv6_udp.pcap, dns_udp_mixed_rcode.pcap, dnssec.pcap, dhcpv6.pcap, dhcp-flow.pcap on the created virtual switch
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
    Then metrics must be correctly generated for flow handler
        And remove the agent .yaml generated on each scenario


@sanity @metric_groups @metrics_flow @root @mocked_interface @auto_provision
Scenario: flow handler type sflow with only by_bytes metric groups enabled
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink with default configuration type already exists
        And a mocked interface is configured with mtu: 65000 and ip: 192.168.100.2/24
        And a virtual switch is configured and is up with type sflow and target \"192.168.100.2:available\"
    When an agent(backend_type:pktvisor, settings: {"input_type":"flow","bind":"192.168.100.2", "port":"switch", "flow_type":"sflow"}) is self-provisioned via a configuration file on port available with 1 agent tags and has status online. [Overwrite default: False. Paste only file: True]
        And pktvisor state is running
        And 1 Agent Group(s) is created with all tags contained in the agent
        And a flow policy flow with tap_selector matching all tag(s) of the tap from an agent, by_bytes metric_groups enabled, top_tos, counters, by_packets, cardinality, top_geo, conversations, top_ports, top_ips_ports, top_interfaces, top_ips metric_groups disabled and settings: default is applied to the group
        And 1 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 1 policies are applied and all has status running
        And run mocked data dns_ipv4_tcp.pcap, dns_ipv4_udp.pcap, dns_ipv6_tcp.pcap, dns_ipv6_udp.pcap, dns_udp_mixed_rcode.pcap, dnssec.pcap, dhcpv6.pcap, dhcp-flow.pcap on the created virtual switch
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
    Then metrics must be correctly generated for flow handler
        And remove the agent .yaml generated on each scenario


@sanity @metric_groups @metrics_flow @root @mocked_interface @auto_provision
Scenario: flow handler type sflow with only top_geo metric groups enabled
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink with default configuration type already exists
        And a mocked interface is configured with mtu: 65000 and ip: 192.168.100.2/24
        And a virtual switch is configured and is up with type sflow and target \"192.168.100.2:available\"
    When an agent(backend_type:pktvisor, settings: {"input_type":"flow","bind":"192.168.100.2", "port":"switch", "flow_type":"sflow"}) is self-provisioned via a configuration file on port available with 1 agent tags and has status online. [Overwrite default: False. Paste only file: True]
        And pktvisor state is running
        And 1 Agent Group(s) is created with all tags contained in the agent
        And a flow policy flow with tap_selector matching all tag(s) of the tap from an agent, top_geo metric_groups enabled, top_tos, counters, by_packets, by_bytes, cardinality, conversations, top_ports, top_ips_ports, top_interfaces, top_ips metric_groups disabled and settings: default is applied to the group
        And 1 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 1 policies are applied and all has status running
        And run mocked data dns_ipv4_tcp.pcap, dns_ipv4_udp.pcap, dns_ipv6_tcp.pcap, dns_ipv6_udp.pcap, dns_udp_mixed_rcode.pcap, dnssec.pcap, dhcpv6.pcap, dhcp-flow.pcap on the created virtual switch
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
    Then metrics must be correctly generated for flow handler
        And remove the agent .yaml generated on each scenario

@sanity @metric_groups @metrics_flow @root @mocked_interface @auto_provision
Scenario: flow handler type sflow with only top_geo and by_bytes metric groups enabled
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink with default configuration type already exists
        And a mocked interface is configured with mtu: 65000 and ip: 192.168.100.2/24
        And a virtual switch is configured and is up with type sflow and target \"192.168.100.2:available\"
    When an agent(backend_type:pktvisor, settings: {"input_type":"flow","bind":"192.168.100.2", "port":"switch", "flow_type":"sflow"}) is self-provisioned via a configuration file on port available with 1 agent tags and has status online. [Overwrite default: False. Paste only file: True]
        And pktvisor state is running
        And 1 Agent Group(s) is created with all tags contained in the agent
        And a flow policy flow with tap_selector matching all tag(s) of the tap from an agent, top_geo, by_bytes metric_groups enabled, top_tos, counters, by_packets, cardinality, conversations, top_ports, top_ips_ports, top_interfaces, top_ips metric_groups disabled and settings: default is applied to the group
        And 1 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 1 policies are applied and all has status running
        And run mocked data dns_ipv4_tcp.pcap, dns_ipv4_udp.pcap, dns_ipv6_tcp.pcap, dns_ipv6_udp.pcap, dns_udp_mixed_rcode.pcap, dnssec.pcap, dhcpv6.pcap, dhcp-flow.pcap on the created virtual switch
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 240 seconds
    Then metrics must be correctly generated for flow handler
        And remove the agent .yaml generated on each scenario


@sanity @metric_groups @metrics_flow @root @mocked_interface @auto_provision
Scenario: flow handler type sflow with only top_geo and by_packets metric groups enabled
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink with default configuration type already exists
        And a mocked interface is configured with mtu: 65000 and ip: 192.168.100.2/24
        And a virtual switch is configured and is up with type sflow and target \"192.168.100.2:available\"
    When an agent(backend_type:pktvisor, settings: {"input_type":"flow","bind":"192.168.100.2", "port":"switch", "flow_type":"sflow"}) is self-provisioned via a configuration file on port available with 1 agent tags and has status online. [Overwrite default: False. Paste only file: True]
        And pktvisor state is running
        And 1 Agent Group(s) is created with all tags contained in the agent
        And a flow policy flow with tap_selector matching all tag(s) of the tap from an agent, top_geo, by_packets metric_groups enabled, top_tos, counters, by_bytes, cardinality, conversations, top_ports, top_ips_ports, top_interfaces, top_ips metric_groups disabled and settings: default is applied to the group
        And 1 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 1 policies are applied and all has status running
        And run mocked data dns_ipv4_tcp.pcap, dns_ipv4_udp.pcap, dns_ipv6_tcp.pcap, dns_ipv6_udp.pcap, dns_udp_mixed_rcode.pcap, dnssec.pcap, dhcpv6.pcap, dhcp-flow.pcap on the created virtual switch
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 240 seconds
    Then metrics must be correctly generated for flow handler
        And remove the agent .yaml generated on each scenario


@sanity @metric_groups @metrics_flow @root @mocked_interface @auto_provision
Scenario: flow handler type sflow with only conversations metric groups enabled
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink with default configuration type already exists
        And a mocked interface is configured with mtu: 65000 and ip: 192.168.100.2/24
        And a virtual switch is configured and is up with type sflow and target \"192.168.100.2:available\"
    When an agent(backend_type:pktvisor, settings: {"input_type":"flow","bind":"192.168.100.2", "port":"switch", "flow_type":"sflow"}) is self-provisioned via a configuration file on port available with 1 agent tags and has status online. [Overwrite default: False. Paste only file: True]
        And pktvisor state is running
        And 1 Agent Group(s) is created with all tags contained in the agent
        And a flow policy flow with tap_selector matching all tag(s) of the tap from an agent, conversations metric_groups enabled, top_tos, counters, by_packets, by_bytes, top_geo, cardinality, top_ports, top_ips_ports, top_interfaces, top_ips metric_groups disabled and settings: default is applied to the group
        And 1 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 1 policies are applied and all has status running
        And run mocked data dns_ipv4_tcp.pcap, dns_ipv4_udp.pcap, dns_ipv6_tcp.pcap, dns_ipv6_udp.pcap, dns_udp_mixed_rcode.pcap, dnssec.pcap, dhcpv6.pcap, dhcp-flow.pcap on the created virtual switch
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
    Then metrics must be correctly generated for flow handler
        And remove the agent .yaml generated on each scenario


@sanity @metric_groups @metrics_flow @root @mocked_interface @auto_provision
Scenario: flow handler type sflow with only conversations and cardinality metric groups enabled
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink with default configuration type already exists
        And a mocked interface is configured with mtu: 65000 and ip: 192.168.100.2/24
        And a virtual switch is configured and is up with type sflow and target \"192.168.100.2:available\"
    When an agent(backend_type:pktvisor, settings: {"input_type":"flow","bind":"192.168.100.2", "port":"switch", "flow_type":"sflow"}) is self-provisioned via a configuration file on port available with 1 agent tags and has status online. [Overwrite default: False. Paste only file: True]
        And pktvisor state is running
        And 1 Agent Group(s) is created with all tags contained in the agent
        And a flow policy flow with tap_selector matching all tag(s) of the tap from an agent, conversations, cardinality metric_groups enabled, top_tos, counters, by_packets, by_bytes, top_geo, top_ports, top_ips_ports, top_interfaces, top_ips metric_groups disabled and settings: default is applied to the group
        And 1 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 1 policies are applied and all has status running
        And run mocked data dns_ipv4_tcp.pcap, dns_ipv4_udp.pcap, dns_ipv6_tcp.pcap, dns_ipv6_udp.pcap, dns_udp_mixed_rcode.pcap, dnssec.pcap, dhcpv6.pcap, dhcp-flow.pcap on the created virtual switch
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
    Then metrics must be correctly generated for flow handler
        And remove the agent .yaml generated on each scenario


@sanity @metric_groups @metrics_flow @root @mocked_interface @auto_provision
Scenario: flow handler type sflow with only conversations and by_bytes metric groups enabled
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink with default configuration type already exists
        And a mocked interface is configured with mtu: 65000 and ip: 192.168.100.2/24
        And a virtual switch is configured and is up with type sflow and target \"192.168.100.2:available\"
    When an agent(backend_type:pktvisor, settings: {"input_type":"flow","bind":"192.168.100.2", "port":"switch", "flow_type":"sflow"}) is self-provisioned via a configuration file on port available with 1 agent tags and has status online. [Overwrite default: False. Paste only file: True]
        And pktvisor state is running
        And 1 Agent Group(s) is created with all tags contained in the agent
        And a flow policy flow with tap_selector matching all tag(s) of the tap from an agent, conversations, by_bytes metric_groups enabled, top_tos, counters, by_packets, top_geo, cardinality, top_ports, top_ips_ports, top_interfaces, top_ips metric_groups disabled and settings: default is applied to the group
        And 1 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 1 policies are applied and all has status running
        And run mocked data dns_ipv4_tcp.pcap, dns_ipv4_udp.pcap, dns_ipv6_tcp.pcap, dns_ipv6_udp.pcap, dns_udp_mixed_rcode.pcap, dnssec.pcap, dhcpv6.pcap, dhcp-flow.pcap on the created virtual switch
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 240 seconds
    Then metrics must be correctly generated for flow handler
        And remove the agent .yaml generated on each scenario


@sanity @metric_groups @metrics_flow @root @mocked_interface @auto_provision
Scenario: flow handler type sflow with only conversations and by_packets metric groups enabled
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink with default configuration type already exists
        And a mocked interface is configured with mtu: 65000 and ip: 192.168.100.2/24
        And a virtual switch is configured and is up with type sflow and target \"192.168.100.2:available\"
    When an agent(backend_type:pktvisor, settings: {"input_type":"flow","bind":"192.168.100.2", "port":"switch", "flow_type":"sflow"}) is self-provisioned via a configuration file on port available with 1 agent tags and has status online. [Overwrite default: False. Paste only file: True]
        And pktvisor state is running
        And 1 Agent Group(s) is created with all tags contained in the agent
        And a flow policy flow with tap_selector matching all tag(s) of the tap from an agent, conversations, by_packets metric_groups enabled, top_tos, counters, by_bytes, top_geo, cardinality, top_ports, top_ips_ports, top_interfaces, top_ips metric_groups disabled and settings: default is applied to the group
        And 1 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 1 policies are applied and all has status running
        And run mocked data dns_ipv4_tcp.pcap, dns_ipv4_udp.pcap, dns_ipv6_tcp.pcap, dns_ipv6_udp.pcap, dns_udp_mixed_rcode.pcap, dnssec.pcap, dhcpv6.pcap, dhcp-flow.pcap on the created virtual switch
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 240 seconds
    Then metrics must be correctly generated for flow handler
        And remove the agent .yaml generated on each scenario


@sanity @metric_groups @metrics_flow @root @mocked_interface @auto_provision
Scenario: flow handler type sflow with only top_ports metric groups enabled
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink with default configuration type already exists
        And a mocked interface is configured with mtu: 65000 and ip: 192.168.100.2/24
        And a virtual switch is configured and is up with type sflow and target \"192.168.100.2:available\"
    When an agent(backend_type:pktvisor, settings: {"input_type":"flow","bind":"192.168.100.2", "port":"switch", "flow_type":"sflow"}) is self-provisioned via a configuration file on port available with 1 agent tags and has status online. [Overwrite default: False. Paste only file: True]
        And pktvisor state is running
        And 1 Agent Group(s) is created with all tags contained in the agent
        And a flow policy flow with tap_selector matching all tag(s) of the tap from an agent, top_ports metric_groups enabled, top_tos, conversations, counters, by_packets, by_bytes, top_geo, cardinality, top_ips_ports, top_interfaces, top_ips metric_groups disabled and settings: default is applied to the group
        And 1 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 1 policies are applied and all has status running
        And run mocked data dns_ipv4_tcp.pcap, dns_ipv4_udp.pcap, dns_ipv6_tcp.pcap, dns_ipv6_udp.pcap, dns_udp_mixed_rcode.pcap, dnssec.pcap, dhcpv6.pcap, dhcp-flow.pcap on the created virtual switch
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
    Then metrics must be correctly generated for flow handler
        And remove the agent .yaml generated on each scenario

@sanity @metric_groups @metrics_flow @root @mocked_interface @auto_provision
Scenario: flow handler type sflow with only top_ports and by_bytes metric groups enabled
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink with default configuration type already exists
        And a mocked interface is configured with mtu: 65000 and ip: 192.168.100.2/24
        And a virtual switch is configured and is up with type sflow and target \"192.168.100.2:available\"
    When an agent(backend_type:pktvisor, settings: {"input_type":"flow","bind":"192.168.100.2", "port":"switch", "flow_type":"sflow"}) is self-provisioned via a configuration file on port available with 1 agent tags and has status online. [Overwrite default: False. Paste only file: True]
        And pktvisor state is running
        And 1 Agent Group(s) is created with all tags contained in the agent
        And a flow policy flow with tap_selector matching all tag(s) of the tap from an agent, top_ports, by_bytes metric_groups enabled, top_tos, conversations, counters, by_packets, top_geo, cardinality, top_ips_ports, top_interfaces, top_ips metric_groups disabled and settings: default is applied to the group
        And 1 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 1 policies are applied and all has status running
        And run mocked data dns_ipv4_tcp.pcap, dns_ipv4_udp.pcap, dns_ipv6_tcp.pcap, dns_ipv6_udp.pcap, dns_udp_mixed_rcode.pcap, dnssec.pcap, dhcpv6.pcap, dhcp-flow.pcap on the created virtual switch
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 240 seconds
    Then metrics must be correctly generated for flow handler
        And remove the agent .yaml generated on each scenario


@sanity @metric_groups @metrics_flow @root @mocked_interface @auto_provision
Scenario: flow handler type sflow with only top_ports and by_packets metric groups enabled
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink with default configuration type already exists
        And a mocked interface is configured with mtu: 65000 and ip: 192.168.100.2/24
        And a virtual switch is configured and is up with type sflow and target \"192.168.100.2:available\"
    When an agent(backend_type:pktvisor, settings: {"input_type":"flow","bind":"192.168.100.2", "port":"switch", "flow_type":"sflow"}) is self-provisioned via a configuration file on port available with 1 agent tags and has status online. [Overwrite default: False. Paste only file: True]
        And pktvisor state is running
        And 1 Agent Group(s) is created with all tags contained in the agent
        And a flow policy flow with tap_selector matching all tag(s) of the tap from an agent, top_ports, by_packets metric_groups enabled, top_tos, conversations, counters, by_bytes, top_geo, cardinality, top_ips_ports, top_interfaces, top_ips metric_groups disabled and settings: default is applied to the group
        And 1 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 1 policies are applied and all has status running
        And run mocked data dns_ipv4_tcp.pcap, dns_ipv4_udp.pcap, dns_ipv6_tcp.pcap, dns_ipv6_udp.pcap, dns_udp_mixed_rcode.pcap, dnssec.pcap, dhcpv6.pcap, dhcp-flow.pcap on the created virtual switch
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 240 seconds
    Then metrics must be correctly generated for flow handler
        And remove the agent .yaml generated on each scenario


@sanity @metric_groups @metrics_flow @root @mocked_interface @auto_provision
Scenario: flow handler type sflow with only top_ips_ports metric groups enabled
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink with default configuration type already exists
        And a mocked interface is configured with mtu: 65000 and ip: 192.168.100.2/24
        And a virtual switch is configured and is up with type sflow and target \"192.168.100.2:available\"
    When an agent(backend_type:pktvisor, settings: {"input_type":"flow","bind":"192.168.100.2", "port":"switch", "flow_type":"sflow"}) is self-provisioned via a configuration file on port available with 1 agent tags and has status online. [Overwrite default: False. Paste only file: True]
        And pktvisor state is running
        And 1 Agent Group(s) is created with all tags contained in the agent
        And a flow policy flow with tap_selector matching all tag(s) of the tap from an agent, top_ips_ports metric_groups enabled, top_tos, conversations, counters, by_packets, by_bytes, top_geo, cardinality, top_ports, top_interfaces, top_ips metric_groups disabled and settings: default is applied to the group
        And 1 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 1 policies are applied and all has status running
        And run mocked data dns_ipv4_tcp.pcap, dns_ipv4_udp.pcap, dns_ipv6_tcp.pcap, dns_ipv6_udp.pcap, dns_udp_mixed_rcode.pcap, dnssec.pcap, dhcpv6.pcap, dhcp-flow.pcap on the created virtual switch
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
    Then metrics must be correctly generated for flow handler
        And remove the agent .yaml generated on each scenario


@sanity @metric_groups @metrics_flow @root @mocked_interface @auto_provision
Scenario: flow handler type sflow with only top_ips_ports and by_bytes metric groups enabled
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink with default configuration type already exists
        And a mocked interface is configured with mtu: 65000 and ip: 192.168.100.2/24
        And a virtual switch is configured and is up with type sflow and target \"192.168.100.2:available\"
    When an agent(backend_type:pktvisor, settings: {"input_type":"flow","bind":"192.168.100.2", "port":"switch", "flow_type":"sflow"}) is self-provisioned via a configuration file on port available with 1 agent tags and has status online. [Overwrite default: False. Paste only file: True]
        And pktvisor state is running
        And 1 Agent Group(s) is created with all tags contained in the agent
        And a flow policy flow with tap_selector matching all tag(s) of the tap from an agent, top_ips_ports, by_bytes metric_groups enabled, top_tos, conversations, counters, by_packets, top_geo, cardinality, top_ports, top_interfaces, top_ips metric_groups disabled and settings: default is applied to the group
        And 1 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 1 policies are applied and all has status running
        And run mocked data dns_ipv4_tcp.pcap, dns_ipv4_udp.pcap, dns_ipv6_tcp.pcap, dns_ipv6_udp.pcap, dns_udp_mixed_rcode.pcap, dnssec.pcap, dhcpv6.pcap, dhcp-flow.pcap on the created virtual switch
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 240 seconds
    Then metrics must be correctly generated for flow handler
        And remove the agent .yaml generated on each scenario


@sanity @metric_groups @metrics_flow @root @mocked_interface @auto_provision
Scenario: flow handler type sflow with only top_ips_ports and by_packets metric groups enabled
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink with default configuration type already exists
        And a mocked interface is configured with mtu: 65000 and ip: 192.168.100.2/24
        And a virtual switch is configured and is up with type sflow and target \"192.168.100.2:available\"
    When an agent(backend_type:pktvisor, settings: {"input_type":"flow","bind":"192.168.100.2", "port":"switch", "flow_type":"sflow"}) is self-provisioned via a configuration file on port available with 1 agent tags and has status online. [Overwrite default: False. Paste only file: True]
        And pktvisor state is running
        And 1 Agent Group(s) is created with all tags contained in the agent
        And a flow policy flow with tap_selector matching all tag(s) of the tap from an agent, top_ips_ports, by_packets metric_groups enabled, top_tos, conversations, counters, by_bytes, top_geo, cardinality, top_ports, top_interfaces, top_ips metric_groups disabled and settings: default is applied to the group
        And 1 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 1 policies are applied and all has status running
        And run mocked data dns_ipv4_tcp.pcap, dns_ipv4_udp.pcap, dns_ipv6_tcp.pcap, dns_ipv6_udp.pcap, dns_udp_mixed_rcode.pcap, dnssec.pcap, dhcpv6.pcap, dhcp-flow.pcap on the created virtual switch
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 240 seconds
    Then metrics must be correctly generated for flow handler
        And remove the agent .yaml generated on each scenario

@sanity @metric_groups @metrics_flow @root @mocked_interface @auto_provision
Scenario: flow handler type sflow with only top_interfaces metric groups enabled
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink with default configuration type already exists
        And a mocked interface is configured with mtu: 65000 and ip: 192.168.100.2/24
        And a virtual switch is configured and is up with type sflow and target \"192.168.100.2:available\"
    When an agent(backend_type:pktvisor, settings: {"input_type":"flow","bind":"192.168.100.2", "port":"switch", "flow_type":"sflow"}) is self-provisioned via a configuration file on port available with 1 agent tags and has status online. [Overwrite default: False. Paste only file: True]
        And pktvisor state is running
        And 1 Agent Group(s) is created with all tags contained in the agent
        And a flow policy flow with tap_selector matching all tag(s) of the tap from an agent, top_interfaces metric_groups enabled, top_tos, conversations, counters, by_packets, by_bytes, top_geo, cardinality, top_ports, top_ips_ports, top_ips metric_groups disabled and settings: default is applied to the group
        And 1 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 1 policies are applied and all has status running
        And run mocked data dns_ipv4_tcp.pcap, dns_ipv4_udp.pcap, dns_ipv6_tcp.pcap, dns_ipv6_udp.pcap, dns_udp_mixed_rcode.pcap, dnssec.pcap, dhcpv6.pcap, dhcp-flow.pcap on the created virtual switch
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
    Then metrics must be correctly generated for flow handler
        And remove the agent .yaml generated on each scenario

@sanity @metric_groups @metrics_flow @root @mocked_interface @auto_provision
Scenario: flow handler type sflow with only top_interfaces and by_bytes metric groups enabled
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink with default configuration type already exists
        And a mocked interface is configured with mtu: 65000 and ip: 192.168.100.2/24
        And a virtual switch is configured and is up with type sflow and target \"192.168.100.2:available\"
    When an agent(backend_type:pktvisor, settings: {"input_type":"flow","bind":"192.168.100.2", "port":"switch", "flow_type":"sflow"}) is self-provisioned via a configuration file on port available with 1 agent tags and has status online. [Overwrite default: False. Paste only file: True]
        And pktvisor state is running
        And 1 Agent Group(s) is created with all tags contained in the agent
        And a flow policy flow with tap_selector matching all tag(s) of the tap from an agent, top_interfaces, by_bytes metric_groups enabled, top_tos, conversations, counters, by_packets, top_geo, cardinality, top_ports, top_ips_ports, top_ips metric_groups disabled and settings: default is applied to the group
        And 1 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 1 policies are applied and all has status running
        And run mocked data dns_ipv4_tcp.pcap, dns_ipv4_udp.pcap, dns_ipv6_tcp.pcap, dns_ipv6_udp.pcap, dns_udp_mixed_rcode.pcap, dnssec.pcap, dhcpv6.pcap, dhcp-flow.pcap on the created virtual switch
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 240 seconds
    Then metrics must be correctly generated for flow handler
        And remove the agent .yaml generated on each scenario


@sanity @metric_groups @metrics_flow @root @mocked_interface @auto_provision
Scenario: flow handler type sflow with only top_interfaces and by_packets metric groups enabled
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink with default configuration type already exists
        And a mocked interface is configured with mtu: 65000 and ip: 192.168.100.2/24
        And a virtual switch is configured and is up with type sflow and target \"192.168.100.2:available\"
    When an agent(backend_type:pktvisor, settings: {"input_type":"flow","bind":"192.168.100.2", "port":"switch", "flow_type":"sflow"}) is self-provisioned via a configuration file on port available with 1 agent tags and has status online. [Overwrite default: False. Paste only file: True]
        And pktvisor state is running
        And 1 Agent Group(s) is created with all tags contained in the agent
        And a flow policy flow with tap_selector matching all tag(s) of the tap from an agent, top_interfaces, by_packets metric_groups enabled, top_tos, conversations, counters, by_bytes, top_geo, cardinality, top_ports, top_ips_ports, top_ips metric_groups disabled and settings: default is applied to the group
        And 1 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 1 policies are applied and all has status running
        And run mocked data dns_ipv4_tcp.pcap, dns_ipv4_udp.pcap, dns_ipv6_tcp.pcap, dns_ipv6_udp.pcap, dns_udp_mixed_rcode.pcap, dnssec.pcap, dhcpv6.pcap, dhcp-flow.pcap on the created virtual switch
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 240 seconds
    Then metrics must be correctly generated for flow handler
        And remove the agent .yaml generated on each scenario


@sanity @metric_groups @metrics_flow @root @mocked_interface @auto_provision
Scenario: flow handler type sflow with only top_ips metric groups enabled
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink with default configuration type already exists
        And a mocked interface is configured with mtu: 65000 and ip: 192.168.100.2/24
        And a virtual switch is configured and is up with type sflow and target \"192.168.100.2:available\"
    When an agent(backend_type:pktvisor, settings: {"input_type":"flow","bind":"192.168.100.2", "port":"switch", "flow_type":"sflow"}) is self-provisioned via a configuration file on port available with 1 agent tags and has status online. [Overwrite default: False. Paste only file: True]
        And pktvisor state is running
        And 1 Agent Group(s) is created with all tags contained in the agent
        And a flow policy flow with tap_selector matching all tag(s) of the tap from an agent, top_ips metric_groups enabled, top_tos, conversations, counters, by_packets, by_bytes, top_geo, cardinality, top_ports, top_ips_ports, top_interfaces metric_groups disabled and settings: default is applied to the group
        And 1 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 1 policies are applied and all has status running
        And run mocked data dns_ipv4_tcp.pcap, dns_ipv4_udp.pcap, dns_ipv6_tcp.pcap, dns_ipv6_udp.pcap, dns_udp_mixed_rcode.pcap, dnssec.pcap, dhcpv6.pcap, dhcp-flow.pcap on the created virtual switch
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
    Then metrics must be correctly generated for flow handler
        And remove the agent .yaml generated on each scenario

@sanity @metric_groups @metrics_flow @root @mocked_interface @auto_provision
Scenario: flow handler type sflow with only top_ips and by_bytes metric groups enabled
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink with default configuration type already exists
        And a mocked interface is configured with mtu: 65000 and ip: 192.168.100.2/24
        And a virtual switch is configured and is up with type sflow and target \"192.168.100.2:available\"
    When an agent(backend_type:pktvisor, settings: {"input_type":"flow","bind":"192.168.100.2", "port":"switch", "flow_type":"sflow"}) is self-provisioned via a configuration file on port available with 1 agent tags and has status online. [Overwrite default: False. Paste only file: True]
        And pktvisor state is running
        And 1 Agent Group(s) is created with all tags contained in the agent
        And a flow policy flow with tap_selector matching all tag(s) of the tap from an agent, top_ips, by_bytes metric_groups enabled, top_tos, conversations, counters, by_packets, top_geo, cardinality, top_ports, top_ips_ports, top_interfaces metric_groups disabled and settings: default is applied to the group
        And 1 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 1 policies are applied and all has status running
        And run mocked data dns_ipv4_tcp.pcap, dns_ipv4_udp.pcap, dns_ipv6_tcp.pcap, dns_ipv6_udp.pcap, dns_udp_mixed_rcode.pcap, dnssec.pcap, dhcpv6.pcap, dhcp-flow.pcap on the created virtual switch
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 240 seconds
    Then metrics must be correctly generated for flow handler
        And remove the agent .yaml generated on each scenario


@sanity @metric_groups @metrics_flow @root @mocked_interface @auto_provision
Scenario: flow handler type sflow with only top_ips and by_packets metric groups enabled
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink with default configuration type already exists
        And a mocked interface is configured with mtu: 65000 and ip: 192.168.100.2/24
        And a virtual switch is configured and is up with type sflow and target \"192.168.100.2:available\"
    When an agent(backend_type:pktvisor, settings: {"input_type":"flow","bind":"192.168.100.2", "port":"switch", "flow_type":"sflow"}) is self-provisioned via a configuration file on port available with 1 agent tags and has status online. [Overwrite default: False. Paste only file: True]
        And pktvisor state is running
        And 1 Agent Group(s) is created with all tags contained in the agent
        And a flow policy flow with tap_selector matching all tag(s) of the tap from an agent, top_ips, by_packets metric_groups enabled, top_tos, conversations, counters, by_bytes, top_geo, cardinality, top_ports, top_ips_ports, top_interfaces metric_groups disabled and settings: default is applied to the group
        And 1 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 1 policies are applied and all has status running
        And run mocked data dns_ipv4_tcp.pcap, dns_ipv4_udp.pcap, dns_ipv6_tcp.pcap, dns_ipv6_udp.pcap, dns_udp_mixed_rcode.pcap, dnssec.pcap, dhcpv6.pcap, dhcp-flow.pcap on the created virtual switch
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 240 seconds
    Then metrics must be correctly generated for flow handler
        And remove the agent .yaml generated on each scenario


@sanity @metric_groups @metrics_flow @root @mocked_interface @auto_provision
Scenario: flow handler type sflow with only top_tos metric groups enabled
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink with default configuration type already exists
        And a mocked interface is configured with mtu: 65000 and ip: 192.168.100.2/24
        And a virtual switch is configured and is up with type sflow and target \"192.168.100.2:available\"
    When an agent(backend_type:pktvisor, settings: {"input_type":"flow","bind":"192.168.100.2", "port":"switch", "flow_type":"sflow"}) is self-provisioned via a configuration file on port available with 1 agent tags and has status online. [Overwrite default: False. Paste only file: True]
        And pktvisor state is running
        And 1 Agent Group(s) is created with all tags contained in the agent
        And a flow policy flow with tap_selector matching all tag(s) of the tap from an agent, top_tos metric_groups enabled, top_ips, by_packets, conversations, counters, by_bytes, top_geo, cardinality, top_ports, top_ips_ports, top_interfaces metric_groups disabled and settings: default is applied to the group
        And 1 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 1 policies are applied and all has status running
        And run mocked data dns_ipv4_tcp.pcap, dns_ipv4_udp.pcap, dns_ipv6_tcp.pcap, dns_ipv6_udp.pcap, dns_udp_mixed_rcode.pcap, dnssec.pcap, dhcpv6.pcap, dhcp-flow.pcap on the created virtual switch
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
    Then metrics must be correctly generated for flow handler
        And remove the agent .yaml generated on each scenario


@sanity @metric_groups @metrics_flow @root @mocked_interface @auto_provision
Scenario: flow handler type sflow with only top_tos and by_bytes metric groups enabled
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink with default configuration type already exists
        And a mocked interface is configured with mtu: 65000 and ip: 192.168.100.2/24
        And a virtual switch is configured and is up with type sflow and target \"192.168.100.2:available\"
    When an agent(backend_type:pktvisor, settings: {"input_type":"flow","bind":"192.168.100.2", "port":"switch", "flow_type":"sflow"}) is self-provisioned via a configuration file on port available with 1 agent tags and has status online. [Overwrite default: False. Paste only file: True]
        And pktvisor state is running
        And 1 Agent Group(s) is created with all tags contained in the agent
        And a flow policy flow with tap_selector matching all tag(s) of the tap from an agent, top_tos, by_bytes metric_groups enabled, top_ips, by_packets, conversations, counters, top_geo, cardinality, top_ports, top_ips_ports, top_interfaces metric_groups disabled and settings: default is applied to the group
        And 1 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 1 policies are applied and all has status running
        And run mocked data dns_ipv4_tcp.pcap, dns_ipv4_udp.pcap, dns_ipv6_tcp.pcap, dns_ipv6_udp.pcap, dns_udp_mixed_rcode.pcap, dnssec.pcap, dhcpv6.pcap, dhcp-flow.pcap on the created virtual switch
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 240 seconds
    Then metrics must be correctly generated for flow handler
        And remove the agent .yaml generated on each scenario


@sanity @metric_groups @metrics_flow @root @mocked_interface @auto_provision
Scenario: flow handler type sflow with only top_tos and by_packets metric groups enabled
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink with default configuration type already exists
        And a mocked interface is configured with mtu: 65000 and ip: 192.168.100.2/24
        And a virtual switch is configured and is up with type sflow and target \"192.168.100.2:available\"
    When an agent(backend_type:pktvisor, settings: {"input_type":"flow","bind":"192.168.100.2", "port":"switch", "flow_type":"sflow"}) is self-provisioned via a configuration file on port available with 1 agent tags and has status online. [Overwrite default: False. Paste only file: True]
        And pktvisor state is running
        And 1 Agent Group(s) is created with all tags contained in the agent
        And a flow policy flow with tap_selector matching all tag(s) of the tap from an agent, top_tos, by_packets metric_groups enabled, top_ips, conversations, counters, by_bytes, top_geo, cardinality, top_ports, top_ips_ports, top_interfaces metric_groups disabled and settings: default is applied to the group
        And 1 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 1 policies are applied and all has status running
        And run mocked data dns_ipv4_tcp.pcap, dns_ipv4_udp.pcap, dns_ipv6_tcp.pcap, dns_ipv6_udp.pcap, dns_udp_mixed_rcode.pcap, dnssec.pcap, dhcpv6.pcap, dhcp-flow.pcap on the created virtual switch
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 240 seconds
    Then metrics must be correctly generated for flow handler
        And remove the agent .yaml generated on each scenario


#### pcap

@sanity @metric_groups @metrics_pcap @root @mocked_interface @auto_provision
Scenario: pcap handler with default metric groups configuration
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink with default configuration type already exists
        And a mocked interface is configured with mtu: 65000 and ip: None
    When an agent(backend_type:pktvisor, settings: {"input_type":"pcap","iface":"mocked"}) is self-provisioned via a configuration file on port available with 1 agent tags and has status online. [Overwrite default: False. Paste only file: True]
        And pktvisor state is running
        And 1 Agent Group(s) is created with all tags contained in the agent
        And a pcap policy pcap with tap_selector matching all tag(s) of the tap from an agent, default metric_groups enabled, default metric_groups disabled and settings: default is applied to the group
        And 1 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 1 policies are applied and all has status running
        And run mocked data dns_ipv4_tcp.pcap, dns_ipv4_udp.pcap, dns_ipv6_tcp.pcap, dns_ipv6_udp.pcap, dns_udp_mixed_rcode.pcap, dnssec.pcap, dhcpv6.pcap on the created virtual interface
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 240 seconds
    Then metrics must be correctly generated for pcap handler
        And remove the agent .yaml generated on each scenario


@sanity @metric_groups @metrics_pcap @root @mocked_interface @auto_provision
Scenario: pcap handler with all metric groups enabled
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink with default configuration type already exists
        And a mocked interface is configured with mtu: 65000 and ip: None
    When an agent(backend_type:pktvisor, settings: {"input_type":"pcap","iface":"mocked"}) is self-provisioned via a configuration file on port available with 1 agent tags and has status online. [Overwrite default: False. Paste only file: True]
        And pktvisor state is running
        And 1 Agent Group(s) is created with all tags contained in the agent
        And a pcap policy pcap with tap_selector matching all tag(s) of the tap from an agent, all metric_groups enabled, none metric_groups disabled and settings: default is applied to the group
        And 1 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 1 policies are applied and all has status running
        And run mocked data dns_ipv4_tcp.pcap, dns_ipv4_udp.pcap, dns_ipv6_tcp.pcap, dns_ipv6_udp.pcap, dns_udp_mixed_rcode.pcap, dnssec.pcap, dhcpv6.pcap on the created virtual interface
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 240 seconds
    Then metrics must be correctly generated for pcap handler
        And remove the agent .yaml generated on each scenario


@sanity @metric_groups @metrics_pcap @root @mocked_interface @auto_provision
Scenario: pcap handler with all metric groups disabled
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink with default configuration type already exists
        And a mocked interface is configured with mtu: 65000 and ip: None
    When an agent(backend_type:pktvisor, settings: {"input_type":"pcap","iface":"mocked"}) is self-provisioned via a configuration file on port available with 1 agent tags and has status online. [Overwrite default: False. Paste only file: True]
        And pktvisor state is running
        And 1 Agent Group(s) is created with all tags contained in the agent
        And a pcap policy pcap with tap_selector matching all tag(s) of the tap from an agent, none metric_groups enabled, all metric_groups disabled and settings: default is applied to the group
        And 1 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 1 policies are applied and all has status running
        And run mocked data dns_ipv4_tcp.pcap, dns_ipv4_udp.pcap, dns_ipv6_tcp.pcap, dns_ipv6_udp.pcap, dns_udp_mixed_rcode.pcap, dnssec.pcap, dhcpv6.pcap on the created virtual interface
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 240 seconds
    Then metrics must be correctly generated for pcap handler
        And remove the agent .yaml generated on each scenario


#### bgp

@sanity @metric_groups @metrics_bgp @root @mocked_interface @auto_provision
Scenario: bgp handler with default metric groups configuration
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink with default configuration type already exists
        And a mocked interface is configured with mtu: 65000 and ip: None
    When an agent(backend_type:pktvisor, settings: {"input_type":"pcap","iface":"mocked"}) is self-provisioned via a configuration file on port available with 1 agent tags and has status online. [Overwrite default: False. Paste only file: True]
        And pktvisor state is running
        And 1 Agent Group(s) is created with all tags contained in the agent
        And a bgp policy pcap with tap_selector matching all tag(s) of the tap from an agent, default metric_groups enabled, default metric_groups disabled and settings: default is applied to the group
        And 1 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 1 policies are applied and all has status running
        And run mocked data bgp.pcap on the created virtual interface
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 240 seconds
    Then metrics must be correctly generated for bgp handler
        And remove the agent .yaml generated on each scenario


@sanity @metric_groups @metrics_bgp @root @mocked_interface @auto_provision
Scenario: bgp handler with all metric groups enabled
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink with default configuration type already exists
        And a mocked interface is configured with mtu: 65000 and ip: None
    When an agent(backend_type:pktvisor, settings: {"input_type":"pcap","iface":"mocked"}) is self-provisioned via a configuration file on port available with 1 agent tags and has status online. [Overwrite default: False. Paste only file: True]
        And pktvisor state is running
        And 1 Agent Group(s) is created with all tags contained in the agent
        And a bgp policy pcap with tap_selector matching all tag(s) of the tap from an agent, all metric_groups enabled, none metric_groups disabled and settings: default is applied to the group
        And 1 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 1 policies are applied and all has status running
        And run mocked data bgp.pcap on the created virtual interface
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 240 seconds
    Then metrics must be correctly generated for bgp handler
        And remove the agent .yaml generated on each scenario


@sanity @metric_groups @metrics_bgp @root @mocked_interface @auto_provision
Scenario: bgp handler with all metric groups disabled
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink with default configuration type already exists
        And a mocked interface is configured with mtu: 65000 and ip: None
    When an agent(backend_type:pktvisor, settings: {"input_type":"pcap","iface":"mocked"}) is self-provisioned via a configuration file on port available with 1 agent tags and has status online. [Overwrite default: False. Paste only file: True]
        And pktvisor state is running
        And 1 Agent Group(s) is created with all tags contained in the agent
        And a bgp policy pcap with tap_selector matching all tag(s) of the tap from an agent, none metric_groups enabled, all metric_groups disabled and settings: default is applied to the group
        And 1 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 1 policies are applied and all has status running
    And run mocked data bgp.pcap on the created virtual interface
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 240 seconds
    Then metrics must be correctly generated for bgp handler
        And remove the agent .yaml generated on each scenario


#### dhcp

@sanity @metric_groups @metrics_dhcp @root @mocked_interface @auto_provision
Scenario: dhcp handler with default metric groups configuration
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink with default configuration type already exists
        And a mocked interface is configured with mtu: 65000 and ip: None
    When an agent(backend_type:pktvisor, settings: {"input_type":"pcap","iface":"mocked"}) is self-provisioned via a configuration file on port available with 1 agent tags and has status online. [Overwrite default: False. Paste only file: True]
        And pktvisor state is running
        And 1 Agent Group(s) is created with all tags contained in the agent
        And a dhcp policy pcap with tap_selector matching all tag(s) of the tap from an agent, default metric_groups enabled, default metric_groups disabled and settings: default is applied to the group
        And 1 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 1 policies are applied and all has status running
        And run mocked data dns_ipv4_tcp.pcap, dns_ipv4_udp.pcap, dns_ipv6_tcp.pcap, dns_ipv6_udp.pcap, dns_udp_mixed_rcode.pcap, dhcpv6.pcap, dhcp-flow.pcap on the created virtual interface
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 240 seconds
    Then metrics must be correctly generated for dhcp handler
        And remove the agent .yaml generated on each scenario


@sanity @metric_groups @metrics_dhcp @root @mocked_interface @auto_provision
Scenario: dhcp handler with all metric groups enabled
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink with default configuration type already exists
        And a mocked interface is configured with mtu: 65000 and ip: None
    When an agent(backend_type:pktvisor, settings: {"input_type":"pcap","iface":"mocked"}) is self-provisioned via a configuration file on port available with 1 agent tags and has status online. [Overwrite default: False. Paste only file: True]
        And pktvisor state is running
        And 1 Agent Group(s) is created with all tags contained in the agent
        And a dhcp policy pcap with tap_selector matching all tag(s) of the tap from an agent, all metric_groups enabled, none metric_groups disabled and settings: default is applied to the group
        And 1 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 1 policies are applied and all has status running
        And run mocked data dns_ipv4_tcp.pcap, dns_ipv4_udp.pcap, dns_ipv6_tcp.pcap, dns_ipv6_udp.pcap, dns_udp_mixed_rcode.pcap, dhcpv6.pcap, dhcp-flow.pcap on the created virtual interface
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 240 seconds
    Then metrics must be correctly generated for dhcp handler
        And remove the agent .yaml generated on each scenario


@sanity @metric_groups @metrics_dhcp @root @mocked_interface @auto_provision
Scenario: dhcp handler with all metric groups disabled
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink with default configuration type already exists
        And a mocked interface is configured with mtu: 65000 and ip: None
    When an agent(backend_type:pktvisor, settings: {"input_type":"pcap","iface":"mocked"}) is self-provisioned via a configuration file on port available with 1 agent tags and has status online. [Overwrite default: False. Paste only file: True]
        And pktvisor state is running
        And 1 Agent Group(s) is created with all tags contained in the agent
        And a dhcp policy pcap with tap_selector matching all tag(s) of the tap from an agent, none metric_groups enabled, all metric_groups disabled and settings: default is applied to the group
        And 1 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 1 policies are applied and all has status running
        And run mocked data dns_ipv4_tcp.pcap, dns_ipv4_udp.pcap, dns_ipv6_tcp.pcap, dns_ipv6_udp.pcap, dns_udp_mixed_rcode.pcap, dhcpv6.pcap, dhcp-flow.pcap on the created virtual interface
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 240 seconds
    Then metrics must be correctly generated for dhcp handler
        And remove the agent .yaml generated on each scenario


#### net v1.0

@sanity @metric_groups @metrics_net @root @mocked_interface @auto_provision
Scenario: net handler with default metric groups configuration
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink with default configuration type already exists
        And a mocked interface is configured with mtu: 65000 and ip: None
    When an agent(backend_type:pktvisor, settings: {"input_type":"pcap","iface":"mocked", "host_spec":"fe80::a00:27ff:fed4:10bb/48,192.168.0.0/24,75.75.75.75/32"}) is self-provisioned via a configuration file on port available with 1 agent tags and has status online. [Overwrite default: False. Paste only file: True]
        And pktvisor state is running
        And 1 Agent Group(s) is created with all tags contained in the agent
        And a net policy pcap with tap_selector matching all tag(s) of the tap from an agent, default metric_groups enabled, default metric_groups disabled and settings: default is applied to the group
        And 1 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 1 policies are applied and all has status running
        And run mocked data dns_ipv4_tcp.pcap, dns_ipv4_udp.pcap, dns_ipv6_tcp.pcap, dns_ipv6_udp.pcap, dns_udp_mixed_rcode.pcap, dns_udp_tcp_random.pcap, ecs.pcap, ipfix.pcap, ecmp.pcap, ipfix.pcap, nf9.pcap, dnssec.pcap, dhcpv6.pcap on the created virtual interface
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 240 seconds
    Then metrics must be correctly generated for net handler
        And remove the agent .yaml generated on each scenario


@sanity @metric_groups @metrics_net @root @mocked_interface @auto_provision
Scenario: net handler with all metric groups enabled
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink with default configuration type already exists
        And a mocked interface is configured with mtu: 65000 and ip: None
    When an agent(backend_type:pktvisor, settings: {"input_type":"pcap","iface":"mocked", "host_spec":"fe80::a00:27ff:fed4:10bb/48,192.168.0.0/24,75.75.75.75/32"}) is self-provisioned via a configuration file on port available with 1 agent tags and has status online. [Overwrite default: False. Paste only file: True]
        And pktvisor state is running
        And 1 Agent Group(s) is created with all tags contained in the agent
        And a net policy pcap with tap_selector matching all tag(s) of the tap from an agent, all metric_groups enabled, none metric_groups disabled and settings: default is applied to the group
        And 1 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 1 policies are applied and all has status running
        And run mocked data dns_ipv4_tcp.pcap, dns_ipv4_udp.pcap, dns_ipv6_tcp.pcap, dns_ipv6_udp.pcap, dns_udp_mixed_rcode.pcap, dns_udp_tcp_random.pcap, ecs.pcap, ipfix.pcap, ecmp.pcap, ipfix.pcap, nf9.pcap, dnssec.pcap, dhcpv6.pcap on the created virtual interface
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 240 seconds
    Then metrics must be correctly generated for net handler
        And remove the agent .yaml generated on each scenario


@sanity @metric_groups @metrics_net @root @mocked_interface @auto_provision
Scenario: net handler with all metric groups disabled
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink with default configuration type already exists
        And a mocked interface is configured with mtu: 65000 and ip: None
    When an agent(backend_type:pktvisor, settings: {"input_type":"pcap","iface":"mocked", "host_spec":"fe80::a00:27ff:fed4:10bb/48,192.168.0.0/24,75.75.75.75/32"}) is self-provisioned via a configuration file on port available with 1 agent tags and has status online. [Overwrite default: False. Paste only file: True]
        And pktvisor state is running
        And 1 Agent Group(s) is created with all tags contained in the agent
        And a net policy pcap with tap_selector matching all tag(s) of the tap from an agent, none metric_groups enabled, all metric_groups disabled and settings: default is applied to the group
        And 1 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 1 policies are applied and all has status running
        And run mocked data dns_ipv4_tcp.pcap, dns_ipv4_udp.pcap, dns_ipv6_tcp.pcap, dns_ipv6_udp.pcap, dns_udp_mixed_rcode.pcap, dns_udp_tcp_random.pcap, ecs.pcap, ipfix.pcap, ecmp.pcap, ipfix.pcap, nf9.pcap, dnssec.pcap, dhcpv6.pcap on the created virtual interface
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
    Then metrics must be correctly generated for net handler
        And remove the agent .yaml generated on each scenario


@sanity @metric_groups @metrics_net @root @mocked_interface @auto_provision
Scenario: net handler with only cardinality metric groups enabled
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink with default configuration type already exists
        And a mocked interface is configured with mtu: 65000 and ip: None
    When an agent(backend_type:pktvisor, settings: {"input_type":"pcap","iface":"mocked", "host_spec":"fe80::a00:27ff:fed4:10bb/48,192.168.0.0/24,75.75.75.75/32"}) is self-provisioned via a configuration file on port available with 1 agent tags and has status online. [Overwrite default: False. Paste only file: True]
        And pktvisor state is running
        And 1 Agent Group(s) is created with all tags contained in the agent
        And a net policy pcap with tap_selector matching all tag(s) of the tap from an agent, cardinality metric_groups enabled, counters, top_geo, top_ips metric_groups disabled and settings: default is applied to the group
        And 1 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 1 policies are applied and all has status running
        And run mocked data dns_ipv4_tcp.pcap, dns_ipv4_udp.pcap, dns_ipv6_tcp.pcap, dns_ipv6_udp.pcap, dns_udp_mixed_rcode.pcap, dns_udp_tcp_random.pcap, ecs.pcap, ipfix.pcap, ecmp.pcap, ipfix.pcap, nf9.pcap, dnssec.pcap, dhcpv6.pcap on the created virtual interface
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 240 seconds
    Then metrics must be correctly generated for net handler
        And remove the agent .yaml generated on each scenario


@sanity @metric_groups @metrics_net @root @mocked_interface @auto_provision
Scenario: net handler with only counters metric groups enabled
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink with default configuration type already exists
        And a mocked interface is configured with mtu: 65000 and ip: None
    When an agent(backend_type:pktvisor, settings: {"input_type":"pcap","iface":"mocked", "host_spec":"fe80::a00:27ff:fed4:10bb/48,192.168.0.0/24,75.75.75.75/32"}) is self-provisioned via a configuration file on port available with 1 agent tags and has status online. [Overwrite default: False. Paste only file: True]
        And pktvisor state is running
        And 1 Agent Group(s) is created with all tags contained in the agent
        And a net policy pcap with tap_selector matching all tag(s) of the tap from an agent, counters metric_groups enabled, cardinality, top_geo, top_ips metric_groups disabled and settings: default is applied to the group
        And 1 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 1 policies are applied and all has status running
        And run mocked data dns_ipv4_tcp.pcap, dns_ipv4_udp.pcap, dns_ipv6_tcp.pcap, dns_ipv6_udp.pcap, dns_udp_mixed_rcode.pcap, dns_udp_tcp_random.pcap, ecs.pcap, ipfix.pcap, ecmp.pcap, ipfix.pcap, nf9.pcap, dnssec.pcap, dhcpv6.pcap on the created virtual interface
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 240 seconds
    Then metrics must be correctly generated for net handler
        And remove the agent .yaml generated on each scenario


@sanity @metric_groups @metrics_net @root @mocked_interface @auto_provision
Scenario: net handler with only top_geo metric groups enabled
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink with default configuration type already exists
        And a mocked interface is configured with mtu: 65000 and ip: None
    When an agent(backend_type:pktvisor, settings: {"input_type":"pcap","iface":"mocked", "host_spec":"fe80::a00:27ff:fed4:10bb/48,192.168.0.0/24,75.75.75.75/32"}) is self-provisioned via a configuration file on port available with 1 agent tags and has status online. [Overwrite default: False. Paste only file: True]
        And pktvisor state is running
        And 1 Agent Group(s) is created with all tags contained in the agent
        And a net policy pcap with tap_selector matching all tag(s) of the tap from an agent, top_geo metric_groups enabled, counters, cardinality, top_ips metric_groups disabled and settings: default is applied to the group
        And 1 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 1 policies are applied and all has status running
        And run mocked data dns_ipv4_tcp.pcap, dns_ipv4_udp.pcap, dns_ipv6_tcp.pcap, dns_ipv6_udp.pcap, dns_udp_mixed_rcode.pcap, dns_udp_tcp_random.pcap, ecs.pcap, ipfix.pcap, ecmp.pcap, ipfix.pcap, nf9.pcap, dnssec.pcap, dhcpv6.pcap on the created virtual interface
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 240 seconds
    Then metrics must be correctly generated for net handler
        And remove the agent .yaml generated on each scenario


@sanity @metric_groups @metrics_net @root @mocked_interface @auto_provision
Scenario: net handler with only top_ips metric groups enabled
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink with default configuration type already exists
        And a mocked interface is configured with mtu: 65000 and ip: None
    When an agent(backend_type:pktvisor, settings: {"input_type":"pcap","iface":"mocked", "host_spec":"fe80::a00:27ff:fed4:10bb/48,192.168.0.0/24,75.75.75.75/32"}) is self-provisioned via a configuration file on port available with 1 agent tags and has status online. [Overwrite default: False. Paste only file: True]
        And pktvisor state is running
        And 1 Agent Group(s) is created with all tags contained in the agent
        And a net policy pcap with tap_selector matching all tag(s) of the tap from an agent, top_ips metric_groups enabled, counters, top_geo, cardinality metric_groups disabled and settings: default is applied to the group
        And 1 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 1 policies are applied and all has status running
        And run mocked data dns_ipv4_tcp.pcap, dns_ipv4_udp.pcap, dns_ipv6_tcp.pcap, dns_ipv6_udp.pcap, dns_udp_mixed_rcode.pcap, dns_udp_tcp_random.pcap, ecs.pcap, ipfix.pcap, ecmp.pcap, ipfix.pcap, nf9.pcap, dnssec.pcap, dhcpv6.pcap on the created virtual interface
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 240 seconds
    Then metrics must be correctly generated for net handler
        And remove the agent .yaml generated on each scenario


#### dns v1.0

@sanity @metric_groups @metrics_dns @root @mocked_interface @auto_provision
Scenario: dns handler with default metric groups configuration
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink with default configuration type already exists
        And a mocked interface is configured with mtu: 65000 and ip: None
    When an agent(backend_type:pktvisor, settings: {"input_type":"pcap","iface":"mocked", "host_spec":"fe80::a00:27ff:fed4:10bb/48,192.168.0.0/24,75.75.75.75/32"}) is self-provisioned via a configuration file on port available with 1 agent tags and has status online. [Overwrite default: False. Paste only file: True]
        And pktvisor state is running
        And 1 Agent Group(s) is created with all tags contained in the agent
        And a dns policy pcap with tap_selector matching all tag(s) of the tap from an agent, default metric_groups enabled, default metric_groups disabled and settings: default is applied to the group
        And 1 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 1 policies are applied and all has status running
        And run mocked data dns_ipv4_tcp.pcap, dns_ipv4_udp.pcap, dns_ipv6_tcp.pcap, dns_ipv6_udp.pcap, dns_udp_mixed_rcode.pcap, dns_udp_tcp_random.pcap, ecs.pcap, dnssec.pcap, dhcpv6.pcap on the created virtual interface
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 240 seconds
    Then metrics must be correctly generated for dns handler
        And remove the agent .yaml generated on each scenario


@sanity @metric_groups @metrics_dns @root @mocked_interface @auto_provision
Scenario: dns handler with all metric groups enabled
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink with default configuration type already exists
        And a mocked interface is configured with mtu: 65000 and ip: None
    When an agent(backend_type:pktvisor, settings: {"input_type":"pcap","iface":"mocked", "host_spec":"fe80::a00:27ff:fed4:10bb/48,192.168.0.0/24,75.75.75.75/32"}) is self-provisioned via a configuration file on port available with 1 agent tags and has status online. [Overwrite default: False. Paste only file: True]
        And pktvisor state is running
        And 1 Agent Group(s) is created with all tags contained in the agent
        And a dns policy pcap with tap_selector matching all tag(s) of the tap from an agent, all metric_groups enabled, none metric_groups disabled and settings: default is applied to the group
        And 1 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 1 policies are applied and all has status running
        And run mocked data dns_ipv4_tcp.pcap, dns_ipv4_udp.pcap, dns_ipv6_tcp.pcap, dns_ipv6_udp.pcap, dns_udp_mixed_rcode.pcap, dns_udp_tcp_random.pcap, ecs.pcap, dnssec.pcap, dhcpv6.pcap on the created virtual interface
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 240 seconds
    Then metrics must be correctly generated for dns handler
        And remove the agent .yaml generated on each scenario


@sanity @metric_groups @metrics_dns @root @mocked_interface @auto_provision
Scenario: dns handler with all metric groups disabled
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink with default configuration type already exists
        And a mocked interface is configured with mtu: 65000 and ip: None
    When an agent(backend_type:pktvisor, settings: {"input_type":"pcap","iface":"mocked", "host_spec":"fe80::a00:27ff:fed4:10bb/48,192.168.0.0/24,75.75.75.75/32"}) is self-provisioned via a configuration file on port available with 1 agent tags and has status online. [Overwrite default: False. Paste only file: True]
        And pktvisor state is running
        And 1 Agent Group(s) is created with all tags contained in the agent
        And a dns policy pcap with tap_selector matching all tag(s) of the tap from an agent, none metric_groups enabled, all metric_groups disabled and settings: default is applied to the group
        And 1 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 1 policies are applied and all has status running
        And run mocked data dns_ipv4_tcp.pcap, dns_ipv4_udp.pcap, dns_ipv6_tcp.pcap, dns_ipv6_udp.pcap, dns_udp_mixed_rcode.pcap, dns_udp_tcp_random.pcap, ecs.pcap, dnssec.pcap, dhcpv6.pcap on the created virtual interface
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
    Then metrics must be correctly generated for dns handler
        And remove the agent .yaml generated on each scenario


@sanity @metric_groups @metrics_dns @root @mocked_interface @auto_provision
Scenario: dns handler with only top_ecs metric groups enabled
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink with default configuration type already exists
        And a mocked interface is configured with mtu: 65000 and ip: None
    When an agent(backend_type:pktvisor, settings: {"input_type":"pcap","iface":"mocked", "host_spec":"fe80::a00:27ff:fed4:10bb/48,192.168.0.0/24,75.75.75.75/32"}) is self-provisioned via a configuration file on port available with 1 agent tags and has status online. [Overwrite default: False. Paste only file: True]
        And pktvisor state is running
        And 1 Agent Group(s) is created with all tags contained in the agent
        And a dns policy pcap with tap_selector matching all tag(s) of the tap from an agent, top_ecs metric_groups enabled, quantiles, top_qnames_details, cardinality, counters, dns_transaction, top_qnames, top_ports metric_groups disabled and settings: default is applied to the group
        And 1 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 1 policies are applied and all has status running
        And run mocked data dns_ipv4_tcp.pcap, dns_ipv4_udp.pcap, dns_ipv6_tcp.pcap, dns_ipv6_udp.pcap, dns_udp_mixed_rcode.pcap, dns_udp_tcp_random.pcap, ecs.pcap, dnssec.pcap, dhcpv6.pcap on the created virtual interface
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 240 seconds
    Then metrics must be correctly generated for dns handler
        And remove the agent .yaml generated on each scenario


@sanity @metric_groups @metrics_dns @root @mocked_interface @auto_provision
Scenario: dns handler with only top_qnames_details metric groups enabled
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink with default configuration type already exists
        And a mocked interface is configured with mtu: 65000 and ip: None
    When an agent(backend_type:pktvisor, settings: {"input_type":"pcap","iface":"mocked", "host_spec":"fe80::a00:27ff:fed4:10bb/48,192.168.0.0/24,75.75.75.75/32"}) is self-provisioned via a configuration file on port available with 1 agent tags and has status online. [Overwrite default: False. Paste only file: True]
        And pktvisor state is running
        And 1 Agent Group(s) is created with all tags contained in the agent
        And a dns policy pcap with tap_selector matching all tag(s) of the tap from an agent, top_qnames_details metric_groups enabled, quantiles, top_ecs, cardinality, counters, dns_transaction, top_qnames, top_ports metric_groups disabled and settings: default is applied to the group
        And 1 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 1 policies are applied and all has status running
        And run mocked data dns_ipv4_tcp.pcap, dns_ipv4_udp.pcap, dns_ipv6_tcp.pcap, dns_ipv6_udp.pcap, dns_udp_mixed_rcode.pcap, dns_udp_tcp_random.pcap, ecs.pcap, dnssec.pcap, dhcpv6.pcap on the created virtual interface
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 240 seconds
    Then metrics must be correctly generated for dns handler
        And remove the agent .yaml generated on each scenario


@sanity @metric_groups @metrics_dns @root @mocked_interface @auto_provision
Scenario: dns handler with only cardinality metric groups enabled
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink with default configuration type already exists
        And a mocked interface is configured with mtu: 65000 and ip: None
    When an agent(backend_type:pktvisor, settings: {"input_type":"pcap","iface":"mocked", "host_spec":"fe80::a00:27ff:fed4:10bb/48,192.168.0.0/24,75.75.75.75/32"}) is self-provisioned via a configuration file on port available with 1 agent tags and has status online. [Overwrite default: False. Paste only file: True]
        And pktvisor state is running
        And 1 Agent Group(s) is created with all tags contained in the agent
        And a dns policy pcap with tap_selector matching all tag(s) of the tap from an agent, cardinality metric_groups enabled, quantiles, top_qnames_details, top_ecs, counters, dns_transaction, top_qnames, top_ports metric_groups disabled and settings: default is applied to the group
        And 1 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 1 policies are applied and all has status running
        And run mocked data dns_ipv4_tcp.pcap, dns_ipv4_udp.pcap, dns_ipv6_tcp.pcap, dns_ipv6_udp.pcap, dns_udp_mixed_rcode.pcap, dns_udp_tcp_random.pcap, ecs.pcap, dnssec.pcap, dhcpv6.pcap on the created virtual interface
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 240 seconds
    Then metrics must be correctly generated for dns handler
        And remove the agent .yaml generated on each scenario


@sanity @metric_groups @metrics_dns @root @mocked_interface @auto_provision
Scenario: dns handler with only counters metric groups enabled
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink with default configuration type already exists
        And a mocked interface is configured with mtu: 65000 and ip: None
    When an agent(backend_type:pktvisor, settings: {"input_type":"pcap","iface":"mocked", "host_spec":"fe80::a00:27ff:fed4:10bb/48,192.168.0.0/24,75.75.75.75/32"}) is self-provisioned via a configuration file on port available with 1 agent tags and has status online. [Overwrite default: False. Paste only file: True]
        And pktvisor state is running
        And 1 Agent Group(s) is created with all tags contained in the agent
        And a dns policy pcap with tap_selector matching all tag(s) of the tap from an agent, counters metric_groups enabled, quantiles, top_qnames_details, cardinality, top_ecs, dns_transaction, top_qnames, top_ports metric_groups disabled and settings: default is applied to the group
        And 1 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 1 policies are applied and all has status running
        And run mocked data dns_ipv4_tcp.pcap, dns_ipv4_udp.pcap, dns_ipv6_tcp.pcap, dns_ipv6_udp.pcap, dns_udp_mixed_rcode.pcap, dns_udp_tcp_random.pcap, ecs.pcap, dnssec.pcap, dhcpv6.pcap on the created virtual interface
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 240 seconds
    Then metrics must be correctly generated for dns handler
        And remove the agent .yaml generated on each scenario


@sanity @metric_groups @metrics_dns @root @mocked_interface @auto_provision
Scenario: dns handler with only dns_transaction metric groups enabled
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink with default configuration type already exists
        And a mocked interface is configured with mtu: 65000 and ip: None
    When an agent(backend_type:pktvisor, settings: {"input_type":"pcap","iface":"mocked", "host_spec":"fe80::a00:27ff:fed4:10bb/48,192.168.0.0/24,75.75.75.75/32"}) is self-provisioned via a configuration file on port available with 1 agent tags and has status online. [Overwrite default: False. Paste only file: True]
        And pktvisor state is running
        And 1 Agent Group(s) is created with all tags contained in the agent
        And a dns policy pcap with tap_selector matching all tag(s) of the tap from an agent, dns_transaction metric_groups enabled, quantiles, top_qnames_details, cardinality, counters, top_ecs, top_qnames, top_ports metric_groups disabled and settings: default is applied to the group
        And 1 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 1 policies are applied and all has status running
        And run mocked data dns_ipv4_tcp.pcap, dns_ipv4_udp.pcap, dns_ipv6_tcp.pcap, dns_ipv6_udp.pcap, dns_udp_mixed_rcode.pcap, dns_udp_tcp_random.pcap, ecs.pcap, dnssec.pcap, dhcpv6.pcap on the created virtual interface
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 240 seconds
    Then metrics must be correctly generated for dns handler
        And remove the agent .yaml generated on each scenario


@sanity @metric_groups @metrics_dns @root @mocked_interface @auto_provision
Scenario: dns handler with only top_qnames metric groups enabled
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink with default configuration type already exists
        And a mocked interface is configured with mtu: 65000 and ip: None
    When an agent(backend_type:pktvisor, settings: {"input_type":"pcap","iface":"mocked", "host_spec":"fe80::a00:27ff:fed4:10bb/48,192.168.0.0/24,75.75.75.75/32"}) is self-provisioned via a configuration file on port available with 1 agent tags and has status online. [Overwrite default: False. Paste only file: True]
        And pktvisor state is running
        And 1 Agent Group(s) is created with all tags contained in the agent
        And a dns policy pcap with tap_selector matching all tag(s) of the tap from an agent, top_qnames metric_groups enabled, quantiles, top_qnames_details, cardinality, counters, dns_transaction, top_ecs, top_ports metric_groups disabled and settings: default is applied to the group
        And 1 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 1 policies are applied and all has status running
        And run mocked data dns_ipv4_tcp.pcap, dns_ipv4_udp.pcap, dns_ipv6_tcp.pcap, dns_ipv6_udp.pcap, dns_udp_mixed_rcode.pcap, dns_udp_tcp_random.pcap, ecs.pcap, dnssec.pcap, dhcpv6.pcap on the created virtual interface
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 240 seconds
    Then metrics must be correctly generated for dns handler
        And remove the agent .yaml generated on each scenario


@sanity @metric_groups @metrics_dns @root @mocked_interface @auto_provision
Scenario: dns handler with only top_ports metric groups enabled
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink with default configuration type already exists
        And a mocked interface is configured with mtu: 65000 and ip: None
    When an agent(backend_type:pktvisor, settings: {"input_type":"pcap","iface":"mocked", "host_spec":"fe80::a00:27ff:fed4:10bb/48,192.168.0.0/24,75.75.75.75/32"}) is self-provisioned via a configuration file on port available with 1 agent tags and has status online. [Overwrite default: False. Paste only file: True]
        And pktvisor state is running
        And 1 Agent Group(s) is created with all tags contained in the agent
        And a dns policy pcap with tap_selector matching all tag(s) of the tap from an agent, top_ports metric_groups enabled, quantiles, top_qnames_details, cardinality, counters, dns_transaction, top_qnames, top_ecs metric_groups disabled and settings: default is applied to the group
        And 1 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 1 policies are applied and all has status running
        And run mocked data dns_ipv4_tcp.pcap, dns_ipv4_udp.pcap, dns_ipv6_tcp.pcap, dns_ipv6_udp.pcap, dns_udp_mixed_rcode.pcap, dns_udp_tcp_random.pcap, ecs.pcap, dnssec.pcap, dhcpv6.pcap on the created virtual interface
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 240 seconds
    Then metrics must be correctly generated for dns handler
        And remove the agent .yaml generated on each scenario


@sanity @metric_groups @metrics_dns @root @mocked_interface @auto_provision
Scenario: dns handler with only histograms metric groups enabled
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink with default configuration type already exists
        And a mocked interface is configured with mtu: 65000 and ip: None
    When an agent(backend_type:pktvisor, settings: {"input_type":"pcap","iface":"mocked", "host_spec":"fe80::a00:27ff:fed4:10bb/48,192.168.0.0/24,75.75.75.75/32"}) is self-provisioned via a configuration file on port available with 1 agent tags and has status online. [Overwrite default: False. Paste only file: True]
        And pktvisor state is running
        And 1 Agent Group(s) is created with all tags contained in the agent
        And a dns policy pcap with tap_selector matching all tag(s) of the tap from an agent, histograms metric_groups enabled, quantiles, top_ports, top_qnames_details, cardinality, counters, dns_transaction, top_qnames, top_ecs metric_groups disabled and settings: default is applied to the group
        And 1 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 1 policies are applied and all has status running
        And run mocked data dns_ipv4_tcp.pcap, dns_ipv4_udp.pcap, dns_ipv6_tcp.pcap, dns_ipv6_udp.pcap, dns_udp_mixed_rcode.pcap, dns_udp_tcp_random.pcap, ecs.pcap, dnssec.pcap, dhcpv6.pcap on the created virtual interface
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
    Then metrics must be correctly generated for dns handler
        And remove the agent .yaml generated on each scenario


@sanity @metric_groups @metrics_dns @root @mocked_interface @auto_provision
Scenario: dns handler with only histograms and dns_transaction metric groups enabled
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink with default configuration type already exists
        And a mocked interface is configured with mtu: 65000 and ip: None
    When an agent(backend_type:pktvisor, settings: {"input_type":"pcap","iface":"mocked", "host_spec":"fe80::a00:27ff:fed4:10bb/48,192.168.0.0/24,75.75.75.75/32"}) is self-provisioned via a configuration file on port available with 1 agent tags and has status online. [Overwrite default: False. Paste only file: True]
        And pktvisor state is running
        And 1 Agent Group(s) is created with all tags contained in the agent
        And a dns policy pcap with tap_selector matching all tag(s) of the tap from an agent, histograms, dns_transaction metric_groups enabled, quantiles, top_ports, top_qnames_details, cardinality, counters, dns_transaction, top_qnames, top_ecs metric_groups disabled and settings: default is applied to the group
        And 1 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 1 policies are applied and all has status running
        And run mocked data dns_ipv4_tcp.pcap, dns_ipv4_udp.pcap, dns_ipv6_tcp.pcap, dns_ipv6_udp.pcap, dns_udp_mixed_rcode.pcap, dns_udp_tcp_random.pcap, ecs.pcap, dnssec.pcap, dhcpv6.pcap on the created virtual interface
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 240 seconds
    Then metrics must be correctly generated for dns handler
        And remove the agent .yaml generated on each scenario


@sanity @metric_groups @metrics_dns @root @mocked_interface @auto_provision
Scenario: dns handler with only quantiles metric groups enabled
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink with default configuration type already exists
        And a mocked interface is configured with mtu: 65000 and ip: None
    When an agent(backend_type:pktvisor, settings: {"input_type":"pcap","iface":"mocked", "host_spec":"fe80::a00:27ff:fed4:10bb/48,192.168.0.0/24,75.75.75.75/32"}) is self-provisioned via a configuration file on port available with 1 agent tags and has status online. [Overwrite default: False. Paste only file: True]
        And pktvisor state is running
        And 1 Agent Group(s) is created with all tags contained in the agent
        And a dns policy pcap with tap_selector matching all tag(s) of the tap from an agent, quantiles metric_groups enabled, histograms, top_ports, top_qnames_details, cardinality, counters, dns_transaction, top_qnames, top_ecs metric_groups disabled and settings: default is applied to the group
        And 1 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 1 policies are applied and all has status running
        And run mocked data dns_ipv4_tcp.pcap, dns_ipv4_udp.pcap, dns_ipv6_tcp.pcap, dns_ipv6_udp.pcap, dns_udp_mixed_rcode.pcap, dns_udp_tcp_random.pcap, ecs.pcap, dnssec.pcap, dhcpv6.pcap on the created virtual interface
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
    Then metrics must be correctly generated for dns handler
        And remove the agent .yaml generated on each scenario


#### dns v2.0

@sanity @metric_groups @metrics_dns_v2 @root @mocked_interface @auto_provision
Scenario: dns handler with default metric groups configuration (v2)
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink with default configuration type already exists
        And a mocked interface is configured with mtu: 65000 and ip: None
    When an agent(backend_type:pktvisor, settings: {"input_type":"pcap","iface":"mocked"}) is self-provisioned via a configuration file on port available with 1 agent tags and has status online. [Overwrite default: False. Paste only file: True]
        And pktvisor state is running
        And 1 Agent Group(s) is created with all tags contained in the agent
        And a dns policy pcap with tap_selector matching all tag(s) of the tap from an agent, default metric_groups enabled, default metric_groups disabled and settings: {"require_version":"2.0"} is applied to the group
        And 1 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 1 policies are applied and all has status running
        And run mocked data dns_ipv4_tcp.pcap, dns_ipv4_udp.pcap, dns_ipv6_tcp.pcap, dns_ipv6_udp.pcap, dns_udp_mixed_rcode.pcap, dns_udp_tcp_random.pcap, ecs.pcap, dnssec.pcap, dhcpv6.pcap on the created virtual interface
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 240 seconds
    Then metrics must be correctly generated for dns-v2 handler
        And remove the agent .yaml generated on each scenario


@sanity @metric_groups @metrics_dns_v2 @root @mocked_interface @auto_provision
Scenario: dns handler with all metric groups enabled (v2)
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink with default configuration type already exists
        And a mocked interface is configured with mtu: 65000 and ip: None
    When an agent(backend_type:pktvisor, settings: {"input_type":"pcap","iface":"mocked"}) is self-provisioned via a configuration file on port available with 1 agent tags and has status online. [Overwrite default: False. Paste only file: True]
        And pktvisor state is running
        And 1 Agent Group(s) is created with all tags contained in the agent
        And a dns policy pcap with tap_selector matching all tag(s) of the tap from an agent, all metric_groups enabled, none metric_groups disabled and settings: {"require_version":"2.0"} is applied to the group
        And 1 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 1 policies are applied and all has status running
        And run mocked data dns_ipv4_tcp.pcap, dns_ipv4_udp.pcap, dns_ipv6_tcp.pcap, dns_ipv6_udp.pcap, dns_udp_mixed_rcode.pcap, dns_udp_tcp_random.pcap, ecs.pcap, dnssec.pcap, dhcpv6.pcap on the created virtual interface
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 240 seconds
    Then metrics must be correctly generated for dns-v2 handler
        And remove the agent .yaml generated on each scenario


@sanity @metric_groups @metrics_dns_v2 @root @mocked_interface @auto_provision
Scenario: dns handler with all metric groups disabled (v2)
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink with default configuration type already exists
        And a mocked interface is configured with mtu: 65000 and ip: None
    When an agent(backend_type:pktvisor, settings: {"input_type":"pcap","iface":"mocked"}) is self-provisioned via a configuration file on port available with 1 agent tags and has status online. [Overwrite default: False. Paste only file: True]
        And pktvisor state is running
        And 1 Agent Group(s) is created with all tags contained in the agent
        And a dns policy pcap with tap_selector matching all tag(s) of the tap from an agent, none metric_groups enabled, all metric_groups disabled and settings: {"require_version":"2.0"} is applied to the group
        And 1 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 1 policies are applied and all has status running
        And run mocked data dns_ipv4_tcp.pcap, dns_ipv4_udp.pcap, dns_ipv6_tcp.pcap, dns_ipv6_udp.pcap, dns_udp_mixed_rcode.pcap, dns_udp_tcp_random.pcap, ecs.pcap, dnssec.pcap, dhcpv6.pcap on the created virtual interface
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
    Then metrics must be correctly generated for dns-v2 handler
        And remove the agent .yaml generated on each scenario


@sanity @metric_groups @metrics_dns_v2 @root @mocked_interface @auto_provision
Scenario: dns handler with only top_ecs metric groups enabled (v2)
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink with default configuration type already exists
        And a mocked interface is configured with mtu: 65000 and ip: None
    When an agent(backend_type:pktvisor, settings: {"input_type":"pcap","iface":"mocked"}) is self-provisioned via a configuration file on port available with 1 agent tags and has status online. [Overwrite default: False. Paste only file: True]
        And pktvisor state is running
        And 1 Agent Group(s) is created with all tags contained in the agent
        And a dns policy pcap with tap_selector matching all tag(s) of the tap from an agent, top_ecs metric_groups enabled, cardinality, counters, top_qnames, top_ports, top_size, xact_times, quantiles, top_qtypes, top_rcodes metric_groups disabled and settings: {"require_version":"2.0"} is applied to the group
        And 1 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 1 policies are applied and all has status running
        And run mocked data dns_ipv4_tcp.pcap, dns_ipv4_udp.pcap, dns_ipv6_tcp.pcap, dns_ipv6_udp.pcap, dns_udp_mixed_rcode.pcap, dns_udp_tcp_random.pcap, ecs.pcap, dnssec.pcap, dhcpv6.pcap on the created virtual interface
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 240 seconds
    Then metrics must be correctly generated for dns-v2 handler
        And remove the agent .yaml generated on each scenario


@sanity @metric_groups @metrics_dns_v2 @root @mocked_interface @auto_provision
Scenario: dns handler with only top_ports metric groups enabled (v2)
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink with default configuration type already exists
        And a mocked interface is configured with mtu: 65000 and ip: None
    When an agent(backend_type:pktvisor, settings: {"input_type":"pcap","iface":"mocked"}) is self-provisioned via a configuration file on port available with 1 agent tags and has status online. [Overwrite default: False. Paste only file: True]
        And pktvisor state is running
        And 1 Agent Group(s) is created with all tags contained in the agent
        And a dns policy pcap with tap_selector matching all tag(s) of the tap from an agent, top_ports metric_groups enabled, top_ecs, cardinality, counters, top_qnames, top_size, xact_times, quantiles, top_qtypes, top_rcodes metric_groups disabled and settings: {"require_version":"2.0"} is applied to the group
        And 1 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 1 policies are applied and all has status running
        And run mocked data dns_ipv4_tcp.pcap, dns_ipv4_udp.pcap, dns_ipv6_tcp.pcap, dns_ipv6_udp.pcap, dns_udp_mixed_rcode.pcap, dns_udp_tcp_random.pcap, ecs.pcap, dnssec.pcap, dhcpv6.pcap on the created virtual interface
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 240 seconds
    Then metrics must be correctly generated for dns-v2 handler
        And remove the agent .yaml generated on each scenario


@sanity @metric_groups @metrics_dns_v2 @root @mocked_interface @auto_provision
Scenario: dns handler with only top_size metric groups enabled (v2)
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink with default configuration type already exists
        And a mocked interface is configured with mtu: 65000 and ip: None
    When an agent(backend_type:pktvisor, settings: {"input_type":"pcap","iface":"mocked"}) is self-provisioned via a configuration file on port available with 1 agent tags and has status online. [Overwrite default: False. Paste only file: True]
        And pktvisor state is running
        And 1 Agent Group(s) is created with all tags contained in the agent
        And a dns policy pcap with tap_selector matching all tag(s) of the tap from an agent, top_size metric_groups enabled, top_ecs, cardinality, counters, top_qnames, top_ports, xact_times, quantiles, top_qtypes, top_rcodes metric_groups disabled and settings: {"require_version":"2.0"} is applied to the group
        And 1 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 1 policies are applied and all has status running
        And run mocked data dns_ipv4_tcp.pcap, dns_ipv4_udp.pcap, dns_ipv6_tcp.pcap, dns_ipv6_udp.pcap, dns_udp_mixed_rcode.pcap, dns_udp_tcp_random.pcap, ecs.pcap, dnssec.pcap, dhcpv6.pcap on the created virtual interface
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 240 seconds
    Then metrics must be correctly generated for dns-v2 handler
        And remove the agent .yaml generated on each scenario


@sanity @metric_groups @metrics_dns_v2 @root @mocked_interface @auto_provision
Scenario: dns handler with only xact_times metric groups enabled (v2)
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink with default configuration type already exists
        And a mocked interface is configured with mtu: 65000 and ip: None
    When an agent(backend_type:pktvisor, settings: {"input_type":"pcap","iface":"mocked"}) is self-provisioned via a configuration file on port available with 1 agent tags and has status online. [Overwrite default: False. Paste only file: True]
        And pktvisor state is running
        And 1 Agent Group(s) is created with all tags contained in the agent
        And a dns policy pcap with tap_selector matching all tag(s) of the tap from an agent, xact_times metric_groups enabled, top_ecs, cardinality, counters, top_qnames, top_ports, top_size, quantiles, top_qtypes, top_rcodes metric_groups disabled and settings: {"require_version":"2.0"} is applied to the group
        And 1 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 1 policies are applied and all has status running
        And run mocked data dns_ipv4_tcp.pcap, dns_ipv4_udp.pcap, dns_ipv6_tcp.pcap, dns_ipv6_udp.pcap, dns_udp_mixed_rcode.pcap, dns_udp_tcp_random.pcap, ecs.pcap, dnssec.pcap, dhcpv6.pcap on the created virtual interface
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 240 seconds
    Then metrics must be correctly generated for dns-v2 handler
        And remove the agent .yaml generated on each scenario


@sanity @metric_groups @metrics_dns_v2 @root @mocked_interface @auto_provision
Scenario: dns handler with only cardinality metric groups enabled (v2)
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink with default configuration type already exists
        And a mocked interface is configured with mtu: 65000 and ip: None
    When an agent(backend_type:pktvisor, settings: {"input_type":"pcap","iface":"mocked"}) is self-provisioned via a configuration file on port available with 1 agent tags and has status online. [Overwrite default: False. Paste only file: True]
        And pktvisor state is running
        And 1 Agent Group(s) is created with all tags contained in the agent
        And a dns policy pcap with tap_selector matching all tag(s) of the tap from an agent, cardinality metric_groups enabled, top_ecs, counters, top_qnames, top_ports, top_size, xact_times, quantiles, top_qtypes, top_rcodes metric_groups disabled and settings: {"require_version":"2.0"} is applied to the group
        And 1 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 1 policies are applied and all has status running
        And run mocked data dns_ipv4_tcp.pcap, dns_ipv4_udp.pcap, dns_ipv6_tcp.pcap, dns_ipv6_udp.pcap, dns_udp_mixed_rcode.pcap, dns_udp_tcp_random.pcap, ecs.pcap, dnssec.pcap, dhcpv6.pcap on the created virtual interface
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 240 seconds
    Then metrics must be correctly generated for dns-v2 handler
        And remove the agent .yaml generated on each scenario


@sanity @metric_groups @metrics_dns_v2 @root @mocked_interface @auto_provision
Scenario: dns handler with only counters metric groups enabled (v2)
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink with default configuration type already exists
        And a mocked interface is configured with mtu: 65000 and ip: None
    When an agent(backend_type:pktvisor, settings: {"input_type":"pcap","iface":"mocked"}) is self-provisioned via a configuration file on port available with 1 agent tags and has status online. [Overwrite default: False. Paste only file: True]
        And pktvisor state is running
        And 1 Agent Group(s) is created with all tags contained in the agent
        And a dns policy pcap with tap_selector matching all tag(s) of the tap from an agent, counters metric_groups enabled, top_ecs, cardinality, top_qnames, top_ports, top_size, xact_times, quantiles, top_qtypes, top_rcodes metric_groups disabled and settings: {"require_version":"2.0"} is applied to the group
        And 1 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 1 policies are applied and all has status running
        And run mocked data dns_ipv4_tcp.pcap, dns_ipv4_udp.pcap, dns_ipv6_tcp.pcap, dns_ipv6_udp.pcap, dns_udp_mixed_rcode.pcap, dns_udp_tcp_random.pcap, ecs.pcap, dnssec.pcap, dhcpv6.pcap on the created virtual interface
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 240 seconds
    Then metrics must be correctly generated for dns-v2 handler
        And remove the agent .yaml generated on each scenario


@sanity @metric_groups @metrics_dns_v2 @root @mocked_interface @auto_provision
Scenario: dns handler with only top_qnames metric groups enabled (v2)
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink with default configuration type already exists
        And a mocked interface is configured with mtu: 65000 and ip: None
    When an agent(backend_type:pktvisor, settings: {"input_type":"pcap","iface":"mocked"}) is self-provisioned via a configuration file on port available with 1 agent tags and has status online. [Overwrite default: False. Paste only file: True]
        And pktvisor state is running
        And 1 Agent Group(s) is created with all tags contained in the agent
        And a dns policy pcap with tap_selector matching all tag(s) of the tap from an agent, top_qnames metric_groups enabled, top_ecs, cardinality, counters, top_ports, top_size, xact_times, quantiles, top_qtypes, top_rcodes metric_groups disabled and settings: {"require_version":"2.0"} is applied to the group
        And 1 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 1 policies are applied and all has status running
        And run mocked data dns_ipv4_tcp.pcap, dns_ipv4_udp.pcap, dns_ipv6_tcp.pcap, dns_ipv6_udp.pcap, dns_udp_mixed_rcode.pcap, dns_udp_tcp_random.pcap, ecs.pcap, dnssec.pcap, dhcpv6.pcap on the created virtual interface
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 240 seconds
    Then metrics must be correctly generated for dns-v2 handler
        And remove the agent .yaml generated on each scenario


@sanity @metric_groups @metrics_dns_v2 @root @mocked_interface @auto_provision
Scenario: dns handler with only quantiles metric groups enabled (v2)
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink with default configuration type already exists
        And a mocked interface is configured with mtu: 65000 and ip: None
    When an agent(backend_type:pktvisor, settings: {"input_type":"pcap","iface":"mocked"}) is self-provisioned via a configuration file on port available with 1 agent tags and has status online. [Overwrite default: False. Paste only file: True]
        And pktvisor state is running
        And 1 Agent Group(s) is created with all tags contained in the agent
        And a dns policy pcap with tap_selector matching all tag(s) of the tap from an agent, quantiles metric_groups enabled, top_ecs, cardinality, counters, top_qnames, top_ports, top_size, xact_times, top_qtypes, top_rcodes metric_groups disabled and settings: {"require_version":"2.0"} is applied to the group
        And 1 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 1 policies are applied and all has status running
        And run mocked data dns_ipv4_tcp.pcap, dns_ipv4_udp.pcap, dns_ipv6_tcp.pcap, dns_ipv6_udp.pcap, dns_udp_mixed_rcode.pcap, dns_udp_tcp_random.pcap, ecs.pcap, dnssec.pcap, dhcpv6.pcap on the created virtual interface
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 240 seconds
    Then metrics must be correctly generated for dns-v2 handler
        And remove the agent .yaml generated on each scenario


@sanity @metric_groups @metrics_dns_v2 @root @mocked_interface @auto_provision
Scenario: dns handler with only top_qtypes metric groups enabled (v2)
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink with default configuration type already exists
        And a mocked interface is configured with mtu: 65000 and ip: None
    When an agent(backend_type:pktvisor, settings: {"input_type":"pcap","iface":"mocked"}) is self-provisioned via a configuration file on port available with 1 agent tags and has status online. [Overwrite default: False. Paste only file: True]
        And pktvisor state is running
        And 1 Agent Group(s) is created with all tags contained in the agent
        And a dns policy pcap with tap_selector matching all tag(s) of the tap from an agent, top_qtypes metric_groups enabled, top_ecs, cardinality, counters, top_qnames, top_ports, top_size, xact_times, quantiles, top_rcodes metric_groups disabled and settings: {"require_version":"2.0"} is applied to the group
        And 1 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 1 policies are applied and all has status running
        And run mocked data dns_ipv4_tcp.pcap, dns_ipv4_udp.pcap, dns_ipv6_tcp.pcap, dns_ipv6_udp.pcap, dns_udp_mixed_rcode.pcap, dns_udp_tcp_random.pcap, ecs.pcap, dnssec.pcap, dhcpv6.pcap on the created virtual interface
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 240 seconds
    Then metrics must be correctly generated for dns-v2 handler
        And remove the agent .yaml generated on each scenario


@sanity @metric_groups @metrics_dns_v2 @root @mocked_interface @auto_provision
Scenario: dns handler with only top_rcodes metric groups enabled (v2)
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink with default configuration type already exists
        And a mocked interface is configured with mtu: 65000 and ip: None
    When an agent(backend_type:pktvisor, settings: {"input_type":"pcap","iface":"mocked"}) is self-provisioned via a configuration file on port available with 1 agent tags and has status online. [Overwrite default: False. Paste only file: True]
        And pktvisor state is running
        And 1 Agent Group(s) is created with all tags contained in the agent
        And a dns policy pcap with tap_selector matching all tag(s) of the tap from an agent, top_rcodes metric_groups enabled, top_ecs, cardinality, counters, top_qnames, top_ports, top_size, xact_times, quantiles, top_qtypes metric_groups disabled and settings: {"require_version":"2.0"} is applied to the group
        And 1 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 1 policies are applied and all has status running
        And run mocked data dns_ipv4_tcp.pcap, dns_ipv4_udp.pcap, dns_ipv6_tcp.pcap, dns_ipv6_udp.pcap, dns_udp_mixed_rcode.pcap, dns_udp_tcp_random.pcap, ecs.pcap, dnssec.pcap, dhcpv6.pcap on the created virtual interface
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 240 seconds
    Then metrics must be correctly generated for dns-v2 handler
        And remove the agent .yaml generated on each scenario


#### net v2.0

@sanity @metric_groups @metrics_net_v2 @root @mocked_interface
Scenario: net handler with default metric groups configuration (v2)
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink with default configuration type already exists
        And a mocked interface is configured with mtu: 65000 and ip: None
    When an agent(backend_type:pktvisor, settings: {"input_type":"pcap","iface":"mocked"}) is self-provisioned via a configuration file on port available with 1 agent tags and has status online. [Overwrite default: False. Paste only file: True]
        And pktvisor state is running
        And 1 Agent Group(s) is created with all tags contained in the agent
        And a net policy pcap with tap_selector matching all tag(s) of the tap from an agent, default metric_groups enabled, default metric_groups disabled and settings: {"require_version":"2.0"} is applied to the group
        And 1 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 1 policies are applied and all has status running
        And run mocked data dns_ipv4_tcp.pcap, dns_ipv4_udp.pcap, dns_ipv6_tcp.pcap, dns_ipv6_udp.pcap, dns_udp_mixed_rcode.pcap, dns_udp_tcp_random.pcap, ecs.pcap, ipfix.pcap, ecmp.pcap, ipfix.pcap, nf9.pcap, dnssec.pcap, dhcpv6.pcap on the created virtual interface
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 240 seconds
    Then metrics must be correctly generated for net-v2 handler
        And remove the agent .yaml generated on each scenario


@sanity @metric_groups @metrics_net_v2 @root @mocked_interface
Scenario: net handler with all metric groups enabled (v2)
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink with default configuration type already exists
        And a mocked interface is configured with mtu: 65000 and ip: None
    When an agent(backend_type:pktvisor, settings: {"input_type":"pcap","iface":"mocked"}) is self-provisioned via a configuration file on port available with 1 agent tags and has status online. [Overwrite default: False. Paste only file: True]
        And pktvisor state is running
        And 1 Agent Group(s) is created with all tags contained in the agent
        And a net policy pcap with tap_selector matching all tag(s) of the tap from an agent, all metric_groups enabled, none metric_groups disabled and settings: {"require_version":"2.0"} is applied to the group
        And 1 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 1 policies are applied and all has status running
        And run mocked data dns_ipv4_tcp.pcap, dns_ipv4_udp.pcap, dns_ipv6_tcp.pcap, dns_ipv6_udp.pcap, dns_udp_mixed_rcode.pcap, dns_udp_tcp_random.pcap, ecs.pcap, ipfix.pcap, ecmp.pcap, ipfix.pcap, nf9.pcap, dnssec.pcap, dhcpv6.pcap on the created virtual interface
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 240 seconds
    Then metrics must be correctly generated for net-v2 handler
        And remove the agent .yaml generated on each scenario


@sanity @metric_groups @metrics_net_v2 @root @mocked_interface
Scenario: net handler with all metric groups disabled (v2)
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink with default configuration type already exists
        And a mocked interface is configured with mtu: 65000 and ip: None
    When an agent(backend_type:pktvisor, settings: {"input_type":"pcap","iface":"mocked"}) is self-provisioned via a configuration file on port available with 1 agent tags and has status online. [Overwrite default: False. Paste only file: True]
        And pktvisor state is running
        And 1 Agent Group(s) is created with all tags contained in the agent
        And a net policy pcap with tap_selector matching all tag(s) of the tap from an agent, none metric_groups enabled, all metric_groups disabled and settings: {"require_version":"2.0"} is applied to the group
        And 1 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 1 policies are applied and all has status running
        And run mocked data dns_ipv4_tcp.pcap, dns_ipv4_udp.pcap, dns_ipv6_tcp.pcap, dns_ipv6_udp.pcap, dns_udp_mixed_rcode.pcap, dns_udp_tcp_random.pcap, ecs.pcap, ipfix.pcap, ecmp.pcap, ipfix.pcap, nf9.pcap, dnssec.pcap, dhcpv6.pcap on the created virtual interface
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
    Then metrics must be correctly generated for net-v2 handler
        And remove the agent .yaml generated on each scenario


@sanity @metric_groups @metrics_net_v2 @root @mocked_interface
Scenario: net handler with only cardinality metric groups enabled (v2)
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink with default configuration type already exists
        And a mocked interface is configured with mtu: 65000 and ip: None
    When an agent(backend_type:pktvisor, settings: {"input_type":"pcap","iface":"mocked"}) is self-provisioned via a configuration file on port available with 1 agent tags and has status online. [Overwrite default: False. Paste only file: True]
        And pktvisor state is running
        And 1 Agent Group(s) is created with all tags contained in the agent
        And a net policy pcap with tap_selector matching all tag(s) of the tap from an agent, cardinality metric_groups enabled, quantiles, counters, top_geo, top_ips metric_groups disabled and settings: {"require_version":"2.0"} is applied to the group
        And 1 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 1 policies are applied and all has status running
        And run mocked data dns_ipv4_tcp.pcap, dns_ipv4_udp.pcap, dns_ipv6_tcp.pcap, dns_ipv6_udp.pcap, dns_udp_mixed_rcode.pcap, dns_udp_tcp_random.pcap, ecs.pcap, ipfix.pcap, ecmp.pcap, ipfix.pcap, nf9.pcap, dnssec.pcap, dhcpv6.pcap on the created virtual interface
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 240 seconds
    Then metrics must be correctly generated for net-v2 handler
        And remove the agent .yaml generated on each scenario


@sanity @metric_groups @metrics_net_v2 @root @mocked_interface
Scenario: net handler with only counters metric groups enabled (v2)
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink with default configuration type already exists
        And a mocked interface is configured with mtu: 65000 and ip: None
    When an agent(backend_type:pktvisor, settings: {"input_type":"pcap","iface":"mocked"}) is self-provisioned via a configuration file on port available with 1 agent tags and has status online. [Overwrite default: False. Paste only file: True]
        And pktvisor state is running
        And 1 Agent Group(s) is created with all tags contained in the agent
        And a net policy pcap with tap_selector matching all tag(s) of the tap from an agent, counters metric_groups enabled, quantiles, cardinality, top_geo, top_ips metric_groups disabled and settings: {"require_version":"2.0"} is applied to the group
        And 1 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 1 policies are applied and all has status running
        And run mocked data dns_ipv4_tcp.pcap, dns_ipv4_udp.pcap, dns_ipv6_tcp.pcap, dns_ipv6_udp.pcap, dns_udp_mixed_rcode.pcap, dns_udp_tcp_random.pcap, ecs.pcap, ipfix.pcap, ecmp.pcap, ipfix.pcap, nf9.pcap, dnssec.pcap, dhcpv6.pcap on the created virtual interface
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 240 seconds
    Then metrics must be correctly generated for net-v2 handler
        And remove the agent .yaml generated on each scenario


@sanity @metric_groups @metrics_net_v2 @root @mocked_interface
Scenario: net handler with only top_geo metric groups enabled (v2)
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink with default configuration type already exists
        And a mocked interface is configured with mtu: 65000 and ip: None
    When an agent(backend_type:pktvisor, settings: {"input_type":"pcap","iface":"mocked"}) is self-provisioned via a configuration file on port available with 1 agent tags and has status online. [Overwrite default: False. Paste only file: True]
        And pktvisor state is running
        And 1 Agent Group(s) is created with all tags contained in the agent
        And a net policy pcap with tap_selector matching all tag(s) of the tap from an agent, top_geo metric_groups enabled, quantiles, counters, cardinality, top_ips metric_groups disabled and settings: {"require_version":"2.0"} is applied to the group
        And 1 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 1 policies are applied and all has status running
        And run mocked data dns_ipv4_tcp.pcap, dns_ipv4_udp.pcap, dns_ipv6_tcp.pcap, dns_ipv6_udp.pcap, dns_udp_mixed_rcode.pcap, dns_udp_tcp_random.pcap, ecs.pcap, ipfix.pcap, ecmp.pcap, ipfix.pcap, nf9.pcap, dnssec.pcap, dhcpv6.pcap on the created virtual interface
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 240 seconds
    Then metrics must be correctly generated for net-v2 handler
        And remove the agent .yaml generated on each scenario


@sanity @metric_groups @metrics_net_v2 @root @mocked_interface
Scenario: net handler with only top_ips metric groups enabled (v2)
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink with default configuration type already exists
        And a mocked interface is configured with mtu: 65000 and ip: None
    When an agent(backend_type:pktvisor, settings: {"input_type":"pcap","iface":"mocked"}) is self-provisioned via a configuration file on port available with 1 agent tags and has status online. [Overwrite default: False. Paste only file: True]
        And pktvisor state is running
        And 1 Agent Group(s) is created with all tags contained in the agent
        And a net policy pcap with tap_selector matching all tag(s) of the tap from an agent, top_ips metric_groups enabled, quantiles, counters, top_geo, cardinality metric_groups disabled and settings: {"require_version":"2.0"} is applied to the group
        And 1 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 1 policies are applied and all has status running
        And run mocked data dns_ipv4_tcp.pcap, dns_ipv4_udp.pcap, dns_ipv6_tcp.pcap, dns_ipv6_udp.pcap, dns_udp_mixed_rcode.pcap, dns_udp_tcp_random.pcap, ecs.pcap, ipfix.pcap, ecmp.pcap, ipfix.pcap, nf9.pcap, dnssec.pcap, dhcpv6.pcap on the created virtual interface
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 240 seconds
    Then metrics must be correctly generated for net-v2 handler
        And remove the agent .yaml generated on each scenario


@sanity @metric_groups @metrics_net_v2 @root @mocked_interface
Scenario: net handler with only quantiles metric groups enabled (v2)
    Given the Orb user has a registered account
        And the Orb user logs in
        And that a sink with default configuration type already exists
        And a mocked interface is configured with mtu: 65000 and ip: None
    When an agent(backend_type:pktvisor, settings: {"input_type":"pcap","iface":"mocked"}) is self-provisioned via a configuration file on port available with 1 agent tags and has status online. [Overwrite default: False. Paste only file: True]
        And pktvisor state is running
        And 1 Agent Group(s) is created with all tags contained in the agent
        And a net policy pcap with tap_selector matching all tag(s) of the tap from an agent, quantiles metric_groups enabled, top_ips, counters, top_geo, cardinality metric_groups disabled and settings: {"require_version":"2.0"} is applied to the group
        And 1 dataset(s) have validity valid and 0 have validity invalid in 30 seconds
        And this agent's heartbeat shows that 1 groups are matching the agent
        And this agent's heartbeat shows that 1 policies are applied and all has status running
        And run mocked data dns_ipv4_tcp.pcap, dns_ipv4_udp.pcap, dns_ipv6_tcp.pcap, dns_ipv6_udp.pcap, dns_udp_mixed_rcode.pcap, dns_udp_tcp_random.pcap, ecs.pcap, ipfix.pcap, ecmp.pcap, ipfix.pcap, nf9.pcap, dnssec.pcap, dhcpv6.pcap on the created virtual interface
        And the container logs contain the message "policy applied successfully" referred to each policy within 30 seconds
        And the container logs that were output after all policies have been applied contain the message "scraped metrics for policy" referred to each applied policy within 180 seconds
        And referred sink must have active state on response within 240 seconds
    Then metrics must be correctly generated for net-v2 handler
        And remove the agent .yaml generated on each scenario