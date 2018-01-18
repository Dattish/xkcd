package xkcd

import (
	"net/http"
	"encoding/json"
	"fmt"
	"os"
	"io"
	"errors"
)

type ComicData struct {
	Year       string `json:"year"`
	Month      string `json:"month"`
	Day        string `json:"day"`
	Num        int    `json:"num"`
	Title      string `json:"title"`
	SafeTitle  string `json:"safe_title"`
	ImageUrl   string `json:"img"`
	Transcript string `json:"transcript"`
	News       string `json:"news"`
	Link       string `json:"link"`
	AltText    string `json:"alt"`
}

type ResponseError struct {
	Status     string
	StatusCode int
}

func (e *ResponseError) Error() string {
	return fmt.Sprintf("%d: %s", e.StatusCode, e.Status)
}

func (comicData ComicData) String() string {
	return fmt.Sprintf("ComicData{Year=%v, Month=%v, Day=%v, Num=%v, Title=%v, SafeTitle=%v, ImageUrl=%v, Transcript=%v, News=%v, Link=%v, AltText=%v}",
		comicData.Year,
		comicData.Month,
		comicData.Day,
		comicData.Num,
		comicData.Title,
		comicData.SafeTitle,
		comicData.ImageUrl,
		comicData.Transcript,
		comicData.News,
		comicData.Link,
		comicData.AltText)
}

func (comicData ComicData) SaveComicImage(prefix string) error {
	return SaveComicImage(comicData.ImageUrl, prefix, comicData.ImageUrl[len("https://imgs.xkcd.com/comics/"):])
}

func GetComicData(comicNumber int) (*ComicData, error) {
	return getComicData(fmt.Sprintf("https://xkcd.com/%v/info.0.json", comicNumber))
}

func GetLatestComicData() (*ComicData, error) {
	return getComicData("https://xkcd.com/info.0.json")
}

func SaveComicImage(url string, prefix string, filename string) error {
	if filename == "" {
		return errors.New("empty filename")
	}

	filePath := prefix + filename
	response, responseErr := http.Get(url)
	if responseErr != nil {
		return responseErr
	}
	defer response.Body.Close()
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return &ResponseError{response.Status, response.StatusCode}
	}

	if _, err := os.Stat(filePath); err == nil {
		return errors.New("comic already exists")
	}

	file, fileErr := os.Create(filePath)
	defer file.Close()
	if fileErr != nil {
		return fileErr
	}

	_, copyErr := io.Copy(file, response.Body)
	if copyErr != nil {
		return copyErr
	}

	return nil
}

func getComicData(url string) (*ComicData, error) {
	response, responseErr := http.Get(url)
	if responseErr != nil {
		return nil, responseErr
	}
	defer response.Body.Close()
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return nil, &ResponseError{response.Status, response.StatusCode}
	}
	var comicData ComicData
	jsonErr := json.NewDecoder(response.Body).Decode(&comicData)
	if jsonErr != nil {
		return nil, jsonErr
	}
	return &comicData, nil
}
