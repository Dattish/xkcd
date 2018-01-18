package main

import (
	"dattri.eu/xkcd"
	"os"
	"strings"
	"log"
	"fmt"
	"time"
)

func getComic(comicNumber int, retrieved chan *xkcd.ComicData, fail chan int) {
	comicData, err := xkcd.GetComicData(comicNumber)

	if err != nil {
		fail <- comicNumber
	} else {
		retrieved <- comicData
	}
}

func saveComic(comicData *xkcd.ComicData, prefix string, notifier chan int, fail chan int) {
	err := comicData.SaveComicImage(prefix)
	if err != nil {
		fail <- comicData.Num
	}
	notifier <- 1
}

func getAllComics(numberOfComics int, prefix string) ([]int) {
	formattedPrefix := prefix
	if !strings.HasSuffix(formattedPrefix, "/") {
		formattedPrefix = formattedPrefix + "/"
	}
	retrieved := make(chan *xkcd.ComicData)
	fail := make(chan int)

	for i := numberOfComics; i > 0; i-- {
		go getComic(i, retrieved, fail)
	}

	fails := make([]int, 0)
	notifier := make(chan int)
	saves := 0
	for counter := 0 ; counter < numberOfComics; {
		select {
		case comicData := <-retrieved:
			go saveComic(comicData, formattedPrefix, notifier, fail)
			counter++
			saves++
		case failedComic := <-fail:
			fails = append(fails, failedComic)
			counter++
		}
	}

	for counter := 0; counter < saves; {
		select {
		case <-notifier:
			counter++
		case failedComic := <-fail:
			fails = append(fails, failedComic)
			counter++
		default:
			time.Sleep(50 * time.Millisecond)
		}
	}

	return fails
}

func main() {
	comicData, cdErr := xkcd.GetLatestComicData()
	if cdErr != nil {
		log.Fatal(cdErr)
	}
	numberOfComics := comicData.Num

	directory := "images"
	dirErr := os.Mkdir(directory, 0777)
	if dirErr != nil && !os.IsExist(dirErr) {
		log.Fatal(dirErr)
	}
	fmt.Printf("Attempting to fetch %v comics.\n",numberOfComics)
	fails := getAllComics(numberOfComics, directory)
	for _, failure := range fails {
		fmt.Printf("Failed to fetch no. %d\n", failure)
	}
	fmt.Printf("Done, successfully fetched %d comics.\n", numberOfComics - len(fails))
}
