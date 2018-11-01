package meshinfo

import (
	"fmt"
)

// HandlerOverflowError defines the error that happens when a handler cannot create new MeshInfos.
type HandlerOverflowError struct{}

func (e *HandlerOverflowError) Error() string {
	return "handler overflow"
}

// FaceCountMissmatchError defines an error that happens when a mesh info operation is done with a not matching current face number.
type FaceCountMissmatchError struct {
	current, new uint32
}

func (e *FaceCountMissmatchError) Error() string {
	return fmt.Sprintf("mesh information face count (%d) does not match with mesh face count (%d)", e.current, e.new)
}

// FaceDataIndexError defines an error that happens when a face data is queried by index but the index is higher than the number of faces
type FaceDataIndexError struct {
	faceNum, index uint32
}

func (e *FaceDataIndexError) Error() string {
	return fmt.Sprintf("could not access face data (%d > %d)", e.index, e.faceNum)
}