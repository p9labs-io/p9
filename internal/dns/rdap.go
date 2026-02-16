package dns

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type RdapConfig struct {
	Services [][][]string `json:"services"`
}

var rdapconfig RdapConfig

//	//TODO: avoid repetitive file download

func RdapServer(tld string) string {
	configPath := createConfig()
	//load map from config to var or read map from file?
}

func getRdapConfig() {

	res, err := http.Get("https://data.iana.org/rdap/dns.json")
	if err != nil {
		log.Fatal(err)
	}
	body, err := io.ReadAll(res.Body)
	res.Body.Close()
	if res.StatusCode > 299 {
		log.Fatalf("Response failed with status code: %d and\nbody: %s\n", res.StatusCode, body)
	}
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(body, &rdapconfig)
	if err != nil {
		fmt.Println("error:", err)
	}
	config := transformToMap(rdapconfig)

	// build file path
	cfgdir, err := os.UserConfigDir()
	if err != nil {
		fmt.Printf("Error getting home directory: %v\n", err)
	}

	configDir := filepath.Join(cfgdir, ".p9")
	filePath := filepath.Join(configDir, "rdap_config")

	err = os.MkdirAll(configDir, 0755)
	if err != nil {
		fmt.Printf("Error creating directory: %v\n", err)
	}
	//write to file
	file, err := os.Create(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(config)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Successfully created: %s\n", filePath)

	//var rawconfig [][]string //got json from url https://data.iana.org/rdap/dns.json, .services field
	//config := make(map[string]string)
	//for _, v := range rawconfig {
	//	config[v[0]] = v[1]
	//}
	//return config

}

// 1 domain as a key
func transformToMap(rawconfig RdapConfig) map[string]string {
	config := make(map[string]string)

	for _, entry := range rawconfig.Services {
		// entry[0] is the slice of TLDs: ["com", "net"]
		// entry[1] is the slice of URLs: ["https://..."]

		tlds := entry[0]

		// Usually, there's only one URL, so we take the first one
		var url string
		if len(entry[1]) > 0 {
			url = entry[1][0]
		}

		for _, tld := range tlds {
			config[tld] = url
		}
	}
	return config
}

func createConfig() string {
	// build file path
	cfgdir, err := os.UserConfigDir()
	if err != nil {
		fmt.Printf("Error getting home directory: %v\n", err)
	}

	configDir := filepath.Join(cfgdir, ".p9")
	filePath := filepath.Join(configDir, "rdap_config")
	fileInfo, err := os.Stat(filePath)
	age := time.Since(fileInfo.ModTime())
	if os.IsNotExist(err) || age.Hours() > 720 {
		getRdapConfig()
	}

	return filePath
}

func loadConfigFromFile(filePath string) (map[string]string, error) {
	// TODO:
	// 1. Open file

	// 2. Create json.Decoder
	// 3. Decode into map[string]string
	// 4. Return the map
}
