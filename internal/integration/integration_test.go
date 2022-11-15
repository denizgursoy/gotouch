//go:build integration_test
// +build integration_test

package integration

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/suite"
)

const echoWithVersion = "github.com/labstack/echo/v4 v4.9.1"

var (
	dependencies = []string{
		echoWithVersion,
		"github.com/spf13/viper",
	}
	file = ""
)

type ZippingTestSuite struct {
	suite.Suite
	workingDir         string
	binaryPath         string
	binaryDir          string
	createdProjectPath string
}

func TestUnzipping(t *testing.T) {
	suite.Run(t, new(ZippingTestSuite))
}

func (z *ZippingTestSuite) SetupSuite() {
	err := os.Chdir("../../")
	getwd, _ := os.Getwd()
	z.binaryDir = getwd
	z.binaryPath = getwd + "/gotouch-" + runtime.GOOS
	z.Nil(err, "could not change directory")
}

func (z *ZippingTestSuite) SetupTest() {
	mkdirTemp, _ := os.MkdirTemp("", "gotouch-test*")
	z.createdProjectPath = mkdirTemp + "/" + "testapp"
	z.workingDir = mkdirTemp
	//exec.Command("open", mkdirTemp).Start()
	fmt.Println("binaryDir          -->" + z.binaryDir)
	fmt.Println("binaryPath         -->" + z.binaryPath)
	fmt.Println("workingDir         -->" + z.workingDir)
	fmt.Println("createdProjectPath -->" + z.createdProjectPath)
}

func (z *ZippingTestSuite) TearDownTest() {
	os.RemoveAll(z.workingDir)
}

func (z *ZippingTestSuite) TestUnzipping() {
	z.setInputFile("testapp.txt")
	z.executeGotouch()

	z.checkDefaultProjectStructure()

	z.checkModuleName("module testapp", dependencies)
}

func (z *ZippingTestSuite) TestGithub() {
	z.setInputFile("github-full-name.txt")
	z.executeGotouch()

	z.checkDefaultProjectStructure()
}

func (z *ZippingTestSuite) TestDelimiter() {
	z.setInputFile("delimiter.txt")
	z.executeGotouch()

	z.checkFileContent("main.go", "delimiter-main.go")
	z.checkFileContent("testapp/main.go", "delimiter-main.go")
}

func (z *ZippingTestSuite) checkDefaultProjectStructure() {
	directories := make([]string, 0)
	directories = append(directories, "api", "build", "cmd", "configs", "deployments", "web")
	directories = append(directories, "init", "internal", "pkg", "configs", "test", "vendor", "cmd/testapp/")
	z.checkDirectoriesExist(directories)

	files := make([]string, 0)
	files = append(files, "cmd/testapp/main.go", "go.mod", "Dockerfile")
	z.checkFilesExist(files)
	z.checkFileContent("Dockerfile", "Dockerfile")
	z.checkFileContent("test.txt", "test.txt")
}

func (z *ZippingTestSuite) checkFileContent(fileName, expectedFile string) {
	actualFilePath := fmt.Sprintf("%s/%s", z.createdProjectPath, fileName)
	actualFileContent, err := os.ReadFile(actualFilePath)
	z.Nil(err)

	expectedFilePath := fmt.Sprintf("%s/internal/testdata/%s", z.binaryDir, expectedFile)
	expectedFileContent, err := os.ReadFile(expectedFilePath)
	z.Nil(err)

	z.EqualValues(expectedFileContent, actualFileContent)
}

func (z *ZippingTestSuite) checkModuleName(expectedModuleName string, dependencies []string) {
	open, err := os.ReadFile(fmt.Sprintf("%s/go.mod", z.createdProjectPath))
	z.Nil(err, "go module file not found")

	moduleContent := string(open)
	split := strings.Split(moduleContent, "\n")

	z.EqualValues(expectedModuleName, split[0], "Module name did not change: expected: %s, actual: %s", expectedModuleName, split[0])

	for _, dependency := range dependencies {
		z.True(strings.Contains(moduleContent, dependency))
	}

}

func (z *ZippingTestSuite) checkDirectoriesExist(directories []string) {
	for _, directory := range directories {
		directoryPath := fmt.Sprintf("%s/%s", z.createdProjectPath, directory)
		stat, err := os.Stat(directoryPath)
		z.Nil(err, "%s does not exists", directory)
		z.True(stat.IsDir(), "%s does not exists", directory)
	}
}

func (z *ZippingTestSuite) checkFilesExist(files []string) {
	for _, file := range files {
		stat, err := os.Stat(fmt.Sprintf("%s/%s", z.createdProjectPath, file))
		z.Nil(err, "%s does not exists", file)
		z.False(stat.IsDir(), "%s does not exists", file)
	}
}

func (z *ZippingTestSuite) CmdExec(args ...string) {

	baseCmd := args[0]
	cmdArgs := args[1:]

	cmd := exec.Command(baseCmd, cmdArgs...)
	cmd.Env = os.Environ()

	env1 := fmt.Sprintf("%s=%s", "TARGET_DIRECTORY", z.workingDir)
	env2 := fmt.Sprintf("%s=%s", "TARGET_FILE", file)

	//fmt.Println(env1 + " " + env2 + " " + z.binaryPath)

	cmd.Env = append(cmd.Env, env1)
	cmd.Env = append(cmd.Env, env2)

	out, err := cmd.Output()
	if err != nil {
		println(err)
	}
	fmt.Println(string(out))
}

func (z *ZippingTestSuite) executeGotouch() {
	z.CmdExec(z.binaryPath)
}

func (z *ZippingTestSuite) setInputFile(fileName string) {
	source := fmt.Sprintf("%s/internal/testdata/%s", z.binaryDir, fileName)
	file = source
}
