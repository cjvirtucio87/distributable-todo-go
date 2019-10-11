package rlogging

type Logger interface {
	Infof()
	Debugf()
	Tracef()
	Warnf()
}
