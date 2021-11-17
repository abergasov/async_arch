package billing

import (
	"context"
	"encoding/json"

	"async_arch/internal/entities"
	"async_arch/internal/logger"

	"github.com/segmentio/kafka-go"
)

type AccountRepository interface {
	CreateAccount(publicID string) error
	CreateAssignAccountTransaction(userPublicID, taskPublicID string) error
	CreateFinishAccountTransaction(userPublicID, taskPublicID string) error
}

type Accounter struct {
	accRepo      AccountRepository
	userStreamer *kafka.Reader
	taskStreamer *kafka.Reader
}

func InitAccounter(accRepo AccountRepository, userStreamer, taskStreamer *kafka.Reader) *Accounter {
	a := &Accounter{
		accRepo:      accRepo,
		userStreamer: userStreamer,
		taskStreamer: taskStreamer,
	}
	go a.watchUserChanges()
	go a.watchTaskChanges()
	return a
}

func (a *Accounter) watchUserChanges() {
	for {
		m, err := a.userStreamer.ReadMessage(context.Background())
		if err != nil {
			logger.Error("error read message from broker", err)
			break
		}

		var usr entities.UserBIEvent
		if err = json.Unmarshal(m.Value, &usr); err != nil {
			logger.Error("error parse message UserBIEvent", err)
			continue
		}
		switch string(m.Key) {
		case entities.UserBICreatedEvent:
			err = a.accRepo.CreateAccount(usr.PublicID)
		}
		if err != nil {
			logger.Error("error parse message", err)
		}
	}
}

func (a *Accounter) watchTaskChanges() {
	for {
		m, err := a.taskStreamer.ReadMessage(context.Background())
		if err != nil {
			logger.Error("error read message from broker taskStreamer", err)
			continue
		}

		var cnt entities.TaskAssignContainer
		if err = json.Unmarshal(m.Value, &cnt); err != nil {
			logger.Error("adawda", err)
		}

		switch string(m.Key) {
		case entities.TaskAssignedBIEvent:
			err = a.accRepo.CreateAssignAccountTransaction(cnt.UserPublicID, cnt.TaskPublicID)
		case entities.TaskFinishBIEvent:
			err = a.accRepo.CreateFinishAccountTransaction(cnt.UserPublicID, cnt.TaskPublicID)
		}
		if err != nil {
			logger.Error("error parse message", err)
		}
	}
}
