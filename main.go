package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	preprompt := "Here is a transcript from an ongoing conversation. Just answer as if you were to continue the converstaion."
	jane := "You are one of two radio hosts, your name is Jane and you co-host with Joe. You are always trying to crack jokes and teas Joe for the shitty commercial radio gig, you are more of a happy go lucky person that just find the whole crapiness of the job fun. Don't prefix the answer with your name."
	joe := "You are one of two radio hosts, your name is Joe and you co-host with Jane. You are always a bit annoyed ad Janes attempts to joke all the time. You dream of being a real jounalist at New York Times but all you got was this lousy commercial radio gig. Don't prefix the answer with your name."

	resp := ""
	memory := ""
	for i := 0; i < 20; i++ {
		// The path to the file you want to write/append to
		filePath := fmt.Sprintf("content/example%d.txt", i)

		for j := 0; j < 10; j++ {
			time.Sleep(2 * time.Second)
			if j%2 == 0 {
				if memory == "" {
					resp = "Heeeeeeeello world!"
				} else {
					resp = callAPI(jane, fmt.Sprintf("%s %s", preprompt, memory))
				}
				memory = fmt.Sprintf("%s \n %s %s \n\n", memory, "Jane", resp)
			} else {
				resp = callAPI(joe, fmt.Sprintf("%s %s", preprompt, memory))
				memory = fmt.Sprintf("%s \n %s %s \n\n", memory, "Joe", resp)
			}
			fmt.Println(resp + "\n\n")

			// Appending to the file (creates the file if it doesn't exist)
			file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				panic(err)
			}
			defer file.Close()

			if _, err := file.WriteString(fmt.Sprintf("%s\n\n", resp)); err != nil {
				panic(err)
			}

		}
	}

}

func callAPI(systemContent string, userContent string) string {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	// Access the API key stored in the .env file
	openAIKey := os.Getenv("OPENAI_API_KEY")

	apiURL := "https://api.openai.com/v1/chat/completions"

	// Create the payload
	payload := Payload{
		Model: "gpt-4-turbo-preview",
		Messages: []struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		}{
			{Role: "system", Content: systemContent},
			{Role: "user", Content: userContent},
		},
	}

	// Marshal the payload to JSON
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		panic(err)
	}

	// Create a new request using http
	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonPayload))
	if err != nil {
		panic(err)
	}

	// Add headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+openAIKey)

	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	// Unmarshal the JSON response into the ApiResponse struct
	var apiResponse ApiResponse
	if err := json.Unmarshal(body, &apiResponse); err != nil {
		panic(err)
	}

	// Check if there are any choices and messages before accessing
	if len(apiResponse.Choices) > 0 && len(apiResponse.Choices[0].Message.Content) > 0 {
		messageContent := apiResponse.Choices[0].Message.Content
		// fmt.Println("Extracted Message Content:", messageContent)
		return messageContent
	}

	// Return a default or error message if no content is available
	return "No message content available"
}

type ApiResponse struct {
	Choices []struct {
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

// Struct to match the JSON payload structure
type Payload struct {
	Model    string `json:"model"`
	Messages []struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	} `json:"messages"`
}
