package main

import (
    "fmt"
    "net/http"
    "html/template"
    "strconv"
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

func initNewsblur() (newsblur.Newsblur, error) {
    // TODO: Retrieve cookie from user response here, instead of logging in to Newsblur with a test account
    var nb newsblur.Newsblur
    err := nb.Login("mbtest1", "mbtest1");
    if err != nil {
        return nb, err
    }

    nb.GetProfile()
    nb.GetFeeds()
    nb.GetFolders()

    return nb, nil
}


func index (w http.ResponseWriter, r *http.Request) {
    /*
    test1 := Test1{Val1: "Val1str", Val2: 2}
    test2 := Test2{Val3: "Val3str", Val4: 4}
    d := make(map[string]interface{})
    d["test1"] = &test1
    //d["test2"] = &test2
    _ = &test2
    rendered, err := mustache.MustRenderFile("templates/index.mustache", d)
    if err != nil {
        panic(err)
    }
    //rendered := mustache.RenderFile("templates/index.mustache", d)
    fmt.Fprintf(w, rendered)
    */

    nb, err := initNewsblur()
    if err != nil {
        panic(err)
    }

    vals := map[string]interface{}{
        "feeds": nb.Profile.Feeds,
        "folder": nb.Profile.Folder,
        "socialFeeds": nb.Profile.Social,
    }
    fmt.Println(vals)

    t := template.Must(template.New("index").ParseFiles("templates/wrapper.html", "templates/index"))
    err = t.Execute(w, vals)
    if err != nil {
        fmt.Println(err)
    }
}


func stories (w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)

    nb, err := initNewsblur()
    if err != nil {
        // TODO: This should not panic
        panic(err)
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
        fmt.Println("Page not set explicitly")
    }

    feed := nb.Feeds[feed_id]
    stories := feed.GetStoryPage(&nb, page, false).Stories

    vals := map[string]interface{}{
        "Stories": stories,
        "feed_id": feed_id,
        "page": page,  // use this instead of feed ID in template to collapse things
    }

    if r.Header.Get("X-Requested-With") == "XMLHttpRequest" {
        vals["notAJAX"] = false
    } else {
        vals["notAJAX"] = true
    }

    t := template.Must(template.New("stories").ParseFiles("templates/wrapper.html", "templates/stories"))
    t.Execute(w, vals)
}


func mark_story_read (w http.ResponseWriter, r *http.Request) {
    nb, err := initNewsblur()
    if err != nil {
        // TODO: This should not panic
        panic(err)
    }

    feed_id, err := strconv.Atoi(r.URL.Query().Get("feed_id"))
    story_id := r.URL.Query().Get("story_id")
    if err != nil {
        panic(err)
    }

    if err = nb.MarkStoryRead(feed_id, story_id); err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        fmt.Fprintf(w, "true")
        return
    }

    fmt.Fprintf(w, "true")
}


func main() {
    _ = fmt.Println
    r := mux.NewRouter()
    r.HandleFunc("/", index)
    r.HandleFunc("/feeds", index)
    r.HandleFunc("/feeds/{feed_id}", stories)
    r.HandleFunc("/stories/mark_read", mark_story_read)

    fmt.Println("Listening for browser connections")
    http.Handle("/", r)

    http.ListenAndServe(":4001", nil)
}
