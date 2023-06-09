package ais

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

var (
	http_timeout = 30 * time.Second
)

const (
	API_KEY_ENV = "OPENAI_API_KEY"
)

var apiKey string

func init() {
	apiKey = os.Getenv(API_KEY_ENV)
	if apiKey == "" {
		log.Fatalf("%s not set\n", API_KEY_ENV)
	}
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type PromptRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

func SendPrompt(prompt, prev_ans string, max_tokens int) ([]string, error) {
	res := []string{}
	url := "https://api.openai.com/v1/chat/completions"

	PromptRequest := PromptRequest{
		Model: "gpt-3.5-turbo",
		Messages: []Message{
			{Role: "assistant", Content: prev_ans},
			{Role: "user", Content: prompt},
		},
	}

	b, err := json.Marshal(PromptRequest)
	if err != nil {
		return res, err
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(b))
	if err != nil {
		return res, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{
		Timeout: http_timeout,
	}
	resp, err := client.Do(req)
	if err != nil {
		return res, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return res, err
	}

	type Choice struct {
		Message Message `json:"message"`
	}

	type respMessage struct {
		Choices []Choice `json:"choices"`
	}

	r := respMessage{}
	err = json.Unmarshal(body, &r)
	if err != nil {
		return res, err
	}
	for i := range r.Choices {
		res = append(res, r.Choices[i].Message.Content)
	}

	if len(res) == 0 {
		return res, errors.New("Communication failed")
	}
	return res, nil
}

type Engine struct {
	ID    string `json:"id"`
	Owner string `json:"owner"`
	Ready bool   `json:"ready"`
}

type ListEnginesResponse struct {
	Data []Engine `json:"data"`
}

func ListEngines() ([]Engine, error) {
	res := []Engine{}
	url := "https://api.openai.com/v1/engines"

	req, err := http.NewRequest("GET", url, bytes.NewBuffer(nil))
	if err != nil {
		return res, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{
		Timeout: http_timeout,
	}
	resp, err := client.Do(req)
	if err != nil {
		return res, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return res, err
	}

	listresp := ListEnginesResponse{}
	err = json.Unmarshal(body, &listresp)
	if err != nil {
		return res, err
	}

	return listresp.Data, nil
}

// QuotaUsage represents the OpenAI API usage quota information
type QuotaUsage struct {
	Used float64
}

// GetOpenAIQuotaUsage retrieves the current usage quota for the OpenAI API
// using the provided API key and returns the quota usage information as a struct
func GetOpenAIQuotaUsage() (QuotaUsage, error) {
	req, err := http.NewRequest("GET", "https://api.openai.com/v1/usage", nil)
	if err != nil {
		return QuotaUsage{}, fmt.Errorf("error creating request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{
		Timeout: http_timeout,
	}
	resp, err := client.Do(req)
	if err != nil {
		return QuotaUsage{}, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	var usage struct {
		UsageUsd float64 `json:"current_usage_usd"`
	}

	err = json.NewDecoder(resp.Body).Decode(&usage)
	if err != nil {
		return QuotaUsage{}, fmt.Errorf("error decoding response: %v", err)
	}

	quotaUsage := QuotaUsage{
		Used: usage.UsageUsd,
	}

	return quotaUsage, nil
}
