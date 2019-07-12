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

// Devices contains devices information
type Devices struct {
	Dev []struct {
		Tag    string `json:"tag"`
		Clicks string `json:"clicks"`
	} `json:"dev"`
	Sys []struct {
		Tag    string `json:"tag"`
		Clicks string `json:"clicks"`
	} `json:"sys"`
	Bro []struct {
		Tag    string `json:"tag"`
		Clicks string `json:"clicks"`
	} `json:"bro"`
}

// Refs contains references information
type Refs []struct {
	Link   string `json:"link"`
	Clicks string `json:"clicks"`
}

// Stats contains all the information about a Stats query
// Devices and Refs are declared as interface because they are served
// both as object or array, depending on url conditions
// i.e A just created url will come packed with empty Devices and Refs,
// that are passed as empty array (for no logical reason, but that's a server side issue)
// An old and already clicked url will come packed with populated Devices and Refs
// that are passed as an object (which would be the normal thing to do)
type Stats struct {
	Status     int         `json:"status"`
	Clicks     string      `json:"clicks"`
	Date       string      `json:"date"`
	Title      string      `json:"title"`
	FullLink   string      `json:"fullLink"`
	ShortLink  string      `json:"shortLink"`
	Facebook   int         `json:"facebook"`
	Twitter    int         `json:"twitter"`
	Pinterest  int         `json:"pinterest"`
	Instagram  int         `json:"instagram"`
	GooglePlus int         `json:"googlePlus"`
	Linkedin   int         `json:"linkedin"`
	Rest       int         `json:"rest"`
	Devices    interface{} `json:"devices"`
	Refs       interface{} `json:"refs"`
}

// Response condense both URL and Stats information
// in order to use just one type in get
type Response struct {
	URL   `json:"url"`
	Stats `json:"stats"`
}

// New returns a new client
func New(APIKey string) (*Client, error) {
	endpoint, _ := url.Parse("https://cutt.ly/api/api.php")
	return &Client{
		Key:     APIKey,
		BaseURL: endpoint}, nil
}

func (c *Client) get(params url.Values) (Response, error) {
	// Add the parameters to the base url
	c.BaseURL.RawQuery = params.Encode()
	// Make the request to the server
	r, err := http.Get(c.BaseURL.String())
	if err != nil {
		return Response{}, err
	}
	defer r.Body.Close()
	// Read the request
	contents, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return Response{}, err
	}
	// Unmarshal the request into Response struct
	var response Response
	err = json.Unmarshal(contents, &response)
	if err != nil {
		return Response{}, err
	}

	return response, nil
}

// Shorten shortens a given URL
func (c *Client) Shorten(longURL string, customName string) (URL, error) {
	// Creates a new set of url parameters
	v := url.Values{}
	// We start with the key (we use Set because it's the first)
	// then we add the url we want to shorten and the custon name
	v.Set("key", c.Key)
	v.Add("short", longURL)
	v.Add("name", customName)

	// Use the internal modular request system in order to get
	// the url response
	r, err := c.get(v)
	if err != nil {
		return URL{}, fmt.Errorf("Impossible to shorten the url: %s", err.Error())
	}

	// Check if the returned status is an actual error (isURL == true)
	err = checkErrorCode(r.URL.Status, true)
	// Just return the URL information
	return r.URL, err
}

// GetStats get the statistics for a given shortened URL
func (c *Client) GetStats(shortURL string) (Stats, error) {
	// Creates a new set of url parameters
	v := url.Values{}
	// We start with the key (we use Set because it's the first)
	// then we add the already shortened url as stats parameter
	v.Set("key", c.Key)
	v.Add("stats", shortURL)

	// Use the internal modular request system in order to get
	// the url response
	r, err := c.get(v)
	if err != nil {
		return Stats{}, fmt.Errorf("Impossible to get the stats: %s", err.Error())
	}

	// Check if the returned status is an actual error (isURL == false)
	err = checkErrorCode(r.Stats.Status, false)

	// IMPORTANT: this is an ugly solution that I must use because the API calls
	// return an empty array if there's no devices and refs data and an object
	// otherwise, if they'll change their API in the future I'll modify this part

	// If devices field is not passed as an empty interface (empty array in this case)
	if _, ok := r.Stats.Devices.([]interface{}); !ok {
		// Overwrite the old Stats.Devices (made out of map[string]interface{})
		// with the right data (of Devices type)
		r.Stats.Devices, err = ForceDevicesToRightType(r.Stats.Devices)
		if err != nil {
			return Stats{}, fmt.Errorf("Impossible to get the stats: %s", err.Error())
		}
	}

	// If devices field is not passed as an empty interface (empty array in this case)
	if _, ok := r.Stats.Refs.([]interface{}); !ok {
		// Overwrite the old Stats.Refs (made out of map[string]interface{})
		// with the right data (of Refs type)
		r.Stats.Refs, err = ForceRefsToRightType(r.Stats.Refs)
		if err != nil {
			return Stats{}, fmt.Errorf("Impossible to get the stats: %s", err.Error())
		}
	}

	// Just return the Stats information
	return r.Stats, err
}
