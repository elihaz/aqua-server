package main

import (
	"io/ioutil"
	"encoding/json"
	"net/http"
	"fmt"
	"encoding/xml"
	"log"
	"os"
)

type FileInformation struct {
	Name string `xml:"name" json:"name"`
	Size int64  `xml:"size" json:"size"`
	Ext  string `xml:"ext" json:"ext"`
}

type FileStats struct {
	NumOfFiles      int
	AverageFileSize int64
	Extensions      map[string]int
	FrequentExt     string
	Max             *FileInformation
}

type Server struct {
	stats      *FileStats
	port       int
	protocol   string
	format     string
	setChannel chan *FileInformation
	getChannel chan *FileStats
}

func newServer(port int, protocol string, format string) *Server {
	server := new(Server)
	stats := new(FileStats)
	stats.Extensions = make(map[string]int)
	stats.Max = new(FileInformation)

	server.stats = stats
	server.port = port
	server.protocol = protocol
	server.format = format

	server.setChannel = make(chan *FileInformation)
	server.getChannel = make(chan *FileStats)

	go server.coordinator()

	return server
}

func (s *Server) start() {

	http.HandleFunc("/get-stats", s.getHandler)
	http.HandleFunc("/update", s.updateHandler)

	if s.protocol == "HTTPS" {

		var serverCrt string
		var serverKey string

		serverCrt, serverKey = getAuthFiles(serverCrt, serverKey)

		e := http.ListenAndServeTLS(fmt.Sprintf(":%d", s.port),
			serverCrt,
			serverKey, nil)
		if e != nil {
			log.Print(e)
		}
	} else {
		e := http.ListenAndServe(fmt.Sprintf(":%d", s.port), nil)
		if e != nil {
			log.Print(e)
		}
	}
}

func getAuthFiles(serverCrt string, serverKey string) (string, string) {
	if os.Getenv("SERVER_CRT") != "" {
		serverCrt = os.Getenv("SERVER_CRT")
	} else {
		serverCrt = "server.crt"
	}
	if os.Getenv("SERVER_KEY") != "" {
		serverKey = os.Getenv("SERVER_KEY")
	} else {
		serverKey = "server.key"
	}
	return serverCrt, serverKey
}

func (s *Server) getHandler(w http.ResponseWriter, request *http.Request) {
	if request.Method == "GET" {
		stats := s.getCurrentStats()
		bytes, e := json.Marshal(stats)
		if e != nil {
			http.Error(w, "internal error", 500)
			return
		}
		w.Header().Set("content-type", "application/json")
		w.Write(bytes)
	} else {
		http.Error(w, "Method not supported", 405)
	}

}

func (s *Server) updateHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" {
		b, err := ioutil.ReadAll(r.Body)
		defer r.Body.Close()

		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		var fileInfo *FileInformation

		if s.format == "JSON" {
			err = json.Unmarshal(b, &fileInfo)
		} else {
			err = xml.Unmarshal(b, &fileInfo)
		}

		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		s.setStats(fileInfo)

		w.Header().Set("content-type", "application/json")
		w.Write([]byte("OK"))
	} else {
		http.Error(w, "Method not supported", 405)
	}

}

func (s *Server) coordinator() {

	var fileInfo *FileInformation
	for {
		select {
		case fileInfo = <-s.setChannel:

			s.addToExtensions(fileInfo)
			s.calcMostFrequent(fileInfo)
			s.calcMaxSize(fileInfo)
			s.updateAverageSize(fileInfo)
			s.incrementNumOfFiles()

		case s.getChannel <- s.stats:

		}
	}
}

func (s *Server) getCurrentStats() *FileStats {
	return <-s.getChannel
}
func (s *Server) setStats(fileInfo *FileInformation) {
	s.setChannel <- fileInfo
}

func (s *Server) addToExtensions(fileInfo *FileInformation) {
	s.stats.Extensions[fileInfo.Ext]++
}

func (s *Server) incrementNumOfFiles() {
	s.stats.NumOfFiles += 1
}

func (s *Server) updateAverageSize(fileInfo *FileInformation) {
	currentTotalSize := s.stats.AverageFileSize * int64(s.stats.NumOfFiles)
	newTotalSize := currentTotalSize + fileInfo.Size
	s.stats.AverageFileSize = newTotalSize / (int64(s.stats.NumOfFiles + 1))
}

func (s *Server) calcMaxSize(fileInfo *FileInformation) {
	if fileInfo.Size > s.stats.Max.Size {
		s.stats.Max.Name = fileInfo.Name
		s.stats.Max.Size = fileInfo.Size
		s.stats.Max.Ext = fileInfo.Ext
	}
}

func (s *Server) calcMostFrequent(fileInfo *FileInformation) {
	if s.stats.Extensions[fileInfo.Ext] > s.stats.Extensions[s.stats.FrequentExt] {
		s.stats.FrequentExt = fileInfo.Ext
	}
}
