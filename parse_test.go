/* Copyright (C) Alex Zhang */

package corgi

import (
    "os"
    "fmt"
    "testing"
)


func testParseFailed(t *testing.T) {
    c, err := New()
    if err != nil {
        t.Fatal("failed to create corgi instance failed")
    }

    plain := "hello $"

    if _, err := c.Parse(plain); err == nil {
        t.Fatal("unexpected successful parsing")

    } else if err.Error() != "invalid variable name" {
        t.Fatalf("unknown failure reason: %s", err.Error())
    }

    plain = "hello $hostname$"
    if _, err := c.Parse(plain); err == nil {
        t.Fatal("unexpected successful parsing")

    } else if err.Error() != "invalid variable name" {
        t.Fatalf("unknown failure reason: %s", err.Error())
    }

    plain = "hello ${hostname"
    errorReason := "unexpected end of string, \"}\" is missing"

    if _, err := c.Parse(plain); err == nil {
        t.Fatal("unexpected successful parsing")

    } else if err.Error() != errorReason {
        t.Fatalf("unknown failure reason: %s", err.Error())
    }

    plain = "hello ${hostname$pid"
    errorReason = "\"}\" for variable \"hostname\" is missing"

    if _, err := c.Parse(plain); err == nil {
        t.Fatal("unexpected successful parsing")

    } else if err.Error() != errorReason {
        t.Fatalf("unknown failure reason: %s", err.Error())
    }

    plain = "hello ${hostname}$pP0_-id$"
    errorReason = "unknown variable \"pP0_\""

    if _, err := c.Parse(plain); err == nil {
        t.Fatal("unexpected successful parsing")

    } else if err.Error() != errorReason {
        t.Fatalf("unknown failure reason: %s", err.Error())
    }

    plain = "hello ${pP0}$"
    errorReason = "unknown variable \"pP0\""

    if _, err := c.Parse(plain); err == nil {
        t.Fatal("unexpected successful parsing")

    } else if err.Error() != errorReason {
        t.Fatalf("unknown failure reason: %s", err.Error())
    }

    plain = "hello $pP0"
    errorReason = "unknown variable \"pP0\""

    if _, err := c.Parse(plain); err == nil {
        t.Fatal("unexpected successful parsing")

    } else if err.Error() != errorReason {
        t.Fatalf("unknown failure reason: %s", err.Error())
    }
}


func testParseComplex(t *testing.T) {
    c, err := New()
    if err != nil {
        t.Fatal("failed to create corgi instance failed")
    }

    if err := c.RegisterNewVariables(variables); err != nil {
        t.Fatalf("failed to register new variables: %s", err.Error())
    }

    text := `Hello, This is ${name}, dollar is $$, Host is $hostname,
             This is welsh corgi, 世界，你好, gender is $gender`

    expected := `Hello, This is alex, dollar is $, Host is %s,
             This is welsh corgi, 世界，你好, gender is male`

    hostname, err := os.Hostname()
    if err != nil {
        t.Fatalf("failed to get hostname: %s", err.Error())
    }

    expected = fmt.Sprintf(expected, hostname)

    plain := parse(t, c, text)
    if plain != expected {
        t.Fatalf("incorrect value, expected \"%s\" but seen \"%s\"",
                 expected, plain)
    }
}


func TestParse(t *testing.T) {
    testParseFailed(t)
    testParseComplex(t)
}
