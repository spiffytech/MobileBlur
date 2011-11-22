# -*- coding: utf-8 -*-

from pprint import pprint
import time

def view():
    print ""
    s = time.time()
    feed = newsblur.feed(request.args[0])
    stories = feed["stories"]
    print time.time() - s

    print feed.keys()

    if not feed.has_key("feed_title"):
        s = time.time()
        feeds = newsblur.feeds(flat=True)["feeds"]
        print time.time() - s

        s = time.time()
        feed = [feed for feed in feeds.itervalues() if feed["id"]==int(request.args[0])][0]
        print time.time() - s

    return dict(stories=stories, feed=feed)


def mark_read():
    if len(request.args) > 0:
        newsblur.mark_feed_as_read(request.args[0])
    redirect(URL("default", "index"))
