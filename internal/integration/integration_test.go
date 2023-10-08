//go:build integration_test

package integration

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/denizgursoy/gotouch/internal/config"
	"github.com/denizgursoy/gotouch/internal/requirements"

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
	inline             bool
}

func TestUnzipping(t *testing.T) {
	suite.Run(t, new(ZippingTestSuite))
}

func (z *ZippingTestSuite) SetupSuite() {
	err := os.Chdir("../../")
	getwd, _ := os.Getwd()
	z.binaryDir = getwd
	z.binaryPath = filepath.Join(getwd, "gotouch-"+runtime.GOOS)
	z.Nil(err, "could not change directory")
}

func (z *ZippingTestSuite) SetupTest() {
	mkdirTemp, _ := os.MkdirTemp(z.T().TempDir(), "gotouch-test*")

	z.createdProjectPath = filepath.Join(mkdirTemp, "testapp")
	z.workingDir = mkdirTemp

	// reset inline
	z.inline = false

	err := os.Chdir(mkdirTemp)
	if err != nil {
		log.Fatalln("could not change directory")
	}

	// exec.Command("open", mkdirTemp).Start()
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

func (z *ZippingTestSuite) TestProjectWithFilesInline() {
	z.setInputFile("project-with-files-inline.txt")
	z.setInline()
	z.executeGotouch()

	z.checkFileContent("Dockerfile", "Dockerfile")
	z.checkFileContent("values.txt", "values.txt")
}

func (z *ZippingTestSuite) TestProjectWithPropertiesYaml() {
	z.setInputFile("project-with-files-and-dependencies.txt")
	z.executeGotouch()

	z.checkFileExists(requirements.PropertiesYamlName, false)
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
	_, err2 := os.Stat(filepath.Join(z.createdProjectPath, fileName))
	if exists {
		z.NoError(err2)
	} else {
		z.Error(err2)
	}
}

func (z *ZippingTestSuite) checkFilesExist(files []string) {
	for _, file := range files {
		z.FileExists(filepath.Join(z.createdProjectPath, file))
	}
}

func (z *ZippingTestSuite) checkFileContent(fileName, expectedFile string) {
	var actualFilePath string
	if z.inline {
		dir, _ := filepath.Split(z.createdProjectPath)
		actualFilePath = filepath.Join(dir, fileName)
	} else {
		actualFilePath = filepath.Join(z.createdProjectPath, fileName)
	}

	expectedFilePath := filepath.Join(z.binaryDir, "internal", "testdata", expectedFile)
	z.checkFileContentsWithAbsPath(actualFilePath, expectedFilePath)
}

func (z *ZippingTestSuite) checkFileContentsWithAbsPath(actualFilePath, expectedFilePath string) {
	actualFileContent, err := os.ReadFile(actualFilePath)
	z.NoError(err)
	expectedFileContent, err := os.ReadFile(expectedFilePath)
	z.NoError(err)
	z.EqualValues(actualFileContent, expectedFileContent)
}

func (z *ZippingTestSuite) checkModuleName(expectedModuleName string, dependencies []string) {
	open, err := os.ReadFile(filepath.Join(z.createdProjectPath, "go.mod"))
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
		directoryPath := filepath.Join(z.createdProjectPath, directory)
		stat, err := os.Stat(directoryPath)
		z.Nil(err, "%s does not exists", directory)
		z.True(stat.IsDir(), "%s does not exists", directory)
	}
}

func (z *ZippingTestSuite) CmdExec(args ...string) {
	cmd := exec.Command(args[0], args[1:]...)

	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr

	cmd.Env = os.Environ()

	env2 := fmt.Sprintf("%s=%s", "TARGET_FILE", file)

	// fmt.Println(env1 + " " + env2 + " " + z.binaryPath)

	cmd.Env = append(cmd.Env, env2)

	err := cmd.Run()
	if err != nil {
		println(err)
	}
}

func (z *ZippingTestSuite) executeGotouch() {
	args := make([]string, 0)
	args = append(args, z.binaryPath, "-f", PropertiesUrl)
	if z.inline {
		args = append(args, "-i")
	}
	z.CmdExec(args...)
}

func (z *ZippingTestSuite) executeGotouchWithArgs(gotouchArgs ...string) {
	args := make([]string, 0)
	args = append(args, z.binaryPath)
	args = append(args, gotouchArgs...)
	z.CmdExec(args...)
}

func (z *ZippingTestSuite) setInputFile(fileName string) {
	file = filepath.Join(z.binaryDir, "internal", "testdata", fileName)
}

func (z *ZippingTestSuite) setInline() {
	z.inline = true
}
