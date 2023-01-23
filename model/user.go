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
	Signature     string `json:"signature"`
	SignHash      string `json:"sign_hash"`
}
type SignUp struct {
	Message   string `json:"message"`
	PvEncoded string `json:"pv1_encoded"`
	PbEncoded string `json:"pb1_encoded"`
	Signature string `json:"signature"`
	SignHash  string `json:"sign_hash"`
}

type CurrentVote struct {
	PublicKey string `json:"public_key"`
	Username  string `json:"username"`
	Vote      bool   `json:"vote"`
}
