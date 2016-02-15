package main

import (
    "fmt"
    "regexp"
    "strings"
    "github.com/PuerkitoBio/goquery"
    "github.com/huandu/xstrings"
    "time"
    "strconv"
)

func enginePost(blogPosts map[string] *BlogPost, href string, whenDone chan bool) {
    doc, err := goquery.NewDocument(appConfiguration["website_root"].StringValue + href)

    if err != nil {
        fmt.Printf("Impossible to reach %s\n", href)
        whenDone <- false
        return
    }

    author := doc.Find(".james-author").Text()
    headline := doc.Find(".james-headline").Text()
    text := doc.Find(".james-content").Text()
    timestamp, _ := doc.Find(".james-date").Attr("data-timestamp")
    parsedTime, _ := strconv.ParseInt(timestamp, 10, 64)
    hash := make(map[string] int)

    post := BlogPost{
        href,
        author,
        headline,
        time.Unix(parsedTime, 0),
        hash,
    }

    buffer := ""
    for i := 0; i < xstrings.Len(text); i++ {
        c := string(text[i])
        if isMatching, _ := regexp.MatchString("[a-zA-Z0-9]", c); isMatching == true {
            buffer += c
        } else if buffer != "" {
            if xstrings.Len(buffer) >= appConfiguration["min_word_length"].IntValue {
                buffer = strings.ToLower(buffer)
                e, exists := hash[buffer]
                if exists {
                    hash[buffer] = e + 1
                } else {
                    hash[buffer] = 1
                }
            }
            buffer = ""
        }
    }

    blogPosts[href] = &post

    whenDone <- true
}

func fetchHome(blogPosts map[string] *BlogPost, whenDone chan bool) {
    doc, err := goquery.NewDocument(appConfiguration["blog_homepage"].StringValue)

    if err != nil {
        fmt.Printf("Impossible to reach homepage")
        whenDone <- false
        return
    }

    var buffer = make(chan bool)
    started := 0
    ended := 0

    go func() {
        for {
            select {
            case <- buffer:
                ended++
                if ended == started {
                    close(buffer)
                    fmt.Printf("DONE\n")
                    whenDone <- true
                }
            }
        }
    }()

    doc.Find(".james-blogposts .james-post").Each(func (i int, s *goquery.Selection) {
        link, isLinkExisting := s.Attr("href")
        _, isPostAlreadyIndexed := blogPosts[link]

        if isLinkExisting && !isPostAlreadyIndexed {
            started++
            go enginePost(blogPosts, link, buffer)
        }
    })

    if started == 0 {
        fmt.Printf("No post to analyse!\n")
        whenDone <- false
    } else {
        fmt.Printf("Analysing %d posts... ", started)
    }
}

func main() {
    fmt.Printf("Starting indexation...\n")

    whenFetchHomeDone := make(chan bool)

    fmt.Printf("Extracting blog posts... ")
    blogPosts := extractBlogPosts()
    fmt.Printf("DONE\n")

    go fetchHome(blogPosts, whenFetchHomeDone)
    hasToSave := <- whenFetchHomeDone

    if hasToSave {
        fmt.Printf("Saving latest data... ")
        saveBlogPosts(blogPosts)
        fmt.Printf("DONE\n");
    }

    fmt.Printf("Indexation OVER\n");
}