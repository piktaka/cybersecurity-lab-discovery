FROM ubuntu:20.04

# Update and install required tools
RUN apt update && apt install -y wget git sudo gcc make


RUN mkdir -p /opt/openssl
WORKDIR /opt/openssl
RUN wget https://www.openssl.org/source/openssl-1.1.1d.tar.gz
RUN tar -xvzf openssl-1.1.1d.tar.gz
WORKDIR /opt/openssl/openssl-1.1.1d
RUN ./config
RUN make
RUN make install

# RUN echo 'deb http://archive.ubuntu.com/ubuntu/ focal main universe' | tee -a /etc/apt/sources.list
# RUN apt update && apt install -y openssl=1.1.1k-0ubuntu2
# Install OpenSSH server
ENV DEBIAN_FRONTEND=noninteractive
RUN mkdir -p /run/sshd
RUN apt install -y openssh-server
COPY sshd_config /etc/sshd/sshd_config
COPY sudoers /etc/sudoers

RUN useradd -m -G sudo -s /bin/bash ubuntu
RUN adduser lablabee --disabled-password
RUN chmod a+w /etc/passwd
RUN chmod a+r /etc/shadow

# curl -L https://github.com/peass-ng/PEASS-ng/releases/latest/download/linpeas.sh | sh
# Add the Bionic repository and its GPG key
# Set root password
RUN echo 'root:lablabee2025#' | chpasswd

RUN echo 'ubuntu:lablabee2025#' | chpasswd

CMD ["/usr/sbin/sshd", "-D"]
