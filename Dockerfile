FROM ubuntu:bionic as builder

RUN mkdir -p /opt/resource

FROM cloudfoundry/run:tiny
COPY --from=builder /opt /opt

COPY ./bin/* /opt/resource/
