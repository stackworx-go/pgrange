package pgrange

import (
	"database/sql"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type tsrangeTest struct {
	Duration *TimestampRange
}

var tsranges = []*TimestampRange{
	NewTimestampRange(
		timestamp(2020, 1, 1, 12, 0, 0),
		timestamp(2020, 1, 1, 13, 0, 0),
	),
	NewTimestampRange(
		timestamp(2020, 1, 1, 12, 0, 0),
		nil,
	),
	NewTimestampRange(
		nil,
		timestamp(2020, 1, 1, 13, 0, 0),
	),
	NewTimestampRange(
		nil,
		nil,
	),
}

func timestamp(year int, month time.Month, day, hour, min, sec int) *time.Time {
	d := time.Date(year, month, day, hour, min, sec, 0, time.UTC)
	return &d
}

func TsRange(t *testing.T, db *sql.DB) {
	for _, duration := range tsranges {
		t.Run(duration.String(), func(t *testing.T) {
			_, err := db.Exec(
				"delete from tsrange_test",
			)

			require.NoError(t, err)

			c := &tsrangeTest{
				Duration: duration,
			}

			_, err = db.Exec(
				"insert into tsrange_test (duration) values ($1)",
				c.Duration,
			)

			require.NoError(t, err)

			rows, err := db.Query("SELECT duration FROM tsrange_test")
			if err != nil {
				panic(err)
			}

			defer rows.Close()

			if rows.Next() {
				var retrievedDuration TimestampRange
				if err := rows.Scan(&retrievedDuration); err != nil {
					panic(err)
				}

				assert.Equal(t, duration, &retrievedDuration)
			} else {
				t.Fatal("expected single result")
			}
		})
	}
}
