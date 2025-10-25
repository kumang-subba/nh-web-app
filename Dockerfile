FROM --platform=$BUILDPLATFORM golang:alpine AS build 
WORKDIR /build
RUN apk add --no-cache make

COPY go.mod go.sum Makefile ./
RUN make deps

COPY . .
RUN make build

FROM alpine
COPY --from=build /build/hello /hello
EXPOSE 3000
ENTRYPOINT ["/hello"]
