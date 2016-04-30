#!/usr/bin/env python


host=[]

rpinsamerange = 20
rpinotherrange = 10

local = False
lip="127.0.0.1"


rgsize = 4
rgrange = 1000
localpostersize = 1

routeHost="10.0.2.20"
routePort="8089"

monitorHost="10.0.2.20"
monitorPort="8000"

rg = {}
rg[1] = ["10.29.101.3"]
rg[2] = ["10.29.101.7"]
rg[3] = ["192.168.4.45"]
rg[4] = ["192.168.4.46"]
rg[5] = ["10.29.114.35"]
rg[6] = ["192.168.0.1"]

rg[7] = ["192.168.0.1"]
rg[8] = ["192.168.0.1"]

rg[9] = ["192.168.0.1"]
rg[10] = ["192.168.0.1"]


if local == True:
    for i in rg:
        rg[i] = [lip]
    routeHost = lip
    monitorHost = lip

rg = {}
for i in range(1, rgsize + 1):
    rg[i] = ["10.0.2." + str(i+20 - 1)]

routeurl = " -routeserverurl=\"http://" + routeHost + ":" + routePort +"/server/\"  "
routrreg = " -routeregisturl=\"http://" + routeHost + ":" + routePort+ "/regist/\" "
rgurl= " -rgsize=" +  str(rgrange) + " "
moniturl = " -monitorhost=\""+ monitorHost + ":" + monitorPort + "\""

parameter = routeurl + routrreg  +  rgurl +  moniturl


print "\n\n\n"


global cmds

cmds = ""

def Pcmd(i):
    print i
    global cmds
    cmds += i + "\n"

#refactor
def SaveCmds(path):
    global cmds
    filew = open(path, 'w')
    filew.write(cmds)
    filew.close( )
    cmds = ""


v ="""

taskset -c 0 redis-server ./redis1.conf
taskset -c 1 redis-server ./redis2.conf
taskset -c 2 redis-server ./redis3.conf
taskset -c 3 redis-server ./redis4.conf

"""
print v


for i in range(1, rgsize+1):
    for ipadd in rg[i]:
        Pcmd ("redis-cli -h " + str(ipadd)  + " -p  " +  str(1500) + " flushall")
SaveCmds("cmd/clean_redis.sh")

print "\n\n\n"




def makeClientName(ipadd):
    return "cmd/client/client" + str(ipadd) +".sh"

def makeServerName(ipadd):
    return "cmd/server/server" + str(ipadd) +".sh"


for i in range(1, rgsize+1):
    for ipadd in rg[i]:
        Pcmd ("cd /home/ec2-user/bin")
        Pcmd( "./client -host " + str(ipadd)+  " -port "+ str( 9500) +  " -minid " + str( 1 + (i-1)*rgrange) +  " -maxid " + str( rgrange * i)  + " -monitor "+ monitorHost + ":8002"   )
        SaveCmds(makeClientName(ipadd))
print "\n\n\n"


for i in range(1, rgsize+1):
    for ipadd in rg[i]:
        Pcmd ("./checkdb.py -m " + str(ipadd) + " -p " +str( 1500) +  " -l " + str( 1 + (i-1)*rgrange) +  " -g " + str( rgrange * i))
SaveCmds("cmd/checkdb.sh")
print "\n\n\n"

for i in range(1, rgsize+1):
    Pcmd ("./savetoredis.py -m " + rg[i][0] + " -p " + str(1500) + " < " +  str(i) + ".txt")
SaveCmds("cmd/savetoredis.sh")
print "\n\n\n"


for i in range(1, rgsize+1):
    Pcmd("curl -kvv \"http://" + routeHost + ":" + routePort + "/regist/redis/"+ str(i) + "/?id=res" + str(i) +"&host=" + str(rg[i][0]) +"&port="+str(1500)+"&cellid=1\"")
SaveCmds("cmd/resign.sh")


def show(head, tail, log):
    #route
    print "\n\n\n"
    index = 1

    Pcmd (head +   " ./router " + tail)
    Pcmd (head +   " ./monitorserver " + tail)

    #Cache

    for i in range(1, rgsize +1):
        logstr = " "
        if log == True:
            logstr = " > ./ca" + str(index) + ".log  2 >&1 "
        Pcmd( head + "  ./cacheserver -hostid=ca" + str(index) + " -servertype=cache -listenaddress=\"" + str(rg[i]) + "\" -listenport=" + str(9600+i) + " -rgid=" + str(i) + " -process=ca "+ parameter  + logstr  + tail)
        index += 1

    print "\n\n\n"
    #gateway
    for i in range(1, rgsize +1):
        logstr = " "
        if log == True:
            logstr = " > ./gw" + str(index) + ".log  2 >&1 "

        Pcmd( head + "  ./gateway -hostid=gw" + str(index) + " -servertype=gw -listenaddress=\"" + str(rg[i]) + "\" -listenport=" + str(9500+i) + " -rgid=" + str(i) + " -process=ca " + parameter + logstr   + tail)
        index += 1
    print "\n\n\n"
    #localposter
    port = 0
    for k in range(1, localpostersize+1):
        for i in range(1, rgsize +1):
            port += 1
            logstr = " "
            if log == True:
                logstr = " > ./localposter" + str(index) + ".log  2 >&1 "


            Pcmd( head + "  ./localposter -hostid=ca" + str(index) + " -servertype=localposter -listenaddress=\"" + str(rg[i]) + "\" -listenport=" + str(9700+port ) + " -rgid=" + str(i) + " -process=ca " + parameter  + logstr   + tail)
            index += 1
        print "\n\n\n"
    #poster
    for i in range(1, rgsize +1):
        logstr = " "
        if log == True:
            logstr = " > ./poster" + str(index) + ".log  2 >&1 "

        Pcmd( head + "  ./poster -hostid=poster" + str(index) + " -servertype=poster -listenaddress=\"" + str(rg[i]) + "\" -listenport=" + str(9800+i) + " -rgid=" + str(i) + " -process=ca " + parameter  + logstr  + tail)
        index += 1
    print "\n\n\n"



def show2(head, tail, log):
    #route
    print "\n\n\n"
    index = 1

    Pcmd ("cd /home/ec2-user/bin")
    Pcmd (head +   " ./router " + tail)
    Pcmd (head +   " ./monitorserver " + tail)

    SaveCmds("cmd/monitor.sh")
    #Cache

    for i in range(1, rgsize +1):
        for ipadd in rg[i]:
            Pcmd ("cd /home/ec2-user/bin")
            logstr = " "
            if log == True:
                logstr = " > ./ca" + str(index) + ".log  2 >&1 "

            Pcmd( head + "  ./cacheserver -hostid=ca" + str(index) + " -servertype=cache -listenaddress=\"" + str(ipadd) + "\" -listenport=" + str(9600) + " -rgid=" + str(i) + " -process=ca "+ parameter  + logstr  + tail)
            index += 1

            logstr = " "
            if log == True:
                logstr = " > ./gw" + str(index) + ".log  2 >&1 "
            Pcmd( head + "  ./gateway -hostid=gw" + str(index) + " -servertype=gw -listenaddress=\"" + str(ipadd) + "\" -listenport=" + str(9500) + " -rgid=" + str(i) + " -process=gw " + parameter + logstr   + tail)
            index += 1

            #poster
            logstr = " "
            if log == True:
                logstr = " > ./poster" + str(index) + ".log  2 >&1 "
            Pcmd( head + "  ./poster -hostid=poster" + str(index) + " -servertype=poster -listenaddress=\"" + str(ipadd) + "\" -listenport=" + str(9800) + " -rgid=" + str(i) + " -process=poster " + parameter  + logstr  + tail)

            index += 1
            #localposter
            port = 0
            for k in range(1, localpostersize+1):
                port += 1
                logstr = " "
                if log == True:
                    logstr = " > ./localposter" + str(index) + ".log  2 >&1 "
                Pcmd( head + "  ./localposter -hostid=local" + str(index) + " -servertype=localposter -listenaddress=\"" + str(ipadd) + "\" -listenport=" + str(9700+index ) + " -rgid=" + str(i) + " -process=local " + parameter  + logstr   + tail)
                index += 1
            print "\n\n\n"
            SaveCmds(makeServerName(ipadd))

def showansible():
    cmd = """ansible all -m copy -a "src=/home/ec2-user/gopath/src/github.com/fakewechat/bin dest=/home/ec2-user"
ansible all -m copy -a "src=/home/ec2-user/gopath/src/github.com/fakewechat/package/sysctl.conf dest=/home/ec2-user/bin"
ansible all -m copy -a "src=/home/ec2-user/gopath/src/github.com/fakewechat/package/psmisc*.rpm dest=/home/ec2-user/bin"

ansible all -a "sudo cp /home/ec2-user/bin/sysctl.conf /etc"
ansible all -a "sudo sysctl -p"
ansible all -m file -a "dest=/home/ec2-user/bin mode=700"
ansible all -m file -a "dest=/home/ec2-user/bin/gateway mode=700"
ansible all -m file -a "dest=/home/ec2-user/bin/redis-server mode=700"
ansible all -m file -a "dest=/home/ec2-user/bin/router mode=700"
ansible all -m file -a "dest=/home/ec2-user/bin/poster mode=700"
ansible all -m file -a "dest=/home/ec2-user/bin/localposter mode=700"
ansible all -m file -a "dest=/home/ec2-user/bin/cacheserver mode=700"
ansible all -m file -a "dest=/home/ec2-user/bin/monitorclient mode=700"
ansible all -m file -a "dest=/home/ec2-user/bin/monitorserver mode=700"
ansible all -m file -a "dest=/home/ec2-user/bin/client mode=700"
ansible all -m file -a "dest=/home/ec2-user/bin/redis-cli mode=700"
ansible all -m file -a "dest=/home/ec2-user/bin/kill.sh mode=700"
ansible all -m file -a "dest=/home/ec2-user/bin/start_redis.sh mode=700"
ansible all -m file -a "dest=/home/ec2-user/bin/stop_redis.sh mode=700"
ansible all -m file -a "dest=/home/ec2-user/bin/redis.conf mode=700"




    """

    Pcmd(cmd)
    for i in range(1, rgsize +1):
        for ipadd in rg[i]:
            cmd = "ansible " + str(ipadd) + " -m copy -a \"src=/home/ec2-user/gopath/src/github.com/fakewechat/bin/" + makeClientName(ipadd) + " dest=/home/ec2-user/bin/client.sh\""
            Pcmd(cmd)
            cmd = "ansible " + str(ipadd) + " -m copy -a \"src=/home/ec2-user/gopath/src/github.com/fakewechat/bin/" + makeServerName(ipadd) + " dest=/home/ec2-user/bin/server.sh\""
            Pcmd(cmd)

    cmd = "ansible all -m file -a \"dest=/home/ec2-user/bin/client.sh mode=700\""
    Pcmd(cmd)

    cmd = "ansible all -m file -a \"dest=/home/ec2-user/bin/server.sh mode=700\""
    Pcmd(cmd)
    SaveCmds("cmd/ansible.sh")



def listHost():
    for i in range(1, rgsize +1):
        for ipadd in rg[i]:
            Pcmd(str(ipadd))
    SaveCmds("cmd/hosts")

def genCmd():
    cmd = "cd ../python"
    Pcmd(cmd)
    cmd = "python gen.py -a " + "1," + str(rgsize * rgrange ) + "  -c " + str(rpinsamerange) +  "   -m " +  str(rpinotherrange) + "   -r " + str( rgrange )
    Pcmd(cmd)
    cmd = "cat ./rp.txt | ./split.py -k " + str(rgrange )
    Pcmd(cmd)
    for i in range (1, rgsize +1 ):
        cmd = "mv " + str(i)+".txt" + " ../bin"
        Pcmd(cmd)
    cmd = "cd ../bin"
    Pcmd(cmd)


    SaveCmds("cmd/gen.sh")


"""
show ("","", False)
show ("screen", "", False)
show ("nohup", " >/dev/null &", False)
show ("nohup", "   &", True)


show2 ("","", False)
show2 ("screen", "", False)
"""

#show ("nohup", " >/dev/null &", False)



show2 ("nohup", " >/dev/null &", False)
#show2 ("nohup", "   &", True)
showansible()
listHost()
genCmd()

