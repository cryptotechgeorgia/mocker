package router

var (
	DefaultResponse ResponseData = ResponseData{
		ContentType: "application/json",
		Payload:     `{"error": {"UserMsg":"route you're looking for does not exist for this project","Description":"","Base":"" }`,
	}

	DefaultNoRouteResponse ResponseData = ResponseData{
		ContentType: "application/json",
		Payload:     `{"error": {"UserMsg":"oops! the route you're looking for does not exist.","Description":"","Base":"" }`,
	}
)

type RequestData struct {
	ContentType string
	Payload     string
}

type ResponseData struct {
	ContentType string
	Payload     string
}

// gives ability to costumize
// request response pair
type PayloadPair struct {
	req  RequestData
	resp ResponseData
}

func NewRequestData(contentType, payload string) RequestData {
	return RequestData{
		ContentType: contentType,
		Payload:     payload,
	}
}

func NewResponseData(contentType, payload string) ResponseData {
	return ResponseData{
		ContentType: contentType,
		Payload:     payload,
	}
}

func NewPayloadPair(req RequestData, resp ResponseData) PayloadPair {
	return PayloadPair{
		req:  req,
		resp: resp,
	}
}

type Request struct {
	method       string
	path         string
	defaultResp  *ResponseData
	payloadPairs []PayloadPair
}

func (r *Request) SetDefaultResponse(data ResponseData) {
	r.defaultResp = &data
}

func NewRequest(method, path string) *Request {
	return &Request{
		method:       method,
		path:         path,
		defaultResp:  &DefaultResponse,
		payloadPairs: []PayloadPair{},
	}
}

type Project struct {
	Name     string
	BaseAddr string
	Requests []Request
}
