package pgrange

import (
	"database/sql/driver"
	"fmt"
	"regexp"
	"strings"
	"time"
)

// DateRange export
type DateRange struct {
	Start          *time.Time
	StartInclusive bool
	End            *time.Time
	EndInclusive   bool
	// ?
	Valid bool
}

func NewDateRange(start, end *time.Time) *DateRange {
	tr := DateRange{
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

const dateLayout = "2006-01-02"

func (r DateRange) Value() (driver.Value, error) {
	var sb strings.Builder

	if r.StartInclusive {
		sb.WriteString("[")
	} else {
		sb.WriteString("(")
	}

	if r.Start != nil {
		sb.WriteString(fmt.Sprintf("'%s'", r.Start.Format(dateLayout)))
	}

	sb.WriteString(",")

	if r.End != nil {
		sb.WriteString(fmt.Sprintf("'%s'", r.End.Format(dateLayout)))
	}

	if r.EndInclusive {
		sb.WriteString("]")
	} else {
		sb.WriteString(")")
	}

	return sb.String(), nil
}

func (r *DateRange) Scan(src interface{}) error {
	switch src := src.(type) {
	case []byte:
		return r.scanString(string(src))
	case string:
		return r.scanString(src)
	case nil:
		return nil
	}

	return fmt.Errorf("pq: cannot convert %T to daterange", src)
}

var dateRegex = regexp.MustCompile(`^([\(\[])(\d{4}-\d{2}-\d{2})?,(\d{4}-\d{2}-\d{2})?([\)\]])$`)

func (r *DateRange) scanString(src string) error {
	result := dateRegex.FindStringSubmatch(src)

	if result == nil {
		return fmt.Errorf("invalid range: %s", src)
	}

	if result[1] == "[" {
		r.StartInclusive = true
	}

	if result[2] != "" {
		start, err := time.ParseInLocation(dateLayout, result[2], time.UTC)

		if err != nil {
			return fmt.Errorf("invalid start date: %s, range: %s", result[2], src)
		}

		r.Start = &start
	}

	if result[3] != "" {
		end, err := time.ParseInLocation(dateLayout, result[3], time.UTC)

		if err != nil {
			return fmt.Errorf("invalid end date: %s, range: %s", result[3], src)
		}

		r.End = &end
	}

	if result[4] == "]" {
		r.EndInclusive = true
	}

	return nil
}

func (r *DateRange) String() string {
	v, err := r.Value()

	if err != nil {
		return err.Error()
	}

	return fmt.Sprint(v)
}
