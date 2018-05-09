// Copyright (C) Alex Zhang

package corgi

import (
    "fmt"
    "strings"
)


const (
    VARIABLE_NO_CACHEABLE = (1 << iota)
    VARIABLE_CHANGEABLE
    VARIABLE_UNKNOWN
)


type VariableSetHandler func(value *VariableValue, ctx interface{}, name string) error
type VariableGetHandler func(value *VariableValue, ctx interface{}, name string) error


// Variable describles a variable.
// Name, variable's name, when the variable is unknown, it is the fixed prefix.
// Set, the set handler, which will be invoked when changeing the variable.
// Get, the get handler, which will be invoked when getting the variable.
// Flags, marks the variable type.
type Variable struct {
    Name   string
    Set    VariableSetHandler
    Get    VariableGetHandler
    Flags  uint
}

// VariableValue describles the variable value.
// Value, the textual variable value.
// Cacheable, marks whether the variable can be cached.
// NotFound, marks whether the variable value is not found.
type VariableValue struct {
    Value     string
    Cacheable bool
    NotFound  bool
}


func (corgi *Corgi) validUnknownVariable(name string) *Variable {
    for prefix, variable := range corgi.unknowns {

        // FIXME implements with a more effective way(like trie?)
        if strings.HasPrefix(name, prefix) == true {
            return variable
        }
    }

    return nil
}


func (corgi *Corgi) variableGet(name string) (string, error) {
    var value     VariableValue
    var variable *Variable
    var ok        bool

    var varName   string = name

    if variable, ok = corgi.variables[name]; ok == false {
        if variable = corgi.validUnknownVariable(name); variable == nil {
            return "", fmt.Errorf("variable \"%s\" not found", name)
        }

        prefix := len(variable.Name)
        varName = name[prefix:]
    }

    if (variable.Flags & VARIABLE_NO_CACHEABLE) == 0 {
        if value, ok := corgi.caches[name]; ok == true {
            // hits the cache
            if value.NotFound == true {
                return "", fmt.Errorf("vlaue of variable \"%s\" not found",
                name)
            }

            return value.Value, nil
        }
    }

    ctx := corgi.Context

    if err := variable.Get(&value, ctx, varName); err != nil {
        return "", err
    }

    if value.Cacheable {
        corgi.caches[name] = &value
    }

    if value.NotFound == true {
        return "", fmt.Errorf("vlaue of variable \"%s\" not found", name)
    }

    return value.Value, nil
}


// RegisterNewVariable Registers a new variable.
// The unique param is the variable that caller wants to register.
// In case of failure, a corresponding error object will be yielded.
func (corgi *Corgi) RegisterNewVariable(variable *Variable) error {
    var name string = variable.Name

    if oldVariable, ok := corgi.variables[name]; ok == true {

        if oldVariable.Flags & VARIABLE_CHANGEABLE == 0 {
            return fmt.Errorf("variable \"%s\" already exists", name)
        }

        // flushes the cache
        delete(corgi.caches, name)

        if variable.Flags & VARIABLE_UNKNOWN == 0 {
            corgi.variables[name] = variable

        } else {
            delete(corgi.variables, name)
            corgi.unknowns[name] = variable
        }

        return nil
    }

    // name is actually the prefix
    if oldVariable, ok := corgi.unknowns[name]; ok == true {

        if oldVariable.Flags & VARIABLE_CHANGEABLE == 0 {
            return fmt.Errorf("variable \"%s\" already exists", name)
        }

        for key, _ := range corgi.caches {

            // FIXME implements with a more effective way
            if strings.HasPrefix(key, name) == true {

                // flushes the cache
                delete(corgi.caches, key)
            }
        }

        if variable.Flags & VARIABLE_UNKNOWN == 0 {
            delete(corgi.unknowns, name)
            corgi.variables[name] = variable

        } else {
            corgi.unknowns[name] = variable
        }

        return nil
    }

    if variable.Flags & VARIABLE_UNKNOWN == 0 {
        corgi.variables[name] = variable

    } else {
        corgi.unknowns[name] = variable
    }

    return nil
}


// RegisterNewVariable Registers a group of variables, this method is just the
// wrapper of Corgi.RegisterNewVariable.
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
