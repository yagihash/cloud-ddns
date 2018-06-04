FROM ubuntu

ARG SA_FILE=no_such_file
ENV DEBIAN_FRONTEND=noninteractive
ENV LANG=en_US.UTF-8

RUN groupadd -g 1000 worker && useradd -g worker -u 500 --create-home worker

RUN apt-get update
RUN apt-get -y install software-properties-common

RUN set -xe \
    && apt-get update \
    && apt-get install -y --no-install-recommends sudo cron curl \
    && rm -rf /var/lib/apt/lists/*

RUN apt-get -y dist-upgrade
RUN curl -sL https://deb.nodesource.com/setup_10.x | sudo -E sh -
RUN apt-get -y install nodejs

COPY . /home/worker/
COPY $SA_FILE /tmp/sa.json
RUN cd /home/worker/ && npm install
RUN chown -R worker:worker /home/worker/

RUN echo '0-59/1 * * * * worker /usr/bin/node /home/worker/main.js' >> /etc/cron.d/ddns

ENTRYPOINT cron -f
