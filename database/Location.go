package database

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"time"
)

type Location struct {
	ID       uuid.UUID
	Location pgtype.Point
	NiceName string
	Time     time.Time
	UserID   uuid.UUID
	Range    int
	DeleteAt bool
}

func (l *Location) Create() error {
	tx, err := db.Begin()
	if err != nil {
		return pgError(err)
	}
	defer tx.Rollback()

	stmt, err := db.Prepare("INSERT INTO location (location, nick_name, time, user_id) VALUES($1, $2, now(), $3)")
	if err != nil {
		return pgError(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(l.Location, l.NiceName, l.UserID)
	if err != nil {
		return pgError(err)
	}

	err = tx.Commit()
	if err != nil {
		return pgError(err)
	}

	return nil
}

func (l *Location) GetLocationList() ([]Location, error) {
	rows, err := db.Query("SELECT id,nick_name, range FROM location WHERE user_id = $1 AND delete_at = false", l.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to query: %w", err)
	}
	defer rows.Close()

	var ls []Location
	for rows.Next() {
		var location Location
		if err := rows.Scan(&location.ID, &location.NiceName, &location.Range); err != nil {
			return nil, fmt.Errorf("failed to row: %w", err)
		}
		ls = append(ls, location)
	}

	return ls, nil
}

func (l *Location) GetUserId() (*uuid.UUID, error) {
	if l.UserID.ID() != 0 {
		return &l.UserID, nil
	}

	stmt, err := db.Prepare("SELECT user_id FROM location WHERE id = $1")
	if err != nil {
		return nil, fmt.Errorf("failed to query: %w", err)
	}
	defer stmt.Close()

	err = stmt.QueryRow(l.ID).Scan(&l.UserID)
	if err != nil {
		return nil, pgError(err)
	}
	return &l.UserID, nil
}

func (l *Location) EditRange(r int) error {
	tx, err := db.Begin()
	if err != nil {
		return pgError(err)
	}
	defer tx.Rollback()

	stmt, err := db.Prepare("UPDATE location SET range = $1 WHERE id = $2")
	if err != nil {
		return pgError(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(r, l.ID)
	if err != nil {
		return pgError(err)
	}

	err = tx.Commit()
	if err != nil {
		return pgError(err)
	}

	return nil
}

func (l *Location) Delete() error {
	tx, err := db.Begin()
	if err != nil {
		return pgError(err)
	}
	defer tx.Rollback()

	stmt, err := db.Prepare("UPDATE location SET delete_at = true WHERE id = $1")
	if err != nil {
		return pgError(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(l.ID)
	if err != nil {
		return pgError(err)
	}

	err = tx.Commit()
	if err != nil {
		return pgError(err)
	}

	return nil
}

func (l *Location) GetNickName() (string, error) {

	stmt, err := db.Prepare("SELECT nick_name FROM location WHERE id = $1")
	if err != nil {
		return "", fmt.Errorf("failed to query: %w", err)
	}
	defer stmt.Close()

	err = stmt.QueryRow(l.ID).Scan(&l.NiceName)
	if err != nil {
		return "", pgError(err)
	}
	return l.NiceName, nil
}

func (l *Location) IsDelete() (bool, error) {
	stmt, err := db.Prepare("SELECT delete_at FROM location WHERE id = $1")
	if err != nil {
		return false, fmt.Errorf("failed to query: %w", err)
	}
	defer stmt.Close()

	err = stmt.QueryRow(l.ID).Scan(&l.DeleteAt)
	if err != nil {
		return false, pgError(err)
	}
	return l.DeleteAt, nil
}
