package pgrange

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"
	"reflect"
	"runtime"
	"strings"
	"testing"
)

func setup(t *testing.T, db *sql.DB) {
	_, err := db.Exec(
		"create table tsrange_test (duration tsrange not null)",
	)

	require.NoError(t, err)

	_, err = db.Exec(
		"create table tstzrange_test (duration tstzrange not null)",
	)

	require.NoError(t, err)

	_, err = db.Exec(
		"create table daterange_test (duration daterange not null)",
	)

	require.NoError(t, err)
}

func teardown(t *testing.T, db *sql.DB) {
	_, err := db.Exec(
		"drop table if exists tsrange_test",
	)

	require.NoError(t, err)

	_, err = db.Exec(
		"drop table if exists tstzrange_test",
	)

	require.NoError(t, err)

	_, err = db.Exec(
		"drop table if exists daterange_test",
	)

	require.NoError(t, err)
}

func Test(t *testing.T) {
	for version, port := range map[string]int{"10": 5430, "11": 5431, "12": 5433} {
		t.Run(version, func(t *testing.T) {
			db, err := sql.Open("postgres", fmt.Sprintf("host=localhost port=%d user=postgres dbname=test password=pass sslmode=disable", port))
			require.NoError(t, err)

			teardown(t, db)

			defer db.Close()
			for _, tt := range tests {
				name := runtime.FuncForPC(reflect.ValueOf(tt).Pointer()).Name()
				t.Run(name[strings.LastIndex(name, ".")+1:], func(t *testing.T) {
					setup(t, db)
					tt(t, db)
					teardown(t, db)
				})
			}
		})
	}
}

var (
	tests = [...]func(*testing.T, *sql.DB){
		TsRange,
		TsTzRange,
		DRange,
	}
)
