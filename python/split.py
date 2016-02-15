#!/usr/bin/env python

import sys
import getopt

class saveline:
    def __init__(self, range, prefix = ""):
        self.index = 1
        self.f = ""
        self.filename = ""
        self.range = range
        self.maxId = 0
        self.prefix = prefix

    def newfile(self):
        if self.prefix == "":
            self.filename ="%d.txt"%(self.index)
        else:
            self.filename = "%s%05d_%08d_%08d.txt"%(self.prefix, self.index,self.maxId+1, self.maxId + self.range)
        self.index = self.index + 1
        self.maxId += self.range
        self.f = open(self.filename,'wb')

    def close(self):
        self.f.close()

    def saveLine(self, line):
        v = line.split(",")
        if len(v) > 1 and v[0] != '':
            index = int(v[0])
            while True:
                if index <= self.maxId:
                    self.f.write(line)
                    break
                else:
                    self.close()
                    self.newfile()


def usage():
    print "split.py -k 1000 -p prefix "


if __name__ == "__main__":
    if len(sys.argv) == 1:
        usage()
        sys.exit(2)
    try:
        opts, args = getopt.getopt(sys.argv[1:], "k:p:h", ["help"])

    except getopt.GetoptError:
        usage()
        sys.exit(2)

    range = 0
    prefix = ""

    for o, a in opts:
        if o in ('-h', '--help'):
            usage()
            sys.exit(1)
        elif o in ('-k' ):
            range = int(a)
        elif o in ('-p'):
            prefix = a

    s = saveline(range, prefix)

    s.newfile()

    while True:
        line = sys.stdin.readline()
        #print line
        if line:
            s.saveLine(line)
        else:
            s.close()
            break
