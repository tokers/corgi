// Copyright (C) Alex Zhang

// Package corgi does the variables interpolation job.
package corgi


const (
    VARIABLE_SLOTS = 16
)


// Corgi is the core struct for user.
// The field Context, holds any type data that the caller wants to save,
// which will be used inside the variable get/set handler.
type Corgi struct {
    variables map[string]*Variable
    unknowns  map[string]*Variable
    caches    map[string]*VariableValue
    Context   interface{}
    Group   []string
}


// New returns an instance of Corgi.
// In case of failure, nil and the corresponding error object will be yielded.
// In case of success, the error object will be nil.
func New() (*Corgi, error) {
    var corgi *Corgi = new(Corgi)

    corgi.variables = make(map[string]*Variable, VARIABLE_SLOTS)
    corgi.unknowns = make(map[string]*Variable, VARIABLE_SLOTS >> 1)
    corgi.caches = make(map[string]*VariableValue, VARIABLE_SLOTS)

    if err := corgi.registerPredefineVariables(); err != nil {
        return nil, err
    }

    return corgi, nil
}
