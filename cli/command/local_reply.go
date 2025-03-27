package command

import (
	"fmt"

	"github.com/handlename/otomo/internal/app/usecase"
	"github.com/handlename/otomo/internal/domain/entity"
	"github.com/handlename/otomo/internal/infra/repository"
	"github.com/handlename/otomo/internal/infra/service"
	"github.com/morikuni/failure/v2"
)

type LocalReply struct {
	Message string `help:"Message to send"`
}

func (o *LocalReply) Run(c *Context) error {
	repoS := &repository.VolatileSession{}
	repoB := repository.NewGeneralBrain(c.Ctx, c.App.BotToken, c.App.AppToken)
	msgr := &service.NopSlack{}

	brain, err := repoB.New(c.Ctx)
	if err != nil {
		return failure.Wrap(err, failure.Message("failed to new brain repository"))
	}

	otomo := entity.NewOtomo(brain)
	inst := entity.NewInstruction("dummy", o.Message)

	uc := usecase.NewReplyToUser(repoS, msgr)
	if err := uc.Run(c.Ctx, otomo, inst); err != nil {
		return failure.Wrap(err)
	}

	fmt.Println(msgr.Memory)

	return nil
}
