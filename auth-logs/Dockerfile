FROM ubuntu:22.04


ENV DEBIAN_FRONTEND=noninteractive

RUN mkdir -p /run/sshd


RUN apt update
ENV DEBIAN_FRONTEND=noninteractive
RUN  apt install -y openssh-server
RUN apt install -y fail2ban

RUN useradd -m -s /bin/bash ubuntu
RUN useradd -m -s /bin/bash lablabee


RUN echo 'root:lablabee2025#' | chpasswd
RUN echo 'ubuntu:lablabee2025@!' | chpasswd
RUN echo 'lablabee:lablabee2025@?' | chpasswd

RUN  apt install -y rsyslog
RUN sed -i '/module(load="imklog" permitnonkernelfacility="on")/s/^/#/' /etc/rsyslog.conf
WORKDIR /home/ubuntu
RUN  echo 'LogLevel INFO' >> /etc/ssh/sshd_config && \
    echo 'SyslogFacility AUTH' >> /etc/ssh/sshd_config
EXPOSE 22

COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"]

