package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/user"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	fast "gopkg.in/ddo/go-fast.v0"
)

var keyAddr string
var Speed float64

type UpdateNodeStats struct {
	Token       string `json:"token"`
	AccountAddr string `json:"account_addr"`
	Info        info   `json:"info"`
}
type info struct {
	Type string `json:"type"`
}

var username string

// Socks5Config info
type Socks5Config struct {
	Server       string   `json:"server"`
	PortPassword portPass `json:"port_password"`
	Timeout      int32    `json:"timeout"`
	Method       string   `json:"method"`
}

//NetworkInfo represent the network info about current socks5 Node
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

type portPass struct {
	PortPassword0 string `json:"4200"`
	PortPassword1 string `json:"4201"`
	PortPassword2 string `json:"4202"`
	PortPassword3 string `json:"4203"`
}

type UserKeystore struct {
	keystoreData keystore `json:"keystore"`
}

type keystore struct {
	Version string `json:"version"`
	ID      string `json:"id"`
	Address string `json:"address"`
	Crypto  crypto `json:"Crypto"`
}

type crypto struct {
	CipherText   string `json:"ciphertext"`
	CipherParams struct {
		IV string `json:"iv"`
	} `json:"cipherparams"`
	Cipher    string `json:"cipher"`
	KDF       string `json:"kdf"`
	KDFParams struct {
		DKLEN int32  `json:"dklen"`
		Salt  string `json:"salt"`
		N     int32  `json:"n"`
		R     int32  `json:"r"`
		P     int32  `json:"p"`
	} `json:"kdfparams"`
	MAC string `json:"mac"`
}

type configData struct {
	Token       string `json:"token"`
	AccountAddr string `json:"account_addr"`
	PricePerGB  int32  `json:"price_per_gb"`
}

//NodeStats represent the info about current socks5 Node
type NodeStats struct {
	AccountAddr string `json:"account_addr"`
	PricePerGB  string `json:"price_per_gb"`
	IPAddr      string `json:"ip"`
	Location    string `json:"location"`
	NetSpeed    string `json:"net_speed"`
	VPNType     string `json:"vpn_type"`
}

type UserAccount struct {
	AccountAddress string `json:"account_addr"`
}

type password struct {
	Password string `json:"password"`
}

type KeyString struct {
	Keystore string `json:"keystore"`
}

type Token struct {
	Token string `json:"token"`
}

type WalletAddr struct {
	AccountAddr string `json:"account_addr"`
}

// runCmd represents the new command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Spin up a new SOCKS5 node in Sentinel's Distributed Network",
	Run: func(cmd *cobra.Command, args []string) {
		StartSocks5Node()
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}

func createConfig() {

	file, err := os.Create("/home/" + username + "/.config.json")
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

//StartSocks5Node start a shadowsocks server as daemon
func StartSocks5Node() {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}
	username = user.Username
	createConfig()
	cmd := "ssserver -c /home/" + username + "/.config.json"
	cmdParts := strings.Fields(cmd)

	startSocks := exec.Command(cmdParts[0], cmdParts[1:]...)

	if err := startSocks.Start(); err != nil {
		// fmt.Errorf("Could Not Start the Shadowsocks server: %v", err)
		log.Fatalf("you are missing shadowsocks server. Please install by running pip install shadowsocks: %v", err)
	}

	configServer := "socks-info-server"
	split := strings.Fields(configServer)
	startInfoServer := exec.Command(split[0], split[1:]...)

	if err := startInfoServer.Start(); err != nil {
		// fmt.Errorf("Could Not Start the Shadowsocks server: %v", err)
		log.Fatalf("you are missing shadowsocks server. Please install by running pip install shadowsocks: %v", err)
	}
	fmt.Println("opened port :8080 for clients go get server config")
	RegisterNode()
	updateNodeInfo()
	fmt.Println(Speed)
	fmt.Println("Started SOCKS5 Server")
}

func RegisterNode() {

	getNewWallet()
	URL := "https://api.sentinelgroup.io/node/register"
	fmt.Println("URL:>", URL)

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

	// var node NodeStatsc

	fastCom := fast.New()

	// init
	err = fastCom.Init()
	if err != nil {
		panic(err)
	}

	// get urls
	urls, err := fastCom.GetUrls()
	if err != nil {
		panic(err)
	}

	// measure
	KbpsChan := make(chan float64)

	go func() {
		for Kbps := range KbpsChan {
			Speed = Kbps

			// fmt.Printf("%.2f %.2f \n", Kbps, Kbps/1000)
		}

		fmt.Println("done")
	}()

	err = fastCom.Measure(urls, KbpsChan)
	if err != nil {
		panic(err)
	}

	some := &NodeStats{
		IPAddr:      netStat.IPAddr,
		Location:    netStat.City + netStat.Country,
		AccountAddr: keyAddr,
		NetSpeed:    floattostrwithprec(Speed, 4),
		PricePerGB:  "200",
		VPNType:     "socks5",
	}

	may, err := json.Marshal(some)

	// var bufffff bytes.Buffer
	// binary.Write(&bufffff, binary.BigEndian, some)

	// .Parse(`{"account_addr":${{accountAddr}}, "price_per_gb":"200", "location": "Hyderabad", "ip": "10.10.10.10", "net_speed":"200"}`)
	// accountAddr := "0x672092107c0b940a8a7041f827c85801f30af9b2"
	// var finData = []byte(fmt.Sprintf("%v", some))
	// var jsonStr = []byte(`{"account_addr":${{accountAddr}}, "price_per_gb":"200", "location": "Hyderabad", "ip": "10.10.10.10", "net_speed":"200"}`)
	req, err := http.NewRequest("POST", URL, bytes.NewBuffer(may))
	// req.Header.Set("X-Custom-Header", "myvalue")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	token := &Token{}
	body, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(body, &token)
	if err != nil {
		panic(err)
	}

	fmt.Println("response Body IS:", string(body))

	file, err := os.Create("config.data")
	if err != nil {
		panic(err)
	}

	bytes, err := json.Marshal(configData{
		Token:       token.Token,
		PricePerGB:  200,
		AccountAddr: keyAddr,
	})
	defer file.Close()
	fmt.Fprintf(file, string(bytes))
}

// Saving Keystore for user and node config

func getNewWallet() {
	URL := "https://api.sentinelgroup.io/node/account"

	pass := &password{
		Password: "witty@123",
	}
	mrsh, _ := json.Marshal(pass)
	req, err := http.NewRequest("POST", URL, bytes.NewBuffer(mrsh))
	if err != nil {
		fmt.Errorf("Could Not Get Wallet Address: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer req.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	// fmt.Printf(string(body))

	accountInfo := &KeyString{}
	walletAddr := &WalletAddr{}
	err = json.Unmarshal(body, &accountInfo)
	if err != nil {
		fmt.Printf("Error here: %v", err)
	}

	err = json.Unmarshal(body, &walletAddr)
	if err != nil {
		fmt.Printf("Error here: %v", err)
	}
	keyAddr = walletAddr.AccountAddr
	fmt.Printf("%v", accountInfo)
	file, err := os.Create("UTCKeystore")
	if err != nil {
		panic(err)
	}
	bytes, err := json.Marshal(accountInfo.Keystore)
	defer file.Close()
	fmt.Fprintf(file, string(bytes))

}

func updateNodeInfo() {
	URL := "https://api.sentinelgroup.io/node/update-nodeinfo"

	file, err := os.Open("config.data")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	b, _ := ioutil.ReadAll(file)
	// if err != nil {
	// 	log.Fatalf("Error in Reading file: %v", err)
	// }
	fileOBJ := &UpdateNodeStats{}
	err = json.Unmarshal(b, &fileOBJ)
	NodeInfo := &UpdateNodeStats{
		Token:       fileOBJ.Token,
		AccountAddr: fileOBJ.AccountAddr,
		Info: info{
			Type: "vpn",
		},
	}

	mrsh, _ := json.Marshal(NodeInfo)
	req, err := http.NewRequest("POST", URL, bytes.NewBuffer(mrsh))
	if err != nil {
		fmt.Errorf("Could Not Get Wallet Address: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer req.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	// err = json.Unmarshal(body, &token)
	// if err != nil {
	// 	panic(err)
	// }

	fmt.Println(string(body))

}

func floattostrwithprec(fv float64, prec int) string {
	return strconv.FormatFloat(fv, 'f', prec, 64)
}
