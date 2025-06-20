// Copyright 2025 Colton Loftus
// SPDX-License-Identifier: AGPL-3.0-only

package main

import (
	"encoding/xml"
	"fmt"

	"github.com/godbus/dbus/v5"
)

const OrcaServiceName = "org.gnome.Orca.Service"
const OrcaObjectPath = "/org/gnome/Orca/Service"
const OrcaCallMethod = "org.gnome.Orca.Module.ExecuteCommand"

func get_modules() ([]string, error) {

	conn, err := dbus.ConnectSessionBus()
	if err != nil {
		return nil, err
	}

	obj := conn.Object(OrcaServiceName, OrcaObjectPath)
	var result []any
	err = obj.Call("org.gnome.Orca.Service.ListModules", 0).Store(&result)
	if err != nil {
		return nil, err
	}
	var modules []string
	for _, item := range result {
		modules = append(modules, item.(string))
	}
	return modules, nil
}

func get_commands_for_module(module string) ([][]string, error) {
	conn, err := dbus.ConnectSessionBus()
	if err != nil {
		return nil, err
	}

	moduleObjectPath := OrcaObjectPath + "/" + module

	obj := conn.Object(OrcaServiceName, dbus.ObjectPath(moduleObjectPath))
	var result [][]interface{}
	err = obj.Call("org.gnome.Orca.Module.ListCommands", 0).Store(&result)
	if err != nil {
		return nil, err
	}
	var resultAsStrings [][]string

	for _, item := range result {
		if item[0] == nil || item[1] == nil || len(item) != 2 {
			return nil, fmt.Errorf("invalid item in result: %v", item)
		}

		resultAsStrings = append(resultAsStrings, []string{item[0].(string), item[1].(string)})
	}

	return resultAsStrings, nil

}

type Node struct {
	Name       string       `xml:"name,attr,omitempty"`
	Interfaces []Interface_ `xml:"interface"`
	Nodes      []Node       `xml:"node"`
}

type Interface_ struct {
	Name       string     `xml:"name,attr"`
	Methods    []Method   `xml:"method"`
	Signals    []Signal   `xml:"signal"`
	Properties []Property `xml:"property"`
}

type Method struct {
	Name string `xml:"name,attr"`
	Args []Arg  `xml:"arg"`
}

type Signal struct {
	Name string `xml:"name,attr"`
	Args []Arg  `xml:"arg"`
}

type Property struct {
	Name   string `xml:"name,attr"`
	Type   string `xml:"type,attr"`
	Access string `xml:"access,attr"`
}

type Arg struct {
	Name      string `xml:"name,attr,omitempty"`
	Type      string `xml:"type,attr"`
	Direction string `xml:"direction,attr,omitempty"`
}

func introspect_root() (Node, error) {
	conn, err := dbus.ConnectSessionBus()
	if err != nil {
		return Node{}, err
	}

	const Introspect = "org.freedesktop.DBus.Introspectable"

	obj := conn.Object(OrcaServiceName, OrcaObjectPath)
	var result string
	err = obj.Call(Introspect+".Introspect", 0).Store(&result)
	if err != nil {
		return Node{}, err
	}

	var root Node
	err = xml.Unmarshal([]byte(result), &root)
	if err != nil {
		return Node{}, err
	}
	return root, nil
}
