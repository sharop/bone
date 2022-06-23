package s

// Constant Error values which can be compared to determine the type of error
const (
	ErrorBoneIndexErrorOnInit Error = iota
)

// Error represents a more strongly typed bleve error for detecting
// and handling specific types of errors.
type Error int

func (e Error) Error() string {
	return errorMessages[e]
}

var errorMessages = map[Error]string{
	ErrorBoneIndexErrorOnInit: "Error initializing core mapping. ",
}