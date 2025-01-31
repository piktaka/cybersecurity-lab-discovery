#!/bin/bash

# Start rsyslog
# /usr/sbin/rsyslogd -d &
rm -f /run/rsyslogd.pid || true
/usr/sbin/rsyslogd -d -n &
# Start SSH server
exec /usr/sbin/sshd -D
