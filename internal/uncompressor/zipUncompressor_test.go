//+build integration

package uncompressor

import (
	"context"
	"fmt"
	"github.com/denizgursoy/gotouch/internal/prompts"
	"github.com/hashicorp/go-uuid"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"io"
	"os"
	"testing"
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

}

func (z *ZippingTestSuite) SetupTest() {
	err := os.Chdir("../../")

	z.containerWorkingDir = "/go/test"
	z.execPath = prompts.GetWd()

	generateUUID, err2 := uuid.GenerateUUID()
	if err2 != nil {
		z.T().Fatal("could not create cnt", err)
	}
	z.mountPath = fmt.Sprintf("%s%s", os.TempDir(), generateUUID)
	z.T().Log("mount:", z.mountPath)
	err = os.Mkdir(z.mountPath, os.ModePerm)
	if err != nil {
		z.T().Fatal("could not create cnt", err)
	}
	binaryName := "gotouch-linux-test"
	sourcePath := fmt.Sprintf("%s/%s", z.execPath, binaryName)
	targetPath := fmt.Sprintf("%s/gotouch", z.mountPath)

	i, err := z.copy(sourcePath, targetPath)
	if err != nil {
		z.T().Fatal("could not create cnt", err, i)
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

	cnt, err := testcontainers.GenericContainer(context.Background(), testcontainers.GenericContainerRequest{
		ContainerRequest: request,
		Started:          true,
	})

	z.T().Log("command:", fmt.Sprintf("docker exec -it %s /bin/bash", cnt.GetContainerID()))

	if cnt == nil {
		z.T().Fatal("could not create cnt", err)
	}

	z.c = cnt
}

func (z *ZippingTestSuite) TestUnzipping() {
	z.moveToDirectory("unzip-test.txt")
	z.executeCommand()
}

func (z *ZippingTestSuite) TestGithub() {
	z.moveToDirectory("github-full-name.txt")
	z.executeCommand()
}

func (z *ZippingTestSuite) executeCommand() {
	sprintf := fmt.Sprintf("%s/gotouch", z.containerWorkingDir)
	command := []string{sprintf}
	i, err := z.c.Exec(context.Background(), command)
	if err != nil {
		z.T().Fatal("could not create cont", err, i)
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
