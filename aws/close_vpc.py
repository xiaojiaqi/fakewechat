#!/usr/bin/env python

import boto.ec2
import boto.vpc

# set the region as you wish
region = 'ap-northeast-1'
keypair = 'testkey'
imgid = 'ami-b1b458b1'
imgid = 'ami-936d9d93'
instancetype = 't2.micro'
aws_access_key =  'xxx'
aws_secret_key = 'xxxx'

vpcCon = boto.vpc.connect_to_region(region, aws_access_key_id= aws_access_key, aws_secret_access_key =aws_secret_key)

print "subnets"    
orgtab = vpcCon.get_all_subnets()
for i in orgtab:
    print i.id
    print vpcCon.delete_subnet(i.id)

print "gateway"    
orgtab = vpcCon.get_all_customer_gateways()
for i in orgtab:
    print i.id
    print vpcCon.delete_customer_gateway(i.id)    

    

print "internet gateway"
gw = vpcCon.get_all_internet_gateways()
for i in gw:
    print i.id
        

print "all vpc"        
vpcs = vpcCon.get_all_vpcs()
for i in vpcs:
    print i.id
    

for i in gw:
    for k in vpcs:
        try:
            print vpcCon.detach_internet_gateway(i.id, k.id)
        except Exception, e:
            pass

for i in gw:
    print i.id
    print vpcCon.delete_internet_gateway(i.id)    

# get_all_security_groups
print "get_all_security_groups"

import boto.ec2
print "begin delete security group"

ec2 = boto.ec2.connect_to_region(region, aws_access_key_id= aws_access_key, aws_secret_access_key =aws_secret_key)
sgroup = ec2.get_all_security_groups()
for i in sgroup:
    print i.id
    print i.name
    if i.name != 'default':
        try:
            ec2.delete_security_group(i.id)
        except Exception, e:
            pass
print "peering_connections"    
orgtab = vpcCon.get_all_vpc_peering_connections()
for i in orgtab:
    print i.id
    print vpcCon.delete_vpc_peering_connection(i.id)  

#print "route"
print "route table"
orgtab = vpcCon.get_all_route_tables()
for i in orgtab:
    print i.id
    print vpcCon.delete_route_table(i.id)    
    
print "all vpc"        
orgtab = vpcCon.get_all_vpcs()
for i in orgtab:
    print i.id
    print vpcCon.delete_vpc(i.id)
    




  
    





    
