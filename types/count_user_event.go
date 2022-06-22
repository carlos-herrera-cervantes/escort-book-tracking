package types

type CountUserEvent struct {
	Accumulator int    `json:"accumulator"`
	Operation   string `json:"operation"`
	UserId      string `json:"userId"`
	UserType    string `json:"userType"`
}
