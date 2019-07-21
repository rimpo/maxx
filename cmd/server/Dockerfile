FROM golang:1.12.7-alpine

RUN mkdir -p /maxx/bin
ADD . /maxx
WORKDIR /maxx

RUN GOBIN=/maxx/bin go install ./...

#Environment
ENV PORT=8443
ENV CERT_FILE=./cert/server.crt
ENV KEY_FILE=./cert/server.key

EXPOSE ${PORT}

CMD ["./bin/server", "--port=${PORT}", "--cert_file=${CERT_FILE}", "--key_file=${KEY_FILE}"]