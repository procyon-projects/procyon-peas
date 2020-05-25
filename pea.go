package peas

type ConfigurablePeaFactory interface {
	PeaFactory
	SharedPeaRegistry
}
