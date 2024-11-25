package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/google/generative-ai-go/genai"
	"github.com/joho/godotenv"
	"google.golang.org/api/option"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

type prompt struct {
	Userprompt string `json:"prompt"`
}

func getgemini(apiKey, userprompt string) (string, error) {
	promptTemplate := `You are Peter Griffin from Family Guy. Respond to the user's message in Peter's voice and style.
	Include his characteristic laugh "hehehe" and make references to Family Guy episodes when relevant.
	Be silly, sometimes inappropriate (but not offensive), and give advice or responses that are humorously questionable
	- just like Peter would. Make sure to stay in character at all times.
	
	User's message: {{prompt}}
	
	Remember to:
	- Use Peter's speaking style and mannerisms
	- Include his laugh "hehehe" naturally in the response
	- Make references to Family Guy episodes or characters when relevant
	- Give advice that's humorous and questionably logical
	- Stay true to Peter's character`

	prompt := strings.ReplaceAll(promptTemplate, "{{prompt}}", userprompt)
	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return "", fmt.Errorf("failed to create client: %w", err)
	}
	defer client.Close()

	model := client.GenerativeModel("gemini-1.5-flash")
	resp, err := model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return "", fmt.Errorf("failed to generate: %w", err)
	}

	if len(resp.Candidates) == 0 || resp.Candidates[0].Content == nil {
		return "", fmt.Errorf("no response generated")
	}

	var parts []string
	for _, part := range resp.Candidates[0].Content.Parts {
		parts = append(parts, fmt.Sprintf("%v", part))
	}

	result := strings.Join(parts, "\n")
	return result, nil
}

func prompthandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method == http.MethodPost {
		var promptdata prompt

		err := json.NewDecoder(r.Body).Decode(&promptdata)
		if err != nil {
			fmt.Println("Error decoding JSON:", err)
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		if promptdata.Userprompt == "" {
			http.Error(w, "Prompt cannot be empty", http.StatusBadRequest)
			return
		}

		apiKey := os.Getenv("API_KEY")
		if apiKey == "" {
			log.Fatal("API_KEY environment variable not set")
			return
		}
		result, err := getgemini(apiKey, promptdata.Userprompt)
		if err != nil {
			fmt.Println("Error generating response:", err)
			http.Error(w, "Error generating response", http.StatusInternalServerError)
			return
		}

		response := map[string]string{
			"result": result,
		}

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			fmt.Println("Error encoding response:", err)
			http.Error(w, "Error encoding response", http.StatusInternalServerError)
			return
		}
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func main() {
	http.HandleFunc("/prompt", prompthandler)
	fmt.Println("Server running on http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("Error starting the server:", err)
	}
}
