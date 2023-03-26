FROM golang:latest as build

WORKDIR /build
COPY . /build

RUN CGO_ENABLED=0 go build -o /vigor -ldflags="-extldflags -static" /build/

FROM gcr.io/distroless/base

WORKDIR /app
COPY --from=build /vigor /app/vigor

CMD ["/app/vigor"]