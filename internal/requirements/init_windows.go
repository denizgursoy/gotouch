package requirements

const InitFileName = "init.bat"

func executeInitFile(address string) {
	commandData := CommandData{
		Command: InitFileName,
	}
	RunCommand(&commandData, str)
}
