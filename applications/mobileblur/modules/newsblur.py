#!/usr/bin/python

"""newsblur.py - An API wrapper library for newsblur.com"""

import cookielib
import simplejson
import urllib
import urllib2

__author__ = 'Dananjaya Ramanayake <dananjaya86@gmail.com>, spiffytech <spiffytechgmail.com>'
__version__ = "0.1"


# Set up cookie handling so we can auth with the Newsblur API
cj = cookielib.LWPCookieJar()
opener = urllib2.build_opener(urllib2.HTTPCookieProcessor(cj))
urllib2.install_opener(opener)

nb_url = "http://newsblur.com/"


def login(username,password):
    '''
    
    Login as an existing user.
    If a user has no password set, you cannot just send any old password. 
    Required parameters, username and password, must be of string type.
    
    '''
    url = nb_url + 'api/login'
    params = urllib.urlencode({'username':username,'password':password})
    results = urllib2.urlopen(url,params).read()
    return simplejson.loads(results)

def logout():
    '''
    
    Logout the currently logged in user.
    
    '''
    url = nb_url + 'api/logout'
    results = urllib2.urlopen(url).read()
    return simplejson.loads(results)

def signup(username,password,email):
    '''
    
    Create a new user.
    All three required parameters must be of type string.
    
    '''
    url = nb_url + 'api/signup'
    params = urllib.urlencode({'signup_username':username,'signup_password':password,'signup_email':email})
    results = urllib2.urlopen(url,params).read()
    return simplejson.loads(results)

def search_feed(address,offset=1):
    '''
    
    Retrieve information about a feed from its website or RSS address.
    Parameter address must be of type string while parameter offset must be an integer.
    Will return a feed.
    
    '''
    url = nb_url + 'rss_feeds/search_feed?%s'
    params = urllib.urlencode({'address':address,'offset':offset})
    results = urllib2.urlopen(url % params).read()
    return simplejson.loads(results)

def feeds(include_favicons=True,flat=False):
    '''
    
    Retrieve a list of feeds to which a user is actively subscribed.
        Includes the 3 unread counts (positive, neutral, negative), as well as optional favicons.

        '''
    
    url = nb_url + 'reader/feeds'
    params = urllib.urlencode({'include_favicons':include_favicons,'flat':flat})
#    print url + " " + url % params
    results = urllib2.urlopen(url).read()
    return simplejson.loads(results)


def favicons(feeds=[1,2,3]):
    '''
    
    Retrieve a list of favicons for a list of feeds. 
    Used when combined with /reader/feeds and include_favicons=false, so the feeds request contains far less data. 
    Useful for mobile devices, but requires a second request. 
    
    '''
    url = nb_url + 'reader/favicons?%s'
    params = urllib.urlencode({'feeds':feeds})
    results = urllib2.urlopen(url % params).read()
    return simplejson.loads(results)
    
def id(id_no):
    '''
    
    Retrieve the original page from a single feed.
    
    '''
    url = nb_url + 'reader/page/%d' % id_no
    results = urllib2.urlopen(url).read()
    return simplejson.loads(results)

def refresh_feeds():
    '''
    
    Up-to-the-second unread counts for each active feed.
        Poll for these counts no more than once a minute.
        
        '''

    url = nb_url + 'reader/refresh_feeds'
    results = urllib2.urlopen(url).read()
    return simplejson.loads(results)

def feeds_trainer(feed_id):
    '''
    
     Retrieves all popular and known intelligence classifiers.
        Also includes user's own classifiers.
        
        '''

    url = nb_url + 'reader/feeds_trainer?%s'
    params = urllib.urlencode({'feed_id':feed_id})
    results = urllib2.urlopen(url % params).read()
    return simplejson.loads(results)

def statistics(id_no):
    '''
    
    If you only want a user's classifiers, use /classifiers/:id.
        Omit the feed_id to get all classifiers for all subscriptions.
        
        '''

    url = nb_url + 'rss_feeds/statistics/%d' % id_no
    results = urllib2.urlopen(url).read()
    return simplejson.loads(results)

def feed_autocomplete(term):
    '''
    
    Get a list of feeds that contain a search phrase.
        Searches by feed address, feed url, and feed title, in that order.
        Will only show sites with 2+ subscribers.
        
        '''
    url = nb_url + 'rss_feeds/feed_autocomplete?%'
    params = urllib.urlencode({'term':term})
    results = urllib2.urlopen(url % params).read()
    return simplejson.loads(results)

def feed(id=1):
    '''
    
    Retrieve stories from a single feed.
    
    '''
    url = nb_url + 'reader/feed/%d' % id
    results = urllib2.urlopen(url).read()
    return simplejson.loads(results)

def starred_stories(page=1):
    '''
    
    Retrieve a user's starred stories.
    
    '''
    url = nb_url + 'reader/starred_stories?%s'
    params = urllib.urlencode({'page':page})
    results = urllib2.urlopen(url % params).read()
    return simplejson.loads(results)

def river_stories(feeds,page=1,read_stories_count=0):
    '''
    
    Retrieve stories from a collection of feeds. This is known as the River of News.
        Stories are ordered in reverse chronological order.
        
        '''

    url = nb_url + 'reader/river_stories?%s'
    params = urllib.urlencode({'feeds':feeds,'page':page,'read_stories_count':read_stories_count})
    results = urllib2.urlopen(url % params).read()
    return simplejson.loads(results)

def mark_story_as_read(story_id,feed_id):
    '''
    
     Mark stories as read.
        Multiple story ids can be sent at once.
        Each story must be from the same feed.
        
        '''

    url = nb_url + 'reader/mark_story_as_read'
    params = urllib.urlencode({'story_id':story_id,'feed_id':feed_id})
    results = urllib2.urlopen(url,params).read()
    return simplejson.loads(results)

def mark_story_as_starred(story_id,feed_id):
    '''
    
    Mark a story as starred (saved).
    
    '''
    url = nb_url + 'reader/mark_story_as_starred'
    params = urllib.urlencode({'story_id':story_id,'feed_id':feed_id})
    results = urllib2.urlopen(url,params).read()
    return simplejson.loads(results)

def mark_all_as_read(days=0):
    '''
    
    Mark all stories in a feed or list of feeds as read.
    
    '''
    url = nb_url + 'reader/mark_all_as_read'
    params = urllib.urlencode({'days':days})
    results = urllib2.urlopen(url,params).read()
    return simplejson.loads(results)

def add_url(url,folder='[Top Level]'):
    '''
    
    Add a feed by its URL. 
    Can be either the RSS feed or the website itself.
    
    '''
    url = nb_url + 'reader/add_url'
    params = urllib.urlencode({'url':url,'folder':folder})
    results = urllib2.urlopen(url,params).read()
    return simplejson.loads(results)


def add_folder(folder,parent_folder='[Top Level]'):
    '''
    
    Add a new folder.
    
    '''
    
    url = nb_url + 'reader/add_folder'
    params = urllib.urlencode({'folder':folder,'parent_folder':parent_folder})
    results = urllib2.urlopen(url,params).read()
    return simplejson.loads(results)

def rename_feed(feed_title,feed_id):
    '''
    
    Rename a feed title. Only the current user will see the new title.
    
    '''
    url = nb_url + 'reader/rename_feed'
    params = urllib.urlencode({'feed_title':feed_title,'feed_id':feed_id})
    results = urllib2.urlopen(url,params).read()
    return simplejson.loads(results)

def delete_feed(feed_id,in_folder):
    '''
    
    Unsubscribe from a feed. Removes it from the folder.
        Set the in_folder parameter to remove a feed from the correct folder, in case the user is subscribed to the feed in multiple folders.

        '''    
    url = nb_url + 'reader/delete_feed'
    params = urllib.urlencode({'feed_id':feed_id,'in_folder':in_folder})
    results = urllib2.urlopen(url,params).read()
    return simplejson.loads(results)

def rename_folder(folder_to_rename,new_folder_name,in_folder):
    '''
    
    Rename a folder.
    
    '''
    url = nb_url + 'reader/rename_folder'
    params = urllib.urlencode({'folder_to_rename':folder_to_rename,'new_folder_name':new_folder_name,'in_folder':in_folder})
    results = urllib2.urlopen(url,params).read()
    return simplejson.loads(results)

def delete_folder(folder_to_delete,in_folder,feed_id):
    '''
    
    Delete a folder and unsubscribe from all feeds inside.
    
    '''
    url = nb_url + 'reader/delete_folder'
    params = urllib.urlencode({'folder_to_delete':folder_to_delete,'in_folder':in_folder,'feed_id':feed_id})
    results = urllib2.urlopen(url,params).read()
    return simplejson.loads(results)


def mark_feed_as_read(feed_id):
    '''
    
    Mark a list of feeds as read.
    
    '''
    url = nb_url + 'reader/mark_feed_as_read'
    params = urllib.urlencode({'feed_id':feed_id})
    results = urllib2.urlopen(url,params).read()
    return simplejson.loads(results)


def save_feed_order(folders):
    '''
    
    Reorder feeds and move them around between folders.
        The entire folder structure needs to be serialized.
        
        '''

    url = nb_url + 'reader/save_feed_order'
    params = urllib.urlencode({'folders':folders})
    results = urllib2.urlopen(url,params).read()
    return simplejson.loads(results)


def classifier(id_no):
    '''
    
        Get the intelligence classifiers for a user's site.
        Only includes the user's own classifiers. 
        Use /reader/feeds_trainer for popular classifiers.
        
        '''

    url = nb_url + 'classifier/%d' % id_no
    results = urllib2.urlopen(url).read()
    return simplejson.loads(results)


def classifier_save(like_type,dislike_type,remove_like_type,remove_dislike_type):
    '''
    
    Save intelligence classifiers (tags, titles, authors, and the feed) for a feed.
    
        '''
    url = nb_url + 'classifier/save'
    params = urllib.urlencode({'like_[TYPE]':like_type,
                   'dislike_[TYPE]':dislike_type,
                    'remove_like_[TYPE]':remove_like_type,
                   'remove_dislike_[TYPE]':remove_dislike_type})
    results = urllib2.urlopen(url,params).read()
    return simplejson.loads(results)


def opml_export():
    '''
    
    Download a backup of feeds and folders as an OPML file.
        Contains folders and feeds in XML; useful for importing in another RSS reader.
        
        '''
    url = nb_url + 'import/opml_export'
    results = urllib2.urlopen(url).read()
    return simplejson.loads(results)



def opml_upload(opml_file):
    '''
    
    Upload an OPML file.
    
    '''
    url = nb_url + 'import/opml_upload'
    f = open(opml_file)
    params = urllib.urlencode({'file':f})
    f.close()
    results = urllib2.urlopen(url,params).read()
    return simplejson.loads(results)
