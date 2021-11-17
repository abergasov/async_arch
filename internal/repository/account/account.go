package account

import (
	"math/rand"

	"async_arch/internal/storage/database"
)

type Account struct {
	conn database.DBConnector
}

func InitAccountRepo(conn database.DBConnector) *Account {
	return &Account{conn: conn}
}

func (a *Account) CreateAssignAccountTransaction(userPublicID, taskPublicID string) error {
	_, err := a.conn.Client().NamedExec("INSERT INTO account_transactions (user_id, amount, transaction_date) VALUES (:user_id, :amount, NOW())", map[string]interface{}{
		"user_id": userPublicID,
		"task_id": taskPublicID,
		"amount":  a.calcCost() * -1,
	})
	return err
}

func (a *Account) CreateFinishAccountTransaction(userPublicID, taskPublicID string) error {
	_, err := a.conn.Client().NamedExec("INSERT INTO account_transactions (user_id, amount, transaction_date) VALUES (:user_id, :amount, NOW())", map[string]interface{}{
		"user_id": userPublicID,
		"task_id": taskPublicID,
		"amount":  a.calcCost(),
	})
	return err
}

func (a *Account) CreateAccount(publicID string) error {
	_, err := a.conn.Client().NamedExec("INSERT INTO account (public_id) VALUES (:public_id)", map[string]interface{}{
		"public_id": publicID,
	})
	return err
}

func (a *Account) calcCost() int64 {
	return int64(rand.Intn(20-1) + 1)
}
