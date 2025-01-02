package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/invopop/jsonschema"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

// GenerateSchema generates a JSON schema for a given generic type T
func GenerateSchema[T any]() interface{} {
	reflector := jsonschema.Reflector{
		AllowAdditionalProperties: false,
		DoNotReference:            true,
	}
	var v T
	schema := reflector.Reflect(v)
	return schema
}

// GetFiltersAndLabelsFromAI interacts with the OpenAI API to generate Gmail filters and labels based on email samples.
func GetFiltersAndLabelsFromAI(openAIKey, openAIHost, openAIModel string, messages Messages) (*Config, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	// Initialize the OpenAI client
	client := openai.NewClient(
		option.WithAPIKey(openAIKey),
		option.WithBaseURL(openAIHost),
	)

	// Convert email data to JSON
	emailDataJSON, err := json.Marshal(messages)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal email data: %v", err)
	}

	// Define the schema for structured outputs
	schemaParam := openai.ResponseFormatJSONSchemaJSONSchemaParam{
		Name:        openai.F("filters_and_labels"),
		Description: openai.F("Generate Gmail filters and labels based on email samples"),
		Schema:      openai.F(GenerateSchema[Config]()),
		Strict:      openai.Bool(true),
	}

	// Define the chat completion request
	message := openai.ChatCompletionNewParams{
		Messages: openai.F([]openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(fmt.Sprintf("Based on the following email data, suggest filters and labels:\n\n%s", emailDataJSON)),
		}),
		Model: openai.F(openAIModel),
		ResponseFormat: openai.F[openai.ChatCompletionNewParamsResponseFormatUnion](
			openai.ResponseFormatJSONSchemaParam{
				Type:       openai.F(openai.ResponseFormatJSONSchemaTypeJSONSchema),
				JSONSchema: openai.F(schemaParam),
			},
		),
	}

	// Call the OpenAI Chat Completions API
	chat, err := client.Chat.Completions.New(ctx, message)
	if err != nil {
		return nil, fmt.Errorf("failed to get OpenAI response: %v", err)
	}

	// Parse the structured response into the Config object
	var response Config
	err = json.Unmarshal([]byte(chat.Choices[0].Message.Content), &response)
	if err != nil {
		return nil, fmt.Errorf("failed to parse OpenAI response: %v", err)
	}

	return &response, nil
}
