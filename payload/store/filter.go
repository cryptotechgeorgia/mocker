package store

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/cryptotechgeorgia/mocker/payload"
)

func applyFilter(f payload.FilterBy, buf *bytes.Buffer) {
	var wc []string

	if f.RequestId != nil {
		wc = append(wc, fmt.Sprintf("request_id = %d", *f.RequestId))
	}

	if len(wc) > 0 {
		buf.WriteString(" WHERE ")
		buf.WriteString(strings.Join(wc, " AND "))
	}
}
