package rest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/Soumil07/gocord"
)

var (
	idRegex = regexp.MustCompile("[0-9]+")
)

func getBody(r io.Reader) []byte {
	buf := new(bytes.Buffer)
	buf.ReadFrom(r)
	return buf.Bytes()
}

type RestManager struct {
	sync.Mutex
	Token       string
	Shard       *gocord.Shard
	GlobalReset time.Time
	buckets     *sync.Map
}

func NewRestManager(token string) *RestManager {
	return &RestManager{
		Token:       token,
		GlobalReset: time.Time{},
		buckets:     &sync.Map{},
	}
}

func (r *RestManager) GloballyRateLimited() bool {
	return time.Now().Before(r.GlobalReset)
}

func (r *RestManager) GetBucket(route string) *Bucket {
	r.Lock()
	defer r.Unlock()

	if bucket, ok := r.buckets.Load(route); ok {
		return bucket.(*Bucket)
	}

	bucket := NewBucket(r, route)
	r.buckets.Store(route, bucket)

	return bucket
}

func (r *RestManager) Do(method string, path string, body []byte, respBody interface{}) error {
	r.Lock()
	defer r.Unlock()
	route := ParseRoute(method, path)
	bucket := r.GetBucket(route)

	resp, err := bucket.Request(method, path, body)
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
