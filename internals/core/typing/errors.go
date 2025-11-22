package typing

import "errors"

var (
	ErrInvalidUser        error = errors.New("invalid user id")
	ErrInvalidWPM         error = errors.New("invalid wpm")
	ErrInvalidTypedWords  error = errors.New("invalid typed words")
	ErrInvalidTotalWords  error = errors.New("invalid total words")
	ErrSomethingWentWrong error = errors.New("something went wrong")
	ErrinvalidTotalErrors error = errors.New("invalid total errors count")
	ErrInsertingData      error = errors.New("error inserting data to DB")
	ErrUpdatingTotalTest  error = errors.New("error in updating test count")
	ErrGettingDataFromDB  error = errors.New("error getting data from DB")
)
