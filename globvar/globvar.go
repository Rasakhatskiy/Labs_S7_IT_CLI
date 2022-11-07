package globvar

var (
	Headers            []string
	Types              []string
	TableName          string
	DBname             string
	TableOperationType int
)

const (
	Create = iota
	Read
	Update
	Delete
)

const (
	REQ_GET    = "GET"
	REQ_POST   = "POST"
	REQ_DELETE = "DELETE"
	REQ_PUT    = "PUT"
)
