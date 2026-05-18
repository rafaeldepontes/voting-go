package model

import "time"

type Poll struct {
	ID        string              `json:"id"`
	Text      string              `json:"text"`
	Options   []Option            `json:"options"`
	Duration  time.Duration       `json:"duration"`
	OwnerID   string              `json:"ownerID"`
	Votes     map[string]struct{} `json:"votes"`
	CreatedAt time.Time           `json:"createdAt"`
}

type PollDto struct {
	ID       string        `json:"id"`
	Text     string        `json:"text"`
	Duration time.Duration `json:"duration"`
}

type PollReq struct {
	Name     string        `json:"name"`
	Options  []string      `json:"options"`
	Duration time.Duration `json:"duration"`
}

type Option struct {
	Text  string `json:"text"`
	ID    int    `json:"id"`
	Votes int    `json:"votes"`
}
