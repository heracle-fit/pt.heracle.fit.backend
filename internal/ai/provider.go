package ai

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/heracle/pt.heracle.fit.go/internal/config"
)

type Provider interface {
	Complete(systemPrompt, userContent string) (string, error)
	CompleteVision(systemPrompt string, description *string, imageData []byte, mimeType string) (string, error)
}

// ── Router ──────────────────────────────────────────────────────────────────────

type AIRouter struct {
	cfg     *config.Config
	openai  *OpenAIProvider
	gemini  *GeminiProvider
	hf      *HuggingFaceProvider
}

func NewAIRouter(cfg *config.Config) *AIRouter {
	return &AIRouter{
		cfg:    cfg,
		openai: &OpenAIProvider{apiKey: cfg.OpenAIKey},
		gemini: &GeminiProvider{apiKey: cfg.GeminiKey},
		hf:     &HuggingFaceProvider{apiKey: cfg.HuggingFaceKey},
	}
}

func (r *AIRouter) GetProvider(providerName string) Provider {
	switch strings.ToLower(providerName) {
	case "gemini":
		return r.gemini
	case "huggingface":
		return r.hf
	default:
		return r.openai
	}
}

func (r *AIRouter) GetModel(providerName, featureDefault string) string {
	switch strings.ToLower(providerName) {
	case "gemini":
		if featureDefault != "" {
			return featureDefault
		}
		return "gemini-1.5-flash"
	case "huggingface":
		if featureDefault != "" {
			return featureDefault
		}
		return "Qwen/Qwen2.5-VL-72B-Instruct"
	default:
		if featureDefault != "" {
			return featureDefault
		}
		return "gpt-4o"
	}
}

func (r *AIRouter) RunFoodAnalysis(description *string, imageData []byte, mimeType string) (map[string]interface{}, error) {
	provider := r.GetProvider(r.cfg.FoodAnalyseProvider)
	model := r.GetModel(r.cfg.FoodAnalyseProvider, r.cfg.FoodAnalyseModel)

	var raw string
	var err error

	if len(imageData) > 0 || description != nil {
		raw, err = provider.CompleteVision(FoodAnalysisPrompt, description, imageData, mimeType)
	} else {
		return nil, fmt.Errorf("at least one of image or description must be provided")
	}

	if err != nil {
		return nil, fmt.Errorf("%s (model: %s): %w", r.cfg.FoodAnalyseProvider, model, err)
	}

	cleaned := cleanAIResponse(raw)
	var result map[string]interface{}
	if err := json.Unmarshal([]byte(cleaned), &result); err != nil {
		return nil, fmt.Errorf("%s returned unexpected format", r.cfg.FoodAnalyseProvider)
	}
	return result, nil
}

func (r *AIRouter) RunDietSuggestion(userContext string) (map[string]interface{}, error) {
	provider := r.GetProvider(r.cfg.DietSuggestionProvider)

	raw, err := provider.Complete(DietSuggestionPrompt, userContext)
	if err != nil {
		return nil, err
	}

	cleaned := cleanAIResponse(raw)
	var result map[string]interface{}
	if err := json.Unmarshal([]byte(cleaned), &result); err != nil {
		return nil, fmt.Errorf("failed to parse AI diet suggestion")
	}
	return result, nil
}

func (r *AIRouter) RunSleepInsight(sleepContext string) (map[string]interface{}, error) {
	provider := r.GetProvider(r.cfg.SleepInsightProvider)

	raw, err := provider.Complete(SleepInsightPrompt, sleepContext)
	if err != nil {
		return nil, err
	}

	cleaned := cleanAIResponse(raw)
	var result map[string]interface{}
	if err := json.Unmarshal([]byte(cleaned), &result); err != nil {
		return nil, fmt.Errorf("failed to parse AI sleep insight")
	}
	return result, nil
}

func cleanAIResponse(raw string) string {
	raw = strings.ReplaceAll(raw, "```json", "")
	raw = strings.ReplaceAll(raw, "```", "")
	return strings.TrimSpace(raw)
}

// ── OpenAI ──────────────────────────────────────────────────────────────────────

type OpenAIProvider struct {
	apiKey string
}

func (p *OpenAIProvider) Complete(systemPrompt, userContent string) (string, error) {
	payload := map[string]interface{}{
		"model": "gpt-4o",
		"messages": []map[string]interface{}{
			{"role": "system", "content": systemPrompt},
			{"role": "user", "content": userContent},
		},
		"max_tokens": 1024,
	}
	return p.makeRequest(payload)
}

func (p *OpenAIProvider) CompleteVision(systemPrompt string, description *string, imageData []byte, mimeType string) (string, error) {
	userContent := []map[string]interface{}{}
	if len(imageData) > 0 {
		if mimeType == "" {
			mimeType = "image/jpeg"
		}
		userContent = append(userContent, map[string]interface{}{
			"type": "image_url",
			"image_url": map[string]string{
				"url":    fmt.Sprintf("data:%s;base64,%s", mimeType, base64.StdEncoding.EncodeToString(imageData)),
				"detail": "auto",
			},
		})
	}
	if description != nil && *description != "" {
		userContent = append(userContent, map[string]interface{}{
			"type": "text", "text": *description,
		})
	}

	payload := map[string]interface{}{
		"model": "gpt-4o",
		"messages": []map[string]interface{}{
			{"role": "system", "content": systemPrompt},
			{"role": "user", "content": userContent},
		},
		"max_tokens": 1024,
	}
	return p.makeRequest(payload)
}

func (p *OpenAIProvider) makeRequest(payload map[string]interface{}) (string, error) {
	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+p.apiKey)

	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	var result map[string]interface{}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return "", fmt.Errorf("failed to parse OpenAI response")
	}

	choices, ok := result["choices"].([]interface{})
	if !ok || len(choices) == 0 {
		return "", fmt.Errorf("no choices in OpenAI response")
	}
	msg, _ := choices[0].(map[string]interface{})["message"].(map[string]interface{})
	content, _ := msg["content"].(string)
	return content, nil
}

// ── Gemini ──────────────────────────────────────────────────────────────────────

type GeminiProvider struct {
	apiKey string
}

func (p *GeminiProvider) Complete(systemPrompt, userContent string) (string, error) {
	model := "gemini-2.5-flash-lite"
	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/%s:generateContent?key=%s", model, p.apiKey)

	payload := map[string]interface{}{
		"contents": []map[string]interface{}{
			{
				"parts": []map[string]interface{}{
					{"text": systemPrompt},
					{"text": userContent},
				},
			},
		},
	}
	return p.makeRequest(url, payload)
}

func (p *GeminiProvider) CompleteVision(systemPrompt string, description *string, imageData []byte, mimeType string) (string, error) {
	model := "gemini-1.5-flash"
	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/%s:generateContent?key=%s", model, p.apiKey)

	parts := []map[string]interface{}{
		{"text": systemPrompt},
	}
	if len(imageData) > 0 {
		if mimeType == "" {
			mimeType = "image/jpeg"
		}
		parts = append(parts, map[string]interface{}{
			"inline_data": map[string]string{
				"mime_type": mimeType,
				"data":      base64.StdEncoding.EncodeToString(imageData),
			},
		})
	}
	if description != nil && *description != "" {
		parts = append(parts, map[string]interface{}{"text": *description})
	}

	payload := map[string]interface{}{
		"contents": []map[string]interface{}{
			{"parts": parts},
		},
	}
	return p.makeRequest(url, payload)
}

func (p *GeminiProvider) makeRequest(url string, payload map[string]interface{}) (string, error) {
	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", url, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	var result map[string]interface{}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return "", fmt.Errorf("failed to parse Gemini response")
	}

	candidates, ok := result["candidates"].([]interface{})
	if !ok || len(candidates) == 0 {
		return "", fmt.Errorf("no candidates in Gemini response: %s", string(respBody))
	}
	content, _ := candidates[0].(map[string]interface{})["content"].(map[string]interface{})
	parts, _ := content["parts"].([]interface{})
	if len(parts) == 0 {
		return "", fmt.Errorf("no parts in Gemini response")
	}
	text, _ := parts[0].(map[string]interface{})["text"].(string)
	return text, nil
}

// ── HuggingFace ─────────────────────────────────────────────────────────────────

type HuggingFaceProvider struct {
	apiKey string
}

func (p *HuggingFaceProvider) Complete(systemPrompt, userContent string) (string, error) {
	payload := map[string]interface{}{
		"model": "Qwen/Qwen2.5-VL-72B-Instruct",
		"messages": []map[string]interface{}{
			{"role": "system", "content": systemPrompt},
			{"role": "user", "content": userContent},
		},
		"max_tokens": 1024,
	}
	return p.makeRequest(payload)
}

func (p *HuggingFaceProvider) CompleteVision(systemPrompt string, description *string, imageData []byte, mimeType string) (string, error) {
	userContent := []map[string]interface{}{}
	if len(imageData) > 0 {
		if mimeType == "" {
			mimeType = "image/jpeg"
		}
		userContent = append(userContent, map[string]interface{}{
			"type": "image_url",
			"image_url": map[string]string{
				"url": fmt.Sprintf("data:%s;base64,%s", mimeType, base64.StdEncoding.EncodeToString(imageData)),
			},
		})
	}
	if description != nil && *description != "" {
		userContent = append(userContent, map[string]interface{}{
			"type": "text", "text": *description,
		})
	}

	payload := map[string]interface{}{
		"model": "Qwen/Qwen2.5-VL-72B-Instruct",
		"messages": []map[string]interface{}{
			{"role": "system", "content": systemPrompt},
			{"role": "user", "content": userContent},
		},
		"max_tokens": 1024,
	}
	return p.makeRequest(payload)
}

func (p *HuggingFaceProvider) makeRequest(payload map[string]interface{}) (string, error) {
	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", "https://router.huggingface.co/v1/chat/completions", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+p.apiKey)

	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	var result map[string]interface{}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return "", fmt.Errorf("HuggingFace: failed to parse response")
	}

	choices, ok := result["choices"].([]interface{})
	if !ok || len(choices) == 0 {
		return "", fmt.Errorf("HuggingFace: no choices: %s", string(respBody))
	}
	msg, _ := choices[0].(map[string]interface{})["message"].(map[string]interface{})
	content, _ := msg["content"].(string)
	return content, nil
}
