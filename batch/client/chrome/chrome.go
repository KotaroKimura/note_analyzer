package chrome

import (
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/sclevine/agouti"
)

var (
	noteID       = os.Getenv("NOTE_ID")
	notePassword = os.Getenv("NOTE_PASSWORD")
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

func (c *Client) Close() {
	c.P.Destroy()
	c.D.Stop()
}

func (c *Client) Login() error {
	if err := c.P.FirstByName("login").Fill(noteID); err != nil {
		return err
	}
	if err := c.P.FirstByName("password").Fill(notePassword); err != nil {
		return err
	}
	if err := c.P.FirstByName("$ctrl.$scope.login_form").Submit(); err != nil {
		return err
	}

	return nil
}

func (c *Client) ScrapeBySelector(content string, selector string) ([]string, error) {
	var result []string

	reader := strings.NewReader(content)
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return result, err
	}
	doc.Find(selector).Each(func(i int, s *goquery.Selection) {
		result = append(result, s.Text())
	})

	return result, nil
}
