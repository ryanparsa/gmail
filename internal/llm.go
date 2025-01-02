package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/invopop/jsonschema"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"time"
)

// Generate schema for the response
func GenerateSchema[T any]() interface{} {
	reflector := jsonschema.Reflector{
		AllowAdditionalProperties: false,
		DoNotReference:            true,
	}
	var v T
	schema := reflector.Reflect(v)
	return schema
}

func GetFiltersAndLabelsFromAI(openAIKey, openAIHost, openAIModel string, m Messages) (*Config, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	// Initialize OpenAI client
	client := openai.NewClient(
		option.WithAPIKey(openAIKey),
		option.WithBaseURL(openAIHost),
	)

	// Prepare email data as a JSON string
	emailDataJSON, err := json.Marshal(m)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal email data: %v", err)
	}

	// Create schema param for structured outputs
	schemaParam := openai.ResponseFormatJSONSchemaJSONSchemaParam{
		Name:        openai.F("filters_and_labels"),
		Description: openai.F("Generate filters and labels based on email samples"),
		Schema:      openai.F(GenerateSchema[Config]()),
		Strict:      openai.Bool(true),
	}

	// Create chat completion parameters
	message := openai.ChatCompletionNewParams{
		Messages: openai.F([]openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(fmt.Sprintf("Based on the following email data, suggest filters and labels:\n\n%s", emailDataJSON)),
		}),
		Model: openai.F(openai.ChatModelGPT4o),
		ResponseFormat: openai.F[openai.ChatCompletionNewParamsResponseFormatUnion](
			openai.ResponseFormatJSONSchemaParam{
				Type:       openai.F(openai.ResponseFormatJSONSchemaTypeJSONSchema),
				JSONSchema: openai.F(schemaParam),
			},
		),
	}

	// Call OpenAI API
	chat, err := client.Chat.Completions.New(ctx, message)
	if err != nil {
		return nil, fmt.Errorf("failed to get OpenAI response: %v", err)
	}

	// Parse structured JSON output
	var response Config
	err = json.Unmarshal([]byte(chat.Choices[0].Message.Content), &response)
	if err != nil {
		return nil, fmt.Errorf("failed to parse OpenAI response: %v", err)
	}

	return &response, nil
}
