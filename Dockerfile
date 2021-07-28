#############
# Build
#############
FROM golang:alpine as build

WORKDIR /showEnv

COPY . /showEnv

RUN go build -o /usr/bin/showenv

CMD ["/usr/bin/showenv"]


#############
# Release
#############
FROM alpine:latest as release

RUN addgroup -g 1000 golang; \
    adduser -H -s /bin/false -u 1000 -G golang -S -D golang

COPY --from=build /usr/bin/showenv /usr/bin

EXPOSE 8080

CMD ["/usr/bin/showenv"]
