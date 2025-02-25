package logger

type Log struct {
	Id      int
	Message string `json:"message"`
}

func (l *Log) CreateLog(map[string]string) {

}
