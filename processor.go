package peas

import (
	"errors"
	"github.com/codnect/goo"
	"sync"
)

type PeaProcessor interface {
	BeforePeaInitialization(peaName string, pea interface{}) (interface{}, error)
	AfterPeaInitialization(peaName string, pea interface{}) (interface{}, error)
}

type PeaProcessors struct {
	processors map[string]PeaProcessor
	mu         sync.RWMutex
}

func NewPeaProcessors() *PeaProcessors {
	return &PeaProcessors{
		make(map[string]PeaProcessor, 0),
		sync.RWMutex{},
	}
}

func (p *PeaProcessors) AddPeaProcessor(processor PeaProcessor) error {
	if processor == nil {
		return errors.New("processor cannot be null")
	}
	p.mu.Lock()
	processorType := goo.GetType(processor)
	if _, ok := p.processors[processorType.GetFullName()]; ok {
		return errors.New("You have already registered this processor : " + processorType.GetFullName())
	}
	p.processors[processorType.GetFullName()] = processor
	p.mu.Unlock()
	return nil
}

func (p *PeaProcessors) RemoveProcessor(processor PeaProcessor) {
	if processor == nil {
		return
	}
	p.mu.Lock()
	processorType := goo.GetType(processor)
	if _, ok := p.processors[processorType.GetFullName()]; ok {
		delete(p.processors, processorType.GetFullName())
	}
	p.mu.Unlock()
}

func (p *PeaProcessors) GetProcessors() []PeaProcessor {
	processors := make([]PeaProcessor, 0)
	p.mu.Lock()
	for _, val := range p.processors {
		processors = append(processors, val)
	}
	p.mu.Unlock()
	return processors
}

func (p *PeaProcessors) GetProcessorsCount() int {
	return len(p.processors)
}

func (p *PeaProcessors) RemoveAllProcessor() {
	p.mu.Lock()
	p.processors = make(map[string]PeaProcessor, 0)
	p.mu.Unlock()
}

type PeaDefinitionRegistryProcessor interface {
	AfterPeaDefinitionRegistryInitialization(registry PeaDefinitionRegistry)
}

type PeaFactoryProcessor interface {
	AfterPeaFactoryInitialization(factory ConfigurablePeaFactory)
}
