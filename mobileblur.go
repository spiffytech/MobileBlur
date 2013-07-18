package main

import (
    "encoding/json"
    "fmt"
    "html/template"
    "net/http"
    "strconv"
    "strings"
    "time"

    "./newsblur"
    mux "github.com/gorilla/mux"
    gocache "github.com/pmylund/go-cache"
)

type AuthMux struct {
    *mux.Router
}

type MyCache struct {
    cache gocache.Cache
}

func (cache *MyCache) get(key string, f func() interface{}, duration time.Duration) interface{} {
    val, found := cache.cache.Get(key)
    if found {
        return val
    }

    val = f()
    cache.cache.Set(key, val, duration)
    return val
}

var cache = gocache.New(2*time.Minute, 30*time.Second)

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
                Name:   "newsblur_sessionid",
                Value:  cookie,
                Path:   "/",
                Domain: ".mbtest.spiffyte.ch",
                MaxAge: 315360000,
            }
            http.SetCookie(w, &c)
            http.Redirect(w, r, "/", http.StatusSeeOther)

            if _, err = r.Cookie("intelligence_threshold"); err != nil {
                c := http.Cookie{
                    Name:   "intelligence_threshold",
                    Value:  strconv.Itoa(-1),
                    Path:   "/",
                    Domain: ".mbtest.spiffyte.ch",
                    MaxAge: 315360000,
                }
                http.SetCookie(w, &c)
            }

            return
        }
        // TODO: Show form again, or redirect to /
    }

    vals := map[string]interface{}{
        "username": username,
        "password": password,
        "error":    err,
    }

    t := template.Must(template.New("login.html").ParseFiles("templates/wrapper.html", "templates/login.html"))
    err = t.Execute(w, vals)
    if err != nil {
        panic(err)
    }
}

func logout(w http.ResponseWriter, r *http.Request) {
    c := http.Cookie{
        Name:    "newsblur_sessionid",
        Value:   "deleted",
        Path:    "/",
        Domain:  ".mbtest.spiffyte.ch",
        Expires: time.Unix(1000000, 0),
    }
    http.SetCookie(w, &c)
    http.Redirect(w, r, "/", http.StatusSeeOther)
}

func index(w http.ResponseWriter, r *http.Request, nb *newsblur.Newsblur) {
    vals := map[string]interface{}{
        "nb":          nb,
        "feeds":       nb.Profile.Feeds,
        "folder":      nb.Profile.Folder,
        "socialFeeds": nb.Profile.SocialFeeds,
    }

    t := template.Must(template.New("index").Funcs(template.FuncMap{"showFeed": showFeed(nb)}).ParseFiles("templates/wrapper.html", "templates/index"))
    err := t.Execute(w, vals)
    if err != nil {
        panic(err)
    }
}

func stories(w http.ResponseWriter, r *http.Request, nb *newsblur.Newsblur) {
    vars := mux.Vars(r)

    feedID := vars["feedID"]

    page, err := strconv.Atoi(r.URL.Query().Get("p"))
    if err != nil {
        page = 1
    }

    feed := nb.Profile.Feeds[feedID]
    stories := feed.GetStoryPage(nb, page, false)
    fmt.Println("Stories")
    fmt.Println(stories)
    if len(stories) == 0 {
        fmt.Fprintf(w, "false")
        return
    }

    vals := map[string]interface{}{
        "nb":       nb,
        "Stories":  stories,
        "feed":     feed,
        "page":     page, // use this instead of feed ID in template to collapse things
        "isSocial": false,
    }

    if r.Header.Get("X-Requested-With") == "XMLHttpRequest" {
        vals["notAJAX"] = false
    } else {
        vals["notAJAX"] = true
    }

    t := template.Must(template.New("stories").Funcs(template.FuncMap{"showStory": showStory(nb)}).ParseFiles("templates/wrapper.html", "templates/stories"))
    err = t.Execute(w, vals)
    if err != nil {
        panic(err)
    }
}

func getStoryContent(w http.ResponseWriter, r *http.Request, nb *newsblur.Newsblur) {
    vars := mux.Vars(r)

    feedID := r.URL.Query().Get("feedID")
    storyID, err := strconv.Atoi(vars["storyID"])
    isSocial, err := strconv.ParseBool(r.URL.Query().Get("isSocial"))
    if err != nil {
        panic(err)
    }

    page, err := strconv.Atoi(r.URL.Query().Get("page"))
    if err != nil {
        page = 1
    }

    var stories []newsblur.StoryInt
    if isSocial {
        feed, ok := nb.Profile.SocialFeeds[feedID]
        if !ok {
            fmt.Fprintf(w, "Feed not found: false: ", feedID)
            return
        }
        stories = feed.GetSocialStoryPage(nb, page, false)
    } else {
        feed, ok := nb.Profile.Feeds[feedID]
        if !ok {
            fmt.Fprintf(w, "Feed not found: false: ", feedID)
            return
        }
        stories = feed.GetStoryPage(nb, page, false)
    }
    if len(stories) == 0 {
        fmt.Fprintf(w, "No stories: false")
        return
    }

    var story newsblur.StoryInt
    found := false
    // TODO: Fix this for social story IDs
    for id, s := range stories {
        if id == storyID {
            story = s
            found = true
            break
        }
    }

    if found == false {
        fmt.Fprintf(w, "Story id invalid: false")
        return
    }

    ret := map[string]string{
        "content": string(story.Content()),
    }

    retj, err := json.Marshal(ret)
    if found == false {
        fmt.Fprintf(w, "false")
        return
    }

    fmt.Fprintf(w, string(retj))
}

func socialStories(w http.ResponseWriter, r *http.Request, nb *newsblur.Newsblur) {
    vars := mux.Vars(r)

    feedID := vars["feedID"]

    page, err := strconv.Atoi(r.URL.Query().Get("p"))
    if err != nil {
        page = 1
        fmt.Println("Page not set explicitly")
    }

    feed := nb.Profile.SocialFeeds[feedID]
    stories := feed.GetSocialStoryPage(nb, page, false)
    if len(stories) == 0 {
        fmt.Fprintf(w, "false")
        return
    }

    vals := map[string]interface{}{
        "nb":       nb,
        "Stories":  stories,
        "feed":     feed,
        "page":     page, // TODO: use this instead of feed ID in template to collapse things
        "isSocial": true,
    }

    if r.Header.Get("X-Requested-With") == "XMLHttpRequest" {
        vals["notAJAX"] = false
    } else {
        vals["notAJAX"] = true
    }

    t := template.Must(template.New("stories").Funcs(template.FuncMap{"showStory": showStory(nb)}).ParseFiles("templates/wrapper.html", "templates/stories"))
    err = t.Execute(w, vals)
    if err != nil {
        panic(err)
    }
}

func markStoryRead(w http.ResponseWriter, r *http.Request, nb *newsblur.Newsblur) {
    feedID := r.URL.Query().Get("feed_id")

    var isSocial bool
    if strings.Contains(feedID, ":") {
        isSocial = true
    } else {
        isSocial = false
    }

    storyID := r.URL.Query().Get("story_id")

    // TODO: Make this support social story (string) IDs. Or better, find the equivalent social story function
    var err error
    if isSocial == true {
        storyFeedID := r.URL.Query().Get("storyFeedID")
        feedID := strings.Split(feedID, ":")[1]
        stories := make(map[string]map[string][]string)
        if stories[feedID] == nil {
            stories[feedID] = make(map[string][]string)
        }
        stories[feedID][storyFeedID] = append(stories[feedID][storyFeedID], storyID)
        err = nb.MarkSocialStoriesRead(stories)
    } else {
        fmt.Println(feedID)
        feed_id_int, err := strconv.Atoi(feedID)
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

func markReadBulk(w http.ResponseWriter, r *http.Request, nb *newsblur.Newsblur) {
    r.ParseForm()

    var err error
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

func markUnread(w http.ResponseWriter, r *http.Request, nb *newsblur.Newsblur) {
    r.ParseForm()
    rawStory := r.Form.Get("story")
    fields := strings.SplitN(rawStory, "-", 3)
    feedID := fields[1]
    storyID := fields[2]

    err := nb.MarkStoryUnread(feedID, storyID)

    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        fmt.Fprintf(w, "true")
        return
    }

    fmt.Fprintf(w, "true")
}

func settings(w http.ResponseWriter, r *http.Request, nb *newsblur.Newsblur) {
    if r.Method == "POST" {
        r.ParseForm()
        if threshold, err := strconv.Atoi(r.Form.Get("threshold")); err == nil {
            setCookie(&w, *r, "threshold", strconv.Itoa(threshold), true)
        }

        if r.Form.Get("showRead") != "" {
            setCookie(&w, *r, "showRead", "true", true)
        } else {
            setCookie(&w, *r, "showRead", "false", true)
        }

        if r.Form.Get("emptyFeeds") != "" {
            setCookie(&w, *r, "emptyFeeds", "true", true)
        } else {
            setCookie(&w, *r, "emptyFeeds", "false", true)
        }

        http.Redirect(w, r, "/settings", http.StatusSeeOther)
    }

    vals := map[string]interface{}{
        "nb":         nb,
        "threshold":  nb.Threshold,
        "showRead":   nb.ShowRead,
        "emptyFeeds": nb.EmptyFeeds,
        "isSocial":   true,
    }

    if r.Header.Get("X-Requested-With") == "XMLHttpRequest" {
        vals["notAJAX"] = false
    } else {
        vals["notAJAX"] = true
    }

    t := template.Must(template.New("settings.html").ParseFiles("templates/wrapper.html", "templates/settings.html"))
    err := t.Execute(w, vals)
    if err != nil {
        panic(err)
    }
}

func PassesThreshold(story newsblur.Story, nb newsblur.Newsblur) bool {
    return story.Score() > nb.Threshold
}

func showFeed(nb *newsblur.Newsblur) func(newsblur.FeedInt) bool {
    return func(feed newsblur.FeedInt) bool {
        numAboveThreshold := 0
        numAboveThreshold += feed.GetPS()
        if nb.Threshold < 1 {
            numAboveThreshold += feed.GetNT()
        }
        if nb.Threshold < 0 {
            numAboveThreshold += feed.GetNG()
        }

        fmt.Println(nb.EmptyFeeds)
        return numAboveThreshold > 0 || nb.EmptyFeeds == true
    }
}

func showStory(nb *newsblur.Newsblur) func(newsblur.StoryInt) bool {
    return func(story newsblur.StoryInt) bool {
        fmt.Println(story)
        fmt.Println(story.ReadStatus())
        return story.Score() >= nb.Threshold && (story.ReadStatus() == 0 || nb.ShowRead == true)
    }
}

func setCookie(w *http.ResponseWriter, r http.Request, name string, defaultValue string, stomp bool) (value string) {
    c, err := r.Cookie(name)
    if err != nil || stomp == true {
        c := http.Cookie{
            Name:   name,
            Value:  defaultValue,
            Path:   "/",
            Domain: ".mbtest.spiffyte.ch",
            MaxAge: 315360000,
        }
        http.SetCookie(*w, &c)
        value = defaultValue
    } else {
        value = c.Value
    }

    return
}

func (mux *AuthMux) HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request, *newsblur.Newsblur)) {
    mux.Router.HandleFunc(pattern, authenticate(handler))
}

func (mux *AuthMux) HandleFuncNoAuth(pattern string, handler func(http.ResponseWriter, *http.Request)) {
    mux.Router.HandleFunc(pattern, handler)
}

func authenticate(handler func(http.ResponseWriter, *http.Request, *newsblur.Newsblur)) func(http.ResponseWriter, *http.Request) {
    return func(w http.ResponseWriter, r *http.Request) {
        var nb newsblur.Newsblur
        if cookie, err := r.Cookie("newsblur_sessionid"); err == nil {
            nb.Cookie = cookie.Value
        } else {
            http.Redirect(w, r, "/login", http.StatusSeeOther)
            fmt.Fprintf(w, "You need to log in")
            return
        }

        threshold := setCookie(&w, *r, "threshold", "-1", false)
        if threshold, err := strconv.Atoi(threshold); err != nil {
            setCookie(&w, *r, "threshold", "-1", true)
            nb.Threshold = -1
        } else {
            nb.Threshold = threshold
        }

        showRead := setCookie(&w, *r, "showRead", "true", false)
        if showRead, err := strconv.ParseBool(showRead); err != nil {
            setCookie(&w, *r, "showRead", "true", true)
            nb.ShowRead = true
        } else {
            nb.ShowRead = showRead
        }

        emptyFeeds := setCookie(&w, *r, "emptyFeeds", "-1", false)
        if emptyFeeds, err := strconv.ParseBool(emptyFeeds); err != nil {
            setCookie(&w, *r, "emptyFeeds", "true", true)
            nb.EmptyFeeds = true
        } else {
            nb.EmptyFeeds = emptyFeeds
        }

        fmt.Println(nb.Threshold, nb.ShowRead, nb.EmptyFeeds)

        nb.GetProfile()

        handler(w, r, &nb)
    }
}

func main() {
    r := AuthMux{mux.NewRouter()}
    r.HandleFunc("/", index)
    r.HandleFuncNoAuth("/login", login)
    r.HandleFuncNoAuth("/logout", logout)
    r.HandleFunc("/feeds", index)
    r.HandleFunc("/feeds/{feedID}", stories)
    r.HandleFunc("/social/{feedID}", socialStories)
    r.HandleFunc("/stories/mark_read", markStoryRead)
    r.HandleFunc("/stories/getContent", getStoryContent)
    r.HandleFunc("/stories/markReadBulk", markReadBulk)
    r.HandleFunc("/stories/markUnread", markUnread)
    r.HandleFunc("/settings", settings)

    http.Handle("/", r)

    fmt.Println("Listening for browser connections")
    http.ListenAndServe(":4001", nil)
}
