package pgrange

import (
	"database/sql"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type daterangeTest struct {
	Duration *DateRange
}

var dateranges = []struct {
	in  *DateRange
	out *DateRange
}{
	{NewDateRange(
		date(2020, 1, 1),
		date(2020, 1, 2),
	), &DateRange{
		Start:          date(2020, 1, 1),
		StartInclusive: true,
		End:            date(2020, 1, 3),
		EndInclusive:   false,
	}},
	{NewDateRange(
		date(2020, 1, 1),
		nil,
	), NewDateRange(
		date(2020, 1, 1),
		nil,
	)},
	{NewDateRange(
		nil,
		date(2020, 1, 1),
	), &DateRange{
		Start:          nil,
		StartInclusive: false,
		End:            date(2020, 1, 2),
		EndInclusive:   false,
	}},
	{NewDateRange(
		nil,
		nil,
	), NewDateRange(
		nil,
		nil,
	)},
}

func date(year int, month time.Month, day int) *time.Time {
	d := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
	return &d
}

func DRange(t *testing.T, db *sql.DB) {
	for _, duration := range dateranges {
		t.Run(duration.in.String(), func(t *testing.T) {
			_, err := db.Exec(
				"delete from daterange_test",
			)

			require.NoError(t, err)

			c := &daterangeTest{
				Duration: duration.in,
			}

			_, err = db.Exec(
				"insert into daterange_test (duration) values ($1)",
				c.Duration,
			)

			require.NoError(t, err)

			rows, err := db.Query("SELECT duration FROM daterange_test")
			if err != nil {
				panic(err)
			}

			defer rows.Close()

			if rows.Next() {
				var retrievedDuration DateRange
				if err := rows.Scan(&retrievedDuration); err != nil {
					panic(err)
				}

				assert.Equal(t, duration.out, &retrievedDuration)
			} else {
				t.Fatal("expected single result")
			}
		})
	}
}
