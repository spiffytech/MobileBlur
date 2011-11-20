newsblur = local_import("newsblur")

threshold = 0
thresholds = ["nt", "ps", "ng"]  # indices -1, 0, 1 for negative, neutral, and positive intelligence filters

def login(username="spiffytech"):
    user = db(db.users.username==username).select().first()
    if user["cookie"] is None:
        results = newsblur.login(user["username"], user["password"])
