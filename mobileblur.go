package main

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
//    "strconv"
//    "strings"
//    "time"
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

func main() {
    feeds := refreshFeeds()
    fmt.Println(feeds)
}

func refreshFeeds() (map[string]Feed) {
    b, err := ioutil.ReadFile("test.json")
    if err != nil {
        panic(err)
    }

    var profile Profile
    json.Unmarshal(b, &profile)
    return profile.Feeds
}
