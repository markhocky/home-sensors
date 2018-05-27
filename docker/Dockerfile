# STEP 1 build executable binary
FROM golang:alpine as builder
COPY . $GOPATH/src/home-sensors/
WORKDIR $GOPATH/src/home-sensors/
#get dependancies
RUN go get -d -v
#build the binary
RUN go build -o /go/bin/sensors
# STEP 2 build a small image
# start from scratch
FROM scratch
# Copy our static executable
COPY --from=builder /go/bin/sensors /go/bin/sensors
ENTRYPOINT ["/go/bin/sensors"]