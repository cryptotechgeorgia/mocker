package store

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/cryptotechgeorgia/mocker/response"
)

func applyFilter(f response.FilterBy, buf *bytes.Buffer) {
	var wc []string

	if f.RequestPayloadId != nil {
		wc = append(wc, fmt.Sprintf("request_payload_id = %d", *f.RequestPayloadId))
	}

	if len(wc) > 0 {
		buf.WriteString(" WHERE ")
		buf.WriteString(strings.Join(wc, " AND "))
	}
}
