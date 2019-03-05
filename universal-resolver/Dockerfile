FROM golang

ENV PORT=8080

COPY main.go $GOPATH/src/github.com/ockam/universal-resolver/

ENV PATH=$PATH:$GOPATH/src/github.com/ockam/universal-resolver/

WORKDIR $GOPATH/src/github.com/ockam/universal-resolver/
RUN go get ./
RUN go build -o main

EXPOSE 8080
ENTRYPOINT ["main"]

