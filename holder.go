package peas

type PeaDefinitionHolder struct {
	peaName       string
	peaDefinition PeaDefinition
	aliases       []string
}

func NewPeaDefinitionHolder(peaName string, peaDefinition PeaDefinition) *PeaDefinitionHolder {
	return NewPeaDefinitionHolderWithAliases(peaName, peaDefinition, nil)
}

func NewPeaDefinitionHolderWithAliases(peaName string, peaDefinition PeaDefinition, aliases []string) *PeaDefinitionHolder {
	if peaName == "" {
		panic("Pea Name must not be empty")
	}
	if peaDefinition == nil {
		panic("Pea Definition must not be nil")
	}
	return &PeaDefinitionHolder{
		peaName,
		peaDefinition,
		aliases,
	}
}

func NewPeaDefinitionHolderWithHolder(peaDefinitionHolder *PeaDefinitionHolder) *PeaDefinitionHolder {
	if peaDefinitionHolder == nil {
		panic("Pea Definition Holder must not be nil")
	}
	return NewPeaDefinitionHolderWithAliases(peaDefinitionHolder.peaName,
		peaDefinitionHolder.peaDefinition,
		peaDefinitionHolder.aliases,
	)
}

func (holder *PeaDefinitionHolder) GetPeaName() string {
	return holder.peaName
}

func (holder *PeaDefinitionHolder) GetPeaDefinition() PeaDefinition {
	return holder.peaDefinition
}

func (holder *PeaDefinitionHolder) GetAliases() []string {
	return holder.aliases
}
