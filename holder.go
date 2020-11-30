package peas

type PeaDefinitionHolder struct {
	peaName       string
	peaDefinition PeaDefinition
}

func NewPeaDefinitionHolder(peaName string, peaDefinition PeaDefinition) *PeaDefinitionHolder {
	if peaName == "" {
		panic("Pea Name must not be empty")
	}
	if peaDefinition == nil {
		panic("Pea Definition must not be nil")
	}
	return &PeaDefinitionHolder{
		peaName,
		peaDefinition,
	}
}

func (holder *PeaDefinitionHolder) GetPeaName() string {
	return holder.peaName
}

func (holder *PeaDefinitionHolder) GetPeaDefinition() PeaDefinition {
	return holder.peaDefinition
}
