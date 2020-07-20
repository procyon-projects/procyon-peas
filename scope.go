package peas

const (
	PeaSharedScope    string = "shared"
	PeaPrototypeScope string = "prototype"
)

type PeaScope interface {
	GetPeaObject(peaName string) interface{}
	RemovePeaObject(peaName string) interface{}
}
