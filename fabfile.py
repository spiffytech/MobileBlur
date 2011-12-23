from fabric.api import *
from fabric.contrib.console import confirm
from fabric.operations import sudo

env.hosts = ["spiffytech@short.csc.ncsu.edu"]

def test():
    print local("git flow help")

def release(version):
    if local("git branch | grep '*'") != "* release/%s" % version:
        print local("git flow release start %s" % version)
        if not confirm("A release has been started and staged locally. Does it behave like it should?"):
            abort("Aborting...")
    else:
        print "Already on branch release/%s" % version

    print local("git flow release finish %s" % version)
    print local("git push github master")

    with cd("apache/mobileblur.spiffyte.ch/docroot"):
        print run("git pull")
        with settings(warn_only = True):
            result = run("httpd -t")
            if result.failed and not ("Apache has errors. Continue anyway?"):
                abort("Aborting...")

            result = sudo("service httpd restart")
            if result.failed and confirm ("Apache didn't start up again! Revert to last release?"):
                print run("git reset --hard HEAD^")
