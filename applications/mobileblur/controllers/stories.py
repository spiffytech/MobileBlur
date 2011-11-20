# -*- coding: utf-8 -*-

from pprint import pprint

def view():
    stories = newsblur.feed(request.vars["feed_id"])["stories"]
    story = [story for story in stories if story["id"]==request.vars["story"]][0]
    return dict(story=story)
