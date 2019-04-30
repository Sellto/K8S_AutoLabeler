apk add gettext \
        openssl \
        go \
        git \
        musl-dev \
        linux-headers \
        ca-certificates \
        eudev \
        eudev-dev
wget https://github.com/Yelp/dumb-init/releases/download/v1.2.0/dumb-init_1.2.0_amd64 -O /usr/local/bin/dumb-init
wget https://storage.googleapis.com/kubernetes-release/release/v1.13.4/bin/linux/amd64/kubectl -O /usr/local/bin/kubectl
chmod a+x /usr/local/bin/kubectl /usr/local/bin/dumb-init
go get -v gopkg.in/yaml.v2
go get -v github.com/jochenvg/go-udev
go build AutoLabeler.go
go clean -i gopkg.in/yaml.v2
go clean -i github.com/jochenvg/go-udev
apk del gettext \
        openssl \
        go \
        git \
        musl-dev \
        linux-headers \
        eudev-dev
