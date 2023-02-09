//go:build integration_test

package integration

import (
	"fmt"
	"github.com/denizgursoy/gotouch/internal/config"
	"log"
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
	file          = ""
	PropertiesUrl = "https://raw.githubusercontent.com/denizgursoy/go-touch-projects/main/test/package.yaml"
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

	err := os.Chdir(mkdirTemp)
	if err != nil {
		log.Fatalln("could not change directory")
	}

	//exec.Command("open", mkdirTemp).Start()
	fmt.Println("binaryDir          -->" + z.binaryDir)
	fmt.Println("binaryPath         -->" + z.binaryPath)
	fmt.Println("workingDir         -->" + z.workingDir)
	fmt.Println("createdProjectPath -->" + z.createdProjectPath)
}

func (z *ZippingTestSuite) TearDownTest() {
	err := os.RemoveAll(z.workingDir)
	if err != nil {
		return
	}
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

	z.checkDirectoriesExist([]string{"testapp/a/testapp"})
}

func (z *ZippingTestSuite) TestGitCheckout() {
	z.setInputFile("git-checkout-main.txt")
	z.executeGotouch()

	z.checkDirectoriesExist([]string{"test"})
}

func (z *ZippingTestSuite) TestGitCheckoutBranch() {
	z.setInputFile("git-checkout-branch.txt")
	z.executeGotouch()

	z.checkFileContent("test-branch-content.txt", "test-branch-content.txt")

	// check init files are called and deleted
	z.checkFileExists("init.sh", false)
	z.checkFileExists("init.bat", false)
	z.checkFileExists("test-linux.txt", true)
	z.checkFileExists("test-windows.txt", false)
}

func (z *ZippingTestSuite) TestMultipleSelectQuestion() {
	z.setInputFile("multiple-select.txt")
	z.executeGotouch()

	z.checkFileExists("app-deployment.yaml", true)
	z.checkFileExists("Dockerfile", true)
}

func (z *ZippingTestSuite) TestProjectWithNoUrl() {
	z.setInputFile("project-with-no-url.txt")
	z.executeGotouch()

	z.checkFileExists("app-deployment.yaml", true)
	z.checkFileExists("Dockerfile", true)
	z.checkFileExists("Makefile", true)
}

func (z *ZippingTestSuite) TestProjectWithFilesAndDependencies() {
	z.setInputFile("project-with-files-and-dependencies.txt")
	z.executeGotouch()

	z.checkFileContent("Dockerfile", "Dockerfile")
	z.checkFileContent("values.txt", "values.txt")
}

func (z *ZippingTestSuite) TestProjectConfig() {
	name, err := config.GetFileName()
	z.Nil(err)

	if _, err = os.Stat(name); err == nil {
		err = os.Remove(name)
	}

	z.executeGotouchWithArgs("config", "set", "url", "test-url")
	expectedFilePath := fmt.Sprintf("%s/internal/testdata/%s", z.binaryDir, config.ConfigFileName)
	z.checkFileContentsWithAbsPath(name, expectedFilePath)

	z.executeGotouchWithArgs("config", "unset", "url")
	expectedFilePath = fmt.Sprintf("%s/internal/testdata/%s", z.binaryDir, config.ConfigFileName+"-empty")
	z.checkFileContentsWithAbsPath(name, expectedFilePath)
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

func (z *ZippingTestSuite) checkFileExists(fileName string, exists bool) {
	actualFilePath := fmt.Sprintf("%s/%s", z.createdProjectPath, fileName)
	_, err2 := os.Stat(actualFilePath)
	if exists {
		z.Nil(err2)
	} else {
		z.NotNil(err2)
	}
}

func (z *ZippingTestSuite) checkFileContent(fileName, expectedFile string) {
	actualFilePath := fmt.Sprintf("%s/%s", z.createdProjectPath, fileName)
	expectedFilePath := fmt.Sprintf("%s/internal/testdata/%s", z.binaryDir, expectedFile)
	z.checkFileContentsWithAbsPath(actualFilePath, expectedFilePath)
}

func (z *ZippingTestSuite) checkFileContentsWithAbsPath(actualFilePath, expectedFilePath string) {
	actualFileContent, err := os.ReadFile(actualFilePath)
	z.Nil(err)
	expectedFileContent, err := os.ReadFile(expectedFilePath)
	z.Nil(err)
	z.EqualValues(actualFileContent, expectedFileContent)
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
	cmd := exec.Command(args[0], args[1:]...)

	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr

	cmd.Env = os.Environ()

	env2 := fmt.Sprintf("%s=%s", "TARGET_FILE", file)

	//fmt.Println(env1 + " " + env2 + " " + z.binaryPath)

	cmd.Env = append(cmd.Env, env2)

	err := cmd.Run()
	if err != nil {
		println(err)
	}
}

func (z *ZippingTestSuite) executeGotouch() {
	args := make([]string, 0)
	args = append(args, z.binaryPath, "-f", PropertiesUrl)
	z.CmdExec(args...)
}

func (z *ZippingTestSuite) executeGotouchWithArgs(gotouchArgs ...string) {
	args := make([]string, 0)
	args = append(args, z.binaryPath)
	args = append(args, gotouchArgs...)
	z.CmdExec(args...)
}

func (z *ZippingTestSuite) setInputFile(fileName string) {
	source := fmt.Sprintf("%s/internal/testdata/%s", z.binaryDir, fileName)
	file = source
}
