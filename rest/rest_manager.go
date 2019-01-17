package rest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

var (
	idRegex = regexp.MustCompile("[0-9]+")
)

func getBody(r io.Reader) []byte {
	body, err := ioutil.ReadAll(r)
	if err != nil {
		panic(err)
	}

	return body
}

type RestManager struct {
	Token string

	global  *int64
	buckets *sync.Map
}

func NewRestManager(token string) *RestManager {
	return &RestManager{
		Token:   token,
		global:  new(int64),
		buckets: &sync.Map{},
	}
}

func (r *RestManager) GloballyRateLimited() bool {
	globalTime := time.Unix(0, atomic.LoadInt64(r.global))
	return time.Now().Before(globalTime)
}

func (r *RestManager) GetBucket(route string) *Bucket {
	if bucket, ok := r.buckets.Load(route); ok {
		return bucket.(*Bucket)
	}

	bucket := NewBucket(r, route)
	r.buckets.Store(route, bucket)

	return bucket
}

func (r *RestManager) Do(method string, path string, body []byte, respBody interface{}, files ...File) error {
	route := ParseRoute(method, path)
	bucket := r.GetBucket(route)

	resp, err := bucket.Request(method, path, body, files...)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if err != nil {
		// sometimes the error is while updating headers, so there is a body
		if resp != nil {
			body := getBody(resp.Body)
			err := json.Unmarshal(body, respBody)
			if err != nil {
				return fmt.Errorf("error while unmarshalling response body: %s", err)
			}

			return err
		}

		return err
	}

	body = getBody(resp.Body)
	err = json.Unmarshal(body, respBody)
	if err != nil {
		return fmt.Errorf("error while unmarshalling response body: %s", err)
	}

	return nil
}

// SimpleRequest creates a simple JSON request to the supplied URL
func SimpleRequest(method string, url string, body []byte, respBody interface{}) error {
	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("error intializing request: %s", err.Error())
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	err = json.Unmarshal(getBody(res.Body), respBody)
	if err != nil {
		return fmt.Errorf("error unmarshalling json: %s", err.Error())
	}

	return nil
}

// ParseRoute parses a route to be used in a bucket
// adapted from the JavaScript version made by PoLLeN#5796
func ParseRoute(method string, route string) string {
	url := strings.Split(route, "?")[0] // query strings don't count
	ids := idRegex.FindAllString(url, -1)
	// if one ID, return the url
	if len(ids) == 1 {
		return url
	}
	url = strings.Replace(url, ids[1], ":id", 1)
	if method == http.MethodDelete && strings.Contains(url, "messages") {
		// message deletes have their own ratelimits
		return fmt.Sprintf("%s %s", http.MethodDelete, url)
	}

	// adding reactions have their own ratelimit across the account
	if strings.Contains(url, "/reactions/:id") {
		return "/channels/messages/:id/reactions"
	}
	return url
}
