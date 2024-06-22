package parser

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/khomart/instagram_recipe_parser/internal/config"
	"github.com/sashabaranov/go-openai"
	"github.com/xfrr/goffmpeg/transcoder"
)

type Parser struct {
	openAIKey string
}

func NewParser(config *config.Config) *Parser {
	return &Parser{
		openAIKey: config.OpenAIApiKey,
	}
}

// analyzeText analyzes text content using the OpenAI API.
func (p *Parser) analyzeText(content string, client *openai.Client) (string, error) {
	resp, err := client.CreateChatCompletion(context.TODO(), openai.ChatCompletionRequest{
		Model: "gpt-4-turbo",
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleUser,
				Content: fmt.Sprintf("Summarize the following content for a recipe instruction:\n\n%s", content),
			},
		},
		MaxTokens: 100,
	})
	if err != nil {
		return "", fmt.Errorf("failed to summarize text content: %v", err)
	}
	return resp.Choices[0].Message.Content, nil
}

// analyzeImage analyzes an image content using the OpenAI API.
func (p *Parser) analyzeImage(filePath string, client *openai.Client) (string, error) {
	// Read the image file
	imgData, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read image: %s", filePath)
	}

	// Encode the image to base64
	base64Img := base64.StdEncoding.EncodeToString(imgData)

	// Use OpenAI API to analyze the image
	resp, err := client.CreateChatCompletion(context.TODO(), openai.ChatCompletionRequest{
		Model: "gpt-4-vision-preview",
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleUser,
				Content: "Analyze the following image and summarize it as a recipe instruction",
			},
			{
				Role: openai.ChatMessageRoleUser,
				MultiContent: []openai.ChatMessagePart{
					{
						ImageURL: &openai.ChatMessageImageURL{
							URL: fmt.Sprintf("data:image/jpeg;base64,%s", base64Img),
						},
					},
				},
			},
		},
	})
	if err != nil {
		return "", fmt.Errorf("failed to analyze image content: %v", err)
	}

	return resp.Choices[0].Message.Content, nil
}

// analyzeAudio analyzes audio content using the OpenAI API.
func (p *Parser) analyzeAudio(filePath string, client *openai.Client) (string, error) {
	req := openai.AudioRequest{
		Model:    openai.Whisper1,
		FilePath: filePath,
		Prompt:   "Analyze the audio content for a recipe instruction",
	}

	// Use OpenAI API to analyze audio content
	resp, err := client.CreateTranscription(context.TODO(), req)
	if err != nil {
		return "", fmt.Errorf("failed to analyze audio content: %v", err)
	}
	return resp.Text, nil
}

// convertVideoToAudio converts a video file to an audio file using FFmpeg.
func (p *Parser) convertVideoToAudio(videoPath, audioPath string) error {
	trans := new(transcoder.Transcoder)
	err := trans.Initialize(videoPath, audioPath)
	if err != nil {
		return fmt.Errorf("failed to initialize transcoder: %v", err)
	}

	done := trans.Run(false)
	err = <-done
	if err != nil {
		return fmt.Errorf("failed to convert video to audio: %v", err)
	}
	return nil
}

// summarizeContent uses OpenAI API to summarize the content of files in the specified folder.
func (p *Parser) SummarizeContent(folderPath string) (string, error) {
	// Initialize OpenAI client
	client := openai.NewClient(p.openAIKey)

	// Read all files in the folder
	files, err := os.ReadDir(folderPath)
	if err != nil {
		return "", fmt.Errorf("failed to read folder: %v", err)
	}

	var combinedContent strings.Builder

	for _, file := range files {
		if !file.IsDir() {
			filePath := filepath.Join(folderPath, file.Name())
			ext := strings.ToLower(filepath.Ext(file.Name()))

			var summary string
			switch ext {
			case ".txt":
				content, err := os.ReadFile(filePath)
				if err != nil {
					return "", fmt.Errorf("failed to read text file %s: %v", filePath, err)
				}
				summary, err = p.analyzeText(string(content), client)
				if err != nil {
					return "", fmt.Errorf("failed to analyze text: %v", err)
				}
			case ".jpg", ".jpeg", ".png":
				summary, err = p.analyzeImage(filePath, client)
				if err != nil {
					return "", fmt.Errorf("failed to analyze image: %v", err)
				}
			case ".mp4", ".avi", ".mkv":
				audioPath := strings.TrimSuffix(filePath, ext) + ".mp3"
				if err := p.convertVideoToAudio(filePath, audioPath); err != nil {
					return "", fmt.Errorf("failed to convert video to audio: %v", err)
				}
				summary, err = p.analyzeAudio(audioPath, client)
				if err != nil {
					return "", fmt.Errorf("failed to analyze audio: %v", err)
				}
			default:
				continue
			}

			combinedContent.WriteString(summary)
			combinedContent.WriteString("\n")
		}
	}

	// Combine and summarize the content into a short recipe instruction
	finalSummary, err := client.CreateChatCompletion(context.TODO(), openai.ChatCompletionRequest{
		Model: "gpt-4-turbo",
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleUser,
				Content: fmt.Sprintf("Combine the following content into a short recipe instruction. Give short 1-2 sentence description, separate list of ingridients and steps how to prepare the dish. \n\n%s", combinedContent.String()),
			},
		},
		MaxTokens: 1000,
	})
	if err != nil {
		return "", fmt.Errorf("failed to create recipe instruction: %v", err)
	}

	return finalSummary.Choices[0].Message.Content, nil
}
