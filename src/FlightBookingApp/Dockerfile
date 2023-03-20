FROM golang:alpine as build_container
WORKDIR /app
COPY ./src/go.mod .
COPY ./src/go.sum .
RUN go mod download
COPY src .
RUN go build -o server

FROM alpine
COPY --from=build_container /app/server /usr/bin
EXPOSE 4200
ENTRYPOINT ["server"]
