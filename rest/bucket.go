package rest

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"sync"
	"time"
)

const (
	API_URL = "https://discordapp.com/api/v6/"
)

type Bucket struct {
	sync.Mutex
	Manager   *RestManager
	Route     string
	Remaining int64
	Limit     int64

	queue      []*http.Request
	busy       bool
	resetTime  time.Time
	httpClient *http.Client
}

type ratelimitedResponse struct {
	Message    string `json:"message"`
	RetryAfter int    `json:"retry_after"`
	Global     bool   `json:"global"`
}

func NewBucket(r *RestManager, route string) *Bucket {
	bucket := &Bucket{
		Mutex:     sync.Mutex{},
		Manager:   r,
		Route:     route,
		Remaining: 1,
		Limit:     1,

		busy:       false,
		resetTime:  time.Time{},
		httpClient: &http.Client{},
	}

	return bucket
}

// Request creates an http request
func (b *Bucket) Request(method string, path string, body []byte) (*http.Response, error) {
	if b.Manager.GloballyRateLimited() {
		<-time.After(time.Until(b.Manager.GlobalReset))
	}

	if b.Remaining < 1 {
		<-time.After(time.Until(b.resetTime))
	}

	b.Lock()
	defer b.Unlock()

	next, _ := http.NewRequest(method, API_URL+path, bytes.NewBuffer(body))
	next.Header.Set("Authorization", "Bot "+b.Manager.Token)
	next.Header.Set("User-Agent", "DiscordBot (https://github.com/Soumil07/gocord, v1)")
	resp, err := b.httpClient.Do(next)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	err = b.UpdateHeaders(resp, path)
	if err != nil {
		return resp, err
	}
	return resp, nil
}

func (b *Bucket) UpdateHeaders(resp *http.Response, path string) error {
	remaining := resp.Header.Get("X-Ratelimit-Remaining")
	limit := resp.Header.Get("X-Ratelimit-Limit")

	if remaining != "" {
		b.Remaining, _ = strconv.ParseInt(remaining, 10, 32)
	}

	if limit != "" {
		b.Limit, _ = strconv.ParseInt(remaining, 10, 32)
	}

	switch {
	// handle ratelimits
	case resp.StatusCode == http.StatusTooManyRequests:
		// TODO:
		// b.Manager.Shard.Cluster.Dispatch("debug", "Ratelimit")
		bytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		var resp ratelimitedResponse
		json.Unmarshal(bytes, &resp)

		reset := time.Now().Add(time.Duration(resp.RetryAfter))
		if resp.Global {
			b.Manager.GlobalReset = reset
		} else {
			b.resetTime = reset
		}

		return nil

	case resp.StatusCode >= 500 && resp.StatusCode <= 600:
		// handle 5xx errors
		<-time.After(5 * time.Second)
		buf := new(bytes.Buffer)
		buf.ReadFrom(resp.Request.Body)
		b.Request(resp.Request.Method, path, buf.Bytes())

	default:
		reset := resp.Header.Get("X-RateLimit-Reset")
		if reset == "" {
			return nil
		}

		resetTime, err := strconv.ParseInt(reset, 10, 32)
		if err != nil {
			return err
		}

		timeSent, err := time.Parse(resp.Header.Get("Date"), time.RFC1123)
		if err != nil {
			timeSent = time.Now()
		}
		b.resetTime = time.Unix(resetTime, 0).Add(time.Now().Sub(timeSent))
		return nil
	}

	return nil
}
