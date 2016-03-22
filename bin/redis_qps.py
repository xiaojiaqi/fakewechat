#!/usr/bin/env python

import sys
import getopt
import redis
import message_pb2
import time
import random

def usage():
    print "redis_qps.py -m host -p port "


if __name__ == "__main__":
    try:
        opts, args = getopt.getopt(sys.argv[1:], "m:p:r:h", ["help"])

    except getopt.GetoptError:
        usage()
        sys.exit(2)

    hostname = "localhost"
    port = 6379
    testcount = 15000*3
    for o, a in opts:
        if o in ('-h', '--help'):
            usage()
            sys.exit(1)
        elif o in ('-m' ):
            hostname = a
        elif o in ('-p'):
            port = int(a)

    r = redis.Redis(host=hostname,port= port  ,db=0)
    la = []
    lb = []
    lc = []
    a = message_pb2.UserInfor()
    astring = a.SerializeToString()
    for i in range (0, testcount):
        la.append(astring)
        lb.append("user#" + str(i))
        lc.append("user###"+ str(i/20))
    timebegin = time.time()
    for i in range(0, testcount):
        r.set(lb[i], astring)
        #r.hset(lc[i], lb[i], astring)
        #m = r.hget(lc[i], lb[i] )
    spendtime = time.time() - timebegin
    print str(spendtime) + " second passed"
    print str(testcount  / spendtime*1.0 ) + " qps "

    pipe = r.pipeline()

    timebegin = time.time()
    for i in range(0, testcount):
        pipe.set(lb[i], astring)
        #pipe.hset(lc[i], lb[i], astring)
        #pipe.hget(lc[i], lb[i] )
    pipe.execute()
    spendtime = time.time() - timebegin
    print str(spendtime) + " second passed"
    print str(testcount  / spendtime*1.0 ) + " qps "
