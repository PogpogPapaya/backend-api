FROM golang:alpine as server-builder
WORKDIR /app
COPY ./go.mod ./
COPY ./go.sum ./
RUN go mod download
COPY ./ ./
RUN go build -o ./server ./main.go

FROM alpine:latest
WORKDIR /app
COPY --from=server-builder /app/server ./
CMD ["/app/server"]
