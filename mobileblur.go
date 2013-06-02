package main

import (
    "errors"
    "fmt"
    "net/http"
    "html/template"
    "strconv"
    "strings"
    "time"

    "./newsblur"
    gocache "github.com/pmylund/go-cache"
    mux "github.com/gorilla/mux"
    //"github.com/hoisie/mustache"
    //"./mustache"
)

type MyCache struct {
    cache gocache.Cache
}

func (cache *MyCache) get(key string, f func() interface{}, duration time.Duration) (interface{}) {
    val, found := cache.cache.Get(key)
    if(found) {
        return val
    }

    val = f()
    cache.cache.Set(key, val, duration)
    return val
}

var cache = gocache.New(2*time.Minute, 30*time.Second)

func initNewsblur(w *http.ResponseWriter, r *http.Request) (newsblur.Newsblur, error) {
    // TODO: Retrieve cookie from user response here, instead of logging in to Newsblur with a test account
    var nb newsblur.Newsblur
        if cookie, err := r.Cookie("newsblur_sessionid"); err == nil {
        nb.Cookie = cookie.Value
    } else {
        http.Redirect(*w, r, "/login", http.StatusSeeOther)
        return nb, errors.New("You need to log in")
    }

    nb.GetProfile()
    nb.GetFeeds()
    nb.GetFolders()

    return nb, nil
}


func index(w http.ResponseWriter, r *http.Request) {
    nb, err := initNewsblur(&w, r)
    if err != nil {
        return
    }

    vals := map[string]interface{}{
        "feeds": nb.Profile.Feeds,
        "folder": nb.Profile.Folder,
        "socialFeeds": nb.Profile.SocialFeeds,
    }

    t := template.Must(template.New("index").ParseFiles("templates/wrapper.html", "templates/index"))
    err = t.Execute(w, vals)
    if err != nil {
        panic(err)
    }
}


func login(w http.ResponseWriter, r *http.Request) {
    var username string
    var password string
    var err error
    if r.Method == "POST" {
        r.ParseForm()
        username = r.Form.Get("username")
        password = r.Form.Get("password")
        var nb newsblur.Newsblur
        cookie, err := nb.Login(username, password)
        if err == nil {
            c := http.Cookie{
                Name: "newsblur_sessionid",
                Value: cookie,
                Path: "/",
                Domain: ".mbtest.spiffyte.ch",
                MaxAge: 315360000,
            }
            http.SetCookie(w, &c)
            http.Redirect(w, r, "/", http.StatusSeeOther)
            return
        }
        // TODO: Show form again, or redirect to /
    }

    vals := map[string]interface{}{
        "username": username,
        "password": password,
        "error": err,
    }

    t := template.Must(template.New("login.html").ParseFiles("templates/wrapper.html", "templates/login.html"))
    err = t.Execute(w, vals)
    if err != nil {
        panic(err)
    }
}


func logout(w http.ResponseWriter, r *http.Request) {
    c := http.Cookie{
        Name: "newsblur_sessionid",
        Value: "deleted",
        Path: "/",
        Domain: ".mbtest.spiffyte.ch",
        Expires: time.Unix(1000000, 0),
    }
    http.SetCookie(w, &c)
    http.Redirect(w, r, "/", http.StatusSeeOther)
}


func stories(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)

    nb, err := initNewsblur(&w, r)
    if err != nil {
        return
    }

    nb.GetFolders()
    nb.GetFeeds()

    feed_id, err := strconv.Atoi(vars["feed_id"])
    if err != nil {
        panic(err)
    }

    page, err := strconv.Atoi(r.URL.Query().Get("p"))
    if err != nil {
        page = 1
    }

    feed := nb.Feeds[feed_id]
    stories := feed.GetStoryPage(&nb, page, false).Stories
    if len(stories) == 0 {
        fmt.Fprintf(w, "false")
        return
    }

    vals := map[string]interface{}{
        "Stories": stories,
        "feed": feed,
        "page": page,  // use this instead of feed ID in template to collapse things
        "isSocial": false,
    }

    if r.Header.Get("X-Requested-With") == "XMLHttpRequest" {
        vals["notAJAX"] = false
    } else {
        vals["notAJAX"] = true
    }

    t := template.Must(template.New("stories").ParseFiles("templates/wrapper.html", "templates/stories"))
    err = t.Execute(w, vals)
    if err != nil {
        panic(err)
    }
}


func socialStories(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)

    nb, err := initNewsblur(&w, r)
    if err != nil {
        return
    }

    feed_id := vars["feed_id"]

    page, err := strconv.Atoi(r.URL.Query().Get("p"))
    if err != nil {
        page = 1
        fmt.Println("Page not set explicitly")
    }

    feed := nb.Profile.SocialFeeds[feed_id]
    stories := feed.GetSocialStoryPage(&nb, page, false).Stories
    if len(stories) == 0 {
        fmt.Fprintf(w, "false")
        return
    }

    vals := map[string]interface{}{
        "Stories": stories,
        "feed": feed,
        "page": page,  // TODO: use this instead of feed ID in template to collapse things
        "isSocial": true,
    }

    if r.Header.Get("X-Requested-With") == "XMLHttpRequest" {
        vals["notAJAX"] = false
    } else {
        vals["notAJAX"] = true
    }

    t := template.Must(template.New("stories").ParseFiles("templates/wrapper.html", "templates/stories"))
    err = t.Execute(w, vals)
    if err != nil {
        panic(err)
    }
}


func markStoryRead(w http.ResponseWriter, r *http.Request) {
    nb, err := initNewsblur(&w, r)
    if err != nil {
        return
    }

    feed_id := r.URL.Query().Get("feed_id")

    var isSocial bool
    if strings.Contains(feed_id, ":") {
        isSocial = true
    } else {
        isSocial = false
    }

    storyID := r.URL.Query().Get("story_id")
    if err != nil {
        panic(err)
    }

    // TODO: Make this support social story (string) IDs. Or better, find the equivalent social story function
    if isSocial == true {
        storyFeedID := r.URL.Query().Get("storyFeedID")
        feedID := strings.Split(feed_id, ":")[1]
        stories := make(map[string]map[string][]string)
        if stories[feedID] == nil {
            stories[feedID] = make(map[string][]string)
        }
        stories[feedID][storyFeedID] = append(stories[feedID][storyFeedID], storyID)
        err = nb.MarkSocialStoriesRead(stories)
    } else {
        feed_id_int, err := strconv.Atoi(feed_id)
        if err != nil {
            panic(err)
        }
        err = nb.MarkStoryRead(feed_id_int, storyID)
    }

    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        fmt.Fprintf(w, "true")
        return
    }

    fmt.Fprintf(w, "true")
}


func markReadBulk(w http.ResponseWriter, r *http.Request) {
    nb, err := initNewsblur(&w, r)
    if err != nil {
        return
    }

    r.ParseForm()
    if r.URL.Query().Get("isSocial") == "true" {
        stories := make(map[string]map[string][]string)
        for story, _ := range r.Form {
            fields := strings.SplitN(story, "-", 4)
            if fields[0] != "story" {
                continue
            }

            feedID := strings.Split(fields[1], ":")[1]
            storyFeedID := fields[2]
            storyID := fields[3]

            if stories[feedID] == nil {
                stories[feedID] = make(map[string][]string)
            }
            stories[feedID][storyFeedID] = append(stories[feedID][storyFeedID], storyID)
        }

        err = nb.MarkSocialStoriesRead(stories)
    } else {
        stories := make(map[string][]string)
        for story, _ := range r.Form {
            fields := strings.SplitN(story, "-", 3)
            if fields[0] != "story" {
                continue
            }

            stories[fields[1]] = append(stories[fields[1]], fields[2])
        }

        err = nb.MarkStoriesReadBulk(stories)
    }

    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        fmt.Fprintf(w, "true")
        return
    }

    fmt.Fprintf(w, "true")
}


func markUnread(w http.ResponseWriter, r *http.Request) {
    nb, err := initNewsblur(&w, r)
    if err != nil {
        // TODO: This should not panic
        return
    }

    r.ParseForm()
    rawStory := r.Form.Get("story")
    fields := strings.SplitN(rawStory, "-", 3)
    feedID := fields[1]
    storyID := fields[2]

    err = nb.MarkStoryUnread(feedID, storyID)

    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        fmt.Fprintf(w, "true")
        return
    }

    fmt.Fprintf(w, "true")
}


func main() {
    r := mux.NewRouter()
    r.HandleFunc("/", index)
    r.HandleFunc("/login", login)
    r.HandleFunc("/logout", logout)
    r.HandleFunc("/feeds", index)
    r.HandleFunc("/feeds/{feed_id}", stories)
    r.HandleFunc("/social/{feed_id}", socialStories)
    r.HandleFunc("/stories/mark_read", markStoryRead)
    r.HandleFunc("/stories/markReadBulk", markReadBulk)
    r.HandleFunc("/stories/markUnread", markUnread)

    fmt.Println("Listening for browser connections")
    http.Handle("/", r)

    http.ListenAndServe(":4001", nil)
}
