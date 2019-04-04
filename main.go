package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi"
	log "github.com/sirupsen/logrus"

	"github.com/spf13/viper"
)

// Config hold the cms api key and other information necessary for the application
type Config struct {
	CMSUrl string
	CMSKey string
	Port   int
	APIKey string
}

// CMSResponse Holds information from a cms response
type CMSResponse struct {
	Error     bool
	Message   string
	UserType  string `json:"user_type"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

// CacheObject represents a cache entry for a cms response
type CacheObject struct {
	RCSID      string
	Expiration time.Time
	CMSValue   CMSResponse
}

// API holds everything needed for the api
type API struct {
	cache  []CacheObject
	Config Config
}

// RCSPostBody is the structure of an rcs id from the frontend
type RCSPostBody struct {
	RCSID string `json:"rcs_id"`
}

// ValidRCSChecker checks if the rcs id is valid
func (a *API) ValidRCSChecker(w http.ResponseWriter, r *http.Request) {
	key := r.Header.Get("Authorization")
	if key != "Token "+a.Config.APIKey {
		w.WriteHeader(403)
		w.Write([]byte("invalid key"))
		log.Info("User tried accessing without api key")
		return
	}

	var postData RCSPostBody
	err := json.NewDecoder(r.Body).Decode(&postData)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
		return
	}
	if postData.RCSID == "" {
		w.WriteHeader(400)
		return
	}
	// Check cache
	for idx, entry := range a.cache {
		if entry.RCSID == postData.RCSID {
			if time.Now().After(entry.Expiration) {
				log.Info("Cache expired for ", entry.RCSID)
				a.cache = append(a.cache[:idx], a.cache[idx+1:]...)
				break
			}
			log.Info("Using cached value")
			WriteJSON(w, entry.CMSValue)
			return
		}
	}

	// wasn't in cache, find it and add it
	req, err := http.NewRequest("GET", "https://cms.union.rpi.edu/api/users/view_rcs/"+postData.RCSID+"/", nil)
	if err != nil {
		log.Fatalf("Failed to create http request")
	}
	log.Info("https://cms.union.rpi.edu/api/users/view_rcs/" + postData.RCSID + "/")
	req.Header.Set("Authorization", "Token "+a.Config.CMSKey)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)

	if err != nil {
		log.Fatal("Couldn't reach cms")
	}

	defer resp.Body.Close()

	var cmsResp CMSResponse

	err = json.NewDecoder(resp.Body).Decode(&cmsResp)
	if err != nil {
		log.Printf("Empty CMS response, is %s a valid RCS ID?\n", "lyonj4")
		cmsResp.Error = true // This is likely not perfect, but if we don't get a json response from cms it is pretty safe to assume it is an invalid rcs
		cmsResp.Message = "Invalid RCS ID"
	}

	c := CacheObject{
		RCSID:      postData.RCSID,
		Expiration: time.Now().Add(48 * time.Hour),
		CMSValue:   cmsResp,
	}
	a.cache = append(a.cache, c)

	WriteJSON(w, c.CMSValue)
	return
}

// WriteJSON writes the data as JSON.
func WriteJSON(w http.ResponseWriter, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	b, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}
	w.Write(b)
	return nil
}

func main() {
	v := viper.New()
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.SetConfigName("config")
	v.AddConfigPath(".")
	v.BindEnv("PORT", "Port")
	v.SetDefault("port", 8080)
	v.BindEnv("CMSURL")
	v.BindEnv("CMSKEY")
	v.BindEnv("APIKEY")

	var config Config
	if err := v.ReadInConfig(); err != nil {
		log.Info("Config file not found, reading from environment")
	}

	err := v.Unmarshal(&config)

	if err != nil {
		log.Error("Failed to decode config")
		return
	}

	if config.CMSKey == "" || config.CMSUrl == "" {
		log.Error("No CMS info provided")
		return
	}

	api := API{
		Config: config,
	}

	r := chi.NewRouter()
	r.Post("/valid", api.ValidRCSChecker)
	log.Info("Serving at ", "0.0.0.0:"+strconv.Itoa(config.Port))

	if err := http.ListenAndServe("0.0.0.0:"+strconv.Itoa(config.Port), r); err != nil {
		log.Error(err.Error())
	}

}
