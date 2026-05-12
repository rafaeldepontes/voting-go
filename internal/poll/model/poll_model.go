package model

type Poll struct {
	ID      string   `json:"id"`
	Text    string   `json:"text"`
	Options []Option `json:"options"`
}

type PollReq struct {
	Name    string   `json:"name"`
	Options []string `json:"options"`
}

type Option struct {
	Text  string `json:"text"`
	ID    int    `json:"id"`
	Votes int    `json:"votes"`
}

type VoteReq struct {
	OptionID int `json:"optionId"`
}
