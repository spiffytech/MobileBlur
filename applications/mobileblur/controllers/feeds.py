# -*- coding: utf-8 -*-

from pprint import pprint

def view():
    stories = newsblur.feed(request.args[0])["stories"]
    feeds = newsblur.feeds(flat=True)["feeds"]
    feed = [feed for feed in feeds.itervalues() if feed["id"]==int(request.args[0])][0]
    return dict(stories=stories, feed=feed)

def mark_read():
    newsblur.mark_feed_as_read(request.vars["feed"])
    redirect(URL("default", "index"))
