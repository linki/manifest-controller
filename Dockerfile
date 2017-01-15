FROM golang:1.7.4-alpine

# make kubectl available in the container
ADD https://storage.googleapis.com/kubernetes-release/release/v1.5.2/bin/linux/amd64/kubectl /
RUN chmod +x /kubectl

# add code and build binary
COPY . /go/src/github.com/linki/manifest-controller
RUN go install -v github.com/linki/manifest-controller

CMD ["--help"]
ENTRYPOINT ["/go/bin/manifest-controller"]
