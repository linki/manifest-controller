FROM registry.opensource.zalan.do/stups/alpine:3.4-2

# make kubectl available in the container
RUN curl -o /kubectl https://storage.googleapis.com/kubernetes-release/release/v1.3.6/bin/linux/amd64/kubectl
RUN chmod +x /kubectl

# add scm-source
ADD scm-source.json /

# add binary
ADD build/linux/barn /

ENTRYPOINT ["/barn"]
