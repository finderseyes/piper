package child

type Data struct {
	A int
}

type Foo interface {
	Bar(*Data) *Data
}
