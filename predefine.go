/* Copyright (C) Alex Zhang */

package corgi

import (
    "os"
    "time"
    "strconv"
)


var predefineVariables []*Variable = []*Variable {
    &Variable {
        Name  : "hostname",
        Get   : predefineVariableHostname,
        Flags : VARIABLE_CHANGEABLE,
    },

    &Variable {
        Name  : "time_local",
        Get   : predefineVariableTimeLocal,
        Flags : VARIABLE_CHANGEABLE,
    },

    &Variable {
        Name  : "pid",
        Get   : predefineVariablePID,
        Flags : VARIABLE_CHANGEABLE,
    },

    &Variable {
        Name  : "pwd",
        Get   : predefineVariablePWD,
        Flags : VARIABLE_CHANGEABLE,
    },
}


func predefineVariableHostname(value *VariableValue, _ interface{}, _ string) error {
    if name, err := os.Hostname(); err != nil {
        value.NotFound = true
        value.Cacheable = false
        return err

    } else {
        value.Value = name
        value.NotFound = false
        value.Cacheable = true
    }

    return nil
}


func predefineVariableTimeLocal(value *VariableValue, _ interface{}, _ string) error {
    value.Value = time.Now().Format("02/Jan/2006:15:04:05 -0700")
    value.Cacheable = false
    value.NotFound = false

    return nil
}


func predefineVariablePID(value *VariableValue, _ interface{}, _ string) error {
    pid := os.Getpid()
    value.Value = strconv.Itoa(pid)

    value.NotFound = false
    value.Cacheable = true

    return nil
}


func predefineVariablePWD(value *VariableValue, _ interface{}, _ string) error {
    value.Cacheable = false

    if dir, err := os.Getwd(); err != nil {
        value.NotFound = true
        return err

    } else {
        value.Value = dir
    }

    value.NotFound = false

    return nil
}
