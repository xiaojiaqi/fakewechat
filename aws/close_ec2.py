#!/usr/bin/env python

import boto.ec2
aws_access_key =  'xxxx'
aws_secret_key = 'xxxx'
region = 'ap-northeast-1'


ec2 = boto.ec2.connect_to_region(region, aws_access_key_id= aws_access_key, aws_secret_access_key =aws_secret_key)
t = ec2.get_all_instances()
for i in t:
    print i.instances
    for k in i.instances:
        print k.id
        ec2.terminate_instances(k.id)