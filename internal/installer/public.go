package installer

// PathManager публичный интерфейс для работы с файлом путей.
type PathManager interface {
	GetPaths() ([]string, error)
	WritePaths(paths []string) error
	FileExists() bool
}

// DefaultPathManager реализация PathManager.
type DefaultPathManager struct{}

// NewPathManager создает новый PathManager.
func NewPathManager() PathManager {
	return &DefaultPathManager{}
}

// GetPaths возвращает пути из файла path.txt.
func (pm *DefaultPathManager) GetPaths() ([]string, error) {
	return ReadPathFile()
}

// WritePaths записывает пути в файл path.txt.
func (pm *DefaultPathManager) WritePaths(paths []string) error {
	return WritePathFile(paths)
}

// FileExists проверяет существует ли файл path.txt.
func (pm *DefaultPathManager) FileExists() bool {
	return PathFileExists()
}

// InitializeApp инициализирует приложение - создает файл path.txt если его нет.
func InitializeApp() error {
	return EnsurePathFileExists()
}
