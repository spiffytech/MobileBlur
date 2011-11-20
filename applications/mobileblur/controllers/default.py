from pprint import pprint

def index():
    raw_feeds = newsblur.feeds(flat=True)["feeds"]
    feeds = {}
    for feed in raw_feeds.itervalues():
        for i in range(threshold, 2):
            if feed[thresholds[i]] > 0:
                feeds[feed["feed_title"]] = feed
                break

    pprint(feeds)
    return dict(feeds=feeds, threshold=threshold)
