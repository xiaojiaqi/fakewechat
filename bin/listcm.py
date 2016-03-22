#!/usr/bin/env python


host=[]

local = False
lip="127.0.0.1"
rgsize = 1
rgrange = 1000
localpostersize = 1

routeHost="10.29.101.3"
routePort="8089"

monitorHost="10.29.101.3"
monitorPort="8000"

rg = {}
rg[1] = ["10.29.101.3","10.29.101.7"]
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

routeurl = " -routeserverurl=\"http://" + routeHost + ":" + routePort +"/server/\"  "
routrreg = " -routeregisturl=\"http://" + routeHost + ":" + routePort+ "/regist/\" "
rgurl= " -rgsize=" +  str(rgrange) + " "
moniturl = " -monitorhost=\""+ monitorHost + ":" + monitorPort + "\""

parameter = routeurl + routrreg  +  rgurl +  moniturl


print "\n\n\n"

def Pcmd(i):
    print i

v ="""

taskset -c 0 redis-server ./redis1.conf
taskset -c 1 redis-server ./redis2.conf
taskset -c 2 redis-server ./redis3.conf
taskset -c 3 redis-server ./redis4.conf

"""
print v


for i in range(1, rgsize+1):
    Pcmd ("redis-cli -p  " +  str(1500+i) + " flushall")

print "\n\n\n"


for i in range(1, rgsize+1):
    for ipadd in rg[i]:
        print "./client -host " + str(ipadd)+  " -port ", 9500+i,  " -minid ", 1 + (i-1)*rgrange,  " -maxid ", rgrange * i

print "\n\n\n"


for i in range(1, rgsize+1):
    for ipadd in rg[i]:
        print "./checkdb.py -m " + str(ipadd) + " -p ", 1500+i,  " -l ", 1 + (i-1)*rgrange,  " -g ", rgrange * i

print "\n\n\n"

for i in range(1, rgsize+1):
    Pcmd ("./savetoredis.py -m " + rg[i][0] + " -p " + str(1500+i) + " < " +  str(i) + ".txt")

print "\n\n\n"


for i in range(1, rgsize+1):
    Pcmd("curl -kvv \"http://" + routeHost + ":" + routePort + "/regist/redis/"+ str(i) + "/?id=res" + str(i) +"&host=" + str(rg[i][0]) +"&port="+str(1500+i)+"&cellid=1\"")



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

    Pcmd (head +   " ./router " + tail)
    Pcmd (head +   " ./monitorserver " + tail)

    #Cache

    for i in range(1, rgsize +1):
        for ipadd in rg[i]: 
            logstr = " "
            if log == True:
                logstr = " > ./ca" + str(index) + ".log  2 >&1 "

            Pcmd( head + "  ./cacheserver -hostid=ca" + str(index) + " -servertype=cache -listenaddress=\"" + str(ipadd) + "\" -listenport=" + str(9600+i) + " -rgid=" + str(i) + " -process=ca "+ parameter  + logstr  + tail)
            index += 1

            logstr = " "
            if log == True:
                logstr = " > ./gw" + str(index) + ".log  2 >&1 "
            Pcmd( head + "  ./gateway -hostid=gw" + str(index) + " -servertype=gw -listenaddress=\"" + str(ipadd) + "\" -listenport=" + str(9500+i) + " -rgid=" + str(i) + " -process=gw " + parameter + logstr   + tail)
            index += 1

            #poster
            logstr = " "
            if log == True:
                logstr = " > ./poster" + str(index) + ".log  2 >&1 "
            Pcmd( head + "  ./poster -hostid=poster" + str(index) + " -servertype=poster -listenaddress=\"" + str(ipadd) + "\" -listenport=" + str(9800+i) + " -rgid=" + str(i) + " -process=poster " + parameter  + logstr  + tail)

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


