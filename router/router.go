package router

import (
	"github.com/labstack/echo/v4"
	"votingblockchain/container"
	"votingblockchain/controller"
)

func Init(e *echo.Echo, container container.Container) {
	health := controller.NewHealthController(container)
	user := controller.NewUserController(container)
	voting := controller.NewVotingController(container)
	e.GET("/health", func(c echo.Context) error { return health.GetHealth(c) })
	e.POST("/login", func(c echo.Context) error { return user.Login(c) })
	e.POST("/signup", func(c echo.Context) error { return user.Signup(c) })
	e.GET("/vote", func(c echo.Context) error { return voting.GetVoteByUsername(c) })
	e.POST("/vote", func(c echo.Context) error { return voting.Vote(c) })
	e.PUT("/voting", func(c echo.Context) error { return voting.UpdateVoting(c) })
	e.GET("/voting/results", func(c echo.Context) error { return voting.GetVotingResults(c) })
	e.GET("/candidates", func(c echo.Context) error { return voting.GetCandidates(c) })
	e.POST("/candidate", func(c echo.Context) error { return voting.AddCandidate(c) })
	e.GET("/end", func(c echo.Context) error { return voting.EndVoting(c) })
}
