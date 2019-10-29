// Copyright (c) 2014 Kelsey Hightower. All rights reserved. Use of this source
// code is governed by the Apache License, Version 2.0 that can be found in the
// LICENSE file. See also https://github.com/kelseyhightower/setup-network-environment.
package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"

	"github.com/docker/libcontainer/netlink"
)

var (
	defaultEnvironmentFilePath = "/etc/network-environment"
	environmentFilePath        string
	helpUsage                  bool
	verboseOutput              bool
)

func init() {
	log.SetFlags(0)
	flag.BoolVar(&helpUsage, "help", false, "print help usage")
	flag.StringVar(&environmentFilePath, "o", defaultEnvironmentFilePath, "environment file")
	flag.BoolVar(&verboseOutput, "verbose", false, "enable verbose output")
}

func main() {
	flag.Parse()
	if helpUsage {
		log.Println("Provide the -o to specify the environment file path or ommit it and use the default.")
		os.Exit(0)
	}
	tempFilePath := environmentFilePath + ".tmp"
	tempFile, err := os.Create(tempFilePath)
	if err != nil {
		log.Fatal(err)
	}
	defer tempFile.Close()
	if err := writeEnvironment(tempFile); err != nil {
		log.Fatal(err)
	}
	os.Rename(tempFilePath, environmentFilePath)
}

func writeEnvironment(w io.Writer) error {
	var buffer bytes.Buffer
	defaultIfaceName, err := getDefaultGatewayIfaceName()
	if err != nil {
		// A default route is not required; log it and keep going.
		log.Println(err)
	}
	interfaces, err := net.Interfaces()
	if err != nil {
		return err
	}
	for _, iface := range interfaces {
		addrs, err := iface.Addrs()
		if err != nil {
			return err
		}
		for _, addr := range addrs {
			ip, _, err := net.ParseCIDR(addr.String())
			// Record IPv4 network settings. Stop at the first IPv4 address
			// found for the interface.
			if err == nil && ip.To4() != nil {
				buffer.WriteString(fmt.Sprintf("%s_IPV4=%s\n", strings.Replace(strings.ToUpper(iface.Name), ".", "_", -1), ip.String()))
				if defaultIfaceName == iface.Name {
					buffer.WriteString(fmt.Sprintf("DEFAULT_IPV4=%s\n", ip.String()))
				}
				break
			}
		}
	}
	if _, err := buffer.WriteTo(w); err != nil {
		return err
	}
	return nil
}

func getDefaultGatewayIfaceName() (string, error) {
	verboseLog("getDefaultGatewayIfaceName started")
	routes, err := netlink.NetworkGetRoutes()
	verboseLog("netlink.NetworkGetRoutes() called")
	if err != nil {
		verboseLog(fmt.Sprintf("netlink.NetworkGetRoutes() error: %s", err))
		return "", err
	}
	verboseLog(fmt.Sprintf("netlink.NetworkGetRoutes() successful. Found %d routes", len(routes)))
	for i, route := range routes {
		ifaceName := "nil"
		if route.Iface != nil {
			ifaceName = route.Iface.Name
		}

		isDefault := "false"
		if route.Default {
			isDefault = "true"
		}

		verboseLog(fmt.Sprintf("Route %d: IFace = %s, Default = %s", i, ifaceName, isDefault))
		if route.Default {
			if route.Iface == nil {
				return "", errors.New("found default route but could not determine interface")
			}
			return route.Iface.Name, nil
		}
	}
	return "", errors.New("unable to find default route")
}

func verboseLog(msg string) {
	if verboseOutput {
		log.Println(msg)
	}
}
