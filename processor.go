package peas

import (
	"errors"
	core "github.com/procyon-projects/procyon-core"
	"sync"
)

type PeaProcessor interface {
	BeforeInitialization(peaName string, pea interface{}) (interface{}, error)
	AfterInitialization(peaName string, pea interface{}) (interface{}, error)
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
	processorType := core.GetType(processor)
	if _, ok := p.processors[processorType.String()]; ok {
		return errors.New("You have already registered this processor : " + processorType.String())
	}
	p.processors[processorType.String()] = processor
	p.mu.Unlock()
	return nil
}

func (p *PeaProcessors) RemoveProcessor(processor PeaProcessor) {
	if processor == nil {
		return
	}
	p.mu.Lock()
	processorType := core.GetType(processor)
	if _, ok := p.processors[processorType.String()]; ok {
		delete(p.processors, processorType.String())
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
