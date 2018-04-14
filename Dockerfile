FROM golang:1.9 as builder

# build directories
RUN mkdir -p /go/src/github.com/rerorero/myghome
WORKDIR /go/src/github.com/rerorero/myghome
ADD . .

# Build
RUN go get github.com/golang/dep/...
RUN dep ensure
RUN CGO_ENABLED=0 GOOS=linux go build -a -tags netgo -installsuffix netgo --ldflags '-extldflags "-static"' -o /myghome ./main.go

# runner container
FROM alpine:latest
COPY --from=builder /myghome /bin/myghome

CMD /bin/myghome --assistant ${GOOGLE_IP} --port ${PORT}
