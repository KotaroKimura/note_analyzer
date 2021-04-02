package main

import (
	"fmt"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"github.com/KotaroKimura/note_analyzer/batch/client/chrome"
	"github.com/KotaroKimura/note_analyzer/batch/client/mysql"
	"github.com/KotaroKimura/note_analyzer/batch/mapper"
	"github.com/KotaroKimura/note_analyzer/batch/model"
)

var (
	noteURL         = os.Getenv("NOTE_URL")
	titleSELECTOR   = os.Getenv("TITLE_SELECTOR")
	viewSELECTOR    = os.Getenv("VIEW_SELECTOR")
	commentSELECTOR = os.Getenv("COMMENT_SELECTOR")
	sukiSELECTOR    = os.Getenv("SUKI_SELECTOR")
)

func main() {
	c, err := chrome.NewClient()
	if err != nil {
		fmt.Printf("Fail to Chrome NewClient: %+v.\n", err)
		return
	}
	defer c.Close()

	conn, err := mysql.NewConn()
	if err != nil {
		fmt.Printf("Fail to MySQL NewConn: %+v.\n", err)
		return
	}
	defer conn.Close()

	if err = c.P.Navigate(noteURL); err != nil {
		fmt.Printf("Fail to Navigate: %+v.\n", err)
		return
	}
	time.Sleep(time.Second)

	if err = c.Login(); err != nil {
		fmt.Printf("Fail to Login: %+v.\n", err)
		return
	}
	time.Sleep(time.Second)

	content, err := c.P.HTML()
	if err != nil {
		fmt.Printf("Fail to get content: %+v.\n", err)
		return
	}

	titles, err := c.ScrapeBySelector(content, titleSELECTOR)
	if err != nil {
		fmt.Printf("Fail to get titles: %+v.\n", err)
		return
	}

	views, err := c.ScrapeBySelector(content, viewSELECTOR)
	if err != nil {
		fmt.Printf("Fail to get views: %+v.\n", err)
		return
	}

	comments, err := c.ScrapeBySelector(content, commentSELECTOR)
	if err != nil {
		fmt.Printf("Fail to get comments: %+v.\n", err)
		return
	}

	sukis, err := c.ScrapeBySelector(content, sukiSELECTOR)
	if err != nil {
		fmt.Printf("Fail to get sukis: %+v.\n", err)
		return
	}

	aMapeer := mapper.NewArticleMapper(conn)
	for i, t := range titles {
		articles, err := aMapeer.FindByTitle(t)
		if err != nil {
			fmt.Printf("Fail to ArticleMapper.FindByTitle: %+v.\n", err)
			return
		}

		if len(articles) == 0 {
			now := time.Now()
			article := &model.Article{Title: t, CreatedAt: now, UpdatedAt: now}

			if err = aMapeer.Insert(article); err != nil {
				fmt.Printf("Fail to ArticleMapper.Insert: %+v.\n", err)
				return
			}
		}
		v := views[i+1]
		c := comments[i+1]
		s := sukis[i+1]

		sql := fmt.Sprintf("INSERT INTO articles (title, view, comment, suki) VALUES (%s, %s, %s, %s)", t, v, c, s)
		fmt.Println(sql)
	}
}
