# -*- coding: utf-8 -*-

from pprint import pprint

def view():
    page = int(request.vars["page"]) if request.vars.has_key("page") else 1
    requested_story_id = request.vars["story"]
    feed_id = request.vars["feed_id"]
    feed = newsblur.feed(feed_id, page=page)
    stories = feed["stories"]
    
    previous_story = None
    requested_story = None
    next_story = None
    for story in range(len(stories)):
        if stories[story]["id"] == requested_story_id:
            requested_story = stories[story]

            stories_page = stories[story+1:]
            filtered_stories = intelligence_filter(stories_page)
            if len(filtered_stories) == 0:
                stories_page = newsblur.feed(feed_id, page=page+1)["stories"]
                filtered_stories = intelligence_filter(stories_page)
            if len(filtered_stories) > 0:
                previous_story = filtered_stories[0]

            stories_page = stories[:story]
            filtered_stories = intelligence_filter(stories_page)
            if len(filtered_stories) == 0:
                stories_page = newsblur.feed(feed_id, page=page-1)["stories"]
                filtered_stories = intelligence_filter(stories_page)
            if len(filtered_stories) > 0 and filtered_stories[-1]["id"] != previous_story["id"]:  # Fetching page 0 returns page 1, complicating things
                next_story = filtered_stories[-1]

            break

    newsblur.mark_story_as_read(requested_story_id, feed_id)


    feed_title = request.vars["feed_title"]
    response.title = (requested_story["story_title"][:15] + "...") if len(requested_story["story_title"]) > 15 else requested_story["story_title"]
    response.title += " on "
    response.title += (feed_title[:15] + "...") if len(feed_title) > 15 else feed_title
    
    return dict(
        previous_story=previous_story,
        requested_story=requested_story, 
        next_story=next_story,
        feed_id=feed["feed_id"],
        feed_title=feed_title,
    )

def mark_unread():
    results = newsblur.mark_story_as_unread(request.vars["story_id"], request.vars["feed_id"])
    redirect(URL("feeds", "view", args=[request.vars["feed_id"]], vars={"page": request.vars["page"]}))
