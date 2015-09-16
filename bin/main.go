package main

import (
	"github.com/docker/machine/libmachine/drivers/plugin"
	"github.com/zchee/docker-machine-hypercore"
)

func main() {
	plugin.RegisterDriver(new(hypercore.Driver))
}
