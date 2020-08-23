package notifier

type Notifier interface {
	Notify(message string) error
}
