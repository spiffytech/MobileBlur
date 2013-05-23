package main

import (
    "fmt"
    "net/http"
    "html/template"
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
    // TODO: Retrieve cookie from user response here, instead of logging in to Newsblur
    var nb newsblur.Newsblur
    err := nb.Login("mbtest1", "mbtest1");
    if err != nil {
        return nb, err
    }

    return nb, nil
}


    type Test1 struct {
        Val1 string
        Val2 int
    }
    type Test2 struct {
        Val3 string
        Val4 int
    }

    func (t *Test2) Dostuff() (string) {
        return "look, it's stuff!"
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

    //feeds := nb.RefreshFeedStories(false)
    feeds := nb.GetFeeds()

    vals := map[string]interface{}{"Feeds": feeds}

    t := template.Must(template.New("index").ParseFiles("templates/index"))
    t.Execute(w, vals)
}


func main() {
    _ = fmt.Println
    r := mux.NewRouter()
    r.HandleFunc("/", index)

    fmt.Println("Listening for browser connections")
    http.Handle("/", r)

    http.ListenAndServe(":4001", nil)
}
