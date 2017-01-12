FROM alpine:3.5

# make kubectl available in the container
ADD https://storage.googleapis.com/kubernetes-release/release/v1.3.6/bin/linux/amd64/kubectl /
RUN chmod +x /kubectl

# add binary
ADD build/linux/manifest-controller /

ENTRYPOINT ["/manifest-controller"]
