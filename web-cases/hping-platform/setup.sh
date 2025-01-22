#!/bin/bash



tc qdisc add dev eth0 root tbf rate 1mbit burst 32kbit latency 400ms
tc qdisc add dev eth0 handle ffff: ingress
tc filter add dev eth0 parent ffff: protocol ip prio 50 u32 match ip src 0.0.0.0/0 police rate 1mbit burst 32kbit drop

/opt/hping-platform-app/hping-platform