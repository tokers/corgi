/* Copyright (C) Alex Zhang */

package corgi

import (
    "fmt"
)


const (
    VARIABLE_NO_CACHEABLE = (1 << iota)
    VARIABLE_CHANGEABLE
)


type VariableSetHandler func(value *VariableValue, ctx, data interface{}) error
type VariableGetHandler func(value *VariableValue, ctx, data interface{}) error

type Variable struct {
    Name   string
    Set    VariableSetHandler
    Get    VariableGetHandler
    Data   interface{}
    Flags  uint
}

type VariableValue struct {
    Value     string
    Cacheable bool
    NotFound  bool
}


func (corgi *Corgi) variableGet(name string) (string, error) {
    var value VariableValue

    if variable, ok := corgi.variables[name]; ok == false {
        return "", fmt.Errorf("variable \"%s\" not found", name)

    } else {
        if (variable.Flags & VARIABLE_NO_CACHEABLE) == 0 {
            if value, ok := corgi.caches[name]; ok == true {
                /* hits the cache */
                if value.NotFound == true {
                    return "", fmt.Errorf("vlaue of variable \"%s\" not found",
                                          name)
                }

                return value.Value, nil
            }
        }

        ctx := corgi.Context

        if err := variable.Get(&value, ctx, variable.Data); err != nil {
            return "", err
        }

        if value.NotFound == true {
            return "", fmt.Errorf("vlaue of variable \"%s\" not found", name)
        }

        if value.Cacheable {
            corgi.caches[name] = &value
        }

        return value.Value, nil
    }
}


func (corgi *Corgi) RegisterNewVariable(variable *Variable) error {
    var name string = variable.Name

    if oldVariable, ok := corgi.variables[name]; ok == true {

        if oldVariable.Flags & VARIABLE_CHANGEABLE == 0 {
            return fmt.Errorf("variable \"%s\" already exists", name)
        }

        /* flushes the cache */

        delete(corgi.caches, name)
    }

    corgi.variables[name] = variable

    return nil
}


func (corgi *Corgi) RegisterNewVariables(variables []*Variable) error {
    var err error

    for _, variable := range variables {
        err = corgi.RegisterNewVariable(variable)
        if err != nil {
            return err
        }
    }

    return nil
}


func (corgi *Corgi) registerPredefineVariables() error {
    return corgi.RegisterNewVariables(predefineVariables)
}
