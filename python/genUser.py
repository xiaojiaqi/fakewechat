#! /usr/bin/python

#from gen import *
from random import *
import sys
import getopt
import genUser
from rg import *

def create_user(count, rg):
    for i in range(1, count+1):
        print str(i)+","+str( i)+ "," + str(i)+"@fakewechat.com,"+ str( 13500000000 +i) + "," + "paswd" + str(i) + "," +  str(getRG(i,rg))

if __name__ == "__main__":
    if len(sys.argv) == 1:
        usage()
        sys.exit(1)
    try:
        opts, args = getopt.getopt(sys.argv[1:], "a:r:", ["help"])
    except getopt.GetoptError:
        usage()
        sys.exit(2)
    count = 0
    rg = 2000000
    for o, a in opts:
        if o in ('-h', '--help'):
            usage()
            sys.exit(1)
        elif o in ('-a' ):
            count = int(a) 
        elif o in ('-r'):
            rg = int(a)
    create_user(count, rg)


