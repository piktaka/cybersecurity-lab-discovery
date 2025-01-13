FROM ubuntu:22.04

# Update and install required tools
RUN apt update && apt install -y wget git sudo gcc make nano


ENV DEBIAN_FRONTEND=noninteractive
RUN mkdir -p /run/sshd
RUN apt install -y openssh-server
COPY sshd_config /etc/sshd/sshd_config
COPY sudoers /etc/sudoers

RUN useradd -m -G sudo -s /bin/bash ubuntu
RUN adduser lablabee --disabled-password
RUN chmod a+w /etc/passwd
RUN chmod a+r /etc/shadow

RUN mkdir -p /opt/linpeas
WORKDIR /opt/linepeas
RUN wget https://github.com/peass-ng/PEASS-ng/releases/latest/download/linpeas.sh
# Add the Bionic repository and its GPG key
# Set root password
RUN echo 'root:lablabee2025#' | chpasswd

RUN echo 'ubuntu:lablabee@#' | chpasswd


WORKDIR /home/ubuntu
EXPOSE 22
CMD ["/usr/sbin/sshd", "-D"]
