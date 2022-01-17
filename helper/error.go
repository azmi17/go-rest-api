package helper

func PanicIfError(err error) {
	if err != nil {
		panic(err) // nanti akan di handle panic-nya..
	}
}
