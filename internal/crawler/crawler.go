package crawler

import (
	"bytes"
	"io"
	"net/http"

	"github.com/pkg/errors"
	"golang.org/x/net/html"
)

type Crawler struct {
	client *http.Client
}

func New(client *http.Client) *Crawler {
	return &Crawler{
		client: client,
	}
}

func (c *Crawler) Fetch(url string) (string, []string, error) {
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)")

	resp, err := c.client.Do(req)
	if err != nil {
		return "", nil, errors.Wrap(err, "http request")
	}
	defer resp.Body.Close()

	b := bytes.NewBuffer(nil)
	_, err = io.Copy(b, resp.Body)
	if err != nil {
		return "", nil, errors.Wrap(err, "copying response body")
	}

	urls, err := ParseURLs(b)
	if err != nil {
		return "", nil, err
	}

	return b.String(), urls, nil
}

func ParseURLs(body io.Reader) ([]string, error) {
	doc, err := html.Parse(body)
	if err != nil {
		return nil, errors.Wrap(err, "parsing html")
	}

	var (
		urls []string
		f    func(*html.Node)
	)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, a := range n.Attr {
				if a.Key == "href" {
					urls = append(urls, a.Val)
					break
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)

	return urls, nil
}
