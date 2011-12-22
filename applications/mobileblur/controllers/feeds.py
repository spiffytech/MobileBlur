# -*- coding: utf-8 -*-

from pprint import pprint
import time

def view():
    print ""
    s = time.time()
    feed_id = request.args[0]
    page = int(request.vars["page"]) if request.vars.has_key("page") else 1
    feed = newsblur.feed(feed_id, page=page)
    stories = feed["stories"]
    print time.time() - s

    if not feed.has_key("feed_title"):
        s = time.time()
        feeds = newsblur.feeds(flat=True)["feeds"]
        print time.time() - s

        s = time.time()
        feed = [feed for feed in feeds.itervalues() if feed["id"]==int(request.args[0])][0]
        print time.time() - s

    response.title = feed["feed_title"]

    return dict(stories=stories, feed=feed, feed_id=feed_id, page=page)


def mark_read():
    if len(request.args) > 0:
        newsblur.mark_feed_as_read(request.args[0])
    redirect(URL("default", "index"))
