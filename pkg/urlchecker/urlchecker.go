package urlchecker

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"sync"
	"time"

	"golang.org/x/sync/errgroup"
)

type URLChecker struct {
	client *http.Client
}

func NewURLChecker() *URLChecker {
	client := http.Client{
		Timeout: 5 * time.Second,
	}
	return &URLChecker{
		client: &client,
	}
}

type PingRequest struct {
	List []string `json:"list"`
}

type PingResponse struct {
	mu   sync.Mutex
	Data map[string]string `json:"data"`
}

func (p *PingResponse) Set(target, value string) {
	p.mu.Lock()
	p.Data[target] = value
	p.mu.Unlock()
}

func (u *URLChecker) ListenAndServe() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/ping-urls", u.PingHandler)
	mux.HandleFunc("/hello", u.HelloHandler)
	return http.ListenAndServe(":3000", mux)
}

func (u *URLChecker) HelloHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = fmt.Fprintf(w, "Hello there.")
}

func (u *URLChecker) PingHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		break
	default:
		BadRequestResponse(w, ErrorResponse{Message: "Wrong http method"})
		return
	}
	var requestBody PingRequest

	err := DecodeRequestBody(r, &requestBody)
	if err != nil {
		BadRequestResponse(w, ErrorResponse{Message: fmt.Sprintf("Error json decode: %v", err)})
		return
	}
	urls := []*url.URL{}
	for _, urlString := range requestBody.List {
		validUrl, err := url.ParseRequestURI(urlString)
		if err != nil {
			BadRequestResponse(w, ErrorResponse{Message: fmt.Sprintf("Invalid url [%s]: %v", urlString, err)})
			return
		}
		urls = append(urls, validUrl)
	}
	result := PingResponse{
		Data: make(map[string]string),
	}
	eg := new(errgroup.Group)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	callback := func(urlStr string) func() error {
		return func() error {
			response, err := u.Ping(&result, urlStr, ctx, cancel)
			if err != nil {
				log.Printf("error creating http request: %s\n", err)
				return err
			}
			result.Set(urlStr, response)
			return nil
		}
	}
	for _, target := range urls {
		eg.Go(callback(target.String()))
	}
	err = eg.Wait()
	if err != nil {
		BadRequestResponse(w, ErrorResponse{Message: fmt.Sprintf("Fetch all urls failed, %v", err)})
		return
	}
	SuccessResponse(w, result.Data)
	return
}

func (u *URLChecker) Ping(result *PingResponse, target string, ctx context.Context, cancelFunc context.CancelFunc) (string, error) {
	req, err := http.NewRequest("GET", target, nil)
	if err != nil {
		cancelFunc()
		return "", err
	}
	req = req.WithContext(ctx)
	res, err := u.client.Do(req)
	if err != nil {
		if os.IsTimeout(err) {
			return "inactive", nil
		}
		cancelFunc()
		return "", err
	}
	switch res.StatusCode {
	case http.StatusOK:
		return "active", nil
	default:
		return "inactive", nil
	}
}
