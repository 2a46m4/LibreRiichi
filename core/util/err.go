package core

func PanicIf(err error) {
    if err != nil {
	panic(err)
    }
}
