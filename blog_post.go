package main

import (
    "time"
    "io/ioutil"
    "os"
    "github.com/ricardolonga/jsongo"
    "github.com/antonholmquist/jason"
)

type BlogPost struct {
    LocalPath string
    Author string
    Headline string
    Date time.Time
    ContentHash map[string] int
    Brief string
}

func resetIndexation() {
    os.Remove(appConfiguration["saved_data_filename"].StringValue)
}

func saveBlogPosts(hash map[string] *BlogPost) {
    outcome := jsongo.Object()

    for _, v1 := range hash {
        o1 := jsongo.Object()
        contentHash := jsongo.Object()

        o1.Put("author", v1.Author)
        o1.Put("headline", v1.Headline)
        o1.Put("date", int32(v1.Date.Unix()))

        for k2, v2 := range v1.ContentHash {
            contentHash.Put(k2, v2)
        }

        o1.Put("contentHash", contentHash)
        o1.Put("brief", v1.Brief)

        outcome.Put(v1.LocalPath, o1)
    }

    ioutil.WriteFile(appConfiguration["saved_data_filename"].StringValue, []byte(outcome.String()), 0644)
}

func extractBlogPosts() map[string] *BlogPost {
    outcome := make(map[string] *BlogPost)
    rawContent, readErr := ioutil.ReadFile(appConfiguration["saved_data_filename"].StringValue)

    if readErr == nil {
        json, _ := jason.NewObjectFromBytes(rawContent)

        for k1, _ := range json.Map() {
            hash := make(map[string] int)
            author, _ := json.GetString(k1, "author")
            headline, _ := json.GetString(k1, "headline")
            date, _ := json.GetInt64(k1, "date")
            contentHash, _ := json.GetObject(k1, "contentHash")
            brief, _ := json.GetString(k1, "brief")

            post := BlogPost {
                k1,
                author,
                headline,
                time.Unix(date, 0),
                hash,
                brief,
            }

            for k2, v2 := range contentHash.Map() {
                n, _ := v2.Int64()
                hash[k2] = int(n)
            }

            outcome[k1] = &post
        }
    }

    return outcome
}

func toJSON(posts []WrappedBlogPost) string {
    outcome := jsongo.Array()

    for _, e := range posts {
        o := jsongo.Object()

        o.Put("Path", appConfiguration["website_root"].StringValue + e.Value.LocalPath)
        o.Put("Author", e.Value.Author)
        o.Put("Headline", e.Value.Headline)
        o.Put("Date", int32(e.Value.Date.Unix()))
        o.Put("Brief", e.Value.Brief)
        o.Put("Weight", e.Weight)

        outcome.Put(o)
    }

    return outcome.String()
}