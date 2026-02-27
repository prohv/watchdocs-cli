package resolver

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"google.golang.org/genai"
)

type DocResult struct {
	Name   string `json:"name"`
	DocURL string `json:"docUrl"`
}

func ResolveDocs(deps []string) ([]DocResult, error) {
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("GEMINI_API_KEY environment variable not set")
	}

	ctx := context.Background()
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  apiKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return nil, err
	}

	prompt := fmt.Sprintf("Return a JSON array of objects with 'name' and 'docUrl' for these librariesl, make sure they are original docs: %v. Only return the JSON, no markdown code blocks.", deps)

	resp, err := client.Models.GenerateContent(ctx, "gemini-3-flash-preview", genai.Text(prompt), nil)
	if err != nil {
		return nil, err
	}

	var results []DocResult
	err = json.Unmarshal([]byte(resp.Candidates[0].Content.Parts[0].Text), &results)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Gemini response: %v", err)
	}

	return results, nil
}