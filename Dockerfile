FROM debian:buster
LABEL maintainer="Ruben Homs <ruben@homs.codes>"

USER root

# Set Environment Variables
ENV DEBIAN_FRONTEND noninteractive

ARG OPENSIPS_VERSION=3.0
ARG OPENSIPS_BUILD=releases

#install basic components
RUN apt update -qq && apt install -y gnupg2 ca-certificates

#add keyserver, repository
RUN apt-key adv --keyserver keyserver.ubuntu.com --recv-keys 049AD65B
RUN echo "deb https://apt.opensips.org buster ${OPENSIPS_VERSION}-${OPENSIPS_BUILD}" >/etc/apt/sources.list.d/opensips.list

RUN apt update -qq && apt install -y opensips curl net-tools procps

RUN apt-get -y install opensips-http-modules

RUN rm -rf /var/lib/apt/lists/*

EXPOSE 5060/udp
EXPOSE 8888/tcp

COPY run.sh /run.sh
COPY opensips.cfg /etc/opensips/opensips.cfg

ENTRYPOINT ["/run.sh"]
