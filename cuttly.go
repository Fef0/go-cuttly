package cuttly

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

// Client contains the API Key and the API Base url
type Client struct {
	Key     string
	BaseURL *url.URL
}

// URL contains information about a URL query
type URL struct {
	Status    int    `json:"status"`
	FullLink  string `json:"fullLink"`
	Date      string `json:"date"`
	ShortLink string `json:"shortLink"`
	Title     string `json:"title"`
}

// Response contains the URL information
type Response struct {
	URL `json:"url"`
}

// New returns a new client
func New(APIKey string) (*Client, error) {
	endpoint, _ := url.Parse("https://cutt.ly/api/api.php")
	return &Client{
		Key:     APIKey,
		BaseURL: endpoint}, nil
}

func (c *Client) get(params url.Values) (Response, error) {
	c.BaseURL.RawQuery = params.Encode()
	r, err := http.Get(c.BaseURL.String())
	if err != nil {
		fmt.Println("Errore requester", err.Error())
		return Response{}, err
	}
	defer r.Body.Close()
	contents, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println("Errore reader", err.Error())
		return Response{}, err
	}
	fmt.Println(string(contents))
	var response Response
	err = json.Unmarshal(contents, &response)
	if err != nil {
		fmt.Println("Errore json", err.Error())
		return Response{}, err
	}

	return response, nil
}

// Shorten shortens a given URL
func (c *Client) Shorten(longURL string, customName string) (URL, error) {
	// Creates a new set of url parameters
	v := url.Values{}
	// The first will be the key, we use Set because it's the first
	// then we'll add the url we want to shorten and the custon name
	v.Set("key", c.Key)
	v.Add("short", longURL)
	v.Add("name", customName)
	// Use the internal modular request system in order to get
	// the url response
	r, err := c.get(v)
	if err != nil {
		return URL{}, fmt.Errorf("Impossible to shorten the url: %s", err.Error())
	}
	// Check if the returned status is an actual error
	err = checkErrorCode(r.URL.Status, true)
	// Just return the URL information
	return r.URL, err
}
