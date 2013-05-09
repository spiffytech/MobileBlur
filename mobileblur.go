package main

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
//    "strconv"
//    "strings"
    "time"
)

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

var _ time.Time
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

func main() {
    feeds := refreshFeeds()
    //fmt.Println(feeds)

    for _, feed := range feeds {
        feed.Refresh()
    }
}

func refreshFeeds() (map[string]Feed) {
    b, err := ioutil.ReadFile("profile.json")
    if err != nil {
        panic(err)
    }

    var profile Profile
    json.Unmarshal(b, &profile)
    return profile.Feeds
}

func (feed *Feed) Refresh() {
    b, err := ioutil.ReadFile("stories.json")
    if err != nil {
        panic(err)
    }

    var storyList StoryList
    json.Unmarshal(b, &storyList)
    for _, story := range storyList.Stories {
        fmt.Println(story.Permalink)
    }
}
