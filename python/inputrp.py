#!/usr/bin/env python

import sys
import getopt
import redis


def usage():
    print "imput.py -f file -h host -p port -d db"

if __name__ == "__main__":
    if len(sys.argv) == 1:
        print "len"
        usage()
        sys.exit(2)
    try:
        opts, args = getopt.getopt(sys.argv[1:], "f:h:p:d:", ["help"])
    except getopt.GetoptError, e:
        print  "1111"
        usage()
        sys.exit(2)

    redishost ="localhost"
    redisport=6379
    redisdb=1
    rpfile=""
    for o, a in opts:
        if o in ('-f' ):
            rpfile = a
        elif o in ('-h'):
            redishost = a
        elif o in ('-p'):
            redisport = int(a)
        elif o in ('-d'):
            redisdb = int(a)

    r =  redis.Redis(host=redishost, port=redisport, db=redisdb)
    while True:
        line = sys.stdin.readline()
        if line:
            v = line.split(",")
            if len(v) >= 3:
                name = "rp" + v[0]
                r.sadd(name, int(v[1])*1000000+int(v[2])*10000) 
        else:
            break
    for i in range(1, 999):
        
        listl= r.smembers('rp' + str(i))
        print "="*20
        for i in listl:
            print i

