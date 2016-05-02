#!/usr/bin/env python
import os



host=[]

rgsize = 1
rgrange = 1000

rpinsamerange = 20
rpinotherrange = 10

localposterinstancesize = 1

local = False
lip="127.0.0.1"

HOME=os.getenv('HOME')

CD_HOME="cd " + HOME + "/bin"


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

def CleanCmd():
    global cmds
    cmds = ""

v ="""

taskset -c 0 redis-server ./redis1.conf
taskset -c 1 redis-server ./redis2.conf
taskset -c 2 redis-server ./redis3.conf
taskset -c 3 redis-server ./redis4.conf

"""
print v

def CleanRedisData():
    CleanCmd()
    for i in range(1, rgsize+1):
        for ipadd in rg[i]:
            Pcmd ("redis-cli -h " + str(ipadd)  + " -p  " +  str(1500) + " flushall")
    SaveCmds("cmd/clean_redis.sh")

print "\n\n\n"


# list host for /etc/ansible/host
def listHost():
    CleanCmd()
    for i in range(1, rgsize +1):
        for ipadd in rg[i]:
            Pcmd(str(ipadd))
    SaveCmds("cmd/hosts")



def makeClientName(ipadd):
    return "cmd/client/client" + str(ipadd) +".sh"

def makeServerName(ipadd):
    return "cmd/server/server" + str(ipadd) +".sh"

def makeCheckName(ipadd):
    return "cmd/server/checkdb" + str(ipadd) +".sh"

def GenClientCmd():
    CleanCmd()
    for i in range(1, rgsize+1):
        for ipadd in rg[i]:
            Pcmd (CD_HOME)
            Pcmd( "./client -host " + str(ipadd)+  " -port "+ str( 9500) +  " -minid " + str( 1 + (i-1)*rgrange) +  " -maxid " + str( rgrange * i)  + " -monitor "+ monitorHost + ":8002"   )
            SaveCmds(makeClientName(ipadd))
    print "\n\n\n"

#gen command for Monitor and router
def GenMonitorCmd(head, tail, log):
    CleanCmd()
    Pcmd (CD_HOME)
    Pcmd (head +   " ./router " + tail)
    Pcmd (head +   " ./monitorserver " + tail)
    SaveCmds("cmd/monitor.sh")


def GenServerCmd(head, tail, log):
    #route
    print "\n\n\n"
    index = 1
    CleanCmd()
    #Cache

    for i in range(1, rgsize +1):
        for ipadd in rg[i]:
            Pcmd (CD_HOME)
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
            for k in range(1, localposterinstancesize+1):
                port += 1
                logstr = " "
                if log == True:
                    logstr = " > ./localposter" + str(index) + ".log  2 >&1 "
                Pcmd( head + "  ./localposter -hostid=local" + str(index) + " -servertype=localposter -listenaddress=\"" + str(ipadd) + "\" -listenport=" + str(9700+index ) + " -rgid=" + str(i) + " -process=local " + parameter  + logstr   + tail)
                index += 1
            print "\n\n\n"
            SaveCmds(makeServerName(ipadd))

def GenDataCmd():
    CleanCmd()
    Pcmd ("cd " + HOME +"/gopath/src/github.com/fakewechat/bin")
    cmd = "cd ../python"
    Pcmd(cmd)
    cmd = "python gen.py -a " + "1," + str(rgsize * rgrange ) + "  -c " + str(rpinsamerange) +  "   -m " +  str(rpinotherrange) + "   -r " + str( rgrange )
    Pcmd(cmd)
    cmd = "cat ./rp.txt | ./split.py -k " + str(rgrange )
    Pcmd(cmd)
    cmd = "cd ../bin"
    Pcmd(cmd)
    SaveCmds("cmd/genData.sh")

def GoCheckDBCmd():
    CleanCmd()
    for i in range(1, rgsize+1):
        for ipadd in rg[i]:
            Pcmd (CD_HOME)
            Pcmd( "./goredis -host " + str(ipadd)+  " -port "+ str( 1500) +  " -minid " + str( 1 + (i-1)*rgrange) +  " -maxid " + str( rgrange * i)  + " -method " +  " check  > ./" + "checkdbresult" + str(ipadd) + ".txt"   )
            SaveCmds(makeCheckName(ipadd))
    print "\n\n\n"

def GoSaveDBCmd():
    CleanCmd()
    Pcmd (CD_HOME)
    Pcmd( "./goredis -host 127.0.0.1 -port "+ str( 1500) +   " -method " +  "save " + " -file db.txt"   )
    SaveCmds("cmd/gosaveDb.sh")
    print "\n\n\n"

def PyCheckDBCmd():
    CleanCmd()
    Pcmd (CD_HOME)
    for i in range(1, rgsize+1):
        for ipadd in rg[i]:
            Pcmd ("./checkdb.py -m " + str(ipadd) + " -p " +str(1500) +  " -l " + str( 1 + (i-1)*rgrange) +  " -g " + str( rgrange * i))
    SaveCmds("cmd/pycheckdb.sh")
    print "\n\n\n"

def PySaveDBCmd():
    CleanCmd()
    for i in range(1, rgsize+1):
        Pcmd ("./savetoredis.py -m " + rg[i][0] + " -p " + str(1500) + " < " +  str(i) + ".txt")
    SaveCmds("cmd/savetoredis_all.sh")
    print "\n\n\n"

def Genresign():
    CleanCmd()
    for i in range(1, rgsize+1):
        Pcmd("curl -kvv \"http://" + routeHost + ":" + routePort + "/regist/redis/"+ str(i) + "/?id=res" + str(i) +"&host=" + str(rg[i][0]) +"&port="+str(1500)+"&cellid=1\"")
    SaveCmds("cmd/resign.sh")

def Stop_Redis():
    CleanCmd()
    cmd = "killall -9 " + "redis-server"
    Pcmd(cmd)
    cmd = "rm -f " + HOME + "/bin/dump.rdb"
    Pcmd(cmd)
    SaveCmds("cmd/stopredis.sh")

def Start_Redis():
    CleanCmd()
    Pcmd (CD_HOME)
    Pcmd("./redis-server  ./redis.conf")
    SaveCmds("cmd/startredis.sh")

def GenAnsible():
    CleanCmd()

    # copy data
    cmd = "ansible all -m copy -a \"src=" + HOME + "/gopath/src/github.com/fakewechat/bin dest=" + HOME + "\" "
    Pcmd(cmd)

    cmd = "ansible all -m copy -a \"src=" + HOME + "/gopath/src/github.com/fakewechat/package/sysctl.conf dest=" + HOME + "/bin\" "
    Pcmd(cmd)
    for i in range(1, rgsize +1):
        for ipadd in rg[i]:
            cmd = "ansible " + str(ipadd) + " -m copy -a \"src=" + HOME + "/gopath/src/github.com/fakewechat/bin/" + makeClientName(ipadd) + " dest=" + HOME + "/bin/client.sh\""
            Pcmd(cmd)
            cmd = "ansible " + str(ipadd) + " -m copy -a \"src=" + HOME + "/gopath/src/github.com/fakewechat/bin/" + makeServerName(ipadd) + " dest=" + HOME + "/bin/server.sh\""
            Pcmd(cmd)
            cmd = "ansible " + str(ipadd) + " -m copy -a \"src=" + HOME + "/gopath/src/github.com/fakewechat/bin/" + makeCheckName(ipadd)  + " dest=" + HOME + "/bin/checkdb.sh\""
            Pcmd(cmd)


            cmd = "ansible " + str(ipadd) + " -m copy -a \"src=" + HOME + "/gopath/src/github.com/fakewechat/python/" + str(i) + ".txt" + " dest=" + HOME + "/bin/db.txt\""
            Pcmd(cmd)

    # chmod
    cmd = "ansible all -a \"sudo cp " + HOME + "/bin/sysctl.conf /etc\""
    Pcmd(cmd)

    cmd = "ansible all  -m raw -a \"chmod 777 -R " + HOME + "/bin/\""
    Pcmd(cmd)



    # setup
    cmd = "ansible all -a \"sudo sysctl -p\""
    Pcmd(cmd)

    cmd = "ansible all  -m raw -a \"systemctl stop firewalld.service\""
    Pcmd(cmd)

    #cmd = "ansible all -a \"sudo rpm -ivh  " + HOME + "/bin/*.rpm\""
    #Pcmd(cmd)

    SaveCmds("cmd/ansible.sh")


def GerDBResultCmd():
    CleanCmd()
    for i in range(1, rgsize+1):
        for ipadd in rg[i]:
            Pcmd( "scp "+ str(ipadd) +":" +HOME + "/bin/checkdbresult*.txt ./"   )
    SaveCmds("cmd/getalldbresult.sh")

    print "\n\n\n"
if __name__ == "__main__":
    #list host for /etc/ansible/hosts
    listHost()
    # Gen client cmd
    GenClientCmd()
    # Gen Server side cmd
    GenServerCmd("nohup", " >/dev/null &", False)

    GenMonitorCmd("nohup", " >/dev/null &", False)

    GoCheckDBCmd()
    GoSaveDBCmd()

    CleanRedisData()

    PyCheckDBCmd()
    PySaveDBCmd()

    Genresign()

    GenDataCmd()

    GenAnsible()

    Start_Redis()
    Stop_Redis()

    GerDBResultCmd()
