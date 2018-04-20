FROM alpine
MAINTAINER Evil Thanos "evilthanos@protonmail.com"

ENV SS_LIBEV_VERSION=3.0.5
ENV SS_DIR=shadowsocks-libev-${SS_LIBEV_VERSION}
ENV SENT_ENV=PROD

ADD sentinel /root/sentinel

ADD app.py server.py run.sh /root/

ADD config /etc

RUN set -ex && \
    echo '@testing http://nl.alpinelinux.org/alpine/edge/testing' >> /etc/apk/repositories && \
    apk add --no-cache libcrypto1.0 \
                        libev \
                        libsodium \
                        mbedtls \
                        pcre \
                        udns \
    && apk add --no-cache --virtual TMP autoconf \
                             automake \
                             build-base \
                             curl \
                             gettext-dev \
                             libev-dev \
                             libsodium-dev \
                             libtool \
                             linux-headers \
                             mbedtls-dev \
                             openssl-dev \
                             pcre-dev \
                             tar \
                             udns-dev && \
    cd /tmp && \
    curl -SL -k https://github.com/shadowsocks/shadowsocks-libev/releases/download/v${SS_LIBEV_VERSION}/shadowsocks-libev-${SS_LIBEV_VERSION}.tar.gz | tar xz && \
    cd $SS_DIR && \
    ./configure --prefix=/usr --disable-documentation && \
    make install

RUN apk add --no-cache ca-certificates easy-rsa mongodb openvpn redis ufw@testing && \
    mkdir -p /data/db && \
    wget -c https://bootstrap.pypa.io/get-pip.py -O /tmp/get-pip.py && \
    python /tmp/get-pip.py && \
    pip install --no-cache-dir falcon gunicorn pymongo raven redis requests speedtest_cli && \
    rm -rf /tmp/* /var/tmp/* /var/cache/apk/* /var/cache/distfiles/* /root/.cache .wget-hsts

COPY ./docker-entrypoint.sh /usr/local/bin/
RUN ln -s usr/local/bin/docker-entrypoint.sh /entrypoint.sh && \
    chmod +x usr/local/bin/docker-entrypoint.sh

EXPOSE 8388/tcp
EXPOSE 8388/udp
EXPOSE 1194/udp 3000

ENTRYPOINT ["sh", "/root/run.sh"]
