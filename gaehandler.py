#!/usr/bin/env python
# -*- coding: utf-8 -*-

"""
This file is part of the web2py Web Framework
Copyrighted by Massimo Di Pierro <mdipierro@cs.depaul.edu>
License: LGPLv3 (http://www.gnu.org/licenses/lgpl.html)
"""

##############################################################################
# Configuration parameters for Google App Engine
##############################################################################
KEEP_CACHED = False    # request a dummy url every 10secs to force caching app
LOG_STATS = False      # web2py level log statistics
APPSTATS = True         # GAE level usage statistics and profiling
DEBUG = False          # debug mode
AUTO_RETRY = True      # force gae to retry commit on failure
#
# Read more about APPSTATS here
#   http://googleappengine.blogspot.com/2010/03/easy-performance-profiling-with.html
# can be accessed from:
#   http://localhost:8080/_ah/stats
##############################################################################
# All tricks in this file developed by Robin Bhattacharyya
##############################################################################


import time
import os
import sys
import logging
import cPickle
import pickle
import wsgiref.handlers
import datetime

path = os.path.dirname(os.path.abspath(__file__))
sys.path = [path]+[p for p in sys.path if not p==path]

sys.modules['cPickle'] = sys.modules['pickle']


from gluon.settings import global_settings
from google.appengine.api.labs import taskqueue
from google.appengine.ext import webapp
from google.appengine.ext.webapp.util import run_wsgi_app


global_settings.web2py_runtime_gae = True
global_settings.db_sessions = True
if os.environ.get('SERVER_SOFTWARE', '').startswith('Devel'):
    (global_settings.web2py_runtime, DEBUG) = \
        ('gae:development', True)
else:
    (global_settings.web2py_runtime, DEBUG) = \
        ('gae:production', False)


import gluon.main


def log_stats(fun):
    """Function that will act as a decorator to make logging"""
    def newfun(env, res):
        """Log the execution time of the passed function"""
        timer = lambda t: (t.time(), t.clock())
        (t0, c0) = timer(time)
        executed_function = fun(env, res)
        (t1, c1) = timer(time)
        log_info = """**** Request: %.2fms/%.2fms (real time/cpu time)"""
        log_info = log_info % ((t1 - t0) * 1000, (c1 - c0) * 1000)
        logging.info(log_info)
        return executed_function
    return newfun


logging.basicConfig(level=logging.INFO)


def wsgiapp(env, res):
    """Return the wsgiapp"""
    if env['PATH_INFO'] == '/_ah/queue/default':
        if KEEP_CACHED:
            delta = datetime.timedelta(seconds=10)
            taskqueue.add(eta=datetime.datetime.now() + delta)
        res('200 OK',[('Content-Type','text/plain')])
        return ['']
    env['PATH_INFO'] = env['PATH_INFO'].encode('utf8')

    #this deals with a problem where GAE development server seems to forget
    # the path between requests
    if global_settings.web2py_runtime == 'gae:development':
        gluon.admin.create_missing_folders()

    return gluon.main.wsgibase(env, res)


if LOG_STATS or DEBUG:
    wsgiapp = log_stats(wsgiapp)


if AUTO_RETRY:
    from gluon.contrib.gae_retry import autoretry_datastore_timeouts
    autoretry_datastore_timeouts()


def main():
    """Run the wsgi app"""
    if APPSTATS:
        run_wsgi_app(wsgiapp)
    else:
        wsgiref.handlers.CGIHandler().run(wsgiapp)

if __name__ == '__main__':
    main()

