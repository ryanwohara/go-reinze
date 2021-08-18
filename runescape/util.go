package runescape

func maybePanic(err error) {
	if err != nil {
		panic(err)
	}
}