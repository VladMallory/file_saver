package archivecore

import "errors"

var (
	ErrNoPaths    = errors.New("no files to archive")
	ErrNoFindFile = errors.New("не удалось определить путь к бинарнику")
)
