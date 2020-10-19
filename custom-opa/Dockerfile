FROM golang:1.15 as build

WORKDIR /go/src/opa
ADD . /go/src/opa

RUN go build -o /go/bin/opa++ .

FROM gcr.io/distroless/base-debian10
COPY --from=build "/go/bin/opa++" /
ENTRYPOINT ["/opa++"]
