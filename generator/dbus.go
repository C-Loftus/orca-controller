package main

import (
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
			return nil, fmt.Errorf("nvalid item in result: %v", item)
		}

		resultAsStrings = append(resultAsStrings, []string{item[0].(string), item[1].(string)})
	}

	return resultAsStrings, nil

}
