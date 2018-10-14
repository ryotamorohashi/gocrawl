package main

import (
	"fmt"
	"log"
	"os"
)

type message struct {
    res  chan *respons
    req  chan *request
    quit chan int
}

type respons struct {
    url string
    err interface{}
}

type request struct {
    url   string
    depth int
}

func newMessage() *message {
    return &message{
        res:  make(chan *respons),
        req:  make(chan *request),
        quit: make(chan int),
    }
}

func (m *message) execute() {
    // number of worker
    count := 0
    urlMap := make(map[string]bool, 500)
    done := false
    for !done {
        select {
        case res := <-m.res:
            if res.err == nil {
				fmt.Printf("%s\n", res.url)
            } else {
                fmt.Fprintf(os.Stderr, "Error %s\n%v\n", res.url, res.err)
            }
        case req := <-m.req:
            if req.depth == 0 {
                break
            }

            if urlMap[req.url] {
                // is getted url
                break
            }
            urlMap[req.url] = true

            count++
            go Crawl(req.url, req.depth, m)
        case <-m.quit:
            count--
            if count == 0 {
                done = true
            }
        }
    }
    log.Println("scraping complete!")
    os.Exit(0)
}
