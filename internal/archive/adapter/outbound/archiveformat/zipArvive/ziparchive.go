package ziparchive

type Archiver struct {
	path []string
}

func NewArchiver(path []string) Archiver {
	return Archiver{
		path: path,
	}
}

func Archive(path []string) (string, error) {
	result := ""
	return result, nil
}
