package data

var absRootDir string

func ChRoot(d string) {
	absRootDir = d
}

func Pwd() string {
	return absRootDir
}
