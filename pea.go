package peas

import (
	"github.com/procyon-projects/goo"
)

type ConfigurablePeaFactory interface {
	SharedPeaRegistry
	PeaFactory
	RegisterTypeAsOnlyReadable(typ goo.Type) error
	ExcludeType(typ goo.Type) error
	AddPeaProcessor(processor PeaProcessor) error
	GetPeaProcessors() []PeaProcessor
	GetPeaProcessorsCount() int
	PreInstantiateSharedPeas()
}

type PeaInitializer interface {
	InitializePea() error
}

type PeaNameGenerator interface {
	GenerateName(peaDefinition PeaDefinition) string
}
