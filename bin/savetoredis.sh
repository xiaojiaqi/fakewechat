#!/bin/bash
set -o -e

cd /home/ec2-user/bin

./goredis -host 127.0.0.1 -port 1500 -file db.txt