package crawler

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/temoto/robotstxt"
)

type robotsChecker struct {
	cache map[string]*robotstxt.RobotsData
	mu    sync.RWMutex
}

func newRobotsChecker() *robotsChecker {
	return &robotsChecker{
		cache: make(map[string]*robotstxt.RobotsData),
	}
}

func (r *robotsChecker) fetchRobotsData(scheme, host string) (*robotstxt.RobotsData, error) {
	robotsURL := fmt.Sprintf("%s://%s/robots.txt", scheme, host)
	resp, err := http.Get(robotsURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return robotstxt.FromResponse(resp)
}

func (r *robotsChecker) isAllowed(scheme, host, path, userAgent string) bool {
	r.mu.RLock()
	robotsData, exists := r.cache[host]
	r.mu.RUnlock()

	if !exists {
		var err error
		robotsData, err = r.fetchRobotsData(scheme, host)
		if err != nil {
			return true // Allow on error
		}

		r.mu.Lock()
		r.cache[host] = robotsData
		r.mu.Unlock()
	}

	return robotsData.TestAgent(path, userAgent)
}
