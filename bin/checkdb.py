#!/usr/bin/env python

import sys
import getopt
import redis
import message_pb2
import time
import random

def usage():
    print "checkdb.py -m host -p port -l 1 -g 10"


if __name__ == "__main__":
    try:
        opts, args = getopt.getopt(sys.argv[1:], "m:p:h:l:g:", ["help"])

    except getopt.GetoptError:
        usage()
        sys.exit(2)

    hostname = "localhost"
    port = 6379
    min = 0
    max = 1
    for o, a in opts:
        if o in ('-h', '--help'):
            usage()
            sys.exit(1)
        elif o in ('-m' ):
            hostname = a
        elif o in ('-p'):
            port = int(a)
        elif o in ('-l'):
            min = int(a)
        elif o in ('-g'):
            max = int(a)
    r = redis.Redis(host=hostname,port= port  ,db=0)
    num = 0
    userid = ""
    datasize = 0
    tnow = time.time()
    for i in range(min, max+1):
        username = "user#" + str(i)
        data = r.get(username)
        a = message_pb2.UserInfor()
        #print data
        a.ParseFromString((data))
        #print a
        #print "user", i, a.SendId, a.SendAckId, a.ReceiveId, len(a.UserMap)*5
        if (a.SendId != a.SendAckId) or (a.ReceiveId != a.SendId):
            print "wrong", "user", i, a.SendId, a.SendAckId, a.ReceiveId, len(a.UserMap)*5
            print a
            os.exit(1)
        else:
            user = "C#inbox" + str(i)
            list = r.lrange(user,start=0,end=-1)
            #print "inbox " , user, " ", len(list) ," want ", len(a.UserMap)*5
            if len(list) != len(a.UserMap)*5:
                break
            for item in list:
                mess = message_pb2.GeneralMessage()
                mess.ParseFromString( item )
                if i != mess.ReceiverId:
                    print "wrong", i, mess.ReceiverId
                    os.exit(1)
                #print mess.messageType,   mess.ReceiverId, mess.SenderId, mess.SendId
            #print "outbox"
            user = "C#outbox" + str(i)
            list = r.lrange(user,start=0,end=-1)
            #print "outbox " , user, " ", len(list) ," want ", len(a.UserMap)*5
            if len(list) != len(a.UserMap)*5:
                print len(list) != len(a.UserMap)*5
                os.exit(1)

            for item in list:
                mess = message_pb2.GeneralMessage()
                mess.ParseFromString( item )
                if i != mess.SenderId:
                    print "wrong", i, mess.SenderId
                    os.exit(1)
                #print mess.messageType, mess.SenderId,  mess.ReceiverId, mess.SendId
    print "spend ", time.time() - tnow
