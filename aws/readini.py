#!/usr/bin/env python


import sys,os
import ConfigParser

_aws_access_key =  'xxxx'
_aws_secret_key = 'xxxx'

class readini:
  def __init__(self, file_path):
    cf = ConfigParser.ConfigParser()
    cf.read(file_path)
    global _aws_secret_key
    global _aws_access_key
    _aws_secret_key = cf.get("aws", "aws_secret_key")
    _aws_access_key = cf.get("aws", "aws_access_key")
    #print _aws_secret_key,  _aws_access_key

def get_aws_access_key():
    global _aws_access_key 
    return _aws_access_key
 

def get_aws_secret_key():
    global _aws_secret_key
    return _aws_secret_key

f = readini("config.ini")

pwd = os.path.split(os.path.realpath(__file__))[0]

pwd += "/package"
sys.path.append(pwd)
