package service

type Info struct {
	Name    string    `json:"name"`
	Service []Service `json:"children"`
}

type Service struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	IsTelemedicine bool   `json:"isTelemedicine"`
}
