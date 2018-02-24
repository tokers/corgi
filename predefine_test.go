/* Copyright (C) Alex Zhang */

package corgi

import (
    "os"
    "fmt"
    "time"
    "testing"
)


func parse(t *testing.T, c *Corgi, plain string) string {
    if cv, err := c.Parse(plain); err != nil {
        t.Fatalf("failed to parse plain string to corgi complex value: %s",
                 err.Error())

    } else {
        if result, err := c.Code(cv); err != nil {
            t.Fatalf("failed to run complex value's code: %s", err.Error())

        } else {
            return result
        }
    }

    return ""
}


func testPID(t *testing.T, c *Corgi) {
    pid := fmt.Sprintf("The process pid is %d", os.Getpid())

    data := parse(t, c, "The process pid is $pid")
    if data != pid {
        t.Fatalf("incorrect value, expected %s but seen %s",
                 pid, data)
    }

    data = parse(t, c, "The process pid is ${pid}")
    if data != pid {
        t.Fatalf("incorrect value, expected %s but seen %s",
                 pid, data)
    }
}


func testHostname(t *testing.T, c *Corgi) {
    name, err := os.Hostname()
    if err != nil {
        t.Fatalf("failed to get hostname: %s", err.Error())
    }

    hostname := fmt.Sprintf("Hostname is %s", name)

    data := parse(t, c, "Hostname is $hostname")
    if data != hostname {
        t.Fatalf("incorrect value, expected %s but seen %s",
                 hostname, data)
    }

    data = parse(t, c, "Hostname is ${hostname}")
    if data != hostname {
        t.Fatalf("incorrect value, expected %s but seen %s",
                 hostname, data)
    }
}


func testTimeLocal(t *testing.T, c *Corgi) {
    layout := "02/Jan/2006:15:04:05 -0700"
    timeLocal := fmt.Sprintf("time_local is [%s]", time.Now().Format(layout))

    data := parse(t, c, "time_local is [$time_local]")
    if data != timeLocal {
        t.Fatalf("incorrect value, expected %s but seen %s",
                 timeLocal, data)
    }

    data = parse(t, c, "time_local is [${time_local}]")
    if data != timeLocal {
        t.Fatalf("incorrect value, expected %s but seen %s",
                 timeLocal, data)
    }
}


func testPWD(t *testing.T, c *Corgi) {
    pwd, err := os.Getwd()
    if err != nil {
        t.Fatalf("failed to get pwd: %s", err.Error())
    }

    pwd = fmt.Sprintf("pwd is %s", pwd)

    data := parse(t, c, "pwd is $pwd")
    if data != pwd {
        t.Fatalf("incorrect value, expected %s but seen %s",
                 pwd, data)
    }

    data = parse(t, c, "pwd is ${pwd}")
    if data != pwd {
        t.Fatalf("incorrect value, expected %s but seen %s",
                 pwd, data)
    }
}


func TestPredefine(t *testing.T) {
    c, err := New()
    if err != nil {
        t.Fatal("create corgi instance failed")
    }

    testPID(t, c)
    testHostname(t, c)
    testTimeLocal(t, c)
    testPWD(t, c)
}
