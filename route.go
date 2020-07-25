package telebot

// todo desc

type HandlerFunc func(Context) error
type MiddlewareFunc func(Context) error

type route struct {
	h  HandlerFunc
	ms []MiddlewareFunc
}
