#!/usr/bin/env python

#from gen import *
from random import *
import sys
import getopt
from rg import *

# simple ranges
class ranges:
    def __init__(self):
        self.min = 0
        self.max = 0

def create_relationship(idrange_1, idrange_2, numberofcryles, rg):
    result = ""
    for i in range(0, (numberofcryles)) :
        id1 = randint(idrange_1.min, idrange_1.max)
        id2 = 0
        t = 0
        while (True):
            id2 = randint(idrange_2.min, idrange_2.max)
            if id2 != a and t < 10: #hard code
                break
            t = t + 1
        id3 = getRG(id2, rg)
        # you can't make friend to yourself
        if id1 == id2:
            continue
        result += str(id1)+","+str(id2)+"," +str(id3) + "\n"
        id3 = getRG(id1,rg)
        result += str(id2)+","+ str(id1)+","+str(id3) + "\n"
        if  i % 9999 == 0:
            print result
            result = ""
    if len(result):
        print result

def usage():
    print "genRelationship.py -a 1,1000 -b 1,1000 -c 100000 -r 2000000"
    print "-a mean rg1, -b means rg2, -c mean replation of two rg, -r mean rg's range"

if __name__ == "__main__":
    if len(sys.argv) == 1:
        usage()
        sys.exit(1)
    try:
        opts, args = getopt.getopt(sys.argv[1:], "a:b:c:r:h", ["help"])
    except getopt.GetoptError:
        usage()
        sys.exit(2)
    id1 = ranges()
    id2 = ranges()
    count = 0
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
        elif o in ('-b'):
            v = a.split(',')
            if len(v) >= 2:
                id2.min = int(v[0])
                id2.max = int(v[1])

        elif o in ('-c'):
            count = int(a)
        elif o in ('-r'):
            rg = int(a)

    create_relationship(id1, id2, count, rg)
