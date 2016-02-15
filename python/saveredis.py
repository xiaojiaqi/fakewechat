#!/usr/bin/env python

import sys
import getopt
import redis
import message_pb2

import random

def usage():
    print "saveredis.py -m host -p port "


if __name__ == "__main__":
    try:
        opts, args = getopt.getopt(sys.argv[1:], "m:p:r:h", ["help"])

    except getopt.GetoptError:
        usage()
        sys.exit(2)

    hostname = "localhost"
    port = 6379
    for o, a in opts:
        if o in ('-h', '--help'):
            usage()
            sys.exit(1)
        elif o in ('-m' ):
            hostname = a
        elif o in ('-p'):
            port = int(a)

    r = redis.Redis(host=hostname,port= port  ,db=0)
    pipe = r.pipeline()
    num = 0
    userid = ""
    maps = {}

    while True:
        line = sys.stdin.readline()
        num += 1
        line = line.strip('\n')
        if line:
           # print line
            v = line.split(",")
            if len(v) >= 3:
                #print v
                newuserid = v[0]
                if newuserid != userid:
                    if userid == "":
                        maps[ v[1] ] = "0"
                        pass
                    else:
                        a = message_pb2.UserInfor()
                        for i in maps:
                            a.UserMap[int(i)].UserId = int(i)
                        a.SendId = random.randint(1,2000)
                        a.SendAckId = random.randint(100,2000)
                        print "user#" + userid #a.SerializeToString()
                        pipe.set("user#" + userid, a.SerializeToString())
                        pipe.execute()
                        maps={}
                        maps[ v[1] ] = "0"

                    userid = newuserid
                else:
                    maps[ v[1] ] = "0" 
        else:
            break

    a = message_pb2.UserInfor()
    for i in maps:
        a.UserMap[int(i)].UserId = int(i)
    print "user#" + userid #,a.SerializeToString()
    pipe.set("user#" + userid, a.SerializeToString())
    pipe.execute()

