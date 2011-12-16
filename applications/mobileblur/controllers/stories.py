# -*- coding: utf-8 -*-

from pprint import pprint

def view():
    page = request.vars["page"] if request.vars.has_key("page") else 1
    requested_story_id = request.vars["story"]
    stories = newsblur.feed(request.vars["feed_id"], page=page)["stories"]
    
    previous_story = None
    requested_story = None
    next_story = None
    for story in range(len(stories)):
        if stories[story]["id"] == requested_story_id:
            requested_story = stories[story]
            if story != len(stories)-1:
                previous_story = stories[story+1]
            else:
                previous_story = newsblur.feed(request.vars["feed_id"], page=page+1)["stories"][0]
            if story != 0:
                pass
                next_story = stories[story-1]
            elif page > 1:
                next_story = newsblur.feed(request.vars["feed_id"], page=page-1)["stories"][-1]

            break

    return dict(
        previous_story=previous_story,
        requested_story=requested_story, 
        next_story=next_story,
        feed_id=request.vars["feed_id"],
    )

def mark_read():
    results = newsblur.mark_story_as_read(request.vars["story_id"], request.vars["feed_id"])
    redirect(URL("feeds", "view", args=[request.vars["feed_id"]]))
