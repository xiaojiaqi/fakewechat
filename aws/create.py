#!/usr/bin/env python

import boto.ec2

ec2 = boto.ec2.connect_to_region('ap-northeast-1')
security_group=[]
security_group.append('subnet-20ac0857')

s=ec2.get_all_security_groups()
for i in s:
    print s

ec2 = boto.ec2.connect_to_region('ap-northeast-1')

i=11
ip = '10.0.0.' + str(i)
interface = boto.ec2.networkinterface.NetworkInterfaceSpecification(subnet_id='subnet-20ac0857', groups=['sg-0454d761'], associate_public_ip_address=True, private_ip_address=ip)
interfaces = boto.ec2.networkinterface.NetworkInterfaceCollection(interface)
t = ec2.run_instances('ami-b1b458b1', instance_type='t2.micro', min_count=1, key_name='testkey', network_interfaces=interfaces )


i=11
ip = '10.0.1.' + str(i)
interface = boto.ec2.networkinterface.NetworkInterfaceSpecification(subnet_id='subnet-17ac0860', groups=['sg-e554d780'] , private_ip_address=ip)
interfaces = boto.ec2.networkinterface.NetworkInterfaceCollection(interface)
t = ec2.run_instances('ami-b1b458b1', instance_type='t2.micro', min_count=1, key_name='testkey', network_interfaces=interfaces )

i=11
ip = '10.0.2.' + str(i)
interface = boto.ec2.networkinterface.NetworkInterfaceSpecification(subnet_id='subnet-12ac0865', groups=['sg-e654d783'] , private_ip_address=ip)
interfaces = boto.ec2.networkinterface.NetworkInterfaceCollection(interface)
t = ec2.run_instances('ami-b1b458b1', instance_type='t2.micro', min_count=1, key_name='testkey', network_interfaces=interfaces )

