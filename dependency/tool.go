package dependency

import (
	"github.com/SananGuliyev/goddd/domain/tool"
	tool2 "github.com/SananGuliyev/goddd/infrastructure/tool"
)

func NewIdGenerator() tool.IdGenerator {
	return tool2.NewUuidGenerator()
}
