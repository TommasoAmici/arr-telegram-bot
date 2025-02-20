package imdb

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var ErrNotIMDBUrl = errors.New("url is not from imdb")
var ErrNoIMDBId = errors.New("imdb url has no id")

func idFromURL(rawURL string) string {
	split := strings.Split(rawURL, "/")
	for _, part := range split {
		if strings.HasPrefix(part, "tt") {
			return part
		}
	}
	return ""
}

func isMovieFromURL(u *url.URL) (bool, error) {
	client := http.Client{
		Timeout: 10 * time.Second,
	}
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return false, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:135.0) Gecko/20100101 Firefox/135.0")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	req.Header.Set("Alt-Used", "www.imdb.com")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("Sec-GPC", "1")
	req.Header.Set("Priority", "u=0, i")
	req.Header.Set("referrer", "https://www.imdb.com/")

	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 400 {
		return false, fmt.Errorf("HTTP request failed with status: %s", resp.Status)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}
	bodyStr := string(body)
	isMovie := !strings.Contains(bodyStr, "episode-guide-text")
	return isMovie, nil
}

type IMDBEntity struct {
	ID      string
	IsMovie bool
}

func LookupIMDB(rawURL string) (*IMDBEntity, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return nil, err
	}
	if u.Hostname() != "imdb.com" && u.Hostname() != "www.imdb.com" {
		return nil, ErrNotIMDBUrl
	}
	id := idFromURL(rawURL)
	if id == "" {
		return nil, ErrNoIMDBId
	}
	isMovie, err := isMovieFromURL(u)
	if err != nil {
		return nil, err
	}
	return &IMDBEntity{ID: id, IsMovie: isMovie}, nil
}
