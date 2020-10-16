package peas

type PeaPreparationError struct {
	peaName string
	message string
}

func NewPeaPreparationError(peaName string, message string) PeaPreparationError {
	return PeaPreparationError{peaName, message}
}

func (error PeaPreparationError) GetPeaName() string {
	return error.peaName
}

func (error PeaPreparationError) GetMessage() string {
	return error.message
}

func (error PeaPreparationError) Error() string {
	return error.peaName + " : " + error.peaName
}

type PeaInPreparationError struct {
	PeaPreparationError
}

func NewPeaInPreparationError(peaName string) PeaInPreparationError {
	return PeaInPreparationError{
		NewPeaPreparationError(peaName,
			"Pea is currently in preparation, maybe it has got circular dependency cycle",
		),
	}
}
