package service

import (
	"encoding/json"
	"net/http"
	datastorage "scrapper/internal/storage"
)

type IService interface {
	Run(port string) error
}

type ServiceStat struct {
	Single    int
	Min       int
	Max       int
	Statistic int
}

type Service struct {
	storage   datastorage.IDataStorage
	statistic ServiceStat
}

func (s *Service) Run(port string) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/", s.index)
	mux.HandleFunc("/min", s.getMin)
	mux.HandleFunc("/max", s.getMax)

	mux.HandleFunc("/admin", s.getStatistics)

	return http.ListenAndServe(":"+port, mux)
}

func (s *Service) index(w http.ResponseWriter, r *http.Request) {
	all_sites := s.storage.GetSiteNameList()
	// Return all info
	if !r.URL.Query().Has("site") {
		s.newResponse(w, all_sites)
		return
	}

	siteName := r.URL.Query().Get("site")
	siteStat, err := s.storage.GetSite(siteName)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		s.newResponse(w, map[string]string{"error": "site not found"})
		return
	}

	s.statistic.Single++
	s.newResponse(w, siteStat)
}

func (s *Service) getMin(w http.ResponseWriter, r *http.Request) {
	minStat := s.storage.GetMin()

	s.statistic.Min++
	s.newResponse(w, minStat)
}

func (s *Service) getMax(w http.ResponseWriter, r *http.Request) {
	maxStat := s.storage.GetMax()

	s.statistic.Max++
	s.newResponse(w, maxStat)
}

func (s *Service) getStatistics(w http.ResponseWriter, r *http.Request) {
	s.statistic.Statistic++
	s.newResponse(w, s.statistic)
}

func (s *Service) newResponse(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func NewService(storage datastorage.IDataStorage) IService {
	return &Service{storage, ServiceStat{0, 0, 0, 0}}
}
