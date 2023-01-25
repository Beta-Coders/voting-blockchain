package model

type User struct {
	Id         int    `json:"id"`
	Username   string `json:"username"`
	Password   string `json:"password"`
	PrivateKey string `json:"private_key"`
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
	PrivateKey    string `json:"private_key"`
}

type CurrentVote struct {
	PublicKey string `json:"public_key"`
	Username  string `json:"username"`
	Vote      bool   `json:"vote"`
}

type SignUpResponse struct {
	Username   string `json:"username"`
	PrivateKey string `json:"private_key"`
}

type VoteResponse struct {
	Username      string `json:"username"`
	CandidateName string `json:"candidate_name"`
}
