package main

import (
    "fmt"
    "regexp"
    "strings"
    "github.com/PuerkitoBio/goquery"
    "github.com/huandu/xstrings"
    "time"
    "strconv"
    "github.com/gin-gonic/gin"
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

func indexBlog() {
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

func compareWords(w1 string, w2 string) float64 {
    var s int
    var prev1 string
    var prev2 string
    score := 0.0

    s1 := xstrings.Len(w1)
    s2 := xstrings.Len(w2)
    if s1 < s2 {
        s = s1
    } else {
        s = s2
    }

    for i := 0; i < s; i++ {
        a := string(w1[i])
        b := string(w2[i])

        if i > 0 {
            if prev1 == prev2 && a == b {
                score += 1
            } else if prev1 == b && prev2 == a {
                score += 0.9
            } else if prev1 == prev2 || a == b {
                score += 0.75
            }
        }

        prev1 = a
        prev2 = b
    }

    return score / float64(s - 1)
}

func computeWeight(input []string, references []string, coeff int) float64 {
    weight := 0.0

    for _, w1 := range input {
        for _, w2 := range references {
            rate := compareWords(w1, w2)
            if rate >= appConfiguration["compare_words_ratio"].FloatValue {
                // A word has been identified
                weight += rate * float64(coeff)
                break
            }
        }
    }

    return weight
}

func computeWeightWithMap(input []string, references map[string] int) float64 {
    weight := 0.0

    for _, w1 := range input {
        for w2, coeff := range references {
            rate := compareWords(w1, w2)
            if rate >= appConfiguration["compare_words_ratio"].FloatValue {
                // A word has been identified
                weight += rate * float64(coeff)
                break
            }
        }
    }

    return weight
}

func search(keywords []string, whenDone chan []BlogPost)  {
    outcome := SortedList {}
    posts := extractBlogPosts()

    for _, e := range posts {
        weight := 0.0

        weight += computeWeight(keywords, strings.Split(e.Headline, " "), appConfiguration["headline_coeff"].IntValue)
        weight += computeWeight(keywords, strings.Split(e.Author, " "), appConfiguration["author_coeff"].IntValue)
        weight += computeWeightWithMap(keywords, e.ContentHash)

        outcome.Push(weight, e)
    }

    whenDone <- outcome.ToArray()
}

func main() {
    router := gin.Default()
    router.LoadHTMLGlob("templates/*")

    router.GET("/", func(c *gin.Context) {
        c.HTML(200, "index.tmpl", gin.H {})
    })

    router.GET("/search", func(c *gin.Context) {
        rawKeywords := strings.Split(c.Query("keywords"), "+")
        keywords := make([]string, len(rawKeywords))
        whenSearchIsDone := make(chan []BlogPost)

        for _, s := range rawKeywords {
            keywords = append(keywords, strings.ToLower(strings.Trim(s, " ")))
        }

        go search(keywords, whenSearchIsDone)

        posts := <- whenSearchIsDone

        c.JSON(200, gin.H {"posts": posts})
    })

    router.POST("/reset", func(c *gin.Context) {
        indexBlog()
        c.JSON(200, gin.H {})
    })

    router.Run()
}