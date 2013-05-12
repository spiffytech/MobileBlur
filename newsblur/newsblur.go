package newsblur

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "net/http"
    "net/url"
    "strconv"
)

type Newsblur struct {
    Cookie string
}

type Feed struct {
    ID int `json:"id"`
    PS int `json:"ps"`
    NT int `json:"nt"`
    NG int `json:"ng"`
    UpdatedSecondsAgo int `json:"updated_seconds_ago"`
    Address string `json:"feed_address"`
    Title string `json:"feed_title"`
    Link string `json:"feed_link"`
}

type Profile struct {
    Folders []interface{} `json:"folders"`
    Feeds map[string]Feed
}

type Intelligence struct {
    Feed int `json:"feed"`
    Tags int `json:"tags"`
    Author int `json:"author"`
    Title int `json:"title"`
}

type Story struct {
    ID string `json:"id"`
    GUID string `json:"guid_hash"`
    //Date time.Time `json:"story_date"`
    Title string `json:"story_title"`
    Content string `json:"story_content"`
    Permalink string `json:"story_permalink"`
    ReadStatus int `json:"read_status"`
    Tags []string `json:"story_tags"`
    HasModifications int `json:"has_modifications"`
    Intelligence Intelligence `json:"intelligence"`
}

type StoryList struct {
    Stories []Story `json:"stories"`
}


var nbURL = "http://www.newsblur.com"


func (feed *Feed) Refresh(nb Newsblur) {
    req := nb.NewRequest("GET", "/reader/feed/" + strconv.Itoa(feed.ID))
    client := http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        panic(err)
    }

    b, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        panic(err)
    }

    var storyList StoryList
    json.Unmarshal(b, &storyList)
    // TODO: Store this stuff in the cache
    for _, story := range storyList.Stories {
        fmt.Println(story.Permalink)
    }
}


func (nb *Newsblur) NewRequest(method, path string) (*http.Request) {
    req, err := http.NewRequest("GET", nbURL + path, nil)
    if err != nil {
        panic(err)
    }
    cookie  := http.Cookie{Name: "newsblur_sessionid", Value: nb.Cookie}
    req.AddCookie(&cookie)
    return req
}


func (nb *Newsblur) Login(username, password string) (string) {
    resp, err := http.PostForm(nbURL + "/api/login", url.Values{"username": {username}, "password": {password}})
    if err != nil {
        panic(err)
    }

    for _, cookie := range resp.Cookies() {
        if cookie.Name == "newsblur_sessionid" {
            nb.Cookie = cookie.Value
            return nb.Cookie
        }
    }

    // TODO: Make the login thing check that the login was actually successful
    // Return an error if it wasn't, instead of just panicing.
    panic("No newsblur_sessionid cookie returned")
}


func (nb *Newsblur) RetrieveProfile() (map[string]Feed) {
    req := nb.NewRequest("GET", "/reader/feeds")
    client := http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        panic(err)
    }

    var profile Profile

    b, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        panic(err)
    }
    json.Unmarshal(b, &profile)

    return profile.Feeds
}
