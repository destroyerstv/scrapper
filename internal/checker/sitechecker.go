package sitechecker

import (
	"context"
	"fmt"
	"io"
	"net/http"
	datastorage "scrapper/internal/storage"
	"time"
)

const (
	SCHEME = "http://"
)

type ISiteChecker interface {
	Run()
}

type SiteChecker struct {
	ctx     context.Context
	client  http.Client
	storage datastorage.IDataStorage
}

func (sc *SiteChecker) Run() {
	for {
		select {
		case <-sc.ctx.Done():
			return
		default:
			sc.CheckStorage(sc.storage.GetSiteNameList())
			time.Sleep(1 * time.Minute)
		}
	}
}

func (sc *SiteChecker) CheckStorage(siteNameList []string) {
	for _, siteName := range siteNameList {
		select {
		case <-sc.ctx.Done():
			return
		default:
			statusCode := 0

			timeVal := sc.Stopwatch(func() {
				statusCode = sc.SendRequest(siteName)
			})

			if statusCode == http.StatusOK {
				sc.storage.SetAvailable(siteName, timeVal)
			} else {
				sc.storage.SetNotAvailable(siteName)
			}
		}
	}
}

func (sc *SiteChecker) Stopwatch(f func()) time.Duration {
	timeStart := time.Now()
	f()
	return time.Duration(time.Since(timeStart).Milliseconds())
}

func (sc *SiteChecker) SendRequest(siteName string) int {
	resp, err := sc.client.Get(SCHEME + siteName)
	if err != nil {
		return 0
	}

	defer resp.Body.Close()
	_, err = io.Copy(io.Discard, resp.Body)
	if err != nil {
		fmt.Println("error body")
		return 0
	}

	return resp.StatusCode
}

func NewSiteChecker(ctx context.Context, client http.Client, storage datastorage.IDataStorage) *SiteChecker {
	return &SiteChecker{ctx, client, storage}
}
