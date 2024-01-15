package lib

func Pointer[T any](element T) *T {
	return &element
}
