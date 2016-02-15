#! /usr/bin/env python
# -*- coding:utf-8 -*-

import urllib
import httplib
import json

""" modify the Variables as you wish """
user = "Admin"
password = "zabbix"
zabbixServer = "127.0.0.1"
hostIp="172.16.11.194"


""" create the user login query object"""
def userLogin():
    global user
    global password
    login={}
    login["jsonrpc"] = "2.0"
    login["method"] = "user.login"
    params = {}
    params["user"] = user
    params["password"] = password
    login["id"] = 1
    login["params"] = params
    login["auth"]= None

    return login


""" create the create host query object"""
def createHost(hostIp, auth):
    host = {}

    host["jsonrpc"] = "2.0"
    host["method"] = "host.create"


    params = {}
    params["host"] = "FakeweChat Server"

    interfaces  = {}
    interfaces["type"] = 1
    interfaces["main"] =  1
    interfaces["useip"] = 1
    interfaces["ip"] = hostIp
    interfaces["dns"] = ""
    interfaces["port"] = "10050"

    params["interfaces"] = interfaces

    params["groups"] = []
    params["groups"].append({"groupid" : "2"})

    host["params"] = params
    host["auth"] = auth
    host["id"]  = 1

    return host


def createItem(hostId, itemkey,auth):
    item = {}
    item["jsonrpc"] = "2.0"
    item["method"] = "item.create"
    params = {}

    params["name"] = itemkey
    params["key_"] = itemkey


    params["hostid"] = hostId
    # we use active sender
    params["type"] = 7
    params["value_type"] = 3
    params["delay"] = 10
    item["params"] = params
    item["auth"] = auth
    item["id"]  = 1

    return item


def createGraph(graphname, itemkeys,auth):
    graph = {}
    graph["jsonrpc"] = "2.0"
    graph["method"] = "graph.create"
    params = {}

    params["name"] = graphname

    params["width"] = 900
    params["height"] = 200
    graph["show_3d"] = 1
    gitems = []
    order = 0
    for i in itemkeys:
        item = {}
        item["itemid"] = i["id"]
        item["color"]  = i["color"]
        item["sortorder"] = str(order)
        item["calc_fnc"] = 7
        order += 1
        gitems.append(item)

    params["gitems"] = gitems
    graph["params"] = params

    graph["auth"] = auth
    graph["id"]  = 1
    return graph

""" create the createscreen query obj"""
def createScreen(screenname, graphkeys,auth):
    screen = {}
    screen["jsonrpc"] = "2.0"
    screen["method"] = "screen.create"
    params = {}

    params["name"] = screenname
    params["hsize"] = 1
    params["vsize"] = 3
    screenitems = []
    order = 0
    for i in graphkeys:
        item = {}
        item["resourcetype"] = 0
        item["resourceid"]  = graphkeys[i]
        #item["rowspan"] = 0
        #item["colspan"] = 0
        item["dynamic"] = 1
        item["x"] = order%1
        item["y"] = order/1
        item["width"] = 900
        order += 1
        screenitems.append(item)

    params["screenitems"] = screenitems
    screen["params"] = params

    screen["auth"] = auth
    screen["id"]  = 1
    return screen

""" send the query and get the response """
def query(zabbixServer, query):

    try:
        requrl = "http://" + zabbixServer + "/zabbix/api_jsonrpc.php"
        headerdata = {"Content-Type":"application/json-rpc"}
        conn = httplib.HTTPConnection(zabbixServer)
        query = json.dumps(query)
        print query
        conn.request(method="POST",url=requrl,body=query, headers = headerdata)
        response = conn.getresponse()
        res= response.read()
        return res
    except Exception,e:
        print "query error", e
        exit()

if __name__ == '__main__':
    global zabbixServer
    # login to server
    loginParams = userLogin()
    js = query(zabbixServer, loginParams)
    print js
    login = json.loads(js)
    print login['result']
    authid = login['result']
    hostParams = createHost(hostIp, authid)

    js = query(zabbixServer, hostParams)
    print js
    host = json.loads(js)
    """ got the host id """
    hostId = host['result']['hostids'][0]

    item = []

    item.append("client.KeepConnection")
    item.append("client.Connection")
    item.append("client.http")
    item.append("client.http.login")
    item.append("client.http.failed")
    item.append("client.SumSebdPack")
    item.append("client.SendPack")
    item.append("client.SumRecvPack")
    item.append("client.RecvPack")

    item.append("longCon.KeepConnection")
    item.append("longCon.Connection")
    item.append("longCon.SumSebdPack")
    item.append("longCon.SendPack")
    item.append("longCon.SumRecvPack")
    item.append("longCon.RecvPack")


    item.append("core.SumSebdPack")
    item.append("core.SendPack")
    item.append("core.SumRecvPack")

    item.append("core.RecvPack")
    item.append("core.SumFailed")
    item.append("core.Failed")

    """ store the itemkey => itemid to itemmap   """
    itemmap = {}

    for i in item:
        itemParams = createItem(hostId, i, authid)
        js = query(zabbixServer, itemParams)
        print js
        itemResult = json.loads(js)
        """ got the itemids """
        itemmap[i] = itemResult["result"]["itemids"][0]
    print itemmap

    graphid = {}
    # keep.Conn
    # 00C800 green
    # FF0000 red
    # 3333FF blue
    color = {}

    color["red"] = "FF0000"
    color["green"] = "00C800"
    color["blue"] = "3333FF"
    """  create keep.Conn graph"""
    graphname="keep.Conn"
    graphmap = {}
    items=[]

    item = {"id":itemmap["client.KeepConnection"],"color":color["red"]}
    items.append(item)

    item = {"id":itemmap["longCon.KeepConnection"],"color":color["green"]}
    items.append(item)
    print items
    graphParams = createGraph(graphname, items, authid)
    print graphParams
    js = query(zabbixServer, graphParams)
    print js
    graphResult = json.loads(js)
    """ got the graphid of keep.Conn          """
    graphmap[graphname] = graphResult["result"]["graphids"][0]


    """ create sendPack graph          """
    graphname="SendPack"
    items=[]
    item = {"id":itemmap["client.SendPack"],"color":color["red"]}
    items.append(item)
    item = {"id":itemmap["longCon.SendPack"],"color":color["green"]}
    items.append(item)
    item = {"id":itemmap["core.SendPack"],"color":color["blue"]}
    items.append(item)
    print items
    graphParams = createGraph(graphname, items, authid)
    print graphParams
    js = query(zabbixServer, graphParams)
    print js
    graphResult = json.loads(js)
    """ got the graphid of SendPack  """
    graphmap[graphname ] = graphResult["result"]["graphids"][0]

    """ create Recvpack graph          """
    graphname="RecvPack"
    items=[]
    item = {"id":itemmap["client.RecvPack"],"color":color["red"]}
    items.append(item)
    item = {"id":itemmap["longCon.RecvPack"],"color":color["green"]}
    items.append(item)
    item = {"id":itemmap["core.RecvPack"],"color":color["blue"]}
    items.append(item)
    print items
    graphParams = createGraph(graphname, items, authid)
    print graphParams
    js = query(zabbixServer, graphParams)
    print js
    graphResult = json.loads(js)
    """ got the graphid of RecvPack"""
    graphmap[graphname ] = graphResult["result"]["graphids"][0]

    print graphmap

    """ create the screen"""
    screenParams = createScreen("screen", graphmap, authid)
    js = query(zabbixServer, screenParams)
    print js
    graphResult = json.loads(js)
