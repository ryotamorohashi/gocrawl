package main

import(
    "log"
    "net/http"
    "net/url"
	"strings"
	"bufio"
	"os"

    "github.com/PuerkitoBio/goquery"
)

func Crawl(url string, depth int, m *message){
	defer func(){ m.quit <- 0 }()

	// get url from www
	urls, err := Fetch(url)

	// send results
	m.res <- &respons{
		url: url,
		err: err,
	}

	if err == nil {
		for _, url := range urls {
			// send new requests
			m.req <- &request {
				url: url,
				depth: depth - 1,
			}
		}
	}
}

func Fetch(u string) (urls []string, err error) {
    baseUrl, err := url.Parse(u)
    if err != nil {
        log.Fatal(err)
    }

    res, err := http.Get(baseUrl.String())
    if err != nil {
        return
    }
    defer res.Body.Close()
    doc, err := goquery.NewDocumentFromReader(res.Body)
    if err != nil {
        return
    }

	urls = make([]string, 0)
	// only for url
    doc.Find(".r").Each(func(_ int, srg *goquery.Selection) {
        srg.Find("a").Each(func(_ int, s *goquery.Selection) {
            href, exists := s.Attr("href")
            if exists {
				// catch url
                reqUrl, err := baseUrl.Parse(href)
                if err == nil {
                    urls = append(urls, reqUrl.String())
                }
            }
        })
	})
	// only for title
	// doc.Find(".g").Each(func(i int, srg *goquery.Selection) {
	// 	title := srg.Find("h3").Text()
	// 	log.Println(title)
    //     if err == nil {
    //         urls = append(urls, title)
    //     }
    // })

    return
}

func SearchWordStdin() (stringInput string) {
	scanner := bufio.NewScanner(os.Stdin)

	scanner.Scan()
	stringInput = scanner.Text()
	stringInput = strings.TrimSpace(stringInput)
	return
}
