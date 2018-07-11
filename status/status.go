package status

// StatusAbiturienta is position abiturient in list
type StatusAbiturienta struct {
	Num             int
	NumWithOriginal int
}

func (s StatusAbiturienta) IsEqual(s2 StatusAbiturienta) bool {
	return (s.Num == s2.Num) && (s.NumWithOriginal == s2.NumWithOriginal)
}

// StatusByName is struct status with name of abiturient
type StatusByName struct {
	Name   string
	Status StatusAbiturienta
}
