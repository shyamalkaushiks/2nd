package gemini

import (
	// "encoding/json"
	"context"
	"fmt"
	"log"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

// GenerateData :
func GenerateData() {

	apiKey := ""

	ctx := context.Background()
	// Access your API key as an environment variable (see "Set up your API key" above)
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	// For text-only input, use the gemini-pro model
	model := client.GenerativeModel("gemini-pro")
	resp, err := model.GenerateContent(ctx, genai.Text("Who is Burger King ?"))
	if err != nil {
		log.Fatal(err)
	}

	// var response map[string]interface{}
	// err = json.Unmarshal(resp, &response)
	// if err != nil {
	// 	fmt.Println("Error unmarshalling response:", err)
	// 	return
	// }

	// fmt.Println(resp.PromptFeedback)
	for _, val := range resp.Candidates {
		// fmt.Println(val.Content)
		for _, content := range val.Content.Parts {
			fmt.Println(content)
		}
	}

}

// GenerateMultiTurnData :
func GenerateMultiTurnData(content, keyword string) string {
	apiKey := "AIzaSyBoerDN0e3L-yYuSItTVH9S7Pi37jwX_c8"

	ctx := context.Background()
	// Access your API key as an environment variable (see "Set up your API key" above)
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	// For text-only input, use the gemini-pro model
	model := client.GenerativeModel("gemini-pro")
	// Initialize the chat
	cs := model.StartChat()
	cs.History = []*genai.Content{
		&genai.Content{
			Parts: []genai.Part{
				genai.Text(content),
			},
			Role: "user",
		},
		&genai.Content{
			Parts: []genai.Part{
				genai.Text("What would you like to know?"),
			},
			Role: "model",
		},
	}

	resp, err := cs.SendMessage(ctx, genai.Text(keyword))
	if err != nil {
		log.Fatal(err)
	}

	// for _, val := range resp.Candidates {
	// 	// fmt.Println(genai.Text(val.Content.Parts))
	// 	for _, content := range val.Content.Parts {
	// 	// 	fmt.Println(content)
	// 	fmt.Println(genai.Text(content))
	// 	// 	contentStr = append(contentStr, "+", genai.Text(content))
	// 	}
	// }

	return fmt.Sprintf("generated response: %s\n", resp.Candidates[0].Content.Parts[0])
}
