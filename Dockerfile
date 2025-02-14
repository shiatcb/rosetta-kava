ARG build_image=golang:1.17.6-bullseye
FROM $build_image as kava-rosetta-build

RUN apt-get update \
      && apt-get install -y git make gcc \
      && rm -rf /var/lib/apt/lists/*

ARG kava_node_version=v0.16.1
ARG kava_rosetta_version=v0.0.3
ENV KAVA_NODE_VERSION=$kava_node_version
ENV KAVA_ROSETTA_VERSION=$kava_rosetta_version

RUN mkdir /app
WORKDIR /app

RUN git clone https://github.com/Kava-Labs/kava.git \
      && cd kava \
      && git checkout $KAVA_NODE_VERSION \
      && make install

RUN git clone https://github.com/Kava-Labs/rosetta-kava.git \
      && cd rosetta-kava \
      && git checkout $KAVA_ROSETTA_VERSION \
      && make install

FROM ubuntu:20.04

RUN apt-get update \
      && apt-get install -y supervisor curl \
      && rm -rf /var/lib/apt/lists/*

RUN mkdir /app \
      && mkdir /app/bin
WORKDIR /app

ENV PATH=$PATH:/app/bin

# copy build binaries from build environemtn
COPY --from=kava-rosetta-build /go/bin/kava /app/bin/kava
COPY --from=kava-rosetta-build /go/bin/rosetta-kava /app/bin/rosetta-kava

# copy config templates to automate setup
COPY --from=kava-rosetta-build /app/rosetta-kava/examples /app/templates

# copy scripts to run services
COPY --from=kava-rosetta-build /app/rosetta-kava/conf/start-services.sh /app/bin/start-services.sh
COPY --from=kava-rosetta-build /app/rosetta-kava/conf/kill-supervisord.sh /app/bin/kill-supervisord.sh
COPY --from=kava-rosetta-build /app/rosetta-kava/conf/supervisord.conf /etc/supervisor/conf.d/rosetta-kava.conf

ENV KAVA_RPC_URL=http://localhost:26657

CMD ["/app/bin/start-services.sh"]
