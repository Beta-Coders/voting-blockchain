package controller

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"strings"
	"votingblockchain/ECC"
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
	username := c.Param("username")
	votes := getVotes(controller)
	for _, vote := range votes {
		if strings.Compare(vote.Username, username) == 0 {
			return c.JSON(http.StatusOK, vote)
		}
	}
	return c.JSON(http.StatusOK, model.Response{Message: "vote not found"})
}

// Vote - create voting
func (controller *VotingController) Vote(c echo.Context) (err error) {
	vote := new(model.Vote)
	userVote := new(model.CurrentVote)
	if err = c.Bind(vote); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	fmt.Println(vote)
	bc := controller.container.GetBC()
	if !bc.CheckVotingInProgress() {
		return echo.NewHTTPError(http.StatusBadRequest, model.Response{Message: "Voting is not in progress"})
	}
	rows, err := controller.container.GetDB().Query("SELECT * FROM admins WHERE username=$1", vote.Username)
	if err != nil {
		fmt.Println(err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())

	}
	defer rows.Close()
	if !rows.Next() {
		return echo.NewHTTPError(http.StatusBadRequest, model.Response{Message: "User not found"})
	}
	rows.Scan(&userVote.PublicKey, &userVote.Username, &userVote.Vote)
	if userVote.Vote {
		return echo.NewHTTPError(http.StatusBadRequest, model.Response{Message: "User already voted"})
	}
	privateKey, _, err := Decode(vote.PrivateKey, userVote.PublicKey)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, model.Response{Message: "malformed private key"})
	}
	_, stringPublicKey := Encode(privateKey, &privateKey.PublicKey)
	if strings.Compare(stringPublicKey, userVote.PublicKey) != 0 {
		return echo.NewHTTPError(http.StatusBadRequest, model.Response{Message: "incorrect private key"})
	}
	sign, signHash := ECC.GenSign(vote.Username, privateKey)
	bc.AddBlock(vote.CandidateName, vote.Username, decodePubKey(userVote.PublicKey), sign, signHash)
	_, err = controller.container.GetDB().Exec("UPDATE admins SET vote=$1 WHERE username=$2", true, vote.Username)
	return c.JSON(http.StatusOK, model.Response{Message: "success"})
}

// GetVotingResults - get voting results
func (controller *VotingController) GetVotingResults(c echo.Context) error {
	votes := getVotes(controller)
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
	votes := getVotes(controller)
	err := controller.container.GetBC().EndVoting()
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	_, err = controller.container.GetDB().Exec("TRUNCATE TABLE admins")
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	_, err = controller.container.GetDB().Exec("TRUNCATE TABLE voting_users")
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	_, err = controller.container.GetDB().Exec("TRUNCATE TABLE candidates")
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, votes)
}

func (controller *VotingController) StartVoting(c echo.Context) error {
	fmt.Println("start voting")
	if controller.container.GetBC().CheckVotingInProgress() {
		fmt.Println("voting already in progress")
		return echo.NewHTTPError(http.StatusBadRequest, model.Response{Message: "Voting already in progress"})
	}
	err := controller.container.GetBC().StartVoting()
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	fmt.Println("voting started")
	return c.JSON(http.StatusOK, model.Response{Message: "success"})
}

func decodePubKey(pemEncodedPub string) *ecdsa.PublicKey {
	blockPub, _ := pem.Decode([]byte(pemEncodedPub))
	x509EncodedPub := blockPub.Bytes
	genericPublicKey, _ := x509.ParsePKIXPublicKey(x509EncodedPub)
	publicKey := genericPublicKey.(*ecdsa.PublicKey)
	return publicKey
}

func getVotes(controller *VotingController) []model.VoteResponse {
	bc := controller.container.GetBC()
	if !bc.CheckVotingInProgress() {
		return make([]model.VoteResponse, 0)
	}
	it := bc.Iterator()
	votes := make([]model.VoteResponse, 0)

	for {
		block := it.Next()
		vote := new(model.VoteResponse)
		vote.CandidateName = string(block.CandidateName)
		vote.Username = string(block.Username)
		votes = append(votes, *vote)
		if len(block.PrevHash) == 0 {
			break
		}
	}
	fmt.Println(votes)
	votes = votes[0 : len(votes)-1]
	return votes
}
