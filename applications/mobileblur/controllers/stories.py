# -*- coding: utf-8 -*-

from pprint import pprint

def view():
    requested_story_id = request.vars["story"]
    stories = newsblur.feed(request.vars["feed_id"])["stories"]
    
    previous_story = None
    requested_story = None
    next_story = None
    for story in range(len(stories)):
        if stories[story]["id"] == requested_story_id:
            requested_story = stories[story]
            try:
                previous_story = stories[story+1]
            except IndexError:
                pass
            if story != 0:
                next_story = stories[story-1]

            break

    return dict(
        previous_story=previous_story,
        requested_story=requested_story, 
        next_story=next_story,
        feed_id=request.vars["feed_id"]
    )

def mark_read():
    results = newsblur.mark_story_as_read(request.vars["story_id"], request.vars["feed_id"])
    redirect(URL("feeds", "view", args=[request.vars["feed_id"]]))
