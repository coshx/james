package main

type SortedListElement struct {
    Next *SortedListElement
    Value *BlogPost
    Weight float64
}

type SortedList struct {
    Head *SortedListElement
    Size int
}

func (this *SortedList) Push(weight float64, post *BlogPost) {
    e := this.Head
    var prev *SortedListElement

    if e == nil {
        this.Head = &SortedListElement {
            nil,
            post,
            weight,
        }
        this.Size++
        return
    }

    for i := 0; i < this.Size; i++ {
        if e.Weight < weight {
            f := SortedListElement {
                e,
                post,
                weight,
            }
            this.Size++

            if e == this.Head {
                this.Head = &f
            } else {
                prev.Next = &f
            }

            return
        }

        prev = e
        e = e.Next
    }

    if e == nil {
        f := SortedListElement {
            nil,
            post,
            weight,
        }
        prev.Next = &f
        this.Size++
    }
}

func (this *SortedList) ToArray() []BlogPost {
    e := this.Head
    outcome := make([]BlogPost, this.Size)

    for i := 0; i < this.Size; i++ {
        outcome = append(outcome, *e.Value)
        e = e.Next
    }

    return outcome
}

