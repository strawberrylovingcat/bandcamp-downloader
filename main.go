package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"github.com/PuerkitoBio/goquery"
	"io"
	"net/http"
	"os"
)

/*
TODO
1. Добавить указание директории для сохранения файлов
2. Добавить указание прогресса скачивания
3. Добавить параллельную загрузку файлов
4. Заполнение метаданных файлов
*/
func main() {
	var urlVar string
	flag.StringVar(&urlVar, "url", "", "Link to download")
	flag.Parse()

	if "" == urlVar {
		panic("No url was passed")
	}
	response, err := http.Get(urlVar)
	defer response.Body.Close()

	if nil != err {
		panic(err)
	}
	body, _ := io.ReadAll(response.Body)

	docParser, _ := goquery.NewDocumentFromReader(bytes.NewReader(body))
	pageData, _ := docParser.Find("script[data-tralbum]").Attr("data-tralbum")
	var tracksRes tracksResponse

	if err := json.Unmarshal([]byte(pageData), &tracksRes); nil != err {
		panic(err)
	}
	homeDir, err := os.UserHomeDir()

	if err != nil {
		panic(err)
	}

	for _, trackElement := range tracksRes.Tracks {
		file, err := os.Create(homeDir + "/" + trackElement.Title + ".mp3")
		defer file.Close()

		if err != nil {
			panic(err)
		}

		response, err := http.Get(trackElement.File.Url)
		defer response.Body.Close()

		if err != nil {
			panic(err)
		}
		io.Copy(file, response.Body)
	}
}

type tracksResponse struct {
	Current current `json:"current"`
	Tracks  []track `json:"trackinfo"`
}

type current struct {
	Title string `json:"title"`
}

type track struct {
	Title string    `json:"title"`
	File  trackFile `json:"file"`
}

type trackFile struct {
	Url string `json:"mp3-128"`
}
