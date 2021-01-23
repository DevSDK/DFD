FROM golang:1.15.7 as builder
WORKDIR /go/src/DFD
COPY . .
RUN go mod download
RUN go build 
RUN chmod 755 DFD
ENTRYPOINT [ "./DFD" ]