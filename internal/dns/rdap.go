/*
 *
 *  * Copyright 2026 P9 Labs
 *  *
 *  * Licensed under the Apache License, Version 2.0 (the "License");
 *  * you may not use this file except in compliance with the License.
 *
 */

package dns

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Lookup

func LookupDomain(domain string) {

	baseURL := RdapServer(domain)
	url := fmt.Sprint(baseURL, "domain/", domain)
	res, err := http.Get(url)
	if err != nil {
		log.Fatalf("Failed to query RDAP server %s: %v", url, err)
	}
	body, err := io.ReadAll(res.Body)
	res.Body.Close()
	if res.StatusCode > 399 {
		log.Fatalf("Response failed with status code: %d and\nbody: %s\n", res.StatusCode, body)
	}
	if err != nil {
		log.Fatal(err)
	}

	type Event struct {
		EventAction string `json:"eventAction"`
		EventDate   string `json:"eventDate"`
	}

	type Nameservers struct {
		LdhName string `json:"ldhName"`
	}

	type SecureDNS struct {
		DelegationSigned bool `json:"delegationSigned"`
		MaxSigLife       uint `json:"maxSigLife"`
	}

	type RdapResponse struct {
		LdhName     string        `json:"ldhName"`
		Nameservers []Nameservers `json:"nameservers"`
		Status      []string      `json:"status"`
		Events      []Event       `json:"events"`
		SecureDNS   *SecureDNS    `json:"secureDNS"`
	}

	var rdapResponse RdapResponse

	err = json.Unmarshal(body, &rdapResponse)
	if err != nil {
		fmt.Println("error:", err)
	}

	output, err := json.MarshalIndent(rdapResponse, "", "  ") // 2-space indent
	if err != nil {
		fmt.Println("error:", err)
	}
	os.Stdout.Write(output)

}

// Get RDAP server url
type RdapConfig struct {
	Services [][][]string `json:"services"`
}

var rdapconfig RdapConfig

func RdapServer(domain string) string {
	config := createConfig()
	tld := extractTLD(domain)
	baseURL := config[tld]
	if baseURL == "" {
		log.Fatalf("TLD '.%s' is not supported by RDAP", tld)
	}
	return baseURL
}

func extractTLD(domain string) string {
	i := strings.LastIndex(domain, ".")
	tld := domain[i+1:]
	return tld
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

func createConfig() map[string]string {
	cfgdir, err := os.UserConfigDir()
	if err != nil {
		log.Fatalf("Error getting config directory: %v", err)
	}

	configDir := filepath.Join(cfgdir, ".p9")
	filePath := filepath.Join(configDir, "rdap_config")

	fileInfo, err := os.Stat(filePath)

	if os.IsNotExist(err) {
		getRdapConfig()
		config, err := loadConfigFromFile(filePath)
		if err != nil {
			log.Fatalf("Can't read newly created config: %v", err)
		}
		return config
	}

	if err != nil {
		log.Fatalf("Can't access file %s: %v", filePath, err)
	}

	age := time.Since(fileInfo.ModTime())
	if age.Hours() > 720 {
		getRdapConfig()
	}

	config, err := loadConfigFromFile(filePath)
	if err != nil {
		log.Fatalf("Can't read config file: %v", err)
	}

	return config
}

func loadConfigFromFile(filePath string) (map[string]string, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var config map[string]string
	err = json.Unmarshal(data, &config) // Convert bytes to map
	if err != nil {
		return nil, err
	}

	return config, nil
}
