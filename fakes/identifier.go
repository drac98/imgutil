package fakes

type Indentifier struct {
	hash string
}

func NewIdentifier(hash string) Indentifier {
	return Indentifier{
		hash: hash,
	}
}

func (f Indentifier) String() string {
	return f.hash
}
