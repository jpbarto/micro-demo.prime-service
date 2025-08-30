FROM golang:alpine

WORKDIR /opt
COPY src/* /opt/prime-service/

WORKDIR /opt/prime-service
RUN go build -o ./bin/prime-service
# Make sure the binary is executable
RUN chmod +x ./bin/prime-service

# Set the entrypoint to run the binary
ENTRYPOINT ["/opt/prime-service/bin/prime-service"]