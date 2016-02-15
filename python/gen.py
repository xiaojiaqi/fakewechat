#! /usr/bin/env python
# -*- coding:utf-8 -*-

import sys
import os
import getopt
import os
import subprocess
import datetime


# simple ranges
class ranges:
    def __init__(self):
        self.min = 0
        self.max = 0

def runCmd(cmd):
    print "run: " + cmd
    returnCode = subprocess.Popen(cmd, shell=True, stdout=subprocess.PIPE, stderr=subprocess.PIPE)
    print returnCode.stdout.read()

def usage():
    print "gen.py -a 1,1000 -c 100 -r 200000 -m 10"
    print "it mean : create a relationship,  user id from 1 to 1000"
    print "every user has 100 relationship in same rg"
    print "every user has 10 relationship cross rg"
    print "1 rg contain 200000 user"

if __name__ == "__main__":
    if len(sys.argv) == 1:
        usage()
        sys.exit(1)
    try:
        opts, args = getopt.getopt(sys.argv[1:], "a:c:r:m:h", ["help"])
    except getopt.GetoptError:
        usage()
        sys.exit(2)
    mixcount = 0
    count = 0
    id1 = ranges()
    rg = 2000000
    for o, a in opts:
        if o in ('-h', '--help'):
            usage()
            sys.exit(1)
        elif o in ('-a' ):
            v = a.split(',')
            if len(v) >= 2:
                id1.min = int(v[0])
                id1.max = int(v[1])
        elif o in ('-c'):
            count = int(a)
        elif o in ('-r'):
            rg = int(a)
        elif o in ('-m'):
            mixcount = int(a)

    v = []
    x = id1.min
    y = id1.max
    #tmp file to restore random data
    tmpfile = "tmp.file"
    # create gendata
    while True:

        v1 = ranges()
        v1.min = x
        if v1.min + rg -1 < y:
            v1.max = int(v1.min*1.0/rg + 1)*rg
        else:
            v1.max = y
        v.append(v1)
        if v1.max >= y:
            break
        x = v1.max+1

    for i in v:
        print i.min, i.max
    cmd = 'rm -rf ' + tmpfile
    runCmd(cmd)
    for i in v:
        for k in v:
            print i.min, i.max, k.min, k.max
            cmd = "python genRelationship.py -a " + str(i.min) + "," + str(i.max) + " -b " + str(k.min) + "," + str(k.max) + " -r " +str(rg)

            if i.min == k.min:
                scount = int(count * 0.5 *(i.max+1-i.min))
                cmd += " -c " + str(scount)
            else:
                scount =  int (mixcount/len(v)) * 1 *(i.max+1-i.min)
                cmd += " -c " + str( scount )
            cmd += " >> " + tmpfile
            print cmd
            runCmd(cmd)
    cmd = "cat " + tmpfile + " | sh ./resort.sh > rp.txt"
    runCmd(cmd)

    cmd = "python genUser.py  -a " + str(i.max)  + " -r " + str(rg) + " > user.txt"
    runCmd(cmd)

    cmd = 'rm -rf ' + tmpfile
    runCmd(cmd)
