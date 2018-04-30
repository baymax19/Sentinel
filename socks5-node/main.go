package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/gorilla/mux"
	. "github.com/than-os/socks-info-server/config"
	. "github.com/than-os/socks-info-server/dbo"
	. "github.com/than-os/socks-info-server/models"
	"gopkg.in/mgo.v2/bson"
)

//NetworkInfo d
type NetworkInfo struct {
	// As          string  `json:"as"`
	City    string `json:"city"`
	Country string `json:"country"`
	// CountryCode string  `json:"countryCode"`
	// ISP         string  `json:"isp"`
	// Lat         float32 `json:"lat"`
	// Lon         float32 `json:"lon"`
	// Org         string  `json:"org"`
	IPAddr string `json:"query"`
	// Region      string  `json:"region"`
	// RegionName  string  `json:"regionName"`
	// Status      string  `json:"status"`
	// TimeZone    string  `json:"timezone"`
	// ZipCode     string  `json:"zip"`
}

// Socks5Config info
type Socks5Config struct {
	Server       string   `json:"server"`
	PortPassword portPass `json:"port_password"`
	Timeout      int32    `json:"timeout"`
	Method       string   `json:"method"`
}

type portPass struct {
	PortPassword0 string `json:"4200"`
	PortPassword1 string `json:"4201"`
	PortPassword2 string `json:"4202"`
	PortPassword3 string `json:"4203"`
}

var config = Config{}
var dbo = Socks5DBO{}

func init() {
	config.Read()

	dbo.Server = config.Server
	dbo.Database = config.Database
	dbo.Connect()
}

func main() {
	fmt.Println("Starting The App...")

	r := mux.NewRouter()

	r.HandleFunc("/nodes", GetAllNodes).Methods("GET")
	r.HandleFunc("/nodes/{id}", GetOneNode).Methods("GET")
	r.HandleFunc("/nodes", RegisterNewNode).Methods("POST")
	r.HandleFunc("/nodes", DeleteNode).Methods("DELETE")
	r.HandleFunc("/nodes", UpdateNode).Methods("PUT")

	if err := http.ListenAndServe(":8080", r); err != nil {
		fmt.Println("FATAL ERROR")
	}

	// CreateSocks5Config()
	// createConfig()
	// startShadowsocksNode()

}

// GetAllNodes will give you all stats for all of the registered nodes
func GetAllNodes(w http.ResponseWriter, r *http.Request) {

	nodes, err := dbo.FindAll()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return

	}

	respondWithJSON(w, http.StatusOK, nodes)
}

// GetOneNode takes a node ID and returns a single node stats
func GetOneNode(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	node, err := dbo.FindById(params["id"])
	if err != nil {

		respondWithError(w, http.StatusBadRequest, "Invalid Node ID")
		return
	}

	respondWithJSON(w, http.StatusOK, node)
}

//RegisterNewNode is used to insert a newly launched
//socks5node in the db
func RegisterNewNode(w http.ResponseWriter, r *http.Request) {

	defer r.Body.Close()

	ipAddr, err := http.Get("http://ip-api.com/json")
	if err != nil {
		os.Stderr.WriteString(err.Error())
		os.Stderr.WriteString("\n")
		os.Exit(1)
	}
	netStat := NetworkInfo{}

	defer ipAddr.Body.Close()
	jsn, err := ioutil.ReadAll(ipAddr.Body)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(jsn, &netStat)
	if err != nil {
		panic(err)
	}

	fmt.Printf("NetStats: %v\n", netStat.IPAddr)

	var node Node

	if err := json.NewDecoder(r.Body).Decode(&node); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Request Payload")
		return
	}

	node.ID = bson.NewObjectId()
	node.IPAddr = netStat.IPAddr
	node.Location = netStat.City + ", " + netStat.Country
	node.PortPasswords.PortPassword0 = "shadowsocks"
	node.PortPasswords.PortPassword1 = "shadowsocks"
	node.PortPasswords.PortPassword2 = "shadowsocks"
	node.PortPasswords.PortPassword3 = "shadowsocks"
	node.Method = "aes-256-cfb"

	if err := dbo.Insert(node); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, node)
}

// UpdateNode Node stats
func UpdateNode(w http.ResponseWriter, r *http.Request) {

	defer r.Body.Close()
	var node Node
	if err := json.NewDecoder(r.Body).Decode(&node); err != nil {

		respondWithError(w, http.StatusBadRequest, "Invalid Request Payload")
		return
	}

	if err := dbo.Update(node); err != nil {

		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, node)
}

// DeleteNode Will Delete a node's info from db
func DeleteNode(w http.ResponseWriter, r *http.Request) {

	defer r.Body.Close()
	var node Node
	if err := json.NewDecoder(r.Body).Decode(&node); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Payload")
		return
	}
	if err := dbo.Delete(node); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"result": "successfuly updated"})
}

func createConfig() {
	file, err := os.Create(".config.json")
	if err != nil {
		panic(err)
	}
	bytes, err := json.Marshal(Socks5Config{
		Server: "0.0.0.0",
		Method: "aes-256-cfb",
		PortPassword: portPass{
			PortPassword0: "shadowsocks",
			PortPassword1: "shadowsocks",
			PortPassword2: "shadowsocks",
			PortPassword3: "shadowsocks",
		},
		Timeout: 300,
	})
	defer file.Close()
	fmt.Fprintf(file, string(bytes))
}

func startShadowsocksNode() {
	cmd := "ssserver -c .config.json"
	cmdParts := strings.Fields(cmd)

	startSocks := exec.Command(cmdParts[0], cmdParts[1:]...)

	if err := startSocks.Start(); err != nil {
		// fmt.Errorf("Could Not Start the Shadowsocks server: %v", err)
		panic(err)
	}

	fmt.Println("Started successfully")

}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	respondWithJSON(w, code, map[string]string{"error": msg})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Server", "golang/1.10")
	w.WriteHeader(code)
	w.Write(response)
}
