from fabric.api import *
from fabric.contrib.console import confirm
from fabric.operations import sudo

env.hosts = ["spiffytech@short.csc.ncsu.edu"]

def test():
    print local("git flow help")

def release(version):
    if local("git branch | grep '*'", capture=True) != "* release/%s" % version:
        print local("git flow release start %s" % version)
    else:
        print "Already on branch release/%s" % version

    if not confirm("A release has been started and staged locally. Does it behave like it should?"):
        abort("Aborting...")

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


def push():
    local("git push --tags github master develop")


def update_web2py(version):
    local("git flow feature start web2py_%s" % version)
    local("wget http://www.web2py.com/examples/static/web2py_src.zip")
    local("dtrx -n web2py_src.zip")
    local("mv web2py_src/web2py/* .")
    local("rmdir web2py_src/web2py")
    local("rmdir web2py_src")
    local("git status")
    print "\n\nGo restart the web2py server and make sure the site behaves properly"
    if confirm("Does the site still act OK?"):
        local("git flow feature finish web2py_%s" % version)
        local("rm -rf web2py_src.zip")  # Do this after we confirm the site works, so we don't have to redownload if something's broken
