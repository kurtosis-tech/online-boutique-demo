package cartstore

type connectionInfo struct {
	username     string
	password     string
	host         string
	port         uint16
	databaseName string
}

func NewConnectionInfo(
	username string,
	password string,
	host string,
	port uint16,
	databaseName string,
) (*connectionInfo, error) {
	return &connectionInfo{
		username:     username,
		password:     password,
		host:         host,
		port:         port,
		databaseName: databaseName,
	}, nil
}
