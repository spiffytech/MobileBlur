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

            # Determine the next oldest story
            stories_page = stories[story+1:]
            filtered_stories = intelligence_filter(stories_page)
            if len(filtered_stories) == 0:
                stories_page = newsblur.feed(feed_id, page=page+1)["stories"]
                filtered_stories = intelligence_filter(stories_page)
            if len(filtered_stories) > 0:
                previous_story = filtered_stories[0]

            # Determine the next newest story
            stories_page = stories[:story]
            filtered_stories = intelligence_filter(stories_page)
            if len(filtered_stories) == 0:
                page = page - 1
                stories_page = newsblur.feed(feed_id, page=page)["stories"]
                filtered_stories = intelligence_filter(stories_page)
            if len(filtered_stories) > 0 and page != 0:  # Fetching page 0 returns page 1, complicating things
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


def intelligence():
    feed_id = request.vars["feed_id"]
    requested_story_id = request.vars["story_id"]
    page = request.vars["page"]

    feed = newsblur.feed(feed_id, page=page)
    classifiers = feed["classifiers"]
    stories = feed["stories"]
    for story in range(len(stories)):
        if stories[story]["id"] == requested_story_id:
            requested_story = stories[story]
    tags = requested_story["story_tags"]

    items = []
    items.append(H4("Title"))
    title = requested_story["story_title"]
    items.extend([
        INPUT(_name="title", _value=title),
        LABEL("Like"),
        INPUT(_type="radio", _name="title_rating", _value="Like"),
        LABEL("Dislike"),
        INPUT(_type="radio", _name="title_rating", _value="Dislike"),
        BR()
    ])

    items.append(H4("Tags"))
    for tag in tags:
        try:
            rating = classifiers["tags"][tag]
            value = "Like" if rating == 1 else "Dislike"
        except KeyError:
            value = ""
        t = [
            LABEL(tag + ": "),
            LABEL("Like"),
            INPUT(_type="radio", _name=tag+"][tag", _value="Like", value=value),
            LABEL("Dislike"),
            INPUT(_type="radio", _name=tag+"][tag", _value="Dislike", value=value),
            BR()
        ]
        items.extend(t)

    author = requested_story["story_authors"]
    items.append(H4("Author"))
    items.extend([
        INPUT(_name="author", _value=author),
        LABEL("Like"),
        INPUT(_type="radio", _name="author_rating", _value="Like"),
        LABEL("Dislike"),
        INPUT(_type="radio", _name="author_rating", _value="Dislike"),
        BR()
    ])
    items.append(INPUT(_type="submit"))
    intel_form = FORM(*items)

    if intel_form.accepts(request, session):
        ratings = {
            "Like": {
                "title": [],
                "tag": [],
                "author": [],
            },
            "Dislike": {
                "title": [],
                "tag": [],
                "author": [],
            }
        }
        for rating_k in ratings:
            for form_k, form_v in intel_form.vars.iteritems():
                if form_k.endswith("][tag"):
                    if form_v == rating_k:
                        ratings[rating_k]["tag"].append(form_k[:-5])  # 5 = len("][tag")
            if intel_form.vars["title_rating"] == rating_k:
                ratings[rating_k]["title"].append(intel_form.vars["title"])
            if intel_form.vars["author_rating"] == rating_k:
                ratings[rating_k]["author"].append(intel_form.vars["author"])
        import ipdb
        ipdb.set_trace()

    return dict(requested_story=feed, intel_form=intel_form)


def mark_unread():
    results = newsblur.mark_story_as_unread(request.vars["story_id"], request.vars["feed_id"])
    session.flash = "Story left unread"
    redirect(URL("feeds", "view", args=[request.vars["feed_id"]], vars={"page": request.vars["page"]}))
