#!/usr/bin/env python
# -*- coding: utf-8 -*-

"""
This file is part of the web2py Web Framework
Copyrighted by Massimo Di Pierro <mdipierro@cs.depaul.edu>
License: LGPLv3 (http://www.gnu.org/licenses/lgpl.html)

This file is based, althought a rewrite, on MIT code from the Bottle web framework.
"""

import os, sys, optparse
path = os.path.dirname(os.path.abspath(__file__))
os.chdir(path)
sys.path = [path]+[p for p in sys.path if not p==path]
import gluon.main
from gluon.fileutils import read_file, write_file

class Servers:
    @staticmethod
    def cgi(app, address=None, **options):
        from wsgiref.handlers import CGIHandler
        CGIHandler().run(app) # Just ignore host and port here

    @staticmethod
    def flup(app,address, **options):
        import flup.server.fcgi
        flup.server.fcgi.WSGIServer(app, bindAddress=address).run()

    @staticmethod
    def wsgiref(app,address,**options): # pragma: no cover
        from wsgiref.simple_server import make_server, WSGIRequestHandler
        class QuietHandler(WSGIRequestHandler):
            def log_request(*args, **kw): pass
        options['handler_class'] = QuietHandler
        srv = make_server(address[0],address[1],app,**options)
        srv.serve_forever()

    @staticmethod
    def cherrypy(app,address, **options):
        from cherrypy import wsgiserver
        server = wsgiserver.CherryPyWSGIServer(address, app)
        server.start()

    @staticmethod
    def rocket(app,address, **options):
        from gluon.rocket import CherryPyWSGIServer
        server = CherryPyWSGIServer(address, app)
        server.start()

    @staticmethod
    def rocket_with_repoze_profiler(app,address, **options):
        from gluon.rocket import CherryPyWSGIServer
        from repoze.profile.profiler import AccumulatingProfileMiddleware
        from gluon.settings import global_settings
        global_settings.web2py_crontype = 'none'
        wrapped = AccumulatingProfileMiddleware(
            app,
            log_filename='wsgi.prof',
            discard_first_request=True,
            flush_at_shutdown=True,
            path = '/__profile__'
            )
        server = CherryPyWSGIServer(address, wrapped)
        server.start()

    @staticmethod
    def paste(app,address,**options):
        from paste import httpserver
        from paste.translogger import TransLogger
        httpserver.serve(app, host=address[0], port=address[1], **options)

    @staticmethod
    def fapws(app,address, **options):
        import fapws._evwsgi as evwsgi
        from fapws import base
        evwsgi.start(address[0],str(address[1]))
        evwsgi.set_base_module(base)
        def app(environ, start_response):
            environ['wsgi.multiprocess'] = False
            return app(environ, start_response)
        evwsgi.wsgi_cb(('',app))
        evwsgi.run()


    @staticmethod
    def gevent(app,address, **options):
        from gevent import monkey; monkey.patch_all()
        from gevent import pywsgi
        from gevent.pool import Pool
        pywsgi.WSGIServer(address, app, spawn = 'workers' in options and Pool(int(option.workers)) or 'default').serve_forever()

    @staticmethod
    def bjoern(app,address, **options):
        import bjoern
        bjoern.run(app, *address)

    @staticmethod
    def tornado(app,address, **options):
        import tornado.wsgi
        import tornado.httpserver
        import tornado.ioloop
        container = tornado.wsgi.WSGIContainer(app)
        server = tornado.httpserver.HTTPServer(container)
        server.listen(address=address[0], port=address[1])
        tornado.ioloop.IOLoop.instance().start()

    @staticmethod
    def twisted(app,address, **options):
        from twisted.web import server, wsgi
        from twisted.python.threadpool import ThreadPool
        from twisted.internet import reactor
        thread_pool = ThreadPool()
        thread_pool.start()
        reactor.addSystemEventTrigger('after', 'shutdown', thread_pool.stop)
        factory = server.Site(wsgi.WSGIResource(reactor, thread_pool, app))
        reactor.listenTCP(address[1], factory, interface=address[0])
        reactor.run()

    @staticmethod
    def diesel(app,address, **options):
        from diesel.protocols.wsgi import WSGIApplication
        app = WSGIApplication(app, port=address[1])
        app.run()

    @staticmethod
    def gnuicorn(app,address, **options):
        import gunicorn.arbiter
        gunicorn.arbiter.Arbiter(address, 4, app).run()

    @staticmethod
    def eventlet(app,address, **options):
        from eventlet import wsgi, listen
        wsgi.server(listen(address), app)


def run(servername,ip,port,softcron=True,logging=False,profiler=None):
    if logging:
        application = gluon.main.appfactory(wsgiapp=gluon.main.wsgibase,
                                            logfilename='httpserver.log',
                                            profilerfilename=profiler)
    else:
        application = gluon.main.wsgibase
    if softcron:
        from gluon.settings import global_settings
        global_settings.web2py_crontype = 'soft'
    getattr(Servers,servername)(application,(ip,int(port)))

def main():
    usage = "python anyserver.py -s tornado -i 127.0.0.1 -p 8000 -l -P"
    try:
        version = read_file('VERSION')
    except IOError:
        version = ''
    parser = optparse.OptionParser(usage, None, optparse.Option, version)
    parser.add_option('-l',
                      '--logging',
                      action='store_true',
                      default=False,
                      dest='logging',
                      help='log into httpserver.log')
    parser.add_option('-P',
                      '--profiler',                      
                      default=False,
                      dest='profiler',
                      help='profiler filename')
    servers = ', '.join(x for x in dir(Servers) if not x[0]=='_')
    parser.add_option('-s',
                      '--server',
                      default='rocket',
                      dest='server',
                      help='server name (%s)' % servers)
    parser.add_option('-i',
                      '--ip',
                      default='127.0.0.1',
                      dest='ip',
                      help='ip address')
    parser.add_option('-p',
                      '--port',
                      default='8000',
                      dest='port',
                      help='port number')
    parser.add_option('-w',
                      '--workers',
                      default='',
                      dest='workers',
                      help='number of workers number')
    (options, args) = parser.parse_args()
    print 'starting %s on %s:%s...' % (options.server,options.ip,options.port)
    run(options.server,options.ip,options.port,logging=options.logging,profiler=options.profiler)

if __name__=='__main__':
    main()
