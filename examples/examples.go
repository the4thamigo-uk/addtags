package examples

type MyStruct1 struct {
	X int
	Y int `existing:"tag"`
}

type MyStruct2 struct {
	Z int
}
