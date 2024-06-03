package services

import (
	"context"
	"strings"
	"time"

	"github.com/smirnoffmg/partner/internal/entities"
	"github.com/smirnoffmg/partner/internal/repositories"

	"github.com/rs/zerolog/log"

	"gorm.io/gorm"

	"github.com/sashabaranov/go-openai"
)

type ChatGPTService struct {
	ChatsRepo   *repositories.ChatsRepo
	client      *openai.Client
	assistant   openai.Assistant
	partnerName string
}

func NewChatGPTService(repo *repositories.ChatsRepo, openAiKey string, openAiAssistantId string, partnerName string) (*ChatGPTService, error) {
	openaiConfig := openai.DefaultConfig(openAiKey)
	openaiConfig.AssistantVersion = "v2"

	openaiClient := openai.NewClientWithConfig(openaiConfig)

	assistant, err := openaiClient.RetrieveAssistant(context.Background(), openAiAssistantId)

	if err != nil {
		log.Error().Err(err).Msg("Failed to retrieve assistant")
		return nil, err
	}

	return &ChatGPTService{
		ChatsRepo: repo,
		client:    openaiClient,
		assistant: assistant,
	}, nil
}

func (s *ChatGPTService) GetOrCreateChat(ctx context.Context, telegramChatID int64) (*entities.Chat, error) {
	chat, err := s.ChatsRepo.FindByTelegramChatID(telegramChatID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			newThread, err := s.client.CreateThread(ctx, openai.ThreadRequest{})

			if err != nil {
				log.Error().Err(err).Msg("Failed to create thread")
				return nil, err
			}

			chat = &entities.Chat{
				PartnerName:    s.partnerName,
				TelegramChatID: telegramChatID,
				ThreadID:       newThread.ID,
			}

			if err := s.ChatsRepo.Create(chat); err != nil {
				log.Error().Err(err).Msg("Failed to create chat")
				return nil, err
			}
		} else {
			log.Error().Err(err).Msg("Failed to find chat")
			return nil, err
		}
	}

	return chat, nil
}

func (s *ChatGPTService) GetAnswer(telegramChatID int64, question string) (string, error) {
	chat, err := s.GetOrCreateChat(context.Background(), telegramChatID)

	if err != nil {
		log.Error().Err(err).Msg("Failed to get or create chat")
		return "", err
	}

	// add question to chat

	_, err = s.client.CreateMessage(context.Background(), chat.ThreadID, openai.MessageRequest{
		Role:    "user",
		Content: question,
	})

	if err != nil {
		log.Error().Err(err).Msg("Failed to create message")
		return "", err
	}

	// run assistant on thread

	response, err := s.client.CreateRun(context.Background(), chat.ThreadID, openai.RunRequest{
		AssistantID: s.assistant.ID,
	})

	if err != nil {
		log.Error().Err(err).Msg("Failed to create run")
		return "", err
	}

	// wait for completion

	for {
		run, err := s.client.RetrieveRun(context.Background(), chat.ThreadID, response.ID)

		if err != nil {
			log.Error().Err(err).Msg("Failed to get run")
			return "", err
		}

		log.Debug().Msgf("Run status: %s", run.Status)

		if run.Status == openai.RunStatusCompleted {
			break
		}

		if run.Status == openai.RunStatusFailed {
			log.Error().Msg(run.LastError.Message)
			return "Sorry, cannot answer right now. Please ask later", nil
		}
		<-time.After(5 * time.Second)
	}

	// get messages

	var limit *int = new(int)
	*limit = 1

	messages, err := s.client.ListMessage(context.Background(), chat.ThreadID, limit, nil, nil, nil)

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

	// update last message id

	return strings.Join(assistantResponse, "\n"), nil

}
