FROM golang:1.13
WORKDIR /go/src/github.com/morfien101/asg-lifecyle-hook-manager
COPY go.* *.go ./
COPY ./hookmanager ./hookmanager/
RUN mkdir /output && CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /output/lifecyclehook .

FROM alpine:latest  
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=0 /output/lifecyclehook .
ENTRYPOINT ["./lifecyclehook"]
CMD ["-h"]
