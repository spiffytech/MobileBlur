import sys

nb_module = local_import("newsblur")
newsblur = nb_module.NewsBlur()

def my_handler(f):
    try:
        return f()
    except nb_module.NewsblurException:
        redirect("http://spiffyte.ch")
response._caller = my_handler

if not request.cookies.has_key("threshold"):
    threshold = 0
    response.cookies["threshold"] = threshold
    response.cookies["threshold"]["path"] = "/"
else:
    threshold = int(request.cookies["threshold"].value)
thresholds = ["nt", "ps", "ng"]  # indices -1, 0, 1 for negative, neutral, and positive intelligence filters

if [request.application, request.controller, request.function] != [request.application, "default", "login"]:
    if "nb_cookie" not in request.cookies.keys():
        redirect(URL("default", "login"))
    else:
        newsblur.cookies["newsblur_sessionid"] = request.cookies["nb_cookie"].value

passes_intelligence = lambda story: True if sum([v for k,v in story["intelligence"].iteritems()]) >= threshold else False
intelligence_filter = lambda stories: [s for s in stories if passes_intelligence(s)]
