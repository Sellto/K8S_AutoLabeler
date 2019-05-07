FROM golang:1.12.4-alpine as dev
RUN apk add git
RUN apk add wget
RUN apk add openssl
RUN apk add git
RUN apk add musl-dev
RUN apk add linux-headers
RUN apk add ca-certificates
RUN apk add eudev
RUN apk add eudev-dev
RUN apk add build-base
RUN wget https://storage.googleapis.com/kubernetes-release/release/v1.13.4/bin/linux/amd64/kubectl -O /usr/local/bin/kubectl
RUN chmod a+x /usr/local/bin/kubectl
RUN go get -v gopkg.in/yaml.v2
RUN go get -v github.com/jochenvg/go-udev
ADD AutoLabeler.go AutoLabeler.go
RUN go build AutoLabeler.go
CMD sh

FROM alpine
COPY --from=dev /go/AutoLabeler .
COPY --from=dev /usr/local/bin/kubectl /usr/local/bin/kubectl
RUN apk add eudev git
CMD udevd --daemon && udevadm control --reload-rules && udevadm trigger && ./AutoLabeler
