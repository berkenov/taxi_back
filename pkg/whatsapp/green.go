package whatsapp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

// GreenAPI implements Messenger using Green API service
type GreenAPI struct {
	apiURL     string
	idInstance string
	apiToken   string
	client     *http.Client
}

// NewGreenAPI creates a new Green API client from environment variables
func NewGreenAPI() *GreenAPI {
	apiURL := os.Getenv("API_URL")
	if apiURL == "" {
		apiURL = "https://api.green-api.com"
	}
	idInstance := os.Getenv("ID_INSTANCE")
	apiToken := os.Getenv("API_TOKEN")

	return &GreenAPI{
		apiURL:     apiURL,
		idInstance: idInstance,
		apiToken:   apiToken,
		client:     &http.Client{},
	}
}

// sendMessageRequest represents Green API sendMessage request body
type sendMessageRequest struct {
	ChatID  string `json:"chatId"`
	Message string `json:"message"`
}

// SendOTP sends a 4-digit OTP code to the specified phone number via Green API
func (g *GreenAPI) SendOTP(phone, code string) error {
	chatID := phone + "@c.us"
	message := fmt.Sprintf("Ваш код подтверждения: %s", code)

	reqBody := sendMessageRequest{
		ChatID:  chatID,
		Message: message,
	}
	body, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("marshal request: %w", err)
	}

	url := fmt.Sprintf("%s/waInstance%s/sendMessage/%s", g.apiURL, g.idInstance, g.apiToken)
	log.Printf("[GreenAPI] Sending OTP to %s: POST %s", phone, url)

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := g.client.Do(req)
	if err != nil {
		log.Printf("[GreenAPI] Request failed: %v", err)
		return fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("[GreenAPI] Unexpected status %d for phone %s", resp.StatusCode, phone)
		return fmt.Errorf("api returned status %d", resp.StatusCode)
	}

	log.Printf("[GreenAPI] OTP sent successfully to %s", phone)
	return nil
}
