FROM ubuntu:20.04

LABEL maintainer="netboxlabs"
LABEL version="1.0.0"

RUN mkdir -p /workspace

WORKDIR /workspace

COPY ./entrypoint.sh /entrypoint.sh

RUN chmod a+x /entrypoint.sh

RUN apt-get update \
    && apt-get install jq git curl -y

ENTRYPOINT ["/entrypoint.sh"]
