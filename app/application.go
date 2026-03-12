package app

import (
	"github.com/Strangebrewer/go-server/job"
	"github.com/Strangebrewer/go-server/recruiter"
	"github.com/Strangebrewer/go-server/thing"
	"github.com/Strangebrewer/go-server/token"
	"github.com/Strangebrewer/go-server/user"
)

type Application struct {
	ThingStore     *thing.ThingStore
	JobStore       *job.JobStore
	RecruiterStore *recruiter.RecruiterStore
	TokenStore     *token.TokenStore
	TokenService   *token.TokenService
	UserStore      *user.UserStore
}
