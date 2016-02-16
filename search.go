package main

import(
    "strings"
    "github.com/huandu/xstrings"
)

func compareWords(input string, reference string) float64 {
    var s int
    var prev1 string
    var prev2 string
    score := 0.0

    s1 := xstrings.Len(input)
    s2 := xstrings.Len(reference)
    if s1 < s2 {
        s = s1
    } else {
        s = s2
    }

    for i := 0; i < s; i++ {
        a := string(input[i])
        b := string(reference[i])

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

    score /= float64(s1 - 1)

    return score
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

func search(keywords []string, whenDone chan []WrappedBlogPost)  {
    outcome := SortedList {}
    posts := extractBlogPosts()

    for _, e := range posts {
        weight := 0.0

        weight += computeWeight(keywords, strings.Split(e.Headline, " "), appConfiguration["headline_coeff"].IntValue)
        weight += computeWeight(keywords, strings.Split(e.Author, " "), appConfiguration["author_coeff"].IntValue)
        weight += computeWeightWithMap(keywords, e.ContentHash)

        if weight >= appConfiguration["minimum_weight"].FloatValue {
            outcome.Push(weight, e)
        }
    }

    whenDone <- (outcome.ToArray())
}