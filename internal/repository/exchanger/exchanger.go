package exchanger

import (
	"time"

	"async_arch/internal/logger"
	"async_arch/internal/storage/database"

	"github.com/google/uuid"
)

type Exchanger struct {
	conn database.DBConnector
}

func InitExchanger(conn database.DBConnector) *Exchanger {
	return &Exchanger{conn: conn}
}

func (e *Exchanger) SetKey(key string) (res uuid.UUID, err error) {
	row := e.conn.Client().QueryRowx(
		"INSERT INTO one_time_key (key_id, key_val, expires) VALUES (gen_random_uuid(), $1, $2) RETURNING key_id",
		key, time.Now().Add(10*time.Minute),
	)
	err = row.Scan(&res)
	if err != nil {
		logger.Error("error set key", err)
	}
	return res, err
}

func (e *Exchanger) GetKey(uuid uuid.UUID) (res string, err error) {
	row := e.conn.Client().QueryRowx("SELECT key_val FROM one_time_key WHERE key_id = $1 AND expires > NOW()", uuid)
	if err = row.Scan(&res); err == nil {
		e.conn.Client().Exec("DELETE FROM one_time_key WHERE key_id = $1", uuid)
	}
	return res, err
}
