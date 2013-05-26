package newsblur

import (
    "encoding/json"
    "errors"
    "fmt"
    "html/template"
    "io/ioutil"
    "net/http"
    "net/url"
    "strconv"
)

var _ = fmt.Println

type Newsblur struct {
    Cookie string
    Feeds map[int]Feed
    Profile Profile
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
    Stories StoryList
}

type Folder struct {
    Name string
    Feeds map[int]Feed
    Folders []Folder
}

type FeedList struct {
    Feeds map[int]Feed
    Folders []Folder
    //Blurblogs map[int]Blurblogs
}

type Profile struct {
    RawFolders []interface{} `json:"folders"`
    Folder Folder
    Feeds map[string]Feed
}


type UserProfile struct {
    Username string `json:"username"`
}

type RealProfile struct {
    UserProfile UserProfile `json:"user_profile"`
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
    Content template.HTML `json:"story_content"`
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


func (feed *Feed) IsStale() (bool) {
    // TODO: Need to flesh this out to check the cache when I actually have a cache mechanism to check
    return true
}


func (nb *Newsblur) GetFolders() (folder Folder) {
    profile := nb.GetProfile()

    fmt.Println(profile.RawFolders)

    //folder.Feeds = getFolderFeeds(folder)
    folder.Folders = getFolderFolders(profile.RawFolders)
    nb.Profile.Folder = folder

    return folder
}


func getFolderFolders(folder []interface{}) (folders []Folder) {
    for _, item := range folder {
        fmt.Println(item)
        switch item.(type) {
            case float64:
            case interface{}:
                for folderName, val := range item.(map[string]interface{}) {
                    folder := Folder{}
                    folder.Name = folderName
                    //folder.Feeds = getFolderFeeds(folder)
                    folder.Folders = getFolderFolders(val.([]interface{}))
                    folders = append(folders, folder)
                }
        }
    }
    fmt.Println(folders)

    return folders
}


func (nb *Newsblur) GetFeeds() (map[int]Feed) {
    feeds := make(map[int]Feed)

    profile := nb.GetProfile()
    for feedID, feed := range profile.Feeds {
        feedIDInt, err := strconv.Atoi(feedID)
        if err != nil {
            panic(err)
        }
        feeds[feedIDInt] = feed
    }
    nb.Feeds = feeds
    return feeds
}


func (feed *Feed) GetStoryPage(nb *Newsblur, page int, force bool) (StoryList) {
    url := url.Values{"page": {strconv.Itoa(page)}}
    renderedURL := "/reader/feed/" + strconv.Itoa(feed.ID) + "?" + url.Encode()
    fmt.Println(renderedURL)
    req := nb.NewRequest("GET", renderedURL)
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

    feed.Stories = storyList
    return storyList
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

func (nb *Newsblur) RetrieveRealProfile() (RealProfile) {
    // TODO: Cache this, keyed on the session cookie value
    req := nb.NewRequest("GET", "/social/profile")
    client := http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        panic(err)
    }

    var profile RealProfile

    b, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        panic(err)
    }
    json.Unmarshal(b, &profile)
    return profile
}

func (nb *Newsblur) GetUsername() (string) {
    profile := nb.RetrieveRealProfile()
    return profile.UserProfile.Username
}


func (nb *Newsblur) Login(username, password string) (error) {
    resp, err := http.PostForm(nbURL + "/api/login", url.Values{"username": {username}, "password": {password}})
    if err != nil {
        panic(err)
    }

    for _, cookie := range resp.Cookies() {
        if cookie.Name == "newsblur_sessionid" {
            nb.Cookie = cookie.Value
            return nil
        }
    }

    return errors.New("No newsblur_sessionid cookie returned by login operation")
}


func (nb *Newsblur) GetProfile() (Profile) {
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
    err = json.Unmarshal(b, &profile)
    if err != nil {
        panic(err)
    }

    nb.Profile = profile
    return profile
}
