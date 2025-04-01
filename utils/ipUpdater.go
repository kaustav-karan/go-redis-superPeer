package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

var (
	serverIP     string
	ipMutex      sync.RWMutex
	ipExpiryTime time.Time
)

type IPResponse struct {
	DeviceID  string `json:"device_id"`
	LocalIP   string `json:"local_ip"`
	ExpiresIn int    `json:"expires_in"`
}

// GetServerIP returns the current server IP in a thread-safe manner
func GetServerIP() string {
	ipMutex.RLock()
	defer ipMutex.RUnlock()
	return serverIP
}

func updateServerIP() error {
	resp, err := http.Get("https://beacon.korpsin.in/lookup?device_id=server")
	if err != nil {
		return fmt.Errorf("failed to fetch IP: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var ipResp IPResponse
	if err := json.NewDecoder(resp.Body).Decode(&ipResp); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	ipMutex.Lock()
	defer ipMutex.Unlock()
	
	serverIP = ipResp.LocalIP
	ipExpiryTime = time.Now().Add(time.Duration(ipResp.ExpiresIn) * time.Second)
	
	return nil
}

func StartIPUpdater() {
	// Initial update
	if err := updateServerIP(); err != nil {
		log.Printf("Failed to initialize server IP: %v", err)
	} else {
		log.Printf("Initial server IP: %s (expires at %v)", serverIP, ipExpiryTime)
	}

	// Start periodic updater
	for {
		ipMutex.RLock()
		sleepDuration := time.Until(ipExpiryTime) - (30 * time.Second)
		ipMutex.RUnlock()

		if sleepDuration > 0 {
			time.Sleep(sleepDuration)
		}

		if err := updateServerIP(); err != nil {
			log.Printf("Failed to update server IP: %v. Retrying in 1 minute...", err)
			time.Sleep(1 * time.Minute)
			continue
		}

		log.Printf("Updated server IP: %s (expires at %v)", serverIP, ipExpiryTime)
	}
}