# -*- coding: utf-8 -*-
# this file is released under public domain and you can use without limitations

from pprint import pprint

newsblur = local_import("newsblur")

username = ""
password = ""
threshold = 0
thresholds = ["nt", "ps", "ng"]  # indices -1, 0, 1 for negative, neutral, ane positive inhelligence filters

def index():
    newsblur.login(username, password)
    raw_feeds = newsblur.feeds(flat=True)["feeds"]
    feeds = {}
    for feed in raw_feeds.itervalues():
        for i in range(threshold, 2):
            if feed[thresholds[i]] > 0:
                feeds[feed["feed_title"]] = feed
                break

    pprint(feeds)
    return dict(feeds=feeds, threshold=threshold)
    
def login():
    pass
