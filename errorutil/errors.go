package errorutil

const ErrorBadParams 					= 2
const ErrorFailedConsulConnection 		= 3
const ErrorFailedReadingResponse 		= 4
const ErrorFailedJsonDecode 			= 5
const ErrorFailedCloning 				= 6

type GonsulError struct {
	Code	int
}

func ExitError(err error, errorCode int, logger *Logger) {
	logger.PrintError(err.Error())
	panic(GonsulError{Code: errorCode})
}