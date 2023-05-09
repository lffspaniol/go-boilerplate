FROM golang:1.20-bullseye as dependencies
WORKDIR  /go/src/app

COPY go.mod .
COPY go.sum .

RUN go mod download -x

FROM dependencies as development
ENV CGO_ENABLED=1

RUN go install -v github.com/cosmtrek/air@latest&& \
    go install -v github.com/go-delve/delve/cmd/dlv@latest

EXPOSE 8080
EXPOSE 40000
CMD [ "air" ]

FROM dependencies as build
ENV CGO_ENABLED=0

COPY . .
RUN CGO_ENABLED=0 GOARCH=amd64 go build -ldflags "-w -s" -o /bin/main ./cmd/server/main.go

FROM gcr.io/distroless/static-debian11 AS production
COPY --from=build /bin/main /bin/main
ENV GIN_MODE=release
EXPOSE 8080
CMD ["/bin/main"]