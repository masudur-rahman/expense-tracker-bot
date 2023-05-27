package pkg

import (
	"path/filepath"
	"runtime"
)

var (
	_, filePath, _, _ = runtime.Caller(0)
	ProjectDirectory  = filepath.Dir(filepath.Dir(filePath))
)
