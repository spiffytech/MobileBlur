package newsblur

import (
    "crypto/sha512"
    "encoding/json"
    "errors"
    "fmt"
    "html/template"
    "io/ioutil"
    "math"
    "net/http"
    "net/http/cookiejar"
    "net/url"
    "strconv"
    "time"
)

type Newsblur struct {
    Cookie string
    Profile Profile
    Threshold int
    ShowRead bool
    EmptyFeeds bool
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

type FeedInt interface {
    GetPS() int
    GetNT() int
    GetNG() int
}
func (feed *SocialFeed) GetPS() int {
    return feed.PS
}
func (feed *SocialFeed) GetNT() int {
    return feed.NT
}
func (feed *SocialFeed) GetNG() int {
    return feed.NG
}
func (feed *Feed) GetPS() int {
    return feed.PS
}
func (feed *Feed) GetNT() int {
    return feed.NT
}
func (feed *Feed) GetNG() int {
    return feed.NG
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
    JSONId string `json:"id"`
    JSONGuid string `json:"guid_hash"`
    JSONTitle string `json:"story_title"`
    //Date time.Time `json:"story_date"`
    JSONPrettyDate string `json:"short_parsed_date"`
    JSONContent template.HTML `json:"story_content"`
    JSONPermalink string `json:"story_permalink"`
    JSONReadStatus int `json:"read_status"`
    JSONTags []string `json:"story_tags"`
    JSONHasModifications int `json:"has_modifications"`
    JSONIntelligence Intelligence `json:"intelligence"`
}

type StoryInt interface {
    Score() int
    ReadStatus() int
    Content() template.HTML
    ID() string
    Title() string
    PrettyDate() string
    Permalink() string
    Intelligence() Intelligence
    HashStory() string
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
    JSONId string `json:"id"`
    JSONGuid string `json:"guid_hash"`
    JSONTitle string `json:"story_title"`
    JSONAuthor string `json:"story_authors"`
    //Date time.Time `json:"story_date"`
    JSONPrettyDate string `json:"short_parsed_date"`
    JSONContent template.HTML `json:"story_content"`
    JSONPermalink string `json:"story_permalink"`
    JSONReadStatus int `json:"read_status"`
    JSONStoryFeedID int `json:"story_feed_id"`
    JSONTags []string `json:"story_tags"`
    JSONHasModifications int `json:"has_modifications"`
    JSONIntelligence Intelligence `json:"intelligence"`
    JSONCommentCount int `json:"comment_count"`
    JSONStories []SocialStory
}


var nbURL = "http://www.newsblur.com"

func (nb *Newsblur) Login(username, password string) (cookie string, error error) {
    resp, err := http.PostForm(
        nbURL + "/api/login",
        url.Values{
            "username": {username},
            "password": {password},
        },
    )
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()

    type Response struct {
        Result string `json:"result"`
        Errors interface{} `json:"errors"`
        Code int
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

    fmt.Println(response)

    if response.Errors != nil {
        return "", errors.New("Could not log you in")
    } else {
        for _, c := range resp.Cookies() {
            if c.Name == "newsblur_sessionid" {
                cookie = c.Value
            }
        }

        if cookie == "" {
            return "", errors.New("Newsblur didn't give us a cookie")
        }

        nb.Cookie = cookie
        return cookie, nil
    }

}

// TODO: Deduplicate this logic
func (story *Story) HashStory() string {
    b := []byte(story.ID())
    hasher := sha512.New()
    hasher.Write(b)
    sha := fmt.Sprintf("%x", hasher.Sum(nil))
    return sha
}
func (story *SocialStory) HashStory() string {
    b := []byte(story.ID())
    hasher := sha512.New()
    hasher.Write(b)
    sha := fmt.Sprintf("%x", hasher.Sum(nil))
    return sha
}

func (feed *Feed) IsStale() (bool) {
    // TODO: Need to flesh this out to check the cache when I actually have a cache mechanism to check
    return true
}

func (nb *Newsblur) getFolders() (folder Folder) {
    folder.Feeds = getFolderFeeds(nb, nb.Profile.RawFolders)
    folder.Folders = getFolderFolders(nb, nb.Profile.RawFolders)
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


func (feed *Feed) GetStoryPage(nb *Newsblur, page int, force bool) ([]StoryInt) {
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

    var ret []StoryInt
    for i, _ := range storyList.Stories {
        ret = append(ret, &storyList.Stories[i])
    }

    return ret
}


func (feed *SocialFeed) GetSocialStoryPage(nb *Newsblur, page int, force bool) ([]StoryInt) {
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

    var ret []StoryInt
    for i, _ := range storyList.Stories {
        ret = append(ret, &storyList.Stories[i])
    }

    return ret
}


func (nb *Newsblur) NewClient() (*NBClient) {
    client := NBClient{}
    client.Client = &http.Client{}
    cookie := http.Cookie{Name: "newsblur_sessionid", Value: nb.Cookie}
    cookieJar, err := cookiejar.New(nil)
    if err != nil {
        panic(err)
    }

    u, _ := url.Parse("http://www.newsblur.com")
    cookieJar.SetCookies(u, []*http.Cookie{&cookie})
    client.Client.Jar = cookieJar
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
        fmt.Println(string(b))
        panic(err)
    }

    profile.SocialFeeds = make(map[string]SocialFeed)
    for _, feed := range profile.RawSocialFeeds {
        profile.SocialFeeds[feed.ID] = feed
    }

    nb.Profile = profile

    nb.getFolders()
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



func (story *Story) ID() (string) {
    return story.JSONId
}
func (story *Story) Title() (string) {
    return story.JSONTitle
}
func (story *Story) PrettyDate() (string) {
    return story.JSONPrettyDate
}
func (story *Story) Permalink() (string) {
    return story.JSONPermalink
}
func (story *Story) Intelligence() (Intelligence) {
    return story.JSONIntelligence
}
func (story *Story) ReadStatus() (int) {
    return story.JSONReadStatus
}
func (story *Story) Content() (template.HTML) {
    return story.JSONContent
}

func (story *SocialStory) ID() (string) {
    return story.JSONId
}
func (story *SocialStory) Title() (string) {
    return story.JSONTitle
}
func (story *SocialStory) PrettyDate() (string) {
    return story.JSONPrettyDate
}
func (story *SocialStory) Permalink() (string) {
    return story.JSONPermalink
}
func (story *SocialStory) Intelligence() (Intelligence) {
    return story.JSONIntelligence
}
func (story *SocialStory) ReadStatus() (int) {
    return story.JSONReadStatus
}
func (story *SocialStory) Content() (template.HTML) {
    return story.JSONContent
}


// TODO: Refactor this to be one function that takes a StoryInt and uses getters instead of properties
func (story Story) Score() (score int) {
    return scoreStory(story.JSONIntelligence)
}
func (story *SocialStory) Score() (score int) {
    return scoreStory(story.JSONIntelligence)
}

func scoreStory(intelligence Intelligence) (score int) {
    maxScore := 0.0
    maxScore = math.Max(maxScore, float64(intelligence.Tags))
    maxScore = math.Max(maxScore, float64(intelligence.Author))
    maxScore = math.Max(maxScore, float64(intelligence.Title))

    minScore := 0.0
    minScore = math.Min(minScore, float64(intelligence.Tags))
    minScore = math.Min(minScore, float64(intelligence.Author))
    minScore = math.Min(minScore, float64(intelligence.Title))

    if maxScore > 0 {
        score = 1;
    } else if minScore < 0 {
        score = -1;
    }

    if (score == 0) {
        if intelligence.Feed > 0 {
            score = 1
        } else if intelligence.Feed < 0 {
            score = -1
        } else {
            score = 0
        }
    }

    return;
}

type NBClient struct {
    *http.Client
}

func (client *NBClient) Get(url string) (resp *http.Response, clientErr error) {
    // Handles Newsblur's rate limiting transparently
    for {
        resp, clientErr = client.Client.Get(url)
        if resp.StatusCode == 429 {
            fmt.Println("Waiting for the rate limit (" + url + ")")
            time.Sleep(1 * time.Second)
        } else {
            fmt.Println("No more rate limit!")
            break
        }
    }

    return resp, clientErr
}
