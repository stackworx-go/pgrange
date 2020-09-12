package pgrange

import (
	"database/sql/driver"
	"fmt"
	"regexp"
	"strings"
	"time"
)

// TimestampRange export
type TimestampRange struct {
	Start          *time.Time
	StartInclusive bool
	End            *time.Time
	EndInclusive   bool
	// ?
	Valid bool
}

func NewTimestampRange(start, end *time.Time) *TimestampRange {
	tr := TimestampRange{
		Start: start,
		End:   end,
	}

	if start != nil {
		tr.StartInclusive = true
	}

	if end != nil {
		tr.EndInclusive = true
	}

	return &tr
}

const timestampLayout = "2006-01-02 15:04:05"

func (r TimestampRange) Value() (driver.Value, error) {
	var sb strings.Builder

	if r.StartInclusive {
		sb.WriteString("[")
	} else {
		sb.WriteString("(")
	}

	if r.Start != nil {
		sb.WriteString(fmt.Sprintf("'%s'", r.Start.Format(timestampLayout)))
	}

	sb.WriteString(",")

	if r.End != nil {
		sb.WriteString(fmt.Sprintf("'%s'", r.End.Format(timestampLayout)))
	}

	if r.EndInclusive {
		sb.WriteString("]")
	} else {
		sb.WriteString(")")
	}

	return sb.String(), nil
}

func (r *TimestampRange) Scan(src interface{}) error {
	switch src := src.(type) {
	case []byte:
		return r.scanString(string(src))
	case string:
		return r.scanString(src)
	case nil:
		return nil
	}

	return fmt.Errorf("pq: cannot convert %T to timestamprange", src)
}

var timestampRegex = regexp.MustCompile(`^([\(\[])("\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}")?,("\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}")?([\)\]])$`)

func (r *TimestampRange) scanString(src string) error {
	result := timestampRegex.FindStringSubmatch(src)

	if result == nil {
		return fmt.Errorf("invalid range: %s", src)
	}

	if result[1] == "[" {
		r.StartInclusive = true
	}

	if result[2] != "" {
		start, err := time.ParseInLocation(timestampLayout, unquote(result[2]), time.UTC)

		if err != nil {
			return fmt.Errorf("invalid start timestamp: %s, range: %s", result[2], src)
		}

		r.Start = &start
	}

	if result[3] != "" {
		end, err := time.ParseInLocation(timestampLayout, unquote(result[3]), time.UTC)

		if err != nil {
			return fmt.Errorf("invalid end timestamp: %s, range: %s", result[3], src)
		}

		r.End = &end
	}

	if result[4] == "]" {
		r.EndInclusive = true
	}

	return nil
}

func (r *TimestampRange) String() string {
	v, err := r.Value()

	if err != nil {
		return err.Error()
	}

	return fmt.Sprint(v)
}
