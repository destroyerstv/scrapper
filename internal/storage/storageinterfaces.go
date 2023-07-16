package datastorage

import "time"

type ISiteStat interface {
	SetAvailable(t time.Duration)
	SetNotAvailable()
	IsAvailable() bool
	GetAvailable() time.Duration
	Eq(t time.Duration) int
}

type IDataStorage interface {
	GetSiteNameList() []string
	GetSite(name string) (ISiteStat, error)
	GetMax() []ISiteStat
	GetMin() []ISiteStat
	GetNotAvailable() []ISiteStat
	SetAvailable(name string, t time.Duration) error
	SetNotAvailable(name string) error
}
