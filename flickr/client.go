package flickr

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

const (
	baseEndpoint    = "https://api.flickr.com/services/"
	authEndpoint    = baseEndpoint + "auth/"
	restEndpoint    = baseEndpoint + "rest/"
	uploadEndpoint  = baseEndpoint + "upload/"
	replaceEndpoint = baseEndpoint + "replace/"
)

type Client struct {
	apiKey, secret string
}

func New(apiKey, secret string) *Client {
	return &Client{apiKey, secret}
}

func (c *Client) send(url string) ([]byte, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	resp := Response{}
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("JSON error: %s", err)
	}
	if resp.Status != "ok" {
		return nil, fmt.Errorf(resp.Message)
	}
	return body, nil
}

type Response struct {
	Status  string `json:"stat"`
	Code    int
	Message string
}

func (c *Client) Do(method string, args Args) ([]byte, error) {
	if args == nil {
		args = Args{}
	}
	args["method"] = method
	args["format"] = "json"
	args["nojsoncallback"] = "1"
	url := URL(c.apiKey, c.secret, restEndpoint, args)
	return c.send(url)
}

func (c *Client) Test(method string, args Args) {
	b, err := c.Do(method, args)
	if err != nil {
		log.Fatalf("%s request failed: %s", method, err)
	}
	pretty := bytes.Buffer{}
	if err := json.Indent(&pretty, b, "", "  "); err != nil {
		log.Fatalf("Invalid JSON: %s", err)
	}
	log.Println(pretty.String())
	os.Exit(0)
}
