#! /usr/bin/python

import sys

def getRG(userid, rg):
    if userid % rg  == 0 :
        return userid /rg
    else:
        return userid/rg + 1

