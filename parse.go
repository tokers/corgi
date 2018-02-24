/* Copyright (C) Alex Zhang */

package corgi

import (
    "fmt"
    "bytes"
    "errors"
)


const (
    VARIABLE_PREFACE  = '$'
    VARIABLE_LBRACKET = '{'
    VARIABLE_RBRACKET = '}'

    PARSE_PLAIN = iota
    PARSE_VARIABLE_PREFACE
    PARSE_VARIABLE

    SCRIPT_PLAIN = iota
    SCRIPT_VARIABLE
)


type scriptCode struct {
    kind  uint
    data  string
}


type ComplexValue struct {
    code  []scriptCode
    size    int
    corgi  *Corgi
}


func isValidVariableCharacter(ch rune) bool {
    if ch >= '0' && ch <= '9' {
        return true
    }

    if ch >= 'a' && ch <= 'z' {
        return true
    }

    if ch >= 'A' && ch <= 'Z' {
        return true
    }

    if ch == '_' {
        return true
    }

    return false
}


func (cv *ComplexValue) append(name string, variable bool) error {
    if variable == false {
        cv.code = append(cv.code, scriptCode{
            kind : SCRIPT_PLAIN,
            data : name,
        })

        cv.size++

        return nil
    }

    /* this is a variable */

    if _, ok := cv.corgi.variables[name]; ok == false {
        return fmt.Errorf("unknown variable \"%s\"", name)
    }

    cv.code = append(cv.code, scriptCode {
        kind : SCRIPT_VARIABLE,
        data : name,
    })

    cv.size++

    return nil
}


func (corgi *Corgi) Parse(text string) (*ComplexValue, error) {
    state   := PARSE_PLAIN
    from    := 0
    to      := -1
    bracket := false

    var cv *ComplexValue = new(ComplexValue)

    cv.corgi = corgi

    for i, ch := range text {

        switch (state) {

        case PARSE_PLAIN:

            if ch != VARIABLE_PREFACE {
                to++
                continue
            }

            state = PARSE_VARIABLE_PREFACE

            if err := cv.append(text[from:to+1], false); err != nil {
                return nil, err
            }

            /* reset */
            from = i + 1

        case PARSE_VARIABLE_PREFACE:

            if ch == VARIABLE_LBRACKET {
                bracket = true
                from = i + 1

            } else {
                from = i
            }

            to = i

            state = PARSE_VARIABLE

        case PARSE_VARIABLE:

            if bracket == true && ch == VARIABLE_RBRACKET {
                state = PARSE_PLAIN
                bracket = false

                if err := cv.append(text[from:to+1], true); err != nil {
                    return nil, err
                }

                from = i + 1
                to = i

                continue
            }

            if isValidVariableCharacter(ch) {
                to++
                continue
            }

            if bracket == true {
                return nil, fmt.Errorf("\"}\" in \"%s\" is missing",
                                       text[from:to])
            }

            state = PARSE_PLAIN

            if err := cv.append(text[from:to+1], true); err != nil {
                return nil, err
            }

            from = i
            to = i
        }
    }

    if state == PARSE_VARIABLE_PREFACE {
        return nil, errors.New("invalid variable name")
    }

    if bracket == true {
        return nil, errors.New("unexpected end of string, \"}\" is missing")
    }

    if from <= to {
        err := cv.append(text[from:to+1], state == PARSE_VARIABLE)
        if err != nil {
            return nil, err
        }
    }

    return cv, nil
}


func (corgi *Corgi) Code(cv *ComplexValue) (string, error) {
    var buffer    bytes.Buffer

    pos := 0

    for {
        if pos == cv.size {
            break
        }

        code := cv.code[pos]
        pos++

        if code.kind == SCRIPT_PLAIN {
            buffer.WriteString(code.data)
            continue
        }

        if result, err := corgi.variableGet(code.data); err != nil {
            return "", err

        } else {
            if n, err := buffer.WriteString(result); err != nil {
                return "", err

            } else if n != len(result) {
                return "", errors.New("incomplete written operation")
            }
        }
    }

    return buffer.String(), nil
}
