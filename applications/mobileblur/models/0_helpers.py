import sys

import newsblur as nb_module
newsblur = nb_module.NewsBlur()
def controller_error_handler(f):
    try:
        return f()
    except nb_module.NewsblurException:
        redirect(URL("errors", "index"))
response._caller = controller_error_handler

if [request.application, request.controller, request.function] not in [
    [request.application, "default", "login"], 
    [request.application, "errors", "index"]
]:
    if "nb_cookie" not in request.cookies.keys():
        redirect(URL("default", "login"))
    else:
        newsblur.cookies["newsblur_sessionid"] = request.cookies["nb_cookie"].value

if not request.cookies.has_key("threshold"):
    threshold = 0
    response.cookies["threshold"] = threshold
    response.cookies["threshold"]["path"] = "/"
else:
    threshold = int(request.cookies["threshold"].value)
thresholds = ["nt", "ps", "ng"]  # indices -1, 0, 1 for negative, neutral, and positive intelligence filters


def get_intelligence_rating(story): 
    rating = sum([v for k,v in story["intelligence"].iteritems()])
    rating = 1 if rating >= 1 else -1 if rating < 0 else 0
    return rating

passes_intelligence = lambda story: True if get_intelligence_rating(story) >= threshold else False
intelligence_filter = lambda stories: [s for s in stories if passes_intelligence(s)]
