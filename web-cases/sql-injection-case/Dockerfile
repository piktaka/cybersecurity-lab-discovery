FROM ubuntu:latest


RUN apt update && apt install -y wget gcc

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
RUN mkdir -p /opt/sql-injection-case /opt/sql-injection-case-app /usr/local/go
COPY . /opt/sql-injection-case 

WORKDIR /opt/sql-injection-case 

# Install Go modules and build the project
ENV CGO_ENABLED=1
RUN go mod tidy
RUN go build -o sql-injection-platform .
RUN mv login.html home.html sql-injection-platform /opt/sql-injection-case-app
WORKDIR /opt/sql-injection-case-app
RUN rm -rf /opt/sql-injection-case
EXPOSE 8080

CMD ["./sql-injection-platform"]
