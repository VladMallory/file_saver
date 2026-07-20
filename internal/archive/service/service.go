package archiveusecase

import (
	archivecore "saveFile/internal/archive/domain"

	"go.uber.org/zap"
)

type PathProvider interface {
	GetPath() (resultPath []string, err error)
}

type Archiver interface {
	Run(path []string, savePathFile string) (string, error)
}

type ArchiveService struct {
	log          *zap.Logger
	pathProvider PathProvider
	archiver     Archiver
}

func NewArchiveService(
	log *zap.Logger,
	pathProvider PathProvider,
	archiver Archiver,
) ArchiveService {
	return ArchiveService{
		log:          log,
		pathProvider: pathProvider,
		archiver:     archiver,
	}
}

func (a ArchiveService) Run(savePathFile string) error {
	path, err := a.pathProvider.GetPath()
	if err != nil {
		return err
	}

	if len(path) == 0 {
		return archivecore.ErrNoPaths
	}

	_, err = a.archiver.Run(path, savePathFile)
	if err != nil {
		return err
	}

	return nil
}
