package postgres

import (
	"github.com/jinzhu/gorm"
	"github.com/ucladevx/BConnect-backend/models"
)

// ChatClient Postgresql client for chat table
type ChatClient struct {
	chat *gorm.DB
}

func (client *ChatClient) create() {
	client.chat.AutoMigrate(&models.Chats{})
}
