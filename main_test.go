package main

import (
	"reflect"
	"strconv"
	"testing"
)

func Test_sortInterfaces(t *testing.T) {
	testCases := []struct {
		interfaces       []string
		sortedInterfaces []string
	}{
		{
			interfaces:       []string{"eth0"},
			sortedInterfaces: []string{"eth0"},
		},
		{
			interfaces:       []string{"eth1"},
			sortedInterfaces: []string{"eth1"},
		},
		{
			interfaces:       []string{"eth0", "eth1"},
			sortedInterfaces: []string{"eth0", "eth1"},
		},
		{
			interfaces:       []string{"eth1", "eth0"},
			sortedInterfaces: []string{"eth0", "eth1"},
		},
		{
			interfaces:       []string{"bond1", "bond0"},
			sortedInterfaces: []string{"bond0", "bond1"},
		},
		{
			interfaces:       []string{"eth1", "bond0"},
			sortedInterfaces: []string{"bond0", "eth1"},
		},
		{
			interfaces:       []string{"ens256", "ens224"},
			sortedInterfaces: []string{"ens224", "ens256"},
		},
		{
			interfaces:       []string{"eth4", "eth1", "eth0"},
			sortedInterfaces: []string{"eth0", "eth1", "eth4"},
		},
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			sortInterfaces(tc.interfaces)

			if !reflect.DeepEqual(tc.interfaces, tc.sortedInterfaces) {
				t.Fatalf("%d: interfaces are not properly sorted expected %#v, got %#v\n", i, tc.sortedInterfaces, tc.interfaces)
			}
		})
	}
}
