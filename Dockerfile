FROM alpine

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
    curl -SL -k https://github.com/shadowsocks/shadowsocks-libev/releases/download/v3.0.5/shadowsocks-libev-3.0.5.tar.gz | tar xz && \
    cd shadowsocks-libev-3.0.5 && \
    ./configure --prefix=/usr --disable-documentation && \
    make install

RUN apk add --no-cache ca-certificates easy-rsa mongodb openvpn redis ufw@testing && \
    mkdir -p /data/db && \
    wget -c https://bootstrap.pypa.io/get-pip.py -O /tmp/get-pip.py && \
    python /tmp/get-pip.py && \
    pip install --no-cache-dir falcon gunicorn pymongo raven redis requests speedtest_cli && \
    rm -rf /tmp/* /var/tmp/* /var/cache/apk/* /var/cache/distfiles/* /root/.cache .wget-hsts

EXPOSE 4200/tcp
EXPOSE 4200/udp
EXPOSE 1194/udp 3000

ENTRYPOINT ["sh", "/root/run.sh"]
