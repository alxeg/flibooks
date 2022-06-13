package inpx

import (
	"github.com/alxeg/flibooks/pkg/inpx/models"
)

type Processor interface {
	ProcessBook(*models.Book) error
	FinishProcessing()
}
