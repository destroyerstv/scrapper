package datastorage

import (
	"fmt"
	"time"
)

type SiteStat struct {
	SiteName        string        `json:"name"`
	IsAvailableSite bool          `json:"is_available"`
	TimeAvailable   time.Duration `json:"time"`
}

func (ss *SiteStat) SetAvailable(t time.Duration) {
	ss.IsAvailableSite = true
	ss.TimeAvailable = t
}

func (ss *SiteStat) SetNotAvailable() {
	ss.IsAvailableSite = false
	ss.TimeAvailable = time.Duration(0 * time.Millisecond)
}

func (ss *SiteStat) IsAvailable() bool {
	return ss.IsAvailableSite
}

func (ss *SiteStat) GetAvailable() time.Duration {
	return ss.TimeAvailable
}

func (ss *SiteStat) Eq(t time.Duration) int {
	var ret int = 0
	switch {
	case ss.TimeAvailable < t:
		return 1
	case ss.TimeAvailable > t:
		return -1
	}

	return ret
}

func (ss *SiteStat) String() string {
	if ss.IsAvailableSite {
		return fmt.Sprintf("%s: %s", ss.SiteName, ss.TimeAvailable.String())
	}

	return fmt.Sprintf("%s: not available", ss.SiteName)
}

func NewSiteStat(name string) ISiteStat {
	return &SiteStat{name, false, time.Duration(0 * time.Millisecond)}
}
