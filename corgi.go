/* Copyright (C) Alex Zhang */

package corgi


type Corgi struct {
    variables map[string]*Variable
    caches    map[string]*VariableValue
    Context   interface{}
}


func New(slot uint) (*Corgi, error) {
    var corgi *Corgi = new(Corgi)

    corgi.variables = make(map[string]*Variable, slot)
    corgi.caches    = make(map[string]*VariableValue, slot)

    if err := corgi.registerPredefineVariables(); err != nil {
        return nil, err
    }

    return corgi, nil
}
