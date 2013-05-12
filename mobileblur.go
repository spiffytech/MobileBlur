package main

import (
    //"encoding/json"
    "fmt"
    //"io/ioutil"
    "./newsblur"
//    "strconv"
//    "strings"
)

type MyFeed struct {
    *newsblur.Feed
}

type FeedList struct {
    Feeds []MyFeed
}

func (feed *MyFeed) IsStale() (bool) {
    // TODO: Need to flesh this out to check the cache when I actually have a cache mechanism to check
    return true
}

func (feedlist *FeedList) Refresh(nb newsblur.Newsblur, force bool) {
    for _, feed := range feedlist.Feeds {
        if feed.IsStale() || force == true {
            feed.Refresh(nb)
        }
    }
}

func retrieveCookie() (string) {
    // TODO: Retrieve cookie from user response here, instead of logging in to Newsblur
    var nb newsblur.Newsblur
    return nb.Login("mbtest1", "mbtest1");
}

func main() {
    var nb newsblur.Newsblur

    nb.Login("mbtest1", "mbtest1");

    feeds := nb.RetrieveProfile()
    fmt.Println(feeds)

    for _, feed := range feeds {
        feed.Refresh(nb)
    }
}
