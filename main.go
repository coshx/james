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
    //"website_root": "http://www.coshx.com",
    "website_root": "http://localhost:4567/",
    //"blog_homepage": "http://www.coshx.com/blog/",
    "blog_homepage": "http://localhost:4567/blog/",
    "min_word_length": "3",
}

// struct Type Article {

// }

var articles = list.New()

func enginePost(href string, whenDone chan bool) {
    doc, err := goquery.NewDocument(href)

    if err != nil {
        fmt.Printf("Impossible to reach %s\n", href)
        whenDone <- false
        return
    }

    //author := doc.Find(".james-author").Text()
    //headline := doc.Find(".james-headline").Text()
    text := doc.Find(".james-content").Text()
    //timestamp, _ := doc.Find(".james-date").Attr("data-timestamp")

    buffer := ""
    var hash = make(map[string] int)
    for i := 0; i < xstrings.Len(text); i++ {
        c := string(text[i])
        if isMatching, _ := regexp.MatchString("[a-zA-Z0-9]", c); isMatching == true {
            buffer += c
        } else if buffer != "" {
            if xstrings.Len(buffer) >= /*appConfiguration["min_word_length"]*/ 3 {
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

func fetchHome(whenDone chan bool) {
    doc, err := goquery.NewDocument(appConfiguration["blog_homepage"])

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
        }
    }()

    doc.Find(".james-blogposts .james-post").Each(func (i int, s *goquery.Selection) {
        link, isExisting := s.Attr("href")
        if isExisting {
            started++
            go enginePost(appConfiguration["website_root"] + link, buffer)
        }
    })
}

func main() {
    var c = make(chan bool)
    go fetchHome(c)
    <- c
}