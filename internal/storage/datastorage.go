package datastorage

import (
	"bufio"
	"errors"
	"os"
	"strings"
	"sync"
	"time"
)

type DataStorage struct {
	sites   map[string]ISiteStat
	maxTime time.Duration
	minTime time.Duration
	lock    *sync.RWMutex
}

func NewDataStorage(filepath string) (IDataStorage, error) {
	file, err := os.Open(filepath)

	if err != nil {
		return nil, err
	}

	sites := make(map[string]ISiteStat)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		siteName := strings.TrimSpace(scanner.Text())
		sites[siteName] = NewSiteStat(siteName)
	}

	maxTime := time.Duration(0 * time.Millisecond)
	minTime := time.Duration(4294967295 * time.Millisecond)
	var lock sync.RWMutex

	return &DataStorage{sites, maxTime, minTime, &lock}, nil
}

func (ds *DataStorage) GetSiteNameList() []string {
	siteNameList := make([]string, len(ds.sites))

	i := 0
	for siteName := range ds.sites {
		siteNameList[i] = siteName
		i++
	}

	return siteNameList
}

func (ds *DataStorage) GetSite(name string) (ISiteStat, error) {
	ds.lock.RLock()
	defer ds.lock.RUnlock()

	if siteStat, ok := ds.sites[name]; ok {
		return siteStat, nil
	}

	return nil, errors.New("no site")
}

func (ds *DataStorage) GetMax() []ISiteStat {
	ds.lock.RLock()
	defer ds.lock.RUnlock()

	sitesStat := make([]ISiteStat, 0)

	for _, siteStat := range ds.sites {

		if !siteStat.IsAvailable() {
			continue
		}

		if siteStat.Eq(ds.maxTime) == 0 {
			sitesStat = append(sitesStat, siteStat)
		}
	}

	return sitesStat
}

func (ds *DataStorage) GetMin() []ISiteStat {
	ds.lock.RLock()
	defer ds.lock.RUnlock()

	sitesStat := make([]ISiteStat, 0)

	for _, siteStat := range ds.sites {

		if !siteStat.IsAvailable() {
			continue
		}

		if siteStat.Eq(ds.minTime) == 0 {
			sitesStat = append(sitesStat, siteStat)
		}
	}

	return sitesStat
}

func (ds *DataStorage) GetNotAvailable() []ISiteStat {
	ds.lock.RLock()
	defer ds.lock.RUnlock()

	sitesStat := make([]ISiteStat, 0)

	for _, siteStat := range ds.sites {

		if !siteStat.IsAvailable() {
			sitesStat = append(sitesStat, siteStat)
		}
	}

	return sitesStat
}

func (ds *DataStorage) SetAvailable(name string, t time.Duration) error {
	ds.lock.Lock()
	defer ds.lock.Unlock()

	if siteStat, ok := ds.sites[name]; ok {
		siteStat.SetAvailable(t)
		if ds.maxTime < t {
			ds.maxTime = t
		}
		if ds.minTime > t {
			ds.minTime = t
		}
		return nil
	}

	return errors.New("no site")
}

func (ds *DataStorage) SetNotAvailable(name string) error {
	ds.lock.Lock()
	defer ds.lock.Unlock()

	if siteStat, ok := ds.sites[name]; ok {
		siteStat.SetNotAvailable()
		return nil
	}

	return errors.New("no site")
}
