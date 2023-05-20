package common

func CriticErrorHandler(err error) {
	if err != nil {
		panic(err)
	}
}
