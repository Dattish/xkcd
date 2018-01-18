package main

import (
	"dattri.eu/xkcd"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
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

func saveComic(comicData *xkcd.ComicData, prefix string, success chan int, fail chan int) {
	err := comicData.SaveComicImage(prefix)
	if err != nil {
		fail <- comicData.Num
	} else {
		success <- comicData.Num
	}
}

func getAllComics(numberOfComics int, prefix string) []int {
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
	comicsToSave := 0
	saveSuccess := make(chan int)
	saveFailure := make(chan int)
	for counter := 0; counter < numberOfComics; {
		select {
		case comicData := <-retrieved:
			go saveComic(comicData, formattedPrefix, saveSuccess, saveFailure)
			counter++
			comicsToSave++
		case failedComic := <-fail:
			fails = append(fails, failedComic)
			counter++
		}
	}

	for counter := 0; counter < comicsToSave; {
		select {
		case <-saveSuccess:
			counter++
		case failedComic := <-saveFailure:
			fails = append(fails, failedComic)
			counter++
		default:
			time.Sleep(50 * time.Millisecond)
		}
	}

	return fails
}

func saveAllComics(dirFlag *string) {
	comicData, cdErr := xkcd.GetLatestComicData()
	if cdErr != nil {
		log.Fatal(cdErr)
	}
	numberOfComics := comicData.Num
	if *dirFlag == "" {
		*dirFlag = "images"
	}
	directory := *dirFlag
	dirErr := os.Mkdir(directory, 0777)
	if dirErr != nil && !os.IsExist(dirErr) {
		log.Fatal(dirErr)
	}
	fmt.Printf("Attempting to fetch %v comics.\n", numberOfComics)
	fails := getAllComics(numberOfComics, directory)
	for _, failure := range fails {
		fmt.Printf("Failed to fetch no. %d\n", failure)
	}
	fmt.Printf("Done, successfully fetched %d comics.\n", numberOfComics-len(fails))
	if len(fails) > 0 {
		fmt.Println(
			"Note that some comics such as https://www.xkcd.com/1663/ cannot be downloaded this way" +
				"\ndue to their json not pointing at an image. These are usually interactive comics.")
	}
}

func saveLatestComic(dirFlag *string) {
	mkDirIfNeeded(dirFlag)
	comicData, cdErr := xkcd.GetLatestComicData()
	fmt.Println("Latest comic is no.", comicData.Num)
	if cdErr != nil {
		log.Fatal(cdErr)
	}
	err := comicData.SaveComicImage(*dirFlag)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Successfully fetched no.", comicData.Num)
}

func saveSpecificComic(dirFlag *string, specificFlag *int) {
	mkDirIfNeeded(dirFlag)

	fmt.Println("Fetching comic no.", *specificFlag)
	comicData, cdErr := xkcd.GetComicData(*specificFlag)
	if cdErr != nil {
		log.Fatal(cdErr)
	}
	err := comicData.SaveComicImage(*dirFlag)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Successfully fetched no.", comicData.Num)
}
func mkDirIfNeeded(dirFlag *string) {
	if *dirFlag != "" {
		dirErr := os.Mkdir(*dirFlag, 0777)
		if dirErr != nil && !os.IsExist(dirErr) {
			log.Fatal(dirErr)
		}
		if !strings.HasSuffix(*dirFlag, "/") {
			*dirFlag = *dirFlag + "/"
		}
	}
}

func main() {

	allFlag := flag.Bool("a", false, "Download all comics")
	dirFlag := flag.String("d", "", "Set the directory to download the comic[s] to")
	latestFlag := flag.Bool("l", false, "Download the latest comic")
	specificFlag := flag.Int("n", 0, "Download a specific comic")
	flag.Parse()

	if flag.NFlag() < 1 || (flag.NFlag() == 1 && *dirFlag != "") {
		saveAllComics(dirFlag)
	} else if *allFlag {
		saveAllComics(dirFlag)
	} else if *latestFlag {
		saveLatestComic(dirFlag)
	} else if *specificFlag > 0 {
		saveSpecificComic(dirFlag, specificFlag)
	}

}
