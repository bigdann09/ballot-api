package models

import "github.com/ballot/internals/database"

type UserArticle struct {
	UserID      uint   `json:"user_id"`
	ArticleSlug string `json:"article_slug"`
	Completed   bool   `json:"completed"`
}

func MarkArticleAsCompleted(userID uint, articleSlug string) error {
	result := database.DB.Model(&UserArticle{}).Create(&UserArticle{
        UserID: userID,
        ArticleSlug: articleSlug,
        Completed: true,
    })
    if result.Error!= nil {
        return result.Error
    }
    return nil
}

func GetAllCompletedArticle(userID uint) ([]*UserArticle, error) {
    var userArticles []*UserArticle
    result := database.DB.Model(&UserArticle{}).Scan(&userArticles)
    if result.Error != nil {
        return userArticles, result.Error
    }
    return userArticles, nil
}
