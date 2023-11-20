package model

type HomeResponse struct {
	Items []Item `json:"items"`
	// 给其他响应预留
}

type ChannelResp struct {
	Streams []StreamProp `json:"streams"`
}

type Item struct {
	Lab
	Team
	SubStatus  string       `json:"sub_status"`
	IsNotStart bool         `json:"is_not_start"`
	Zhibos     []ZhiboProp  `json:"zhibos"`
}

type ZhiboProp struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type LiveResponse struct {
	Team    Team         `json:"team"`
	Lab     Lab          `json:"lab"`
	Streams []StreamProp `json:"streams"`
}

type StreamProp struct {
	M3U8URL string `json:"m3u8_url"`
	Name    string `json:"name"`
}

type Lab struct {
	LabEvent        string `json:"lab_event"`
	LabEventBGColor string `json:"lab_event_bg_color"`
	LabBC           string `json:"lab_bc"`
	LabJQ           string `json:"lab_jq"`
	LabTime         string `json:"lab_time"`
	BF              string `json:"bf"`
	ID              string `json:"id"`
	Kind            Kind   `json:"kind"`
	Hot             bool   `json:"hot"`
}

type Team struct {
	TeamHome     string `json:"team_home"`
	TeamHomeIcon string `json:"team_home_icon"`
	TeamAway     string `json:"team_away"`
	TeamAwayIcon string `json:"team_away_icon"`
	Bf string `json:"bf"`
}

type Kind string

const (
	All        Kind = "all"
	Basketball Kind = "basketball"
	Football   Kind = "football"
	Other      Kind = "other"
	ChannelKey = "channel:%v"
	StreamKey = "stream:%v"
	LiveKey = "live:%v"
)
