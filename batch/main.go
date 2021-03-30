package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/sclevine/agouti"
)

var (
	noteID       = os.Getenv("NOTE_ID")
	notePassword = os.Getenv("NOTE_PASSWORD")
	noteURL      = os.Getenv("NOTE_URL")
)

type Client struct {
	D *agouti.WebDriver
	P *agouti.Page
}

func NewClient() (*Client, error) {
	d := agouti.ChromeDriver(
		agouti.ChromeOptions(
			"args", []string{
				"--headless",
				"--disable-gpu",
				"--no-sandbox",
				"--window-size=1280,800",
			},
		),
	)
	err := d.Start()
	if err != nil {
		return nil, err
	}

	p, err := d.NewPage()
	if err != nil {
		return nil, err
	}

	return &Client{
		D: d,
		P: p,
	}, nil
}

func (p *Client) Close() {
	p.P.Destroy()
	p.D.Stop()
}

func Login(p *agouti.Page) error {
	if err := p.FirstByName("login").Fill(noteID); err != nil {
		return err
	}
	if err := p.FirstByName("password").Fill(notePassword); err != nil {
		return err
	}
	if err := p.FirstByName("$ctrl.$scope.login_form").Submit(); err != nil {
		return err
	}

	return nil
}

func main() {
	c, err := NewClient()
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

	if err = Login(c.P); err != nil {
		fmt.Printf("%+v.\n", err)
		return
	}
	time.Sleep(time.Second)

	content, err := c.P.HTML()
	if err != nil {
		fmt.Printf("%+v.\n", err)
		return
	}

	reader := strings.NewReader(content)
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		fmt.Printf("%+v.\n", err)
		return
	}

	doc.Find(".o-statsContent__tableLabel").Each(func(i int, s *goquery.Selection) {
		article := s.Text()
		fmt.Printf("Review %d: %s\n", i, article)
	})
}
