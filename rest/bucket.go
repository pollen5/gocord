package rest

import (
	"bytes"
	"io"
	"net/http"
	"strconv"
	"sync"
	"sync/atomic"
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

		resetTime:  time.Time{},
		httpClient: &http.Client{},
	}

	return bucket
}

// Request creates an http request
func (b *Bucket) Request(method string, path string, body []byte, files ...io.Reader) (*http.Response, error) {
	if b.Manager.GloballyRateLimited() {
		<-time.After(time.Until(time.Unix(0, atomic.LoadInt64(b.Manager.global))))
	}

	if b.Remaining < 1 {
		<-time.After(time.Until(b.resetTime))
	}

	b.Lock()
	defer b.Unlock()

	next, _ := http.NewRequest(method, API_URL+path, bytes.NewBuffer(body))
	next.Header.Set("Authorization", "Bot "+b.Manager.Token)
	next.Header.Set("User-Agent", "DiscordBot (https://github.com/Soumil07/gocord, v1)")

	next.Header.Set("Content-Type", "application/json")
	resp, err := b.httpClient.Do(next)
	if err != nil {
		return nil, err
	}

	err = b.UpdateHeaders(resp, path)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

func (b *Bucket) UpdateHeaders(resp *http.Response, path string) error {
	remaining := resp.Header.Get("X-Ratelimit-Remaining")
	reset := resp.Header.Get("X-Ratelimit-Reset")
	global := resp.Header.Get("X-RateLimit-Global")
	retryAfter := resp.Header.Get("Retry-After")

	if retryAfter != "" {
		parsed, _ := strconv.ParseInt(retryAfter, 10, 64)
		resetTime := time.Now().Add(time.Duration(parsed) * time.Millisecond)
		if global != "" {
			atomic.StoreInt64(b.Manager.global, resetTime.UnixNano())
		} else {
			b.resetTime = resetTime
		}
	} else if reset != "" {
		dTime, err := http.ParseTime(resp.Header.Get("Date"))
		if err != nil {
			return err
		}

		unixTime, err := strconv.ParseInt(reset, 10, 64)
		if err != nil {
			return err
		}

		b.resetTime = time.Now().Add(time.Unix(unixTime, 0).Sub(dTime) + time.Millisecond*250)
	}

	if remaining != "" {
		parsedRemaining, err := strconv.ParseInt(remaining, 10, 32)
		if err != nil {
			return err
		}
		b.Remaining = parsedRemaining
	}

	return nil
}
