FROM golang:latest as build
WORKDIR /build
COPY ./go.mod ./go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /vigor -a -ldflags '-extldflags "-static"' . 

FROM gcr.io/distroless/base
COPY --from=build /vigor /vigor
CMD ["/vigor"]