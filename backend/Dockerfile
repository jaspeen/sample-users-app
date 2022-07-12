FROM golang:1.18-alpine as build-stage

COPY . /build
RUN cd /build && go build cmd/server/server.go


FROM alpine
COPY --from=build-stage /build/server /server
ENTRYPOINT /server