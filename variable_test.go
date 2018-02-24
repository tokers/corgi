/* Copyright (C) Alex Zhang */

package corgi

import (
    "fmt"
    "time"
    "strconv"
    "testing"
)

var count = 0
var flag = false


var variables []*Variable = []*Variable {
    &Variable {
        Name  : "name",
        Get   : variableGet,
        Flags : VARIABLE_CHANGEABLE,
    },

    &Variable {
        Name  : "gender",
        Get   : variableGet,
        Flags : VARIABLE_CHANGEABLE,
    },

    &Variable {
        Name  : "height",
        Get   : variableGetCacheable,
        Flags : VARIABLE_CHANGEABLE,
    },

    &Variable {
        Name  : "weight",
        Get   : variableGet,
    },

    &Variable {
        Name  : "none",
        Get   : variableGetCacheable,
    },

    &Variable {
        Name  : "nil",
        Get   : variableGetCacheable,
        Flags : VARIABLE_CHANGEABLE,
    },

    &Variable {
        Name  : "while_",
        Get   : variableGet,
    },
}


func variableGet(value *VariableValue, _ interface{}, name string) error {
    if name == "name" {
        value.NotFound = false
        value.Cacheable = false
        value.Value = "alex"
        return nil
    }

    if name == "gender" {
        value.NotFound = false
        value.Cacheable = false
        value.Value = "male"
        return nil
    }

    if name == "weight" {
        value.NotFound = false
        value.Cacheable = false
        value.Value = "140"
        return nil
    }

    if name == "height" {
        count++

        value.NotFound = false
        value.Cacheable = false
        value.Value = strconv.Itoa(170 + count)
        return nil
    }

    if name == "hahah" {
        value.NotFound = false
        value.Cacheable = false
        value.Value = "hihihi"
        return nil
    }

    value.NotFound = true

    return nil
}


func variableGetCacheable(value *VariableValue, _ interface{}, name string) error {
    if name == "height" {
        count++

        value.NotFound = false
        value.Cacheable = true
        value.Value = strconv.Itoa(170 + count)
        return nil
    }

    if name == "none" {
        if flag == true {
            value.Value = "abc"
            value.NotFound = false
            fmt.Println("asdasdsd")

        } else {
            value.NotFound = true
        }

        value.Cacheable = true
        flag = !flag

        return nil
    }

    value.NotFound = true
    value.Cacheable = true

    return nil
}


func testVariableRegister(t *testing.T) {
    c, err := New()
    if err != nil {
        t.Fatal("failed to create corgi instance failed")
    }

    if err := c.RegisterNewVariables(variables); err != nil {
        t.Fatalf("failed to register new variables: %s", err.Error())
    }

    err = c.RegisterNewVariable(&Variable {
        Name  : "weight",
        Get   : variableGet,
    })

    if err == nil {
        t.Fatalf("unexpected successful register")
    }

    if err.Error() != "variable \"weight\" already exists" {
        t.Fatalf("unknown failure reason: %s", err.Error())
    }

    err = c.RegisterNewVariable(&Variable {
        Name : "while_",
        Get : variableGet,
        Flags : VARIABLE_UNKNOWN,
    })

    if err == nil {
        t.Fatalf("unexpected successful register")
    }

    if err.Error() != "variable \"while_\" already exists" {
        t.Fatalf("unknown failure reason: %s", err.Error())
    }

    /* the old "nil" is known, replaces it with the unknown one */
    err = c.RegisterNewVariable(&Variable {
        Name : "nil",
        Get : variableGet,
        Flags : VARIABLE_UNKNOWN|VARIABLE_CHANGEABLE,
    })

    if err != nil {
        t.Fatalf("failed to register variable \"nil\": %s", err.Error())
    }

    /* the old "nil" is unknown, replaces it with the known one */
    err = c.RegisterNewVariable(&Variable {
        Name : "nil",
        Get : variableGet,
        Flags : VARIABLE_CHANGEABLE,
    })

    if err != nil {
        t.Fatalf("failed to register variable \"nil\": %s", err.Error())
    }

    timeLocal := time.Now().Format("02/Jan/2006:15:04:05 -0700")
    text := "name is ${name}, gender is $gender, ${time_local}, ok!"
    expected := fmt.Sprintf("name is alex, gender is male, %s, ok!",
                            timeLocal)

    plain := parse(t, c, text)
    if plain != expected {
        t.Fatalf("incorrect value, expected \"%s\" but seen \"%s\"",
                 expected, plain)
    }
}


func testVariableCache(t *testing.T) {
    c, err := New()
    if err != nil {
        t.Fatal("failed to create corgi instance failed")
    }

    if err := c.RegisterNewVariables(variables); err != nil {
        t.Fatalf("failed to register new variables: %s", err.Error())
    }

    plain, err := strconv.Atoi(parse(t, c, "${height}"))
    if err != nil {
        t.Fatalf("failed to convert \"%s\" to integer: %s", plain, err.Error())
    }

    if plain != 171 {
        t.Fatalf("incorrect value, expected \"171\" but seen \"%s\"", plain)
    }

    plain, err = strconv.Atoi(parse(t, c, "${height}"))
    if err != nil {
        t.Fatalf("failed to convert \"%s\" to integer: %s", plain, err.Error())
    }

    if plain != 171 {
        t.Fatalf("incorrect value, expected \"171\" but seen \"%s\"", plain)
    }
}


func testVariableChange(t *testing.T) {
    c, err := New()
    if err != nil {
        t.Fatal("failed to create corgi instance failed")
    }

    if err := c.RegisterNewVariables(variables); err != nil {
        t.Fatalf("failed to register new variables: %s", err.Error())
    }

    plain, err := strconv.Atoi(parse(t, c, "${height}"))
    if err != nil {
        t.Fatalf("failed to convert \"%s\" to integer: %s", plain, err.Error())
    }

    if plain != 172 {
        t.Fatalf("incorrect value, expected \"171\" but seen \"%d\"", plain)
    }

    err = c.RegisterNewVariable(&Variable {
        Name  : "height",
        Get   : variableGet,
        Flags : VARIABLE_CHANGEABLE,
    })

    if err != nil {
        t.Fatalf("failed to register new variable: %s", err.Error())
    }

    plain, err = strconv.Atoi(parse(t, c, "${height}"))
    if err != nil {
        t.Fatalf("failed to convert \"%s\" to integer: %s", plain, err.Error())
    }

    if plain != 173 {
        t.Fatalf("incorrect value, expected \"173\" but seen \"%d\"", plain)
    }

    plain, err = strconv.Atoi(parse(t, c, "${height}"))
    if err != nil {
        t.Fatalf("failed to convert \"%s\" to integer: %s", plain, err.Error())
    }

    if plain != 174 {
        t.Fatalf("incorrect value, expected \"174\" but seen \"%d\"", plain)
    }
}


func testVariableValueNotFound(t *testing.T) {
    c, err := New()
    if err != nil {
        t.Fatal("failed to create corgi instance failed")
    }

    if err := c.RegisterNewVariables(variables); err != nil {
        t.Fatalf("failed to register new variables: %s", err.Error())
    }

    plain := "${none}"

    if cv, err := c.Parse(plain); err != nil {
        t.Fatalf("failed to parse plain string to corgi complex value: %s",
                 err.Error())

    } else {
        if _, err := c.Code(cv); err == nil {
            t.Fatal("unexpected successful parsing")

        } else if err.Error() != "vlaue of variable \"none\" not found" {
            t.Fatalf("unknown failure reason: %s", err.Error())
        }

        /* not found from the cache */
        if _, err := c.Code(cv); err == nil {
            t.Fatal("unexpected successful parsing")

        } else if err.Error() != "vlaue of variable \"none\" not found" {
            t.Fatalf("unknown failure reason: %s", err.Error())
        }
    }

    plain = "${env_xxxxx}"
    if cv, err := c.Parse(plain); err != nil {
        t.Fatalf("failed to parse plain string to corgi complex value: %s",
                 err.Error())

    } else {
        if _, err := c.Code(cv); err == nil {
            t.Fatal("unexpected successful parsing")

        } else if err.Error() != "vlaue of variable \"env_xxxxx\" not found" {
            t.Fatalf("unknown failure reason: %s", err.Error())
        }

        /* not found from the cache */
        if _, err := c.Code(cv); err == nil {
            t.Fatal("unexpected successful parsing")

        } else if err.Error() != "vlaue of variable \"env_xxxxx\" not found" {
            t.Fatalf("unknown failure reason: %s", err.Error())
        }
    }
}


func TestVariable(t *testing.T) {
    testVariableRegister(t)
    testVariableCache(t)
    testVariableChange(t)
    testVariableValueNotFound(t)
}
