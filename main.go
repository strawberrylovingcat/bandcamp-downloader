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
1. Добавить указание прогресса скачивания
2. Добавить параллельную загрузку файлов
3. Заполнение метаданных файлов
*/
func main() {
	var urlVar string
	var pathVar string
	flag.StringVar(&urlVar, "url", "", "Link to download")
	flag.StringVar(&pathVar, "path", "", "Path for downloaded files. Default is working directory")
	flag.Parse()

	if "" == urlVar {
		panic("No url was passed")
	}

	if "" == pathVar {
		tempPathVar, err := os.Getwd()

		if err != nil {
			panic(err)
		}
		pathVar = tempPathVar
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

	for _, pack := range tracksRes.Packages {
		if pack.AlbumId == tracksRes.Current.Id {
			pathVar += "/" + pack.DownloadArtist + "/" + pack.DownloadTitle
		}
	}
	if err := os.MkdirAll(pathVar, 0755); nil != err {
		panic(err)
	}

	for _, trackElement := range tracksRes.Tracks {
		file, err := os.Create(pathVar + "/" + trackElement.Title + ".mp3")
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
	Current  current        `json:"current"`
	Tracks   []track        `json:"trackinfo"`
	Packages []musicPackage `json:"packages"`
}

type current struct {
	Title string `json:"title"`
	Id    int    `json:"id"`
}

type track struct {
	Title string    `json:"title"`
	File  trackFile `json:"file"`
}

type trackFile struct {
	Url string `json:"mp3-128"`
}

type musicPackage struct {
	AlbumId        int    `json:"album_id"`
	DownloadTitle  string `json:"download_title"`
	DownloadArtist string `json:"download_artist"`
}
