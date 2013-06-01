package newsblur

import (
    "crypto/sha512"
    "encoding/json"
    "errors"
    "fmt"
    "html/template"
    "io/ioutil"
    "net/http"
    "net/http/cookiejar"
    "net/url"
    "strconv"
)

type Newsblur struct {
    Cookie string
    Feeds map[int]Feed
    Profile Profile
}

type NBCoreResponse struct {
    Result string `json:"result"`
    Code int `json:"code"`
    Errors []string `json:"errors"`
    Payload interface{} `json:"payload"`
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
    Feeds []Feed
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
    RawSocialFeeds []SocialFeed `json:"social_feeds"`
    SocialFeeds map[string]SocialFeed
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
    Title string `json:"story_title"`
    //Date time.Time `json:"story_date"`
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

type SocialStoryList struct {
    Stories []SocialStory `json:"stories"`
}

type SocialFeed struct {
    Title string `json:"feed_title"`
    ID string `json:"id"`
    SocialID int `json:"subscription_user_id"`
    Link string `json:"feed_link"`
    PS int `json:"ps"`
    NT int `json:"nt"`
    NG int `json:"ng"`
    Stories []SocialStory
}

type SocialStory struct {
    ID string `json:"id"`
    GUID string `json:"guid_hash"`
    Title string `json:"story_title"`
    Author string `json:"story_authors"`
    //Date time.Time `json:"story_date"`
    Content template.HTML `json:"story_content"`
    Permalink string `json:"story_permalink"`
    ReadStatus int `json:"read_status"`
    StoryFeedID int `json:"story_feed_id"`
    Tags []string `json:"story_tags"`
    HasModifications int `json:"has_modifications"`
    Intelligence Intelligence `json:"intelligence"`
    CommentCount int `json:"comment_count"`
    Stories []SocialStory
}


var nbURL = "http://www.newsblur.com"


// TODO: Deduplicate this logic
func (story *Story) HashStory() string {
    b := []byte(story.ID)
    hasher := sha512.New()
    hasher.Write(b)
    sha := fmt.Sprintf("%x", hasher.Sum(nil))
    return sha
}
func (story *SocialStory) HashStory() string {
    b := []byte(story.ID)
    hasher := sha512.New()
    hasher.Write(b)
    sha := fmt.Sprintf("%x", hasher.Sum(nil))
    return sha
}

func (feed *Feed) IsStale() (bool) {
    // TODO: Need to flesh this out to check the cache when I actually have a cache mechanism to check
    return true
}


func (nb *Newsblur) GetFolders() (folder Folder) {
    profile := nb.Profile

    //folder.Feeds = getFolderFeeds(folder)
    folder.Feeds = getFolderFeeds(nb, profile.RawFolders)
    folder.Folders = getFolderFolders(nb, profile.RawFolders)
    nb.Profile.Folder = folder

    return folder
}


func getFolderFeeds(nb *Newsblur, folder []interface{}) (feeds []Feed) {
    for _, item := range folder {
        switch item.(type) {
            case float64:
                feedID := strconv.Itoa(int(item.(float64)))
                // TODO: Find a way to grab a reference to the feed, instead of a copy
                feeds = append(feeds, nb.Profile.Feeds[feedID])
            case interface{}:
        }
    }
    return feeds
}


func getFolderFolders(nb *Newsblur, folder []interface{}) (folders []Folder) {
    for _, item := range folder {
        switch item.(type) {
            case float64:
            case interface{}:
                for folderName, val := range item.(map[string]interface{}) {
                    folder := Folder{}
                    folder.Name = folderName
                    //folder.Feeds = getFolderFeeds(folder)
                    folder.Feeds = getFolderFeeds(nb, val.([]interface{}))
                    folder.Folders = getFolderFolders(nb, val.([]interface{}))
                    folders = append(folders, folder)
                }
        }
    }

    return folders
}


// TODO: Do I need things like this function since I refactored for folder support?
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
    u, err := url.Parse(nbURL + "/reader/feed/" + strconv.Itoa(feed.ID))
    if err != nil {
        panic(err)
    }
    q := u.Query()
    q.Set("page", strconv.Itoa(page))
    u.RawQuery = q.Encode()
    client := nb.NewClient()
    resp, err := client.Get(u.String())
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


func (feed *SocialFeed) GetSocialStoryPage(nb *Newsblur, page int, force bool) (SocialStoryList) {
    u, err := url.Parse(nbURL + "/social/stories/" + strconv.Itoa(feed.SocialID) + "/")  // Trailing slash is necessary or Newsblur 404s
    if err != nil {
        panic(err)
    }
    q := u.Query()
    q.Set("page", strconv.Itoa(page))
    u.RawQuery = q.Encode()
    client := nb.NewClient()
    resp, err := client.Get(u.String())
    if err != nil {
        panic(err)
    }

    b, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        panic(err)
    }

    var storyList SocialStoryList
    json.Unmarshal(b, &storyList)
    // TODO: Store this stuff in the cache

    feed.Stories = storyList.Stories
    return storyList
}


func (nb *Newsblur) NewClient() (*http.Client) {
    client := http.Client{}
    cookie := http.Cookie{Name: "newsblur_sessionid", Value: nb.Cookie}
    cookieJar, err := cookiejar.New(nil)
    if err != nil {
        panic(err)
    }

    u, _ := url.Parse("http://www.newsblur.com")
    cookieJar.SetCookies(u, []*http.Cookie{&cookie})
    client.Jar = cookieJar
    return &client
}


func (nb *Newsblur) RetrieveRealProfile() (RealProfile) {
    // TODO: Cache this, keyed on the session cookie value
    client := nb.NewClient()
    resp, err := client.Get(nbURL + "/social/profile")
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
    client := nb.NewClient()
    resp, err := client.Get(nbURL + "/reader/feeds")
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

    profile.SocialFeeds = make(map[string]SocialFeed)
    for _, feed := range profile.RawSocialFeeds {
        profile.SocialFeeds[feed.ID] = feed
    }

    nb.Profile = profile
    return profile
}


func (nb *Newsblur) MarkStoryRead(feedID int, storyID string) (error) {
    client := nb.NewClient()
    resp, err := client.PostForm(
        nbURL + "/reader/mark_story_as_read",
        url.Values{"feed_id": {strconv.Itoa(feedID)}, "story_id": {storyID}},
    )
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    type Response struct {
        Result string `json:"result"`
    }
    var response Response

    b, err := ioutil.ReadAll(resp.Body)
    fmt.Println(string(b))
    if err != nil {
        panic(err)
    }
    err = json.Unmarshal(b, &response)
    if err != nil {
        panic(err)
    }

    if response.Result == "ok" {
        return nil
    } else {
        return errors.New("Feed not could not be marked read")
    }
}


func (nb *Newsblur) MarkStoriesReadBulk(stories map[string][]string) (error) {
    b, err := json.Marshal(stories)
    if err != nil {
        panic(err)
    }

    client := nb.NewClient()
    resp, err := client.PostForm(
        nbURL + "/reader/mark_feed_stories_as_read",
        url.Values{"feeds_stories": {string(b)}},
    )
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    type Response struct {
        Result string `json:"result"`
    }
    var response Response

    b, err = ioutil.ReadAll(resp.Body)
    if err != nil {
        panic(err)
    }
    err = json.Unmarshal(b, &response)
    if err != nil {
        panic(err)
    }

    if response.Result == "ok" {
        return nil
    } else {
        return errors.New("Feeds not could not be marked read")
    }
}


func (nb *Newsblur) MarkSocialStoriesRead(stories map[string]map[string][]string) (error) {
    b, err := json.Marshal(stories)
    if err != nil {
        panic(err)
    }

    client := nb.NewClient()
    resp, err := client.PostForm(
        nbURL + "/reader/mark_social_stories_as_read",
        url.Values{"users_feeds_stories": {string(b)}},
    )
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    type Response struct {
        Result string `json:"result"`
    }
    var response Response

    b, err = ioutil.ReadAll(resp.Body)
    fmt.Println(string(b))
    if err != nil {
        panic(err)
    }
    err = json.Unmarshal(b, &response)
    if err != nil {
        panic(err)
    }

    if response.Result == "ok" {
        return nil
    } else {
        return errors.New("Feed not could not be marked read")
    }
}


func (nb *Newsblur) MarkStoryUnread(feedID, storyID string) (error) {
    client := nb.NewClient()
    resp, err := client.PostForm(
        nbURL + "/reader/mark_story_as_unread",
        url.Values{
            "feed_id": {feedID},
            "story_id": {storyID},
        },
    )
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    type Response struct {
        Result string `json:"result"`
    }
    var response Response

    b, err := ioutil.ReadAll(resp.Body)
    fmt.Println(string(b))
    if err != nil {
        panic(err)
    }
    err = json.Unmarshal(b, &response)
    if err != nil {
        panic(err)
    }

    if response.Result == "ok" {
        return nil
    } else {
        return errors.New("Feeds not could not be marked read")
    }
}
