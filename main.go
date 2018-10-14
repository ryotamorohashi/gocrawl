package main

import(
	"flag"
    "log"
    "strings"
    "net/http"
)

func main() {
	log.Println("検索するワードを入力してください：")
	input_word := SearchWordStdin()
    var word = flag.String("w", input_word, "今回検索されるワードはこちらです→")
    flag.Parse()
	*word = strings.Replace(*word, " ", "+", -1)
	//? of &num=? is number of query 
    firstURL := "https://www.google.co.jp/search?q=" + string(*word) + "&num=100"
    log.Println("検索URL：", firstURL)
    m := newMessage()
    go m.execute()
    m.req <- &request{
        url:   firstURL,
        depth: 2,
    }

    if err := http.ListenAndServe(":8080", nil); err != nil {
        log.Fatal("ListenAndSearver:", err)
    }
}