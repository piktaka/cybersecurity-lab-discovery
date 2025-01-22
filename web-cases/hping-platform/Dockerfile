FROM ubuntu:latest

# Install dependencies
RUN apt update && apt install -y wget

# Download and install Go
RUN wget https://golang.org/dl/go1.23.4.linux-amd64.tar.gz && \
    tar -C /usr/local -xvzf go1.23.4.linux-amd64.tar.gz && \
    rm go1.23.4.linux-amd64.tar.gz

# Setup environment variables
ENV GOPATH=/opt/go
ENV GOROOT=/usr/local/go
ENV PATH=$PATH:$GOPATH/bin:$GOROOT/bin
ENV GO111MODULE=auto

# Create platform directory and copy files
RUN mkdir -p /opt/hping-platform /opt/hping-platform-app /usr/local/go
COPY . /opt/hping-platform

WORKDIR /opt/hping-platform

# Install Go modules and build the project
RUN apt install -y gcc iproute2
RUN mkdir -p /opt/scripts
COPY setup.sh /opt/scripts
RUN chmod +x /opt/scripts/setup.sh
ENV CGO_ENABLED=1
RUN go mod tidy
ENV CGO_ENABLED=1
RUN go build -o hping-platform .
RUN mv login.html feed.html hping-platform /opt/hping-platform-app
WORKDIR /opt/hping-platform-app
RUN rm -rf /opt/hping-platform
EXPOSE 8080
# Set the entry point
CMD ["/opt/scripts/setup.sh"]
