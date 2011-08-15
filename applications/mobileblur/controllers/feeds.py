# -*- coding: utf-8 -*-
# this file is released under public domain and you can use without limitations

from pprint import pprint

newsblur = local_import("newsblur")

username = ""
password = ""
threshold = 0
thresholds = ["nt", "ps", "ng"]  # indices -1, 0, 1 for negative, neutral, ane positive inhelligence filters

def view():
    newsblur.login(username, password)
    return 
