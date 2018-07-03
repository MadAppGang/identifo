package identifo

//Service represents any external service
type Service interface {
	Connect() (Session, error)
}

//Session is a Service's session created with context
type Session interface {
}
