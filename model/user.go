package model

type User struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type Response struct {
	// example: error message
	Message string `json:"msg"`
}

type Candidate struct {
	Id            int    `json:"id"`
	CandidateName string `json:"candidate_name"`
	PartyName     string `json:"party_name"`
}
