from fabric.api import *
from fabric.contrib.console import confirm
from fabric.operations import sudo

env.hosts = ["spiffytech@short.csc.ncsu.edu"]

def test():
    print local("git flow help")

def release(version):
    workflow(version, hotfix=False)

def hotfix(version):
    workflow(version, hotfix=True)


def workflow(version, hotfix):
    release_or_hotfix = "hotfix" if hotfix is True else "release"

    if local("git branch | grep '*'", capture=True) != "* %s/%s" % (release_or_hotfix, version):
        print local("git flow %s start %s" % (release_or_hotfix, version))
    else:
        print "Already on branch %s/%s" % (release_or_hotfix, version)

    if not confirm("A %s has been started and staged locally. Does it behave like it should?" % (release_or_hotfix)):
        abort("Aborting...")

    print local("git flow %s finish %s" % (release_or_hotfix, version))

    push()
    _update_remote_docroot("apache/mobileblur.spiffyte.ch/docroot")


def _update_remote_docroot(docroot):
    with cd(docroot):
        print run("git pull")
        with settings(warn_only = True):
            result = run("httpd -t")
            if result.failed and not ("Apache has errors. Continue anyway?"):
                abort("Aborting...")

            result = sudo("service httpd restart")
            if result.failed and confirm ("Apache didn't start up again! Revert to last release?"):
                print run("git reset --hard HEAD^")


def stage():
    push()
    _update_remote_docroot("apache/mobileblur-staging.spiffyte.ch/docroot")


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
