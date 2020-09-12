package pgrange

import (
	"database/sql"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type tstzrangeTest struct {
	Duration *TimestampTzRange
}

var tstzranges = []*TimestampTzRange{
	NewTimestampTzRange(
		timestampz(2020, 1, 1, 12, 0, 0, time.Local),
		timestampz(2020, 1, 1, 13, 0, 0, time.Local),
	),
	NewTimestampTzRange(
		timestampz(2020, 1, 1, 12, 0, 0, time.Local),
		nil,
	),
	NewTimestampTzRange(
		nil,
		timestampz(2020, 1, 1, 13, 0, 0, time.Local),
	),
	NewTimestampTzRange(
		nil,
		nil,
	),
}

func timestampz(year int, month time.Month, day, hour, min, sec int, loc *time.Location) *time.Time {
	d := time.Date(year, month, day, hour, min, sec, 0, loc)
	return &d
}

func TsTzRange(t *testing.T, db *sql.DB) {
	for _, duration := range tstzranges {
		t.Run(duration.String(), func(t *testing.T) {
			_, err := db.Exec(
				"delete from tstzrange_test",
			)

			require.NoError(t, err)

			c := &tstzrangeTest{
				Duration: duration,
			}

			_, err = db.Exec(
				"insert into tstzrange_test (duration) values ($1)",
				c.Duration,
			)

			require.NoError(t, err)

			rows, err := db.Query("SELECT duration FROM tstzrange_test")
			if err != nil {
				panic(err)
			}

			defer rows.Close()

			if rows.Next() {
				var retrievedDuration TimestampTzRange
				if err := rows.Scan(&retrievedDuration); err != nil {
					panic(err)
				}

				assert.Equal(t, &TimestampTzRange{
					Start:          utc(duration.Start),
					StartInclusive: duration.StartInclusive,
					End:            utc(duration.End),
					EndInclusive:   duration.EndInclusive,
				}, &retrievedDuration)
			} else {
				t.Fatal("expected single result")
			}
		})
	}
}

func utc(t *time.Time) *time.Time {
	if t != nil {
		utcT := t.UTC()
		return &utcT
	}

	return nil
}
