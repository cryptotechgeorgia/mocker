package store

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/cryptotechgeorgia/mocker/request"
)

func applyFilter(f request.FilterBy, buf *bytes.Buffer) {
	var wc []string

	if f.ProjectId != nil {
		wc = append(wc, fmt.Sprintf("project_id = %d", *f.ProjectId))
	}

	if f.ResponseId != nil {
		wc = append(wc, fmt.Sprintf("response_id = %d", *f.ResponseId))
	}

	if len(wc) > 0 {
		buf.WriteString(" WHERE ")
		buf.WriteString(strings.Join(wc, " AND "))
	}

}
