// Copyright 2025 Colton Loftus
// SPDX-License-Identifier: AGPL-3.0-only

package main

import (
	"encoding/xml"
	"fmt"

	"github.com/godbus/dbus/v5"
)

// Constants used for dbus calls
const (
	OrcaServiceName          = "org.gnome.Orca.Service"
	OrcaObjectPath           = "/org/gnome/Orca/Service"
	OrcaCallMethod           = "org.gnome.Orca.Module.ExecuteCommand"
	OrcaListRuntimeGetters   = "org.gnome.Orca.Module.ListRuntimeGetters"
	OrcaListRuntimeSetters   = "org.gnome.Orca.Module.ListRuntimeSetters"
	OrcaExecuteRuntimeGetter = "org.gnome.Orca.Module.ExecuteRuntimeGetter"
	OrcaExecuteRuntimeSetter = "org.gnome.Orca.Module.ExecuteRuntimeSetter"
)

// Get the name of every module registered inside ORca
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

// Contains information about a dbus command
type CommandInfo struct {
	// the name of the command as returned from dbus introspection
	commandName string
	// the custom or unique name of the command; used in generated code
	// to prevent collisions with other commands of the same name
	customName string
	// the description of the command as returned from dbus introspection
	description string
	// whether or not the command is a setter
	// and require an input argument
	isSetter bool
	// whether or not the command is a getter
	// and returns a value
	isGetter bool
}

// List all the runtime getters for a module; i.e. commands which return a value
func list_runtime_getters(module string) ([]CommandInfo, error) {
	conn, err := dbus.ConnectSessionBus()
	if err != nil {
		return nil, err
	}

	moduleObjectPath := OrcaObjectPath + "/" + module

	obj := conn.Object(OrcaServiceName, dbus.ObjectPath(moduleObjectPath))
	var result [][]any
	err = obj.Call(OrcaListRuntimeGetters, 0).Store(&result)
	if err != nil {
		return nil, err
	}
	var setterInfo []CommandInfo
	for _, item := range result {
		if item[0] == nil || item[1] == nil || len(item) != 2 {
			return nil, fmt.Errorf("invalid item in result: %v", item)
		}
		setterInfo = append(setterInfo, CommandInfo{
			commandName: item[0].(string),
			description: item[1].(string),
			isGetter:    true,
		})
	}
	return setterInfo, nil
}

// List all the runtime setters for a module; i.e. commands which take an argument
func list_runtime_setters(module string) ([]CommandInfo, error) {
	conn, err := dbus.ConnectSessionBus()
	if err != nil {
		return nil, err
	}

	moduleObjectPath := OrcaObjectPath + "/" + module

	obj := conn.Object(OrcaServiceName, dbus.ObjectPath(moduleObjectPath))
	var result [][]any
	err = obj.Call(OrcaListRuntimeSetters, 0).Store(&result)
	if err != nil {
		return nil, err
	}
	var setterInfo []CommandInfo
	for _, item := range result {
		if item[0] == nil || item[1] == nil || len(item) != 2 {
			return nil, fmt.Errorf("invalid item in result: %v", item)
		}
		setterInfo = append(setterInfo, CommandInfo{
			customName:  fmt.Sprintf("Set%s", item[0].(string)),
			commandName: item[0].(string),
			description: item[1].(string),
			isSetter:    true,
		})
	}
	return setterInfo, nil
}

// List all the commands for a module; i.e. commands which semantically
// change the state of the module somehow but are not runtime setters/getters
func get_commands_for_module(module string) ([]CommandInfo, error) {
	conn, err := dbus.ConnectSessionBus()
	if err != nil {
		return nil, err
	}

	moduleObjectPath := OrcaObjectPath + "/" + module

	obj := conn.Object(OrcaServiceName, dbus.ObjectPath(moduleObjectPath))
	var result [][]any
	err = obj.Call("org.gnome.Orca.Module.ListCommands", 0).Store(&result)
	if err != nil {
		return nil, err
	}
	var resultAsCommandInfo []CommandInfo

	for _, item := range result {
		if item[0] == nil || item[1] == nil || len(item) != 2 {
			return nil, fmt.Errorf("invalid item in result: %v", item)
		}

		resultAsCommandInfo = append(resultAsCommandInfo, CommandInfo{
			commandName: item[0].(string),
			description: item[1].(string),
		})
	}

	return resultAsCommandInfo, nil

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
