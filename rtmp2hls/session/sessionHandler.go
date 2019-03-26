package session

import (
	"encoding/json"
	log "github.com/cihub/seelog"
	"net/http"
	"os"
	"regexp"
	"strings"
)

const my_ip = "127.168.0.1"
const regex_uuid = "\\b[0-9a-f]{8}\\b-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-\\b[0-9a-f]{12}\\b"
const regex_last_url = "[^/]*$"

/* /session GET request */
func NewSessionHandler(w http.ResponseWriter, r *http.Request) {

	if http.MethodGet == r.Method {

		resp := Sessions.NewSession(my_ip)
		js, err := json.Marshal(resp)
		if err != nil {

			log.Error("Url: " + r.RequestURI + " new session problems")
			w.WriteHeader(http.StatusInternalServerError)

			return
		}

		log.Info(js)
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)

	} else {

		log.Error("Url: " + r.RequestURI + " method: " + r.Method + " not available")
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

/* /session/uuid DELETE request */
func DeleteSessionHandler(w http.ResponseWriter, r *http.Request) {

	if http.MethodDelete == r.Method {

		re := regexp.MustCompile(regex_uuid)
		uuid := re.FindString(r.RequestURI)


		if len(uuid) == 0 {

			log.Error("Url: " + r.RequestURI + " method: " + r.Method + " bad request")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		log.Info(uuid)
		Sessions.DeleteSession(uuid)
		w.WriteHeader(http.StatusOK)

	} else {

		log.Error("Url: " + r.RequestURI + " method: " + r.Method + " not available")
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

/* get media over http */
func HLSSessionHandler(w http.ResponseWriter, r *http.Request) {

	if http.MethodGet == r.Method {

		directory := determineDirectoryByURL(r)
		re := regexp.MustCompile(regex_last_url)
		fileName := re.FindString(r.RequestURI)
		fileLoc := directory + "/" + fileName
		log.Debug("request file: " + fileLoc)

		if _, err := os.Stat(fileLoc); os.IsNotExist(err) {

			log.Error("file does not exist: " + fileLoc)
			w.WriteHeader(http.StatusNotFound)
		} else {

			if strings.Contains(fileName, ".ts") {

				serveTs(w, r, fileLoc)
			} else if strings.Contains(fileName, ".m3u8") {

				serveM3u8(w, r, fileLoc)
			} else {

				log.Error("Wrong filetype")
				w.WriteHeader(http.StatusNotFound)
			}
		}

	} else {

		log.Error("Url: " + r.RequestURI + " method: " + r.Method + " not available")
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func serveM3u8(w http.ResponseWriter, r *http.Request, path string) {

	w.Header().Add("Content-Type", "application/x-mpegURL")
	http.ServeFile(w, r, path)
}

func serveTs(w http.ResponseWriter, r *http.Request, path string) {

	w.Header().Add("Content-Type", "video/MP2T")
	http.ServeFile(w, r, path)
}

/* determine directory on disc */
func determineDirectoryByURL(r *http.Request) string {

	keys := Sessions.GetAvailableKeys()
	path := ""

	for _, key := range keys {

		if strings.Contains(r.RequestURI, key) {

			path += key
		}
	}

	return path
}