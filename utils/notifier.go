package utils

type Urgency byte

const (
	UrLow Urgency = iota
	UrNormal
	UrCritical
)

type NotifyContent struct {
	Title   string
	Text    string
	Url     string
	Icon    string
	GroupId string
	urgency Urgency
}

type Notifier interface {
	push(c NotifyContent)
}
