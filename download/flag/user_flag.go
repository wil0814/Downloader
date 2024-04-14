package flag

type ProtocolType string

const (
	ProtocolHTTP ProtocolType = "http"
	ProtocolFTP  ProtocolType = "ftp"
)

type UserFlag struct {
	Path        string
	Protocol    ProtocolType
	FileName    string
	Concurrency int
}
