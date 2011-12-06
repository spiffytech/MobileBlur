# -*- coding: utf-8 -*-

from pprint import pprint

def view():
    stories = newsblur.feed(request.vars["feed_id"])["stories"]
    story = [story for story in stories if story["id"]==request.vars["story"]][0]
    return dict(story=story, feed_id=request.vars["feed_id"])

def mark_read():
    results = newsblur.mark_story_as_read(request.vars["story_id"], request.vars["feed_id"])
    redirect(URL("feeds", "view", args=[request.vars["feed_id"]]))
