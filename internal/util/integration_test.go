//+build integration_test

package util

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
)

type ZippingTestSuite struct {
	suite.Suite
	c                   testcontainers.Container
	mountPath           string
	containerWorkingDir string
	folderName          string
}

func TestUnzipping(t *testing.T) {
	suite.Run(t, new(ZippingTestSuite))
}

func (z *ZippingTestSuite) SetupSuite() {
	z.folderName = "testapp"
	err := os.Chdir("../../")
	z.Nil(err, "could not change directory")
}

func (z *ZippingTestSuite) SetupTest() {
	z.containerWorkingDir = "/go/test"

	temp, uuidErr := ioutil.TempDir("", "gotouch-test*")
	z.Nil(uuidErr, "could not create directory")

	z.mountPath = temp

	z.T().Log("mount:", z.mountPath)
	getwd := getWorkingDirectory()

	binaryName := "gotouch-linux-test"
	sourcePath := fmt.Sprintf("%s/%s", getwd, binaryName)
	targetPath := fmt.Sprintf("%s/gotouch", z.mountPath)

	_, err := z.copy(sourcePath, targetPath)
	z.Nil(err, "could not copy the binary")

	request := testcontainers.ContainerRequest{
		Image: "golang:latest",
		Cmd:   []string{"sleep", "600000"},
		Mounts: testcontainers.ContainerMounts{
			{
				Source: testcontainers.GenericBindMountSource{
					HostPath: z.mountPath,
				},
				Target:   testcontainers.ContainerMountTarget(z.containerWorkingDir),
				ReadOnly: false,
			},
		},
	}

	cnt, uuidErr := testcontainers.GenericContainer(context.Background(), testcontainers.GenericContainerRequest{
		ContainerRequest: request,
		Started:          true,
	})

	z.NotNil(cnt, "Make sure docker is running")

	z.T().Log("commander:", fmt.Sprintf("docker exec -it %s /bin/bash", cnt.GetContainerID()))

	z.c = cnt
}

func getWorkingDirectory() string {
	getwd, _ := os.Getwd()
	return getwd
}

func (z *ZippingTestSuite) TestUnzipping() {
	z.moveToDirectory("testapp.txt")
	z.executeCommand()

	z.checkDefaultProjectStructure()

	z.checkModuleName("module testapp")
}

func (z *ZippingTestSuite) TestGithub() {
	z.moveToDirectory("github-full-name.txt")
	z.executeCommand()

	z.checkDefaultProjectStructure()
	z.checkModuleName("module g.c/dg/testapp")
}

func (z *ZippingTestSuite) checkDefaultProjectStructure() {
	directories := make([]string, 0)
	directories = append(directories, "api", "build", "cmd", "configs", "deployments", "web")
	directories = append(directories, "init", "internal", "pkg", "configs", "test", "vendor")
	z.checkDirectoriesExist(directories)

	files := make([]string, 0)
	files = append(files, "cmd/main.go", "go.mod", "Dockerfile")
	z.checkFilesExist(files)
	z.checkFileContent("Dockerfile", "Dockerfile")
}

func (z *ZippingTestSuite) checkFileContent(fileName, expectedFile string) {
	actualFilePath := fmt.Sprintf("%s/%s/%s", z.mountPath, z.folderName, fileName)
	actualFileContent, err := ioutil.ReadFile(actualFilePath)
	z.Nil(err)

	expectedFilePath := fmt.Sprintf("%s/internal/testdata/%s", getWorkingDirectory(), expectedFile)
	expectedFileContent, err := ioutil.ReadFile(expectedFilePath)
	z.Nil(err)

	z.EqualValues(expectedFileContent, actualFileContent)
}

func (z *ZippingTestSuite) checkModuleName(expectedModuleName string) {
	open, err := os.ReadFile(fmt.Sprintf("%s/%s/go.mod", z.mountPath, z.folderName))
	z.Nil(err, "go module file not found")

	split := strings.Split(string(open), "\n")
	if split[0] != expectedModuleName {
		z.T().Fatalf("Module name did not change: expected: %s, actual: %s", expectedModuleName, split[0])
	}
}

func (z *ZippingTestSuite) checkDirectoriesExist(directories []string) {
	for _, directory := range directories {
		directoryPath := fmt.Sprintf("%s/%s/%s", z.mountPath, z.folderName, directory)
		if stat, err := os.Stat(directoryPath); err != nil || !stat.IsDir() {
			z.T().Fatalf("%s does not exists", directory)
		}
	}
}

func (z *ZippingTestSuite) checkFilesExist(files []string) {
	for _, file := range files {
		if stat, err := os.Stat(fmt.Sprintf("%s/%s/%s", z.mountPath, z.folderName, file)); err != nil || stat.IsDir() {
			z.T().Fatalf("%s does not exists", file)
		}
	}
}

func (z *ZippingTestSuite) executeCommand() {
	sprintf := fmt.Sprintf("%s/gotouch", z.containerWorkingDir)
	commander := []string{sprintf}
	i, err := z.c.Exec(context.Background(), commander)
	if err != nil {
		z.T().Fatal("could not execute commander", err, i)
	}
}

func (z *ZippingTestSuite) moveToDirectory(fileName string) {
	source := fmt.Sprintf("%s/internal/testdata/%s", getWorkingDirectory(), fileName)
	target := fmt.Sprintf("%s/input.txt", z.mountPath)
	i, err := z.copy(source, target)
	if err != nil {
		fmt.Println(err, i)
	}
}

func (z *ZippingTestSuite) copy(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	err = os.Chmod(dst, os.ModePerm)
	return nBytes, err
}
