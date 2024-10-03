package utils

import "time"

type UserAPI struct {
	ID            uint      `json:"id,omitempty"`
	TGID          int64     `json:"tg_id,omitempty"`
	Username      string    `json:"username,omitempty"`
	WalletAddress string    `json:"wallet_address,omitempty"`
	TGPremium     string    `json:"tg_premium"`
	Token         string    `json:"token,omitempty"`
	ReferralPoint uint64    `json:"referral_point"`
	TaskPoint     uint64    `json:"task_point"`
}

type PointAPI struct {
	ID            uint   `json:"id"`
	UserID        uint64 `json:"user_id,omitempty"`
	ReferralPoint uint64 `json:"referral_point"`
	TaskPoint     uint64 `json:"task_point"`
}

type TaskAPI struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	Type      string `json:"type"`
	Link      string `json:"link,omitempty"`
	Point     int64  `json:"point"`
	Completed bool   `json:"completed"`
}

// news struct
type Source struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Article struct {
	ID          uint   `json:"id"`
	Source      Source `json:"source"`
	Author      string `json:"author"`
	Title       string `json:"title"`
	Slug        string `json:"slug"`
	Description string `json:"description"`
	URL         string `json:"url"`
	URLToImage  string `json:"urlToImage"`
	PublishedAt string `json:"publishedAt"`
	Content     string `json:"content"`
}

type Data struct {
	Status       string    `json:"status"`
	TotalResults uint      `json:"totalResults"`
	Articles     []Article `json:"articles"`
}

type Response struct {
	RefereshedAt time.Time `json:"refereshedAt"`
	Articles     []Article `json:"articles"`
}

// Request structs

type NewUser struct {

}

type ReferralCreateApiRequest struct {
	Token     string `json:"token"`
	TGID      int64  `json:"tg_id"`
	TGPremium bool   `json:"tg_premium"`
}

type TaskCreateApiRequest struct {
	Name  string `json:"name"`
	Type  string `json:"type"`
	Link  string `json:"link,omitempty"`
	Point int64  `json:"point"`
}