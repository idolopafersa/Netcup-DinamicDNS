package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/lpernett/godotenv"
)

const apiurl = "https://ccp.netcup.net/run/webservice/servers/endpoint.php?JSON"

type DnsRecord struct {
	ID          string `json:"id"`
	Hostname    string `json:"hostname"`
	RecordType  string `json:"type"`
	Destination string `json:"destination"`
}

func InfoDomain(customerNumber, apiKey, apiSessionID, domainName string) (string, error) {
	data := map[string]interface{}{
		"action": "infoDnsRecords",
		"param": map[string]interface{}{
			"domainname":     domainName,
			"customernumber": customerNumber,
			"apikey":         apiKey,
			"apisessionid":   apiSessionID,
		},
	}

	jsonData, _ := json.Marshal(data)

	resp, err := http.Post(apiurl, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("error sending request: %s", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var response struct {
		ResponseData struct {
			DnsRecords []DnsRecord `json:"dnsrecords"`
		} `json:"responsedata"`
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		return "", fmt.Errorf("error parsing JSON: %s", err)
	}

	for _, record := range response.ResponseData.DnsRecords {

		if record.RecordType == "A" && (record.Hostname == "@" || record.Hostname == domainName) {
			return record.Destination, nil
		}
	}

	return "", fmt.Errorf("A record not found for domain %s", domainName)
}

func UpdateDomain(customerNumber, apiKey, apiSessionID, domainName, newIP string) error {

	data := map[string]interface{}{
		"action": "infoDnsRecords",
		"param": map[string]interface{}{
			"domainname":     domainName,
			"customernumber": customerNumber,
			"apikey":         apiKey,
			"apisessionid":   apiSessionID,
		},
	}

	jsonData, _ := json.Marshal(data)

	resp, err := http.Post(apiurl, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("error sending info request: %s", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var response struct {
		ResponseData struct {
			DnsRecords []DnsRecord `json:"dnsrecords"`
		} `json:"responsedata"`
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		return fmt.Errorf("error parsing JSON: %s", err)
	}

	var rootid int
	var wildcardRecordID int

	for _, record := range response.ResponseData.DnsRecords {
		if record.RecordType == "A" && (record.Hostname == "@" || record.Hostname == domainName) {
			rootid, _ = strconv.Atoi(record.ID)

		}
		if record.RecordType == "A" && (record.Hostname == "*" || record.Hostname == domainName) {
			wildcardRecordID, _ = strconv.Atoi(record.ID)

		}
	}

	var dnsRecordsToUpdate []map[string]interface{}

	if rootid != 0 {
		dnsRecordsToUpdate = append(dnsRecordsToUpdate, map[string]interface{}{
			"id": rootid, "hostname": "@", "type": "A", "destination": newIP, "deleterecord": false,
		})
	}
	if wildcardRecordID != 0 {
		dnsRecordsToUpdate = append(dnsRecordsToUpdate, map[string]interface{}{
			"id": wildcardRecordID, "hostname": "*", "type": "A", "destination": newIP, "deleterecord": false,
		})
	}

	if len(dnsRecordsToUpdate) == 0 {
		return fmt.Errorf("No existing DNS records found to update")
	}

	updateData := map[string]interface{}{
		"action": "updateDnsRecords",
		"param": map[string]interface{}{
			"domainname":     domainName,
			"customernumber": customerNumber,
			"apikey":         apiKey,
			"apisessionid":   apiSessionID,
			"dnsrecordset": map[string]interface{}{
				"dnsrecords": dnsRecordsToUpdate,
			},
		},
	}

	updateJsonData, _ := json.Marshal(updateData)

	updateResp, err := http.Post(apiurl, "application/json", bytes.NewBuffer(updateJsonData))
	if err != nil {
		return fmt.Errorf("error sending update request: %s", err)
	}
	defer updateResp.Body.Close()

	updateBody, _ := io.ReadAll(updateResp.Body)

	var updateResponse struct {
		Status string `json:"status"`
	}
	err = json.Unmarshal(updateBody, &updateResponse)
	if err != nil {
		return fmt.Errorf("error parsing update response: %s", err)
	}

	if updateResponse.Status != "success" {
		return fmt.Errorf("failed to update DNS record")
	}

	fmt.Println("DNS record updated successfully.")
	return nil
}

func Login(customerNumber, apikey, apipassword string) (string, error) {
	data := map[string]interface{}{
		"action": "login",
		"param": map[string]string{
			"customernumber": customerNumber,
			"apikey":         apikey,
			"apipassword":    apipassword,
		},
	}

	jsonData, _ := json.Marshal(data)

	resp, err := http.Post(apiurl, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("error sending request: %s", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var response struct {
		ResponseData struct {
			APISessionID string `json:"apisessionid"`
		} `json:"responsedata"`
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		return "", fmt.Errorf("error parsing JSON: %s", err)
	}

	return response.ResponseData.APISessionID, nil
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	customernumber := os.Getenv("customernumber")
	apikey := os.Getenv("apikey")
	apipassword := os.Getenv("apipassword")
	domain := os.Getenv("domain")

	resp, err := http.Get("https://api.ipify.org")
	if err != nil {
		fmt.Printf("Could not get public IP: %s\n", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response: %s\n", err)
		return
	}

	publicip := string(body)
	fmt.Println("Public IP:", publicip)

	apisession, err := Login(customernumber, apikey, apipassword)
	if err != nil {
		log.Fatalf("Login failed: %s", err)
	}

	currentIP, err := InfoDomain(customernumber, apikey, apisession, domain)
	if err != nil {
		log.Fatalf("Error retrieving domain info: %s", err)
	}

	fmt.Println("Current DNS A record:", currentIP)

	if currentIP != publicip {
		fmt.Println("Updating DNS record...")
		err = UpdateDomain(customernumber, apikey, apisession, domain, publicip)
		if err != nil {
			log.Fatalf("Failed to update DNS record: %s", err)
		}
	} else {
		fmt.Println("No update needed. IP is already correct.")
	}
}
