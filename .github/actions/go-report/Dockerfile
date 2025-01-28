FROM golang:1.22-alpine

LABEL maintainer="netboxlabs"
LABEL version="1.0.0"

RUN mkdir -p /workspace

COPY ./entrypoint.sh /entrypoint.sh
COPY ./github-commenter_linux_amd64 /github-commenter_linux_amd64

RUN apk add git make curl jq && \
mv /github-commenter_linux_amd64 /github-commenter && \
chmod +x /github-commenter /entrypoint.sh

RUN git clone https://github.com/gojp/goreportcard.git && \
cd goreportcard && \
make install && \
go install ./cmd/goreportcard-cli && \
mv /go/bin/goreportcard-cli /usr/local/bin/goreportcard-cli

WORKDIR /workspace

CMD ["/entrypoint.sh"]
