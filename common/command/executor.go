package command

import (
	"errors"
	"fmt"

	"cardreader/common/status"
)

type Executor struct {
}

type Response struct {
	status  status.Status
	content []byte
}

func (r Response) Status() status.Status {
	return r.status
}

func (r Response) Content() []byte {
	return r.content
}

func (r *Response) AppendContent(c []byte) {
	r.content = append(r.content, c...)
}

func ParseResponse(raw []byte) (r *Response, err error) {
	var rawLen = len(raw)
	if rawLen < 2 {
		return nil, errors.New("too short")
	}
	var st status.Status
	if st, err = status.Parse(raw[rawLen-2:]); err != nil {
		return nil, fmt.Errorf("failed to parse status: %w", err)
	}
	var response = &Response{
		status:  st,
		content: raw[:rawLen-2],
	}
	if st.Level() > status.Normal {
		return response, st
	}
	return response, nil
}
