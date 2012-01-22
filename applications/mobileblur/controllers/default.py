from pprint import pprint
import random
import simplejson
import time

random.seed(time.time())

def index():
    raw_feeds = newsblur.feeds(flat=True)["feeds"]
    feeds = []
    for feed in raw_feeds.itervalues():
        for i in range(threshold, 2):
            if feed[thresholds[i]] > 0:
                feeds.append(feed)
                break

    feeds.sort(key=lambda f: str.lower(f["feed_title"]))
    return dict(feeds=feeds, threshold=threshold)


def login():
    demo = False
    demo_accounts = ["mbtest1", "mbtest2", "mbtest3", "mbtest4"]
#    demo_accounts = ["mbtest1"]

    login_form = SQLFORM.factory(
        Field("username", requires=IS_NOT_EMPTY()),
        Field("password", "password")
    )

    if len(request.args) > 0 and request.args[0] == "demo":
        demo = True
        account = random.choice(demo_accounts)
        login_form.vars["username"] = account
        login_form.vars["password"] = account

    if login_form.accepts(request) or demo is True:
        try:
            results = newsblur.login(login_form.vars["username"], login_form.vars["password"])
            response.cookies["nb_cookie"] = newsblur.cookies["newsblur_sessionid"]
            response.cookies["nb_cookie"]["path"] = "/"
            session.flash = "Welcome to MobileBlur!"
            redirect(URL("index"))
        except ValueError as ex:
            login_form.insert(-1, ex.message)
            login_form._class = "alert-message block-message error"
            response.flash = "Oh noes! You made a mistake! Try again."

    return dict(login_form=login_form)


def logout():
    response.cookies["nb_cookie"] = ""
    response.cookies["nb_cookie"]["expires"] = -10  
    response.cookies["nb_cookie"]["path"] = "/"
    session.flash = "You have been logged out"
    redirect(URL("default", "login"))


def settings():
    threshold_form = SQLFORM.factory(
        Field(
            "threshold", 
            "integer", 
            requires=IS_IN_SET([
                (-1, "Negative"),
                (0, "Neutral"),
                (1, "Positive")
            ]),
            default=threshold,
            widget=SQLFORM.widgets.radio.widget
        ),
        _name="threshold_form"
    )
    if threshold_form.process(formname="threshold_form").accepted:
        response.cookies["threshold"] = threshold_form.vars.threshold
        response.cookies["threshold"]["path"] = "/"
        session.flash = "Settings saved! New threshold is '%s'" % thresholds[int(response.cookies["threshold"].value)]
        redirect(URL("index"))

    add_feed_form = SQLFORM.factory(
        Field("feed_url", requires=IS_URL()),
        _name="add_feed_form"
    )
    if add_feed_form.process(formname="add_feed_form").accepted:
        resp = newsblur.add_url(add_feed_form.vars.feed_url)
        if resp["result"] == "ok":
            response.flash = "Successfully added feed"
            redirect(URL("index"))
        else:
            response.flash = "Something at Newsblur went wrong while adding that feed"

    return dict(threshold_form=threshold_form, add_feed_form=add_feed_form)
