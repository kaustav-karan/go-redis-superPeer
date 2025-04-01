package utils

import (
	"bytes"
	"errors"
	"log"
	"net"
	"net/http"
	"time"
	"encoding/json"
)

const (
	beaconURL        = "https://beacon.korpsin.in/beacon"
	deviceID         = "superpeer"
	retryInterval    = 15 * time.Second
	maxRetryDuration = 1 * time.Hour
	resendInterval   = 15 * time.Minute
)

func getLocalIP() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}
	for _, addr := range addrs {
		if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				return ipNet.IP.String(), nil
			}
		}
	}
	return "", errors.New("no local IP address found")
}

func sendBeaconRequest(localIP string) error {
	requestBody := map[string]string{
		"device_id": deviceID,
		"local_ip":  localIP,
	}
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return err
	}

	resp, err := http.Post(beaconURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New("server responded with status: " + resp.Status)
	}
	return nil
}

func StartBeaconSender() {
	for {
		localIP, err := getLocalIP()
		if err != nil {
			log.Printf("Error retrieving local IP: %v", err)
			time.Sleep(resendInterval)
			continue
		}

		startTime := time.Now()
		for {
			err := sendBeaconRequest(localIP)
			if err == nil {
				log.Printf("Successfully sent beacon signal with IP: %s", localIP)
				break
			}

			if time.Since(startTime) > maxRetryDuration {
				log.Printf("Server unreachable for 1 hour: %v", err)
				break
			}
			log.Printf("Failed to send beacon signal, retrying in 15 seconds: %v", err)
			time.Sleep(retryInterval)
		}

		time.Sleep(resendInterval)
	}
}