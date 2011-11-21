newsblur = local_import("newsblur")
newsblur = newsblur.NewsBlur()

threshold = 0
thresholds = ["nt", "ps", "ng"]  # indices -1, 0, 1 for negative, neutral, and positive intelligence filters

if [request.application, request.controller, request.function] != [request.application, "default", "login"]:
    if "nb_cookie" not in request.cookies.keys():
        redirect(URL("default", "login"))
    else:
        newsblur.cookies["newsblur_sessionid"] = request.cookies["nb_cookie"].value
