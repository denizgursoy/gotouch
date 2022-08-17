//+build integration

package uncompressor

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"io/ioutil"
	"os"
	"os/exec"
	"testing"
)

type ZippingTestSuite struct {
	suite.Suite
	Container testcontainers.Container
}

func TestUnzipping(t *testing.T) {
	suite.Run(t, new(ZippingTestSuite))
}

func (z *ZippingTestSuite) SetupSuite() {

}

func (z *ZippingTestSuite) SetupTest() {

	err := os.Chdir("../../")
	getwd, err := os.Getwd()
	if err != nil {
		println(err)
	}
	getwd = fmt.Sprintf("%s", getwd)

	command := exec.Command("make")
	err = command.Run()

	request := testcontainers.ContainerRequest{
		Image: "alpine:3.16.2",
		Cmd:   []string{"sleep", "60000"},
		Mounts: testcontainers.ContainerMounts{
			{
				Source: testcontainers.GenericBindMountSource{
					HostPath: getwd + "/gotouch-linux",
				},
				Target:   "/tmp/gotouch-linux",
				ReadOnly: true,
			},
		},
	}

	cont, err := testcontainers.GenericContainer(context.Background(), testcontainers.GenericContainerRequest{
		ContainerRequest: request,
		Started:          true,
	})

	if cont != nil {
		i, err := cont.Exec(context.Background(), []string{"/tmp/gotouch-linux", "-s", "test", "-p", "0"})
		logs, err := cont.Logs(context.Background())
		all, err := ioutil.ReadAll(logs)
		fmt.Println(all, i)
		if err != nil {
			z.T().Fatal("could not create cont", err)
		}
	}

	z.Container = cont
}

func (z *ZippingTestSuite) TestUnzipping() {

	fmt.Println("asds")
}
