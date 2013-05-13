package main

import (
    //"encoding/json"
    "fmt"
    //"io/ioutil"
    "./newsblur"
//    "strconv"
//    "strings"
    "time"

    gocache "github.com/pmylund/go-cache"

)

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


func main() {
    nb, err := initNewsblur()
    if err != nil {
        panic(err)
    }

    nb.RefreshFeedStories(false)
    _ = fmt.Println
}
