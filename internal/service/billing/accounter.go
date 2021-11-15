package billing

import (
	"context"
	"encoding/json"

	"async_arch/internal/entities"
	"async_arch/internal/logger"

	"github.com/segmentio/kafka-go"
)

type AccountRepository interface {
	ChangeAccount(publicID string, amount int64) error
	CreateAccount(publicID string) error
}

type Accounter struct {
	accRepo      AccountRepository
	userStreamer *kafka.Reader
}

func InitAccounter(accRepo AccountRepository, userStreamer *kafka.Reader) *Accounter {
	a := &Accounter{
		accRepo:      accRepo,
		userStreamer: userStreamer,
	}
	go a.watchUserChanges()
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
