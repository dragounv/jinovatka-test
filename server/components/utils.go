package components

import (
	"jinovatka/entities"
)

func prettyPrintCaptureState(state entities.CaptureState) string {
	switch state {
	case entities.NotEnqueued:
		return "Nezařazeno"
	case entities.Pending:
		return "Čeká na sklizení"
	case entities.DoneSuccess:
		return "Úspěšně sklizeno"
	case entities.DoneFailure:
		return "Chyba při sklizni"
	}
	return "Neznámý stav"
}
