package inout

type SubscribeIn struct {
	Channels []string `json:"channels"`
}

type SubscribeOut []string
