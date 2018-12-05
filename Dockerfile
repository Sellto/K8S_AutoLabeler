FROM alpine
RUN apk --no-cache add gettext ca-certificates openssl \
    && wget https://github.com/Yelp/dumb-init/releases/download/v1.2.0/dumb-init_1.2.0_amd64 -O /usr/local/bin/dumb-init \
    && wget https://storage.googleapis.com/kubernetes-release/release/v1.11.2/bin/linux/amd64/kubectl -O /usr/local/bin/kubectl \
    && chmod a+x /usr/local/bin/kubectl /usr/local/bin/dumb-init \
    && apk --no-cache del ca-certificates openssl
RUN apk add eudev
RUN apk add eudev-dev
RUN apk add go
RUN apk add git
RUN apk add musl-dev
RUN apk add linux-headers
RUN go get -v gopkg.in/yaml.v2
RUN go get -v github.com/jochenvg/go-udev
VOLUME /dev /dev
VOLUME /home/selltom/TFE/K8S_Files/3S /source
CMD udevd --daemon --debug && sh
