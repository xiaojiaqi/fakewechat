#!/usr/bin/env python

import boto.vpc

# set the region as you wish
region = 'ap-northeast-1'
keypair = 'testkey'
imgid = 'ami-b1b458b1'
imgid = 'ami-936d9d93'
instancetype = 't2.micro'
aws_access_key =  'xxxx'
aws_secret_key = 'xxxx'



vpcCon = boto.vpc.connect_to_region(region, aws_access_key_id= aws_access_key, aws_secret_access_key =aws_secret_key)

# list route table
print "list route table"
#routetab = vpcCon.create_route_table(vpcid)

orgtab = vpcCon.get_all_route_tables()
for i in orgtab:
    print i.id



vpcs = vpcCon.get_all_vpcs()
for i in vpcs:
    print "vpc ", i.id, i.is_default

        
# create vpc
print "begin create vpc"
vpc = vpcCon.create_vpc("10.0.0.0/16")
print "vpc ", vpc.id, " was created"

vpcid = vpc.id


# create internet gateway

gw = vpcCon.create_internet_gateway()

print "gw ", gw.id , " was created"

gwid = gw.id

# attach internet gateway

b = vpcCon.attach_internet_gateway(gwid, vpcid)
if b == True:
    print "attach gateway success"
else:
    print "attach gateway failes"
    exit()
    

# list route table
print "list route table"
#routetab = vpcCon.create_route_table(vpcid)

#need refactor 
newtab = vpcCon.get_all_route_tables()
routetableid = ""
for i in newtab:
    print i.id
    got = False
    for k in range(0, len(orgtab)):
        if i.id == orgtab[k].id:
            got = True
            break
    if got == False:
        routetableid = i.id
        break
if routetableid == "":
    print "can't find newtable id "
    exit()
print "new route table id is ", routetableid
       


#create route
print "create route"
b = vpcCon.create_route(routetableid, "0.0.0.0/0", gateway_id=gwid)


# create subnet
print "begin create subnet"

net1 = vpcCon.create_subnet(vpcid, "10.0.1.0/24")
net2 = vpcCon.create_subnet(vpcid, "10.0.2.0/24")
net3 = vpcCon.create_subnet(vpcid, "10.0.3.0/24")
net4 = vpcCon.create_subnet(vpcid, "10.0.4.0/24")


print "subnet ", net1.id, " was created"
print "subnet ", net2.id, " was created"
print "subnet ", net3.id, " was created"
print "subnet ", net4.id, " was created"



import boto.ec2
print "begin create security group"

ec2 = boto.ec2.connect_to_region(region, aws_access_key_id= aws_access_key, aws_secret_access_key =aws_secret_key)
net1_group = ec2.create_security_group('net1', 'Our net1 Group', vpc_id=vpcid)
net2_group = ec2.create_security_group('net2', 'Our net2 Group', vpc_id=vpcid)
net3_group = ec2.create_security_group('net3', 'Our net3 Group', vpc_id=vpcid)
net4_group = ec2.create_security_group('net4', 'Our net3 Group', vpc_id=vpcid)


print "begin add authorize to security"
b = net1_group.authorize("tcp", from_port=0,to_port=65535, cidr_ip="0.0.0.0/0")
b = net1_group.authorize("udp", from_port=0,to_port=65535, cidr_ip="0.0.0.0/0")

b = net2_group.authorize("tcp", from_port=0,to_port=65535, cidr_ip="0.0.0.0/0")
b = net2_group.authorize("udp", from_port=0,to_port=65535, cidr_ip="0.0.0.0/0")

b = net3_group.authorize("tcp", from_port=0,to_port=65535, cidr_ip="0.0.0.0/0")
b = net3_group.authorize("udp", from_port=0,to_port=65535, cidr_ip="0.0.0.0/0")

b = net4_group.authorize("tcp", from_port=0,to_port=65535, cidr_ip="0.0.0.0/0")
b = net4_group.authorize("udp", from_port=0,to_port=65535, cidr_ip="0.0.0.0/0")


print " security group ", net1_group.id, " was created"
print " security group ", net2_group.id, " was created"
print " security group ", net3_group.id, " was created"
print " security group ", net4_group.id, " was created"

#exit()

# create ec2 instance

print "begin create ec2 instance"

ec2 = boto.ec2.connect_to_region(region, aws_access_key_id= aws_access_key, aws_secret_access_key =aws_secret_key)

i=11
ip = '10.0.1.' + str(i)
interface = boto.ec2.networkinterface.NetworkInterfaceSpecification(subnet_id=net1.id, groups=[net1_group.id], associate_public_ip_address=True, private_ip_address=ip)
interfaces = boto.ec2.networkinterface.NetworkInterfaceCollection(interface)
t = ec2.run_instances(imgid, instance_type=instancetype, min_count=1, key_name=keypair, network_interfaces=interfaces )
print t

exit()

i=11
ip = '10.0.2.' + str(i)
interface = boto.ec2.networkinterface.NetworkInterfaceSpecification(subnet_id=net2.id, groups=[net2_group.id] , private_ip_address=ip)
interfaces = boto.ec2.networkinterface.NetworkInterfaceCollection(interface)
t = ec2.run_instances(imgid, instance_type=instancetype, min_count=1, key_name=keypair, network_interfaces=interfaces )
print t

i=11
ip = '10.0.3.' + str(i)
interface = boto.ec2.networkinterface.NetworkInterfaceSpecification(subnet_id=net3.id, groups=[net3_group.id] , private_ip_address=ip)
interfaces = boto.ec2.networkinterface.NetworkInterfaceCollection(interface)
t = ec2.run_instances(imgid, instance_type=instancetype, min_count=1, key_name=keypair, network_interfaces=interfaces )
print t
