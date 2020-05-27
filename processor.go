package peas

type PeaProcessor interface {
	BeforeInitialization(peaName string, pea interface{})
	AfterInitialization(peaName string, pea interface{})
}
