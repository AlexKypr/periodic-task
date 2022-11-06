#Build stage
FROM golang:1.19.0 AS builder
WORKDIR /server
COPY . .
RUN go build -o main main.go

#Run stage
FROM golang
WORKDIR /server
COPY --from=builder /server/main .

EXPOSE 8080
ENTRYPOINT ["/server/main"]
#CMD ["/server/main"]