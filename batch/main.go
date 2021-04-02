package main

import (
	"fmt"
	"os"
	"time"

	"github.com/KotaroKimura/note_analyzer/batch/chrome"
)

var (
	noteURL = os.Getenv("NOTE_URL")
	titleSELECTOR = os.Getenv("TITLE_SELECTOR")
	viewSELECTOR = os.Getenv("VIEW_SELECTOR")
	commentSELECTOR = os.Getenv("COMMENT_SELECTOR")
	sukiSELECTOR = os.Getenv("SUKI_SELECTOR")
)

func main() {
	c, err := chrome.NewClient()
	if err != nil {
		fmt.Printf("%+v.\n", err)
		return
	}
	defer c.Close()

	if err = c.P.Navigate(noteURL); err != nil {
		fmt.Printf("%+v.\n", err)
		return
	}
	time.Sleep(time.Second)

	if err = c.Login(); err != nil {
		fmt.Printf("%+v.\n", err)
		return
	}
	time.Sleep(time.Second)

	content, err := c.P.HTML()
	if err != nil {
		fmt.Printf("%+v.\n", err)
		return
	}

	titles, err := c.ScrapeBySelector(content, titleSELECTOR)
	if err != nil {
		fmt.Printf("%+v.\n", err)
		return
	}

	views, err := c.ScrapeBySelector(content, viewSELECTOR)
	if err != nil {
		fmt.Printf("%+v.\n", err)
		return
	}

	comments, err := c.ScrapeBySelector(content, commentSELECTOR)
	if err != nil {
		fmt.Printf("%+v.\n", err)
		return
	}

	sukis, err := c.ScrapeBySelector(content, sukiSELECTOR)
	if err != nil {
		fmt.Printf("%+v.\n", err)
		return
	}

	for i, t := range titles {
		v := views[i+1]
		c := comments[i+1]
		s := sukis[i+1]

		sql := fmt.Sprintf("INSERT INTO articles (title, view, comment, suki) VALUES (%s, %s, %s, %s)", t, v, c, s)
		fmt.Println(sql)
	}
}
