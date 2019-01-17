package rest

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

const (
	API_URL = "https://discordapp.com/api/v6/"
)

var (
	quoteEscaper = strings.NewReplacer("\\", "\\\\", `"`, "\\\"")
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

type File struct {
	Name        string
	Reader      io.Reader
	ContentType string
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
func (b *Bucket) Request(method string, path string, body []byte, files ...File) (*http.Response, error) {
	if b.Manager.GloballyRateLimited() {
		<-time.After(time.Until(time.Unix(0, atomic.LoadInt64(b.Manager.global))))
	}

	if b.Remaining < 1 {
		<-time.After(time.Until(b.resetTime))
	}

	b.Lock()
	defer b.Unlock()

	var req *http.Request

	if len(files) > 0 {
		buf := &bytes.Buffer{}
		bodywriter := multipart.NewWriter(buf)

		var p io.Writer

		h := make(textproto.MIMEHeader)
		h.Set("Content-Disposition", `form-data; name="payload_json"`)
		h.Set("Content-Type", "application/json")

		p, err := bodywriter.CreatePart(h)
		if err != nil {
			panic(err)
		}

		if _, err = p.Write(body); err != nil {
			panic(err)
		}

		for i, file := range files {
			h := make(textproto.MIMEHeader)
			h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="file%d"; filename="%s"`, i, quoteEscaper.Replace(file.Name)))
			contentType := file.ContentType
			if contentType == "" {
				contentType = "application/octet-stream"
			}
			h.Set("Content-Type", contentType)

			p, err = bodywriter.CreatePart(h)
			if err != nil {
				panic(err)
			}

			if _, err = io.Copy(p, file.Reader); err != nil {
				fmt.Printf("%#v", file.Reader)
				panic(err)
			}
		}

		err = bodywriter.Close()
		if err != nil {
			panic(err)
		}

		req, err = http.NewRequest(method, API_URL+path, bytes.NewBuffer(buf.Bytes()))
		if err != nil {
			panic(err)
		}

		req.Header.Set("Content-Type", bodywriter.FormDataContentType())
	} else {
		var err error
		req, err = http.NewRequest(method, API_URL+path, bytes.NewBuffer(body))
		if err != nil {
			return nil, err
		}

		req.Header.Set("Content-Type", "application/json")
	}

	req.Header.Set("Authorization", "Bot "+b.Manager.Token)
	req.Header.Set("User-Agent", "DiscordBot (https://github.com/Soumil07/gocord, v1)")

	resp, err := b.httpClient.Do(req)
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
