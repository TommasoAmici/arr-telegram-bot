package tvdb

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Data struct {
	XMLName xml.Name `xml:"Data"`
	Series  struct {
		IMDBID   string `xml:"IMDB_ID"`
		ID       int64  `xml:"id"`
		SeriesID int64  `xml:"seriesid"`
	} `xml:"Series"`
}

func LookupFromIMDBId(imdbID string) (*Data, error) {
	url := "https://thetvdb.com/api/GetSeriesByRemoteID.php?imdbid=" + imdbID
	client := http.Client{
		Timeout: 10 * time.Second,
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/xml")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 400 {
		return nil, fmt.Errorf("HTTP request failed with status: %s", resp.Status)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var data Data
	err = xml.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}
	if data.Series.IMDBID != imdbID {
		return nil, fmt.Errorf("imdb ID doesn't match request")
	}
	return &data, nil
}
