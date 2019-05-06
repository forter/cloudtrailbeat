FROM golang:latest AS builder
MAINTAINER Forter RnD

WORKDIR /go/src/github.com/forter/cloudtrailbeat
RUN mkdir -p /config
RUN apt-get update && \
    apt-get install -y \
    git gcc g++ binutils make python2.7 python-pip && \
    pip install virtualenv
RUN mkdir -p ${GOPATH}/src/github.com/elastic && git clone https://github.com/elastic/beats ${GOPATH}/src/github.com/elastic/beats
COPY . /go/src/github.com/forter/cloudtrailbeat/
RUN make setup
RUN go run mage.go build
# ---

FROM scratch
COPY --from=build /go/bin/cloudtrailbeat /cloudtrailbeat
VOLUME  /config/config.yml
ENTRYPOINT [ "/cloudtrailbeat" ]
CMD [ "--help" ]