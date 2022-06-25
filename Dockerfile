FROM scratch
COPY go-generator-cli /
ENTRYPOINT ["/go-generator-cli"]