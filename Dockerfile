FROM golang:onbuild

RUN mkdir /app
ADD . /app/