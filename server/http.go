package server

import (
	"encoding/json"
	"log"
	"net/http"
	"recycling-service/pkg/models"
	"recycling-service/pkg/models/sqlserver"

	"github.com/gorilla/mux"
)

type httpServer struct {
	Guidelines map[string]*[]models.MaterialGuidelineResults
}

// Used to received requests for community recycling info
type RequestInfo struct {
	CommunityID string `json:"communityID"`
}

type SQLResult struct {
	MID          int32  `json:"mID"`
	CommunityID  string `json:"communityID"`
	Category     string `json:"category"`
	YesNo        string `json:"yesNo"`
	CategoryType string `json:"categoryType"`
	Material     string `json:"material"`
}

type GuidelinesResponse struct {
	Guidelines []SQLResult `json:"guidelines"`
}

func (s *httpServer) communityGuidelines(w http.ResponseWriter, r *http.Request) {
	log.Println("HIT")

	var reqInfo RequestInfo
	err := json.NewDecoder(r.Body).Decode(&reqInfo)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Println(reqInfo)
	guidelines, _ := s.Guidelines[reqInfo.CommunityID]
	if len(*guidelines) == 0 {
		http.Error(w, "No results found for your location", http.StatusNotFound)
	}
	results := []SQLResult{}

	for _, g := range *guidelines {
		results = append(
			results,
			SQLResult{
				g.MID, g.CommunityID, g.Category, g.YesNo, g.CategoryType, g.Material,
			},
		)
	}
	guidelinesResponse := GuidelinesResponse{Guidelines: results}
	// body, err := json.Marshal(&guidelinesResponse)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(&guidelinesResponse)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// Returns an *http.Server that listens at location `addr` with routes registered.
func NewHTTPServer(addr string) *http.Server {

	db, err := sqlserver.InitDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	mgr := sqlserver.MGRModel{DB: db}

	communityGuidelines, err := mgr.GetAll()
	if err != nil {
		log.Fatal(err)
	}
	httpServer := &httpServer{Guidelines: communityGuidelines}

	r := mux.NewRouter()

	r.HandleFunc("/", httpServer.communityGuidelines).Methods("POST")
	r.HandleFunc("/ping", httpServer.Ping).Methods("GET")

	return &http.Server{
		Addr:    addr,
		Handler: r,
	}
}

func (s *httpServer) Ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("It worked"))
}
