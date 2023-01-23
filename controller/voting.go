package controller

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"votingblockchain/container"
	"votingblockchain/model"
)

type VotingController struct {
	container container.Container
}

func NewVotingController(container container.Container) *VotingController {
	return &VotingController{container: container}
}

// GetVoteByUsername - get voting by username
func (controller *VotingController) GetVoteByUsername(c echo.Context) error {
	return nil
}

// Vote - create voting
func (controller *VotingController) Vote(c echo.Context) (err error) {
	vote := new(model.Vote)
	curVote := new(model.CurrentVote)
	if err = c.Bind(vote); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	bc := controller.container.GetBC()
	rows, err := controller.container.GetDB().Query("SELECT * FROM admins WHERE username=$1", vote.Username)
	if err != nil {
		fmt.Println(err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())

	}
	defer rows.Close()
	if !rows.Next() {
		return echo.NewHTTPError(http.StatusBadRequest, model.Response{Message: "User not found"})
	}
	rows.Scan(curVote.PublicKey, curVote.Username, curVote.Vote)
	if curVote.Vote {
		return echo.NewHTTPError(http.StatusBadRequest, model.Response{Message: "User already voted"})
	}
	if curVote.PublicKey != vote.Pubkey {
		return echo.NewHTTPError(http.StatusBadRequest, model.Response{Message: "Invalid public key"})
	}
	bc.AddBlock(vote.CandidateName, vote.Username, decodePubKey(vote.Pubkey), []byte(vote.Signature), []byte(vote.SignHash))
	_, err = controller.container.GetDB().Exec("UPDATE admins SET vote=$1 WHERE username=$2", true, vote.Username)
	return c.JSON(http.StatusOK, model.Response{Message: "success"})
}

// UpdateVoting - update voting
func (controller *VotingController) UpdateVoting(c echo.Context) (err error) {
	return nil
}

// GetVotingResults - get voting results
func (controller *VotingController) GetVotingResults(c echo.Context) error {
	bc := controller.container.GetBC()
	it := bc.Iterator()
	votes := make([]model.Vote, 0)

	for {
		block := it.Next()
		vote := new(model.Vote)
		vote.CandidateName = string(block.CandidateName)
		vote.Username = string(block.Username)
		votes = append(votes, *vote)
		if len(block.PrevHash) == 0 {
			break
		}
	}
	votes = votes[0 : len(votes)-1]
	return c.JSON(http.StatusOK, votes)
}

// GetCandidates - get candidate list
func (controller *VotingController) GetCandidates(c echo.Context) error {
	rows, err := controller.container.GetDB().Query("select * from candidates")
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	defer rows.Close()
	candidates := make([]model.Candidate, 0)
	for rows.Next() {
		candidate := model.Candidate{}
		err := rows.Scan(&candidate.Id, &candidate.CandidateName, &candidate.PartyName)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		candidates = append(candidates, candidate)
	}
	return c.JSON(http.StatusOK, candidates)

}

func (controller *VotingController) AddCandidate(c echo.Context) error {
	candidate := &model.Candidate{}
	if err := c.Bind(candidate); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	rows, err := controller.container.GetDB().Query("select * from candidates where party_name=$1", candidate.PartyName)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	defer rows.Close()
	for rows.Next() {
		return echo.NewHTTPError(http.StatusBadRequest, model.Response{Message: "Candidate already exists"})
	}
	_, err = controller.container.GetDB().Exec("insert into candidates(candidate_name,party_name) values($1,$2)", candidate.CandidateName, candidate.PartyName)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, model.Response{Message: "success"})
}

func (controller *VotingController) EndVoting(c echo.Context) error {
	return nil
}

func decodePubKey(pemEncodedPub string) *ecdsa.PublicKey {
	blockPub, _ := pem.Decode([]byte(pemEncodedPub))
	x509EncodedPub := blockPub.Bytes
	genericPublicKey, _ := x509.ParsePKIXPublicKey(x509EncodedPub)
	publicKey := genericPublicKey.(*ecdsa.PublicKey)
	return publicKey
}
