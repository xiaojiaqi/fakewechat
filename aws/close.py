#!/usr/bin/env python

import boto.ec2

ec2 = boto.ec2.connect_to_region('ap-northeast-1')
t = ec2.get_all_instances()
for i in t:
    print i.instances
    for k in i.instances:
        print k.id
        ec2.terminate_instances(k.id)