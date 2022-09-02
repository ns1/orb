#!/usr/bin/env bash
#
# entry point for orb-agent
#

# orb agent binary location. by default, matches orb-agent container (see Dockerfile)
orb_agent_bin="${ORB_AGENT_BIN:-/usr/local/bin/orb-agent}"

# support generating API and MQTT addresses with one host name in ORB_CLOUD_ADDRESS
if [[ -n "${ORB_CLOUD_ADDRESS}" ]]; then
  ORB_CLOUD_API_ADDRESS="https://${ORB_CLOUD_ADDRESS}"
  ORB_CLOUD_MQTT_ADDRESS="tls://${ORB_CLOUD_ADDRESS}:8883"
  export ORB_CLOUD_API_ADDRESS ORB_CLOUD_MQTT_ADDRESS
fi

# support generating simple default pktvisor PCAP taps

tmpfile=$(mktemp /tmp/orb-agent-pktvisor-conf.XXXXXX)
trap 'rm -f "$tmpfile"' EXIT

#Add defaults
(
cat <<END
version: "1.0"

visor:
  taps:
END
) > "$tmpfile"

# NetFlow tap
if [ "${PKTVISOR_NETFLOW_BIND_ADDRESS}" = '' ]; then
  PKTVISOR_NETFLOW_BIND_ADDRESS='0.0.0.0'
fi
if [ "${PKTVISOR_NETFLOW_PORT_DEFAULT}" = '' ]; then
  PKTVISOR_NETFLOW_PORT_DEFAULT='2055'
fi
if [ "${PKTVISOR_NETFLOW}" = 'true' ]; then
(
cat <<END
    default_netflow:
      input_type: flow
      config:
        flow_type: netflow
        port: "$PKTVISOR_NETFLOW_PORT_DEFAULT"
        bind: "$PKTVISOR_NETFLOW_BIND_ADDRESS"
END
) >> "$tmpfile"

  export ORB_BACKENDS_PKTVISOR_CONFIG_FILE="$tmpfile"
fi

# SFlow tap
if [ "${PKTVISOR_SFLOW_BIND_ADDRESS}" = '' ]; then
  PKTVISOR_SFLOW_BIND_ADDRESS='0.0.0.0'
fi
if [ "${PKTVISOR_SFLOW_PORT_DEFAULT}" = '' ]; then
  PKTVISOR_SFLOW_PORT_DEFAULT='6343'
fi
if [ "${PKTVISOR_SFLOW}" = 'true' ]; then
(
cat <<END
    default_sflow:
      input_type: flow
      config:
        flow_type: sflow
        port: "$PKTVISOR_SFLOW_PORT_DEFAULT"
        bind: "$PKTVISOR_SFLOW_BIND_ADDRESS"
END
) >> "$tmpfile"

  export ORB_BACKENDS_PKTVISOR_CONFIG_FILE="$tmpfile"
fi

# DNS TAP
if [ "${PKTVISOR_DNSTAP_BIND_ADDRESS}" = '' ]; then
  PKTVISOR_DNSTAP_BIND_ADDRESS='0.0.0.0'
fi
if [ "${PKTVISOR_DNSTAP_PORT_DEFAULT}" = '' ]; then
  PKTVISOR_DNSTAP_PORT_DEFAULT='6000'
fi

if [ "${PKTVISOR_DNSTAP}" = 'true' ]; then
(
cat <<END
    default_dnstap:
      input_type: dnstap
      config:
        tcp: "${PKTVISOR_DNSTAP_BIND_ADDRESS}:${PKTVISOR_DNSTAP_PORT_DEFAULT}"

END
) >> "$tmpfile"

  export ORB_BACKENDS_PKTVISOR_CONFIG_FILE="$tmpfile"
fi

# simplest: specify just interface, creates tap named "default_pcap"
# PKTVISOR_PCAP_IFACE_DEFAULT=en0
# special case: if the iface is "mock", then use "mock" pcap source
if [ "$PKTVISOR_PCAP_IFACE_DEFAULT" = 'mock' ]; then
  MAYBE_MOCK='pcap_source: mock'
fi
if [[ -n "${PKTVISOR_PCAP_IFACE_DEFAULT}" ]]; then
(
cat <<END
    default_pcap:
      input_type: pcap
      config:
        iface: "$PKTVISOR_PCAP_IFACE_DEFAULT"
        $MAYBE_MOCK
END
) >>"$tmpfile"

  export ORB_BACKENDS_PKTVISOR_CONFIG_FILE="$tmpfile"
fi

# or specify pair of TAPNAME:IFACE
# TODO allow multiple, split on comma
# PKTVISOR_PCAP_IFACE_TAPS=default_pcap:en0
if [ $# -eq 0 ]; then
  exec "$orb_agent_bin" run
else
  exec "$orb_agent_bin" "$@"
fi
