package services

import (
	"context"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/sashabaranov/go-openai"
	"github.com/smirnoffmg/partner/internal/repositories"
)

const retryTimeout = 5 * time.Second

type IAssistant interface {
	GetAnswer(chatID int64, question string) (string, error)
}

type ChatGPTService struct {
	chatsRepo repositories.IChatRepo
	client    *openai.Client
	assistant *openai.Assistant
}

func NewChatGPTService(repo repositories.IChatRepo, openAiKey string, openAiAssistantID string) (*ChatGPTService, error) {
	openaiConfig := openai.DefaultConfig(openAiKey)
	openaiConfig.AssistantVersion = "v2"

	openaiClient := openai.NewClientWithConfig(openaiConfig)

	assistant, err := openaiClient.RetrieveAssistant(context.Background(), openAiAssistantID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to retrieve assistant")

		return nil, err
	}

	return &ChatGPTService{
		chatsRepo: repo,
		client:    openaiClient,
		assistant: &assistant,
	}, nil
}

func (s *ChatGPTService) getThreadID(chatID int64) (string, error) {
	var threadID string

	chat, err := s.chatsRepo.GetOrCreate(chatID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get or create chat")

		return "", err
	}

	if chat.ThreadID == "" {
		thread, err := s.client.CreateThread(context.Background(), openai.ThreadRequest{})
		if err != nil {
			log.Error().Err(err).Msg("Failed to create thread")

			return "", err
		}

		err = s.chatsRepo.Update(chatID, map[string]interface{}{"thread_id": thread.ID})
		if err != nil {
			log.Error().Err(err).Msg("Failed to update chat")

			return "", err
		}

		threadID = thread.ID
	} else {
		threadID = chat.ThreadID
	}

	return threadID, nil
}

func (s *ChatGPTService) sendMessage(threadID string, question string) error {
	_, err := s.client.CreateMessage(context.Background(), threadID, openai.MessageRequest{
		Role:    "user",
		Content: question,
	})
	if err != nil {
		log.Error().Err(err).Msg("Failed to send message to ChatGPT")

		return err
	}

	return nil
}

func (s *ChatGPTService) runAndWait(threadID string) error {
	response, err := s.client.CreateRun(context.Background(), threadID, openai.RunRequest{
		AssistantID: s.assistant.ID,
	})
	if err != nil {
		log.Error().Err(err).Msg("Failed to create run")

		return err
	}

	// wait for completion
	for {
		run, err := s.client.RetrieveRun(context.Background(), threadID, response.ID)
		if err != nil {
			log.Error().Err(err).Msg("Failed to get run")

			return err
		}

		if run.Status == openai.RunStatusCompleted {
			break
		}

		if run.Status == openai.RunStatusFailed {
			log.Error().Msg(run.LastError.Message)

			return errChatGPTRunFailed
		}

		<-time.After(retryTimeout)
	}

	return nil
}

func (s *ChatGPTService) getLastMessageInThread(threadID string) (string, error) {
	limit := new(int)
	*limit = 1

	messages, err := s.client.ListMessage(context.Background(), threadID, limit, nil, nil, nil)
	if err != nil {
		log.Error().Err(err).Msg("Failed to list messages")

		return "", err
	}

	if len(messages.Messages) == 0 {
		log.Error().Msg("No messages in response")

		return "Sorry, I have nothing to response", nil
	}

	// get assistant response

	var assistantResponse []string

	for _, message := range messages.Messages {
		if message.Role == "assistant" {
			for _, content := range message.Content {
				if content.Type == "text" {
					assistantResponse = append(assistantResponse, content.Text.Value)
				}
			}
		}
	}

	return strings.Join(assistantResponse, "\n"), nil
}

func (s *ChatGPTService) GetAnswer(chatID int64, question string) (string, error) {
	// ensure we have thread
	threadID, err := s.getThreadID(chatID)
	if err != nil {
		log.Error().Err(err).Msg("Cannot get thread ID")

		return "", err
	}

	// add question to chat
	if err = s.sendMessage(threadID, question); err != nil {
		return "", err
	}

	// run assistant on thread
	if err = s.runAndWait(threadID); err != nil {
		return "", err
	}

	// get messages
	return s.getLastMessageInThread(threadID)
}
