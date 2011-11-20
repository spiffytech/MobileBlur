from pprint import pprint
import simplejson

def index():
    raw_feeds = newsblur.feeds(flat=True)["feeds"]
    feeds = {}
    for feed in raw_feeds.itervalues():
        for i in range(threshold, 2):
            if feed[thresholds[i]] > 0:
                feeds[feed["feed_title"]] = feed
                break

    return dict(feeds=feeds, threshold=threshold)


def login():
    login_form = SQLFORM.factory(
        Field("username", requires=IS_NOT_EMPTY()),
        Field("password", "password", requires=IS_NOT_EMPTY())
    )
    if login_form.accepts(request):
        try:
            results = newsblur.login(login_form.vars["username"], login_form.vars["password"])
            response.cookies["nb_cookie"] = newsblur.cookies["newsblur_sessionid"]
            response.cookies["nb_cookie"]["path"] = "/"
            redirect(URL("index"))
        except Exception as ex:
            login_form.insert(-1, ex.message)

    return dict(login_form=login_form)


def logout():
    response.cookies["nb_cookie"] = ""
    response.cookies["nb_cookie"]["expires"] = -10
    response.cookies["nb_cookie"]["path"] = "/"
    redirect(URL("index"))
