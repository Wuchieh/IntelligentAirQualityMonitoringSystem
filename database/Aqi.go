package database

import (
	"fmt"
	"github.com/jackc/pgx/v5/pgtype"
	"strconv"
	"time"
)

type Aqi struct {
	Id       int
	Location pgtype.Point
	Aqi      float64
	Time     time.Time
}

// AqiSearch 其中一項為-1 就是檢索全部
type AqiSearch struct {
	Lat   float64 `json:"lat"`
	Lng   float64 `json:"lng"`
	Range int     `json:"range"`
	Limit int
}

func GetAqi(s AqiSearch) ([]Aqi, error) {
	var dq string

	switch {
	case s.Lng == -1 || s.Lat == -1 || s.Range == -1:
		dq = "SELECT id, location, aqi, time FROM aqi"
	}

	if s.Limit > 0 {
		dq += " LIMIT " + strconv.Itoa(s.Limit)
	}

	rows, err := db.Query(dq)
	if err != nil {
		return nil, fmt.Errorf("failed to query: %w", err)
	}
	defer rows.Close()

	var aqis []Aqi
	for rows.Next() {
		var aqi Aqi
		var locationBytes []byte
		if err := rows.Scan(&aqi.Id, &locationBytes, &aqi.Aqi, &aqi.Time); err != nil {
			return nil, fmt.Errorf("failed row: %w", err)
		}

		if err := aqi.Location.UnmarshalJSON(locationBytes); err != nil {
			return nil, fmt.Errorf("failed to UnmarshalJSON location: %w", err)
		}
		aqis = append(aqis, aqi)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return aqis, nil
}
