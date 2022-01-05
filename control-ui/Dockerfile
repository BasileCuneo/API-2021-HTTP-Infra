from golang:1.16-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

copy templates/* ./templates/

COPY *.go ./

RUN go build -o /control

ENTRYPOINT ["/control"]