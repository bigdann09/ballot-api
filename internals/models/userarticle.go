package models

import (
	"github.com/ballot/internals/database"
)

type UserArticle struct {
	UserID      uint   `json:"user_id"`
	ArticleSlug string `json:"article_slug"`
	Completed   bool   `json:"completed"`
}

func MarkArticleAsCompleted(userID uint, articleSlug string) error {
	result := database.DB.Model(&UserArticle{}).Create(&UserArticle{
		UserID:      userID,
		ArticleSlug: articleSlug,
		Completed:   true,
	})
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func CheckCompletedArticle(userID uint, slug string) bool {
	var userArticle UserArticle
	result := database.DB.Model(&UserArticle{}).Where("user_id =? AND article_slug =?", userID, slug).Scan(&userArticle)
	if result.Error != nil || (userArticle.UserID == 0 && userArticle.ArticleSlug == "") {
		return false
	}
	return true
}
