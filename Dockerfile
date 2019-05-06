FROM golang:stretch AS builder
MAINTAINER Forter RnD

WORKDIR /go/src/github.com/forter/cloudtrailbeat
RUN mkdir -p /config
RUN apt-get update && \
    apt-get install -y \
    git gcc g++ binutils make
RUN mkdir -p ${GOPATH}/src/github.com/elastic && git clone https://github.com/elastic/beats ${GOPATH}/src/github.com/elastic/beats
COPY . /go/src/github.com/forter/cloudtrailbeat/
RUN make
RUN chmod +x cloudtrailbeat
# ---

FROM golang:stretch
COPY --from=builder /go/src/github.com/forter/cloudtrailbeat/cloudtrailbeat /bin/cloudtrailbeat
VOLUME  /config/beat.yml
ENTRYPOINT [ "/bin/cloudtrailbeat" ]
CMD [ "-c /config/beat.yml" ]