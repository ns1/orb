ARG PKTVISOR_TAG=latest-develop
FROM golang:1.22-alpine AS builder

WORKDIR /go/src/github.com/orbcommunity/orb
COPY go.mod .
RUN go mod tidy
COPY . .
RUN apk update && apk add make build-base git
RUN mkdir /tmp/build && CGO_ENABLED=1 make agent_bin && mv build/orb-agent /tmp/build/orb-agent

FROM orbcommunity/pktvisor:${PKTVISOR_TAG}

RUN mkdir /opt/orb
COPY --from=builder /tmp/build/orb-agent /usr/local/bin/orb-agent
COPY --from=builder /go/src/github.com/orbcommunity/orb/agent/docker/agent_default.yaml /opt/orb/agent_default.yaml
COPY --from=builder /go/src/github.com/orbcommunity/orb/agent/docker/orb-agent-entry.sh /usr/local/bin/orb-agent-entry.sh
COPY --from=builder /go/src/github.com/orbcommunity/orb/agent/docker/run-agent.sh /run-agent.sh

RUN chmod a+x /run-agent.sh

ENTRYPOINT [ "/usr/local/bin/orb-agent-entry.sh" ]
