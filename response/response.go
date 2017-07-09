package response

type Response interface {
	Render() ([]byte, error)
	ContentType() string
	Status() int
}
