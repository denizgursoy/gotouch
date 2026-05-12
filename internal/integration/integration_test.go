//go:build integration_test

package integration

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/creack/pty"
	"github.com/denizgursoy/gotouch/internal/config"
	"github.com/denizgursoy/gotouch/internal/requirements"

	"github.com/stretchr/testify/suite"
)

// Key constants for PTY input simulation.
const (
	keyEnter = "\r"
	keyDown  = "\x1b[B"
	keySpace = " "
	keyCtrlU = "\x15"
	keyY     = "y"
	keyN     = "n"
)

const echoWithVersion = "github.com/labstack/echo/v4 v4.9.1"

var (
	dependencies = []string{
		echoWithVersion,
		"github.com/spf13/viper",
	}
	PropertiesUrl = "https://raw.githubusercontent.com/denizgursoy/go-touch-projects/main/test/package.yaml"
	ansiRegex     = regexp.MustCompile(`\x1b\[[0-9;]*[a-zA-Z]|\x1b\][^\x07]*\x07|\x1b\[\?[0-9;]*[hlsr]|\x1b[=>]`)
)

// ptyRunner manages a gotouch subprocess attached to a pseudo-terminal.
type ptyRunner struct {
	ptmx   *os.File
	cmd    *exec.Cmd
	output *bytes.Buffer
	mu     sync.Mutex
	t      *testing.T
}

func stripANSI(s string) string {
	return ansiRegex.ReplaceAllString(s, "")
}

// waitForOutput polls the PTY output buffer until text appears or timeout is reached.
func (r *ptyRunner) waitForOutput(text string, timeout time.Duration) {
	r.t.Helper()
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		r.mu.Lock()
		output := stripANSI(r.output.String())
		r.mu.Unlock()
		if strings.Contains(output, text) {
			// Clear buffer to avoid matching stale content on next wait
			r.mu.Lock()
			r.output.Reset()
			r.mu.Unlock()
			return
		}
		time.Sleep(50 * time.Millisecond)
	}
	r.mu.Lock()
	raw := r.output.String()
	r.mu.Unlock()
	r.t.Fatalf("timeout (%s) waiting for %q in output:\n%s", timeout, text, stripANSI(raw))
}

// sendKeys writes raw key bytes to the PTY.
func (r *ptyRunner) sendKeys(keys ...string) {
	for _, key := range keys {
		_, _ = r.ptmx.WriteString(key)
		time.Sleep(30 * time.Millisecond)
	}
}

// sendText types each character individually.
func (r *ptyRunner) sendText(text string) {
	for _, ch := range text {
		_, _ = r.ptmx.WriteString(string(ch))
		time.Sleep(10 * time.Millisecond)
	}
}

// selectIndex navigates down n times in a select list and presses Enter.
func (r *ptyRunner) selectIndex(n int) {
	for range n {
		r.sendKeys(keyDown)
	}
	r.sendKeys(keyEnter)
}

// enterText clears the current input (Ctrl+U), types the text, and presses Enter.
func (r *ptyRunner) enterText(text string) {
	r.sendKeys(keyCtrlU)
	r.sendText(text)
	r.sendKeys(keyEnter)
}

// answerYes sends 'y' to a yes/no confirm prompt.
func (r *ptyRunner) answerYes() {
	r.sendKeys(keyY)
}

// answerNo sends 'n' to a yes/no confirm prompt.
func (r *ptyRunner) answerNo() {
	r.sendKeys(keyN)
}

// multiSelect toggles the given indices (0-based) in a multi-select list and submits.
// It starts at index 0 and navigates down, pressing Space on the specified indices.
func (r *ptyRunner) multiSelect(selections ...int) {
	selected := make(map[int]bool)
	maxIdx := 0
	for _, idx := range selections {
		selected[idx] = true
		if idx > maxIdx {
			maxIdx = idx
		}
	}

	for i := 0; i <= maxIdx; i++ {
		if selected[i] {
			r.sendKeys(keySpace)
		}
		if i < maxIdx {
			r.sendKeys(keyDown)
		}
	}
	r.sendKeys(keyEnter)
}

// wait waits for the subprocess to exit.
func (r *ptyRunner) wait() {
	_ = r.cmd.Wait()
	_ = r.ptmx.Close()
}

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

// startGotouchWithPTY starts the gotouch binary attached to a PTY for keystroke simulation.
func (z *ZippingTestSuite) startGotouchWithPTY(args ...string) *ptyRunner {
	cmd := exec.Command(z.binaryPath, args...)
	cmd.Dir = z.workingDir
	cmd.Env = os.Environ()

	winSize := &pty.Winsize{Rows: 40, Cols: 120}
	ptmx, err := pty.StartWithSize(cmd, winSize)
	z.Require().NoError(err, "failed to start PTY")

	runner := &ptyRunner{
		ptmx:   ptmx,
		cmd:    cmd,
		output: &bytes.Buffer{},
		t:      z.T(),
	}

	// Read PTY output in background
	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := ptmx.Read(buf)
			if n > 0 {
				runner.mu.Lock()
				runner.output.Write(buf[:n])
				runner.mu.Unlock()

				// Show PTY output in test console when running with -v
				if testing.Verbose() {
					os.Stdout.Write(buf[:n])
				}
			}
			if err != nil {
				return
			}
		}
	}()

	return runner
}

// --- Test Methods ---

func (z *ZippingTestSuite) TestUnzipping() {
	runner := z.startGotouchWithPTY("-f", PropertiesUrl)

	runner.waitForOutput("Select Project Type", 30*time.Second)
	runner.selectIndex(1) // Default Project Layout

	runner.waitForOutput("Module Name", 5*time.Second)
	runner.enterText("testapp")

	runner.waitForOutput("Do you want Dockerfile", 5*time.Second)
	runner.answerYes()

	runner.waitForOutput("Kubernetes", 5*time.Second)
	runner.answerNo()

	runner.waitForOutput("makefile", 5*time.Second)
	runner.answerNo()

	runner.wait()

	z.checkDefaultProjectStructure()
	z.checkModuleName("module testapp", dependencies)
}

func (z *ZippingTestSuite) TestGithub() {
	runner := z.startGotouchWithPTY("-f", PropertiesUrl)

	runner.waitForOutput("Select Project Type", 30*time.Second)
	runner.selectIndex(1) // Default Project Layout

	runner.waitForOutput("Module Name", 5*time.Second)
	runner.enterText("g.c/dg/testapp")

	runner.waitForOutput("Do you want Dockerfile", 5*time.Second)
	runner.answerYes()

	runner.waitForOutput("Kubernetes", 5*time.Second)
	runner.answerNo()

	runner.waitForOutput("makefile", 5*time.Second)
	runner.answerNo()

	runner.wait()

	z.checkDefaultProjectStructure()
}

func (z *ZippingTestSuite) TestDelimiter() {
	runner := z.startGotouchWithPTY("-f", PropertiesUrl)

	runner.waitForOutput("Select Project Type", 30*time.Second)
	runner.selectIndex(2) // Delimiter test project

	runner.waitForOutput("Module Name", 5*time.Second)
	runner.enterText("g.c/dg/testapp")

	runner.waitForOutput("Do you want Dockerfile", 5*time.Second)
	runner.answerYes()

	runner.waitForOutput("Kubernetes", 5*time.Second)
	runner.answerNo()

	runner.waitForOutput("makefile", 5*time.Second)
	runner.answerNo()

	runner.wait()

	z.checkFileContent("main.go", "delimiter-main.go")
	z.checkFileContent("testapp/main.go", "delimiter-main.go")
	z.checkDirectoriesExist([]string{"testapp/a/testapp"})
}

func (z *ZippingTestSuite) TestGitCheckout() {
	runner := z.startGotouchWithPTY("-f", PropertiesUrl)

	runner.waitForOutput("Select Project Type", 30*time.Second)
	runner.selectIndex(3) // gitcheckout main

	runner.waitForOutput("Module Name", 5*time.Second)
	runner.enterText("g.c/dg/testapp")

	runner.wait()

	z.checkDirectoriesExist([]string{"test"})
}

func (z *ZippingTestSuite) TestGitCheckoutBranch() {
	runner := z.startGotouchWithPTY("-f", PropertiesUrl)

	runner.waitForOutput("Select Project Type", 30*time.Second)
	runner.selectIndex(4) // gitcheckout branch

	runner.waitForOutput("Module Name", 5*time.Second)
	runner.enterText("g.c/dg/testapp")

	runner.wait()

	z.checkFileContent("test-branch-content.txt", "test-branch-content.txt")

	// check init files are called and deleted
	z.checkFileExists("init.sh", false)
	z.checkFileExists("init.bat", false)
	z.checkFileExists("test-linux.txt", true)
	z.checkFileExists("test-windows.txt", false)
}

func (z *ZippingTestSuite) TestMultipleSelectQuestion() {
	runner := z.startGotouchWithPTY("-f", PropertiesUrl)

	runner.waitForOutput("Select Project Type", 30*time.Second)
	runner.selectIndex(5) // multiple select

	runner.waitForOutput("Module Name", 5*time.Second)
	runner.enterText("g.c/dg/testapp")

	runner.waitForOutput("Select choices", 5*time.Second)
	runner.multiSelect(0, 2) // dockerfile + deployment

	runner.wait()

	z.checkFileExists("app-deployment.yaml", true)
	z.checkFileExists("Dockerfile", true)
}

func (z *ZippingTestSuite) TestProjectWithNoUrl() {
	runner := z.startGotouchWithPTY("-f", PropertiesUrl)

	runner.waitForOutput("Select Project Type", 30*time.Second)
	runner.selectIndex(6) // Project with no url

	runner.waitForOutput("Module Name", 5*time.Second)
	runner.enterText("g.c/dg/testapp")

	runner.waitForOutput("Do you want Dockerfile", 5*time.Second)
	runner.answerYes()

	runner.waitForOutput("Kubernetes", 5*time.Second)
	runner.answerYes()

	runner.waitForOutput("makefile", 5*time.Second)
	runner.answerYes()

	runner.wait()

	z.checkFileExists("app-deployment.yaml", true)
	z.checkFileExists("Dockerfile", true)
	z.checkFileExists("Makefile", true)
}

func (z *ZippingTestSuite) TestProjectWithFilesAndDependencies() {
	runner := z.startGotouchWithPTY("-f", PropertiesUrl)

	runner.waitForOutput("Select Project Type", 30*time.Second)
	runner.selectIndex(7) // Project with files and dependencies

	runner.waitForOutput("Module Name", 5*time.Second)
	runner.enterText("g.c/dg/testapp")

	runner.wait()

	z.checkFileContent("Dockerfile", "Dockerfile")
	z.checkFileContent("values.txt", "values.txt")
}

func (z *ZippingTestSuite) TestProjectWithFilesInline() {
	z.inline = true
	runner := z.startGotouchWithPTY("-f", PropertiesUrl, "-i")

	runner.waitForOutput("Select Project Type", 30*time.Second)
	runner.selectIndex(7) // Project with files and dependencies

	runner.waitForOutput("Module Name", 5*time.Second)
	runner.enterText("g.c/dg/testapp")

	runner.wait()

	z.checkFileContent("Dockerfile", "Dockerfile")
	z.checkFileContent("values.txt", "values.txt")
}

func (z *ZippingTestSuite) TestProjectWithPropertiesYaml() {
	runner := z.startGotouchWithPTY("-f", PropertiesUrl)

	runner.waitForOutput("Select Project Type", 30*time.Second)
	runner.selectIndex(7) // Project with files and dependencies

	runner.waitForOutput("Module Name", 5*time.Second)
	runner.enterText("g.c/dg/testapp")

	runner.wait()

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

func (z *ZippingTestSuite) TestLocalDirectory() {
	packageYaml := z.generateLocalPackageYaml("local-template", false)
	runner := z.startGotouchWithPTY("-f", packageYaml)

	// Single project in YAML → no project selection prompt
	runner.waitForOutput("Module Name", 30*time.Second)
	runner.enterText("testapp")

	runner.wait()

	z.checkFileExists("main.go", true)
	z.checkFileExists("README.md", true)
	z.checkFileExists("config.yaml", true)
}

func (z *ZippingTestSuite) TestLocalCompressedFile() {
	packageYaml := z.generateLocalPackageYaml("local-template.tar.gz", true)
	runner := z.startGotouchWithPTY("-f", packageYaml)

	// Single project in YAML → no project selection prompt
	runner.waitForOutput("Module Name", 30*time.Second)
	runner.enterText("testapp")

	runner.wait()

	z.checkFileExists("main.go", true)
	z.checkFileExists("README.md", true)
	z.checkFileExists("config.yaml", true)
}

// --- Helper Methods ---

func (z *ZippingTestSuite) generateLocalPackageYaml(templateName string, isCompressed bool) string {
	localPath := filepath.Join(z.binaryDir, "internal", "testdata", templateName)

	yamlContent := fmt.Sprintf(`- name: Local Template Project
  localPath: "%s"
`, localPath)

	tmpFile, err := os.CreateTemp(z.workingDir, "package-*.yaml")
	z.Require().NoError(err)

	_, err = tmpFile.WriteString(yamlContent)
	z.Require().NoError(err)

	err = tmpFile.Close()
	z.Require().NoError(err)

	return tmpFile.Name()
}

func (z *ZippingTestSuite) executeGotouchWithArgs(gotouchArgs ...string) {
	args := make([]string, 0)
	args = append(args, z.binaryPath)
	args = append(args, gotouchArgs...)

	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Env = os.Environ()

	err := cmd.Run()
	if err != nil {
		println(err)
	}
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
