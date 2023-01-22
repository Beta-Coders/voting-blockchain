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

type Vote struct {
	Username      string `json:"username"`
	CandidateName string `json:"candidate_name"`
	Pubkey        string `json:"pubkey"`
}
