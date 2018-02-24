Name
====

Corgi - Implements the interpolation in Golang

Table of Contents
=================
* [name](#name)
* [status](#status)
* [synopsis](#synopsis)
* [Definition of variables](definition-of-variables)
* [Package](#package)
  * [Constants](#constants) 
  * [Functions](#functions)
     * [New](#new)
  * [Types](#types)
     * [Corgi](#corgi)
     * [Variable](#variable)
     * [VariableValue](#variablevalue)
     * [ComplexValue](#complexvalue)
     * [VariableSetHandler](variablesethandler)
     * [VariableGetHandler](variablegethandler)
  * [Methods](#methods)
     * [Corgi.RegisterNewVariable](#corgiregisternewvariable)
     * [Corgi.RegisterNewVariables](#corgiregisternewvariables)
     * [Corgi.Parse](#corgiparse)
     * [Corgi.Code](#corgicode)
  * [Builtin Variables](#builtin-variables)
* [Auther](#auther)
* [TODO](#todo)
* [Copyright and License](copyright-and-license)

Status
======

This package is still under expermental.

Synopsis
========

```go
package main

import (
    "log"

    "github.com/tokers/corgi"
)


func getMyVar(value *corgi.VariableValue, _ interface{}, name string) error {
    if name == "name" {
        value.Value = "tokers"
        value.NotFound = false
        value.Cacheable = true

    } else if name == "gender" {
        value.Value = "male"
        value.NotFound = false
        value.Cacheable = true

    } else {
        value.NotFound = true /* not found this */
        value.Cacheable = false
    }

    return nil
}

func main() {
    c, err := corgi.New()
    if err != nil {
        log.Fatalf("failed to new the corgi instance: %s", err.Error())
    }

    myVar1 := corgi.Variable {
        Name : "name",
        Get : getMyVar,
    }

    myVar2 := corgi.Variable {
        Name : "gender",
        Get : getMyVar,
    }

    if err := c.RegisterNewVariable(&myVar1); err != nil {
        log.Fatalf("failed to register variable \"name\": %s", err.Error())
    }

    if err := c.RegisterNewVariable(&myVar2); err != nil {
        log.Fatalf("failed to register variable \"gender\": %s", err.Error())
    }

    text := "I am ${name}, who is a $gender :)"
    expected := "I am tokers, who is a male :)"

    if cv, err := c.Parse(text); err != nil {
        log.Fatalf("failed to parse text: %s", err.Error())

    } else {
        if plain, err := c.Code(cv); err != nil {
            log.Fatalf("failed to interprete text: %s", err.Error())

        } else if plain != expected {
            log.Fatalf("unexpected value, expected \"%s\" but seen \"%s\"",
                       expected, plain)
        }
    }
}
```

Definition of variables
========================

The variables that corgi knows has the prefix `$`, followed by a sequence(the sequence can be wrappered by the curly brackets), such as `$name`, `${weather}`.

The valid characters in variable is the English letters(`a` to `z` and `A` to `Z`), the Arabic numerals(`0` to `9`) and The underscore(`_`).

There is a special variable, which name is unknown, i.e. with a fixed prefix and a floating body. This variable can be used to represent a group of variables, for instance, `env_PATH`, `env_HOSTNAME`, and etc etc etc. 

`$` is used as a variable's preface, so using `$$` if the literal meaning is expected.

Package
=======

Constants
---------

```go
const (
	VARIABLE_NO_CACHEABLE = (1 << iota)
	VARIABLE_CHANGEABLE
	VARIABLE_UNKNOWN
)
```

* `VARIABLE_NO_CACHEABLE`, marks that a variable cannot be cached
* `VARIABLE_CHANGEABLE`, marks that a variable can be changed(by calling the method [Corgi.RegisterNewVariable](#corgiregisternewvariable))
* `VARIABLE_UNKNOWN`, marks that this variable is unknown

Functions
---------

### New

*syntax*: **func New() (*Corgi, error)**

`New` returns an instance of [Corgi](#Corgi).

In case of failure, `nil` and the corresponding error object will be yielded.

In case of success, the error object will be `nil`.

Types
-----

### Corgi

```go
type Corgi struct {
	Context   interface{}
    // contains filtered or unexported fields
}
```

The filed `Context`, holds any type data that the caller wants to save, which will be used inside the variable get/set handler.

### Variable

```go
type Variable struct {
	Name   string
	Set    VariableSetHandler
	Get    VariableGetHandler
	Flags  uint
```

* `Name`, variable's name, when the variable is unknown, it is the fixed prefix
* `Set`, the set handler, which will be invoked when changeing the variable
* `Get`, the get handler, which will be invoked when getting the variable
* `Flags`, marks the variable type

### VariableValue

```go
type VariableValue struct {
	Value     string
	Cacheable bool
	NotFound  bool
```

* `Value`, the textual variable value
* `Cacheable`, marks whether the variable can be cached
* `NotFound`, marks whether the variable value is not found

### ComplexValue

```go
type ComplexValue struct {
	// contains filtered or unexported fields
}
```

The type `ComplexValue` is used to describe the result of [Corgi.Parse](#corgi.parse).

### VariableSetHandler

*syntax*: **type VariableSetHandler func(value *VariableValue, ctx interface{}, name string) error**

The prototype of the set handler, and this is not currently in use.

### VariableGetHandler

*syntax*: **type VariableGetHandler func(value *VariableValue, ctx interface{}, name string) error**

The prototype of the get handler.

The first param, `value`, is used to store the get result, one should either set the `value.Value` to the proper value or just set `value.NotFound` to `true`. Optionally, one can set the `value.Chachable` to `true`, if cache is expected(and next time the cache will be hit, so this handler will not be called).

The second param, `ctx`, is the one set in the `Corgi` object. One can always store a Context data for this.

The last param, `name` is the name of this variable, for the unknown variable, `name` represents the part of floating body.

In case of failure, one should return a corresponding error object to advertise the failure.

Methods
-------

### Corgi.RegisterNewVariable

*syntax*: **func (corgi *Corgi) RegisterNewVariable(variable *Variable) error**

`RegisterNewVariable` Registers a new variable.

The unique param is the variable that caller wants to register.

In case of failure, a corresponding error object will be yielded.

### Corgi.RegisterNewVariables

*syntax*: **func (corgi *Corgi) RegisterNewVariable(variable *Variable) error**

`RegisterNewVariable` Registers a group of variables, this method is just the wrapper of [Corgi.RegisterNewVariable](#corgiregisternewvariable).

### Corgi.Parse

*syntax*: **func (corgi *Corgi) Parse(text string) (*ComplexValue, error)**

`Parse` parses the textual data to the intermediate representation, i.e. the instance of type [ComplexValue](#complexvalue).

In case of failure, a corresponding error object will be yielded.

### Corgi.Code

*syntax*: **func (corgi *Corgi) Code(cv *ComplexValue) (string, error)**

`Code` interpretes the intermediate representation to the final result.

The param `cv` is the one generated by [Corgi.Parse](#corgiparse)

In case of failure, a corresponding error object will be yielded.

Builtin Variables
-----------------

The package corgi contains some pre-defined variables.

* `$hostname`, the host name reported by the kernel
* `$time_local`, current time in the form of [common log format](https://en.wikipedia.org/wiki/Common_Log_Format)
* `$pid`, the process id of the caller
* `$pwd`, the working directory of the caller process
* `$env_NAME`, the environment variables, e.g. `$env_PATH`, `$env_HOME`

Auther
======

Alex Zhang(张超) zchao1995@gmail.com, UPYUN Inc.

TODO
====

* regex capture group variables
* methods for flushing vairable caches

Copyright and License
=====================

Copyright (c) 2018, Alex Zhang
All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:

* Redistributions of source code must retain the above copyright notice, this
  list of conditions and the following disclaimer.

* Redistributions in binary form must reproduce the above copyright notice,
  this list of conditions and the following disclaimer in the documentation
  and/or other materials provided with the distribution.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
