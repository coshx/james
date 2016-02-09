package main

import (
    "fmt"
    "container/list"
    "regexp"
    "strings"
    "github.com/PuerkitoBio/goquery"
    "github.com/huandu/xstrings"
)

var appConfiguration = map[string] string {
    "blog_root": "http://www.coshx.com",
    "blog_homepage": "http://www.coshx.com/blog/",
}

var articles = list.New()

func enginePage(href string, whenDone chan bool) {
    doc, err := goquery.NewDocument(href)

    if err != nil {
        fmt.Printf("Impossible to reach %s\n", href)
        whenDone <- false
        return
    }

    text := doc.Find(".post-content").Text()

    buffer := ""
    var hash = make(map[string] int)
    for i := 0; i < xstrings.Len(text); i++ {
        c := string(text[i])
        if isMatching, _ := regexp.MatchString("[a-zA-Z0-9]", c); isMatching == true {
            buffer += c
        } else if buffer != "" {
            if xstrings.Len(buffer) >= 4 {
                buffer = strings.ToLower(buffer)
                e, exists := hash[buffer]
                if exists {
                    hash[buffer] = e + 1
                } else {
                    hash[buffer] = 0
                }
            }
            buffer = ""
        }
    }

    articles.PushBack(hash)
    whenDone <- true
}

func getHome(whenDone chan bool) {
    doc, _ := goquery.NewDocument(appConfiguration["blog_homepage"])

    var buffer = make(chan bool)
    started := 0
    ended := 0
    go func() {
        select {
        case <- buffer:
            ended++
            if ended == started {
                close(buffer)
                if articles.Len() > 0 {
                    element := articles.Front()
                    for {
                        fmt.Printf("%v\n", element)
                        element = element.Next()
                        if element == nil {
                            break
                        }
                    }
                }
                whenDone <- true
            }
        }
    }()

    doc.Find(".blogposts a").Each(func (i int, s *goquery.Selection) {
        link, _ := s.Attr("href")
        started++
        go enginePage(appConfiguration["blog_root"] + link, buffer)
    })
}

func main() {
    var c = make(chan bool)
    go getHome(c)
    <- c
}