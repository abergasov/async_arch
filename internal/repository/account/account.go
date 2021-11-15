package account

import "async_arch/internal/storage/database"

type Account struct {
	conn database.DBConnector
}

func InitAccountRepo(conn database.DBConnector) *Account {
	return &Account{conn: conn}
}

func (a *Account) ChangeAccount(publicID string, amount int64) error {
	_, err := a.conn.Client().NamedExec("UPDATE account SET amount = :amount WHERE public_id = :public_id", map[string]interface{}{
		"public_id": publicID,
		"amount":    amount,
	})
	return err
}

func (a *Account) CreateAccount(publicID string) error {
	_, err := a.conn.Client().NamedExec("INSERT INTO account (public_id) VALUES (:public_id)", map[string]interface{}{
		"public_id": publicID,
	})
	return err
}
