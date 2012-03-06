#!/usr/bin/python

"""newsblur.py - An API wrapper library for newsblur.com"""

import simplejson
import requests

__author__ = 'Dananjaya Ramanayake <dananjaya86@gmail.com>, spiffytech <spiffytech@gmail.com>'
__version__ = "1.0"

nb_url = "http://www.newsblur.com/"


class NewsblurException(IOError):
    pass


class NewsBlur():
    def __init__(self):
        self.cookies = {}

    
    def _nb_get(self, url, payload=None):
        if payload is None:
            payload = {}

        try:
            results = requests.get(nb_url + url, params=payload, cookies=self.cookies)
        except:
            raise NewsblurException("Can't reach Newsblur right now")

        if results.status_code != 200:
            raise NewsblurException("Newsblur returned error code " + str(results.status_code) + ", " + results.content)

        decoded_results = simplejson.loads(results.content)
        return decoded_results
        

    def _nb_post(self, url, payload=None):
        if payload is None:
            payload = {}

        try:
            results = requests.post(nb_url + url, data=payload, cookies=self.cookies)
        except:
            raise NewsblurException("Can't reach Newsblur right now")

        if results.status_code != 200:
            raise NewsblurException("Newsblur returned error code " + str(results.status_code) + ", " + results.content)
        self.cookies = results.cookies

        decoded_results = simplejson.loads(results.content)
        return decoded_results


    def login(self, username, password):
        '''
        Login as an existing user.
        If a user has no password set, you cannot just send any old password. 
        Required parameters, username and password, must be of string type.
        '''

        url = "api/login"

        results = self._nb_post(url, payload={"username": username, "password": password})
        if results["authenticated"] is False:
            raise ValueError("The newsblur credentials you provided are invalid")

        return results

    def logout(self):
        '''
        Logout the currently logged in user.
        '''

        url = "api/logout"
        results = self._nb_get(url)
        return results

    def signup(self, username, password, email):
        '''
        Create a new user.
        All three required parameters must be of type string.
        '''

        url = "api/signup"
        payload = {"signup_username": username, "signup_password": password, "signup_email": email}
        results = self._nb_post(url, payload=payload)
        return results


    def search_feed(self, address, offset=1):
        '''
        
        Retrieve information about a feed from its website or RSS address.
        Parameter address must be of type string while parameter offset must be an integer.
        Will return a feed.
        
        '''

        url = "rss_feeds/search_feed"
        payload = {"address": address, "offset": offset}
        results = self._nb_get(url, payload=payload)
        return results


    def feeds(self, include_favicons=True, flat=False):
        '''
        Retrieve a list of feeds to which a user is actively subscribed.
        Includes the 3 unread counts (positive, neutral, negative), as well as optional favicons.
        '''

        url = "reader/feeds"
        payload = {"include_favicons": include_favicons, "flat": flat}
        results = self._nb_get(url, payload=payload)
        return results


    def favicons(self, feeds=[1,2,3]):
        '''
        Retrieve a list of favicons for a list of feeds. 
        Used when combined with /reader/feeds and include_favicons=false, so the feeds request contains far less data. 
        Useful for mobile devices, but requires a second request. 
        '''
        
        url = "reader/favicons"
        payload = {"feeds": feeds}
        results = self._nb_get(url, payload=payload)
        return results
        

    def id(self, id_no):
        '''
        Retrieve the original page from a single feed.
        '''
        
        url = "reader/page/%d" % id_no
        results = self._nb_get(url)
        return results
        

    def refresh_feeds(self, ):
        '''
        Up-to-the-second unread counts for each active feed.
        Poll for these counts no more than once a minute.
        '''

        url = "reader/refresh_feeds"
        results = self._nb_get(url)
        return results


    def feeds_trainer(self, feed_id):
        '''
        Retrieves all popular and known intelligence classifiers.
        Also includes user's own classifiers.
        '''

        url = "reader/feeds_trainer"
        payload = {"feed_id": feed_id}
        results = self._nb_get(url, payload=payload)
        return results


    def statistics(self, id_no):
        '''
        If you only want a user's classifiers, use /classifiers/:id.
        Omit the feed_id to get all classifiers for all subscriptions.
        '''

        url = "rss_feeds/statistics/%d" % id_no
        results = self._nb_get(url)
        return results


    def feed_autocomplete(self, term):
        '''
        Get a list of feeds that contain a search phrase.
        Searches by feed address, feed url, and feed title, in that order.
        Will only show sites with 2+ subscribers.
        '''

        url = "rss_feeds/feed_autocomplete"
        payload = {"term": term}
        results = self._nb_get(url, payload=payload)
        return results


    def feed(self, id, page=1):
        '''
        Retrieve stories from a single feed.
        '''

        url = "reader/feed/%s" % id
        payload = {"page": page}
        
        results = self._nb_get(url, payload=payload)
        content = results

        for story in range(len(content["stories"])):
            content["stories"][story]["page"] = page
        return content


    def starred_stories(self, page=1):
        '''
        Retrieve a user's starred stories.
        '''
        
        url = "reader/starred_stories"
        payload = {"page": page}
        results = self._nb_get(url, payload=payload)
        return results


    def river_stories(self, feeds, page=1, read_stories_count=0):
        '''
        Retrieve stories from a collection of feeds. This is known as the River of News.
        Stories are ordered in reverse chronological order.
        '''

        url = "reader/river_stories"
        payload = {"feeds": feeds, "page": page, "read_stories_count": read_stories_count}
        results = self._nb_get(url, payload=payload)
        return results


    def mark_story_as_read(self, story_id, feed_id):
        '''
        Mark stories as read.
        Multiple story ids can be sent at once.
        Each story must be from the same feed.
        '''

        url = "reader/mark_story_as_read"
        payload = {"story_id": story_id, "feed_id": feed_id}
        results = self._nb_post(url, payload=payload)
        return results


    def mark_story_as_unread(self, story_id, feed_id):
        '''
        Mark stories as read.
        Multiple story ids can be sent at once.
        Each story must be from the same feed.
        '''

        url = "reader/mark_story_as_unread"
        payload = {"story_id": story_id, "feed_id": feed_id}
        results = self._nb_post(url, payload=payload)
        return results


    def mark_story_as_starred(self, story_id, feed_id):
        '''
        Mark a story as starred (saved).
        '''
        
        url = "reader/mark_story_as_starred"
        payload = {"story_id": story_id, "feed_id": feed_id}
        results = self._nb_post(url, payload=payload)
        return results


    def mark_all_as_read(self, days=0):
        '''
        Mark all stories in *all* feeds read.
        '''
        
        url = "reader/mark_all_as_read"
        payload = {"days": days}
        results = self._nb_post(url, payload=payload)
        return results


    def add_url(self, feed_url, folder='[Top Level]'):
        '''
        Add a feed by its URL. 
        Can be either the RSS feed or the website itself.
        '''

        url = "reader/add_url"
        feed_url = feed_url.strip("/")
        payload = {"url": feed_url, "folder": folder}
        results = self._nb_post(url, payload=payload)
        return results


    def add_folder(self, folder, parent_folder='[Top Level]'):
        '''
        Add a new folder.
        '''
        
        url = "reader/add_folder"
        payload = {"folder": folder, "parent_folder": parent_folder}
        results = self._nb_post(url, payload=payload)
        return results


    def rename_feed(self, feed_title, feed_id):
        '''
        Rename a feed title. Only the current user will see the new title.
        '''
        
        url = "reader/rename_feed"
        payload = {"feed_title": feed_title, "feed_id": feed_id}
        results = self._nb_post(url, payload=payload)
        return results


    def delete_feed(self, feed_id, in_folder):
        '''
        Unsubscribe from a feed. Removes it from the folder.
        Set the in_folder parameter to remove a feed from the correct folder, in case the user is subscribed to the feed in multiple folders.
        '''    

        url = "reader/delete_feed"
        payload = {"feed_id": feed_id, "in_folder": in_folder}
        results = self._nb_post(url, payload=payload)
        return results


    def rename_folder(self, folder_to_rename, new_folder_name, in_folder):
        '''
        Rename a folder.
        '''
        
        url = "reader/rename_folder"
        payload = {"folder_to_rename": folder_to_rename, "new_folder_name": new_folder_name, "in_folder": in_folder}
        results = self._nb_post(url, payload=payload)
        return results


    def delete_folder(self, folder_to_delete, in_folder, feed_id):
        '''
        Delete a folder and unsubscribe from all feeds inside.
        '''
        
        url = "reader/delete_folder"
        payload = {"folder_to_delete": folder_to_delete, "in_folder": in_folder, "feed_id": feed_id}
        results = requests.post(url, data=payload)
        return results


    def mark_feed_as_read(self, feed_id):
        '''
        Mark a list of feeds as read.
        '''
        
        url = "reader/mark_feed_as_read"
        payload = {"feed_id": feed_id}
        results = self._nb_post(url, payload=payload)
        return results


    def save_feed_order(self, folders):
        '''
        Reorder feeds and move them around between folders.
        The entire folder structure needs to be serialized.
        '''

        url = "reader/save_feed_order"
        payload = {"folders": folders}
        results = self._nb_post(url, payload=payload)
        return results


    def classifier(self, id_no):
        '''
        Get the intelligence classifiers for a user's site.
        Only includes the user's own classifiers. 
        Use /reader/feeds_trainer for popular classifiers.
        '''

        url = "classifier/%d' % id_n"
        results = self._nb_get(url)
        return results


    def classifier_save(self, feed_id, likes, dislikes):
        '''
        Save intelligence classifiers (tags, titles, authors, and the feed) for a feed.

        Expects dicts in the form of:

        likes = {
            "tag": ["food", "candy", "fun"],
            "author": ["A. Blogger"],
            "title": ["Awesome Title"]
        }
        '''

        payload = {}
        for _type in likes:
            payload["like_" + _type] = likes[_type]
        for _type in dislikes:
            payload["dislike_" + _type] = dislikes[_type]

        payload["feed_id"] = feed_id
        
        url = "classifier/save"
        results = self._nb_post(url, payload=payload)
        return results


    def opml_export(self, ):
        '''
        Download a backup of feeds and folders as an OPML file.
        Contains folders and feeds in XML; useful for importing in another RSS reader.
        '''
        
        url = "import/opml_export"
        results = self._nb_get(url)
        return results



    def opml_upload(self, opml_file):
        '''
        Upload an OPML file.
        '''
        
        url = "import/opml_upload"
        f = open(opml_file)
        payload = {"file": f}
        f.close()
        results = self._nb_post(url, payload=payload)
        return results
