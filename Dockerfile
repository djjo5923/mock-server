FROM golang:1.14-alpine AS build

WORKDIR /src/
COPY main.go go.* /src/
RUN CGO_ENABLED=0 go build -o /bin/mock-server

FROM scratch
COPY --from=build /bin/mock-server /bin/mock-server
ENTRYPOINT ["/bin/mock-server"]
