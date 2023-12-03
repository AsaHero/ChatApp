package entity

var (
	ErrorConflict = NewErrConflict("object")
	ErrorNotFound = NewErrNotFound("object")
)

// error not found
type ErrNotFound struct {
	name string
}

func (e *ErrNotFound) Error() string {
	return e.name + " not found"
}

func NewErrNotFound(text string) *ErrNotFound {
	return &ErrNotFound{text}
}

// error conflict
type ErrConflict struct {
	name string
}

func (e *ErrConflict) Error() string {
	return e.name + " already exist"
}

func NewErrConflict(text string) *ErrConflict {
	return &ErrConflict{text}
}
