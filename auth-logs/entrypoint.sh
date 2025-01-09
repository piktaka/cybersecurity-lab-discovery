#!/bin/bash

# Start rsyslog
/usr/sbin/rsyslogd &

# Start SSH server
exec /usr/sbin/sshd -D
