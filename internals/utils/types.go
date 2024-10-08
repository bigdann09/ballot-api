package utils

import "time"

type UserAPI struct {
	ID            uint   `json:"id,omitempty"`
	TGID          int64  `json:"tg_id,omitempty"`
	Username      string `json:"username,omitempty"`
	WalletAddress string `json:"wallet_address,omitempty"`
	TGPremium     string `json:"tg_premium"`
	Token         string `json:"token,omitempty"`
	ReferralPoint uint64 `json:"referral_point"`
	TaskPoint     uint64 `json:"task_point"`
}

type PointAPI struct {
	ID            uint   `json:"id"`
	UserID        uint   `json:"user_id,omitempty"`
	ReferralPoint uint64 `json:"referral_point"`
	TaskPoint     uint64 `json:"task_point"`
}

type TaskAPI struct {
	ID        uint   `json:"id"`
	UUID      string `json:"uuid"`
	Name      string `json:"name"`
	Type      string `json:"type"`
	Link      string `json:"link,omitempty"`
	Point     int64  `json:"point"`
	Completed bool   `json:"completed"`
	Validate  bool   `json:"validate"`
	Duration  string `json:"duration"`
}

type CandidateAPI struct {
	ID         uint      `json:"id"`
	Name       string    `json:"name"`
	Votes      uint      `json:"votes"`
	LastVoteAt time.Time `json:"last_vote_at"`
}

type VoteAPI struct {
	ID          uint      `json:"id"`
	UserID      uint      `json:"user_id"`
	CandidateID uint      `json:"candidate_id"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
}

type ReferralAPI struct {
	ID       uint `json:"id"`
	Referrer uint `json:"referrer"`
	Referee  uint `json:"referee"`
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
	TGID          int64  `json:"tg_id"`
	Username      string `json:"username"`
	TGPremium     bool   `json:"tg_premium"`
	WalletAddress string `json:"wallet_address"`
}

type ReferralCreateApiRequest struct {
	Username  string `json:"username"`
	Token     string `json:"token"`
	TGID      int64  `json:"tg_id"`
	TGPremium bool   `json:"tg_premium"`
}

type TaskCreateApiRequest struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Link     string `json:"link,omitempty"`
	Point    int64  `json:"point"`
	Validate bool   `json:"validate"`
	Duration string `json:"duration"`
}

type MakeVote struct {
	Candidate string `json:"candidate"`
}
