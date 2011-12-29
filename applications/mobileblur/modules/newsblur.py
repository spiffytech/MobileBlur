#!/usr/bin/python

"""newsblur.py - An API wrapper library for newsblur.com"""

import simplejson
import requests

__author__ = 'Dananjaya Ramanayake <dananjaya86@gmail.com>, spiffytech <spiffytechgmail.com>'
__version__ = "0.1"

nb_url = "http://www.newsblur.com/"

class NewsBlur():
    def __init__(self):
        self.cookies = {}

    def login(self, username,password):
        '''
        Login as an existing user.
        If a user has no password set, you cannot just send any old password. 
        Required parameters, username and password, must be of string type.
        '''

        url = nb_url + 'api/login'
        results = requests.post(url, data={"username": username, "password": password})
        self.cookies = results.cookies
        results = simplejson.loads(results.content)
        if results["authenticated"] is False:
            raise ValueError("The newsblur credentials you provided are invalid")
        return results

    def logout(self, ):
        '''
        Logout the currently logged in user.
        '''

        url = nb_url + 'api/logout'
        results = requests.get(url, cookies=self.cookies)
        return simplejson.loads(results.content)

    def signup(self, username,password,email):
        '''
        Create a new user.
        All three required parameters must be of type string.
        '''

        url = nb_url + 'api/signup'
        payload = {'signup_username':username,'signup_password':password,'signup_email':email}
        results = requests.post(url, data=payload, cookies=self.cookies)
        return simplejson.loads(results.content)

    def search_feed(self, address,offset=1):
        '''
        
        Retrieve information about a feed from its website or RSS address.
        Parameter address must be of type string while parameter offset must be an integer.
        Will return a feed.
        
        '''

        url = nb_url + 'rss_feeds/search_feed'
        payload = {'address':address,'offset':offset}
        results = results.get(url, data=payload, cookies=self.cookies)
        return simplejson.loads(results.content)

    def feeds(self, include_favicons=True,flat=False):
        '''
        Retrieve a list of feeds to which a user is actively subscribed.
        Includes the 3 unread counts (positive, neutral, negative), as well as optional favicons.
        '''
        
        url = nb_url + 'reader/feeds'
        payload = {'include_favicons':include_favicons,'flat':flat}
        results = requests.get(url, params=payload, cookies=self.cookies)
        return simplejson.loads(results.content)


    def favicons(self, feeds=[1,2,3]):
        '''
        Retrieve a list of favicons for a list of feeds. 
        Used when combined with /reader/feeds and include_favicons=false, so the feeds request contains far less data. 
        Useful for mobile devices, but requires a second request. 
        '''
        
        url = nb_url + 'reader/favicons'
        payload = {'feeds':feeds}
        results = requests.get(url, params=payload, cookies=self.cookies)
        return simplejson.loads(results.content)
        
    def id(self, id_no):
        '''
        Retrieve the original page from a single feed.
        '''
        
        url = nb_url + 'reader/page/' % id_no
        payload = {}
        results = requests.get(url, params=payload, cookies=self.cookies)
        return simplejson.loads(results.content)

    def refresh_feeds(self, ):
        '''
        Up-to-the-second unread counts for each active feed.
        Poll for these counts no more than once a minute.
        '''

        url = nb_url + 'reader/refresh_feeds'
        results = requests.get(url, cookies=self.cookies)
        return simplejson.loads(results.content)

    def feeds_trainer(self, feed_id):
        '''
        Retrieves all popular and known intelligence classifiers.
        Also includes user's own classifiers.
        '''

        url = nb_url + 'reader/feeds_trainer'
        payload = {'feed_id':feed_id}
        results = requests.get(url, params=payload, cookies=self.cookies)
        return simplejson.loads(results.content)

    def statistics(self, id_no):
        '''
        If you only want a user's classifiers, use /classifiers/:id.
        Omit the feed_id to get all classifiers for all subscriptions.
        '''

        url = nb_url + 'rss_feeds/statistics/%d' % id_no
        results = requests.get(url, cookies=self.cookies)
        return simplejson.loads(results.content)

    def feed_autocomplete(self, term):
        '''
        Get a list of feeds that contain a search phrase.
        Searches by feed address, feed url, and feed title, in that order.
        Will only show sites with 2+ subscribers.
        '''

        url = nb_url + 'rss_feeds/feed_autocomplete?%'
        payload = {'term':term}
        results = requests.get(url, params=payload, cookies=self.cookies)
        return simplejson.loads(results.content)

    def feed(self, id, page=1):
        '''
        Retrieve stories from a single feed.
        '''

        url = nb_url + 'reader/feed/%s' % id
        payload = {"page": page}
        results = requests.get(url, params=payload, cookies=self.cookies)
        content = simplejson.loads(results.content)
        for story in range(len(content["stories"])):
            content["stories"][story]["page"] = page
        return content

    def starred_stories(self, page=1):
        '''
        Retrieve a user's starred stories.
        '''
        
        url = nb_url + 'reader/starred_stories'
        payload = {'page':page}
        results = requests.get(url, params=payload, cookies=self.cookies)
        return simplejson.loads(results.content)

    def river_stories(self, feeds,page=1,read_stories_count=0):
        '''
        Retrieve stories from a collection of feeds. This is known as the River of News.
        Stories are ordered in reverse chronological order.
        '''

        url = nb_url + 'reader/river_stories'
        payload = {'feeds':feeds,'page':page,'read_stories_count':read_stories_count}
        results = urllib2.urlopen(url, data=payload, cookies=self.cookies)
        return simplejson.loads(results.content)

    def mark_story_as_read(self, story_id,feed_id):
        '''
        Mark stories as read.
        Multiple story ids can be sent at once.
        Each story must be from the same feed.
        '''

        url = nb_url + 'reader/mark_story_as_read'
        payload = {'story_id':story_id,'feed_id':feed_id}
        results = requests.post(url, data=payload, cookies=self.cookies)
        return simplejson.loads(results.content)

    def mark_story_as_unread(self, story_id,feed_id):
        '''
        Mark stories as read.
        Multiple story ids can be sent at once.
        Each story must be from the same feed.
        '''

        url = nb_url + 'reader/mark_story_as_unread'
        payload = {'story_id':story_id,'feed_id':feed_id}
        results = requests.post(url, data=payload, cookies=self.cookies)
        return simplejson.loads(results.content)

    def mark_story_as_starred(self, story_id,feed_id):
        '''
        Mark a story as starred (saved).
        '''
        
        url = nb_url + 'reader/mark_story_as_starred'
        payload = {'story_id':story_id,'feed_id':feed_id}
        results = requests.post(url, data=payload, cookies=self.cookies)
        return simplejson.loads(results.content)

    def mark_all_as_read(self, days=0):
        '''
        Mark all stories in *all* feeds read.
        '''
        
        url = nb_url + 'reader/mark_all_as_read'
        payload = {'days':days}
        results = requests.post(url, data=payload, cookies=self.cookies)
        return simplejson.loads(results.content)

    def add_url(self, feed_url,folder='[Top Level]'):
        '''
        Add a feed by its URL. 
        Can be either the RSS feed or the website itself.
        '''

        url = nb_url + 'reader/add_url'
        feed_url = feed_url.strip("/")
        payload = {'url': feed_url, 'folder': folder}
        results = requests.post(url, data=payload, cookies=self.cookies)
        print results.content
        return simplejson.loads(results.content)


    def add_folder(self, folder,parent_folder='[Top Level]'):
        '''
        Add a new folder.
        '''
        
        url = nb_url + 'reader/add_folder'
        payload = {'folder':folder,'parent_folder':parent_folder}
        results = requests.post(url, data=payload, cookies=self.cookies)
        return simplejson.loads(results.content)

    def rename_feed(self, feed_title,feed_id):
        '''
        Rename a feed title. Only the current user will see the new title.
        '''
        
        url = nb_url + 'reader/rename_feed'
        payload = {'feed_title':feed_title,'feed_id':feed_id}
        results = requests.post(url, data=payload, cookies=self.cookies)
        return simplejson.loads(results.content)

    def delete_feed(self, feed_id,in_folder):
        '''
        Unsubscribe from a feed. Removes it from the folder.
        Set the in_folder parameter to remove a feed from the correct folder, in case the user is subscribed to the feed in multiple folders.
        '''    

        url = nb_url + 'reader/delete_feed'
        payload = {'feed_id':feed_id,'in_folder':in_folder}
        results = requests.post(url, data=payload, cookies=self.cookies)
        return simplejson.loads(results.content)

    def rename_folder(self, folder_to_rename,new_folder_name,in_folder):
        '''
        Rename a folder.
        '''
        
        url = nb_url + 'reader/rename_folder'
        payload = {'folder_to_rename':folder_to_rename,'new_folder_name':new_folder_name,'in_folder':in_folder}
        results = requests.post(url, data=payload, cookies=self.cookies)
        return simplejson.loads(results.content)

    def delete_folder(self, folder_to_delete,in_folder,feed_id):
        '''
        Delete a folder and unsubscribe from all feeds inside.
        '''
        
        url = nb_url + 'reader/delete_folder'
        payload = {'folder_to_delete':folder_to_delete,'in_folder':in_folder,'feed_id':feed_id}
        results = requests.post(url, data=payload, cookies=self.cookies)
        return simplejson.loads(results.content)


    def mark_feed_as_read(self, feed_id):
        '''
        Mark a list of feeds as read.
        '''
        
        url = nb_url + 'reader/mark_feed_as_read'
        payload = {'feed_id':feed_id}
        results = requests.post(url, data=payload, cookies=self.cookies)
        return simplejson.loads(results.content)


    def save_feed_order(self, folders):
        '''
        Reorder feeds and move them around between folders.
        The entire folder structure needs to be serialized.
        '''

        url = nb_url + 'reader/save_feed_order'
        payload = {'folders':folders}
        results = requests.post(url, data=payload, cookies=self.cookies)
        return simplejson.loads(results.content)


    def classifier(self, id_no):
        '''
        Get the intelligence classifiers for a user's site.
        Only includes the user's own classifiers. 
        Use /reader/feeds_trainer for popular classifiers.
        '''

        url = nb_url + 'classifier/%d' % id_no
        results = requests.get(url)
        return simplejson.loads(results.content)


    def classifier_save(self, like_type,dislike_type,remove_like_type,remove_dislike_type):
        '''
        Save intelligence classifiers (tags, titles, authors, and the feed) for a feed.
        '''
        
        url = nb_url + 'classifier/save'
        payload = {'like_[TYPE]':like_type,
                       'dislike_[TYPE]':dislike_type,
                        'remove_like_[TYPE]':remove_like_type,
                       'remove_dislike_[TYPE]':remove_dislike_type}
        results = requests.post(url, data=payload, cookies=self.cookies)
        return simplejson.loads(results.content)


    def opml_export(self, ):
        '''
        Download a backup of feeds and folders as an OPML file.
        Contains folders and feeds in XML; useful for importing in another RSS reader.
        '''
        
        url = nb_url + 'import/opml_export'
        results = requests.get(url)
        return simplejson.loads(results.content)



    def opml_upload(self, opml_file):
        '''
        Upload an OPML file.
        '''
        
        url = nb_url + 'import/opml_upload'
        f = open(opml_file)
        payload = {'file':f}
        f.close()
        results = requests.post(url, data=payload, cookies=self.cookies)
        return simplejson.loads(results.content)
