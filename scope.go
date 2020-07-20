package peas

const (
	SharedScope    string = "shared"
	PrototypeScope string = "prototype"
)

type PeaScope interface {
	GetPeaObject(peaName string) interface{}
	RemovePeaObject(peaName string) interface{}
}
