/* Copyright (C) Alex Zhang */

package corgi


const (
    VARIABLE_SLOTS = 16
)


type Corgi struct {
    variables map[string]*Variable
    caches    map[string]*VariableValue
    Context   interface{}
}


func New() (*Corgi, error) {
    var corgi *Corgi = new(Corgi)

    corgi.variables = make(map[string]*Variable, VARIABLE_SLOTS)
    corgi.caches    = make(map[string]*VariableValue, VARIABLE_SLOTS)

    if err := corgi.registerPredefineVariables(); err != nil {
        return nil, err
    }

    return corgi, nil
}
