package ftp

//type FTPUrlParser struct{}
//type FTPUrlInfo struct {
//	Username string
//	Password string
//	Host     string
//	Port     string
//	Path     string
//	Filename string
//}
//
//func NewFTPURLParser() *FTPUrlParser {
//	return &FTPUrlParser{}
//}
//func (p *FTPUrlParser) Parse(ftpPath string) (*FTPUrlInfo, error) {
//	u, err := url.Parse(ftpPath)
//	if err != nil {
//		return nil, err
//	}
//
//	if u.Scheme != "ftp" {
//		return nil, fmt.Errorf("unsupported scheme: %s", u.Scheme)
//	}
//
//	username := ""
//	password := ""
//	if u.User != nil {
//		username = u.User.Username()
//		password, _ = u.User.Password()
//	}
//
//	host := u.Hostname()
//	port := u.Port()
//	if port == "" {
//		port = "21"
//	}
//
//	path := u.Path
//	filename := path[strings.LastIndex(path, "/")+1:]
//	path = path[:strings.LastIndex(path, "/")+1]
//
//	return &FTPUrlInfo{
//		Username: username,
//		Password: password,
//		Host:     host,
//		Port:     port,
//		Path:     path,
//		Filename: filename,
//	}, nil
//}
