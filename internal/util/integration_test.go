//+build integration

package util

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/denizgursoy/gotouch/internal/manager"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
)

type ZippingTestSuite struct {
	suite.Suite
	c                   testcontainers.Container
	execPath            string
	mountPath           string
	containerWorkingDir string
}

func TestUnzipping(t *testing.T) {
	suite.Run(t, new(ZippingTestSuite))
}

func (z *ZippingTestSuite) SetupSuite() {
	err := os.Chdir("../../")
	if err != nil {
		z.T().Fatal("could not change directory", err)
	}
}

func (z *ZippingTestSuite) SetupTest() {
	z.containerWorkingDir = "/go/test"
	z.execPath = manager.GetInstance().GetWd()

	temp, uuidErr := ioutil.TempDir("", "gotouch-test*")
	if uuidErr != nil {
		z.T().Fatal("could not create directory", temp, uuidErr)
	}

	z.mountPath = temp

	z.T().Log("mount:", z.mountPath)

	binaryName := "gotouch-linux-test"
	sourcePath := fmt.Sprintf("%s/%s", z.execPath, binaryName)
	targetPath := fmt.Sprintf("%s/gotouch", z.mountPath)

	if _, err := z.copy(sourcePath, targetPath); err != nil {
		z.T().Fatal("could not copy the binary", err)
	}

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

	z.T().Log("commander:", fmt.Sprintf("docker exec -it %s /bin/bash", cnt.GetContainerID()))

	if cnt == nil {
		z.T().Fatal("could not create cnt", uuidErr)
	}

	z.c = cnt
}

func (z *ZippingTestSuite) TestUnzipping() {
	z.moveToDirectory("unzip-test.txt")
	z.executeCommand()

	folderName := "unzip"

	z.checkDefaultProjectStructure(folderName)

	z.checkModuleName("module unzip", folderName)
}

func (z *ZippingTestSuite) TestGithub() {
	z.moveToDirectory("github-full-name.txt")
	z.executeCommand()

	folderName := "call"

	z.checkDefaultProjectStructure(folderName)
	z.checkModuleName("module g.c/dg/call", folderName)
}

func (z *ZippingTestSuite) checkDefaultProjectStructure(folderName string) {
	directories := make([]string, 0)
	directories = append(directories, "api", "build", "cmd", "configs", "deployments", "web")
	directories = append(directories, "init", "internal", "pkg", "configs", "test", "vendor")
	z.checkDirectoriesExist(directories, folderName)

	files := make([]string, 0)
	files = append(files, "cmd/main.go", "go.mod")
	z.checkFilesExist(files, folderName)
}

func (z *ZippingTestSuite) checkModuleName(expectedModuleName, folderName string) {
	open, err := os.ReadFile(fmt.Sprintf("%s/%s/go.mod", z.mountPath, folderName))
	if err != nil {
		z.T().Fatal("go module file not found")
	}
	split := strings.Split(string(open), "\n")
	if split[0] != expectedModuleName {
		z.T().Fatalf("Module name did not change: expected: %s, actual: %s", expectedModuleName, split[0])
	}
}

func (z *ZippingTestSuite) checkDirectoriesExist(directories []string, folderName string) {
	for _, directory := range directories {
		if stat, err := os.Stat(fmt.Sprintf("%s/%s/%s", z.mountPath, folderName, directory)); err != nil || !stat.IsDir() {
			z.T().Fatalf("%s does not exists", directory)
		}
	}
}

func (z *ZippingTestSuite) checkFilesExist(files []string, folderName string) {
	for _, file := range files {
		if stat, err := os.Stat(fmt.Sprintf("%s/%s/%s", z.mountPath, folderName, file)); err != nil || stat.IsDir() {
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
	source := fmt.Sprintf("%s/internal/testdata/%s", z.execPath, fileName)
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
