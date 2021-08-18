package main

func maybePanic(err error) {
	if err != nil {
		panic(err)
	}
}
