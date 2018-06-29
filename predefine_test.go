// Copyright (C) Alex Zhang

package corgi

import (
    "os"
    "fmt"
    "time"
    "math"
    "strconv"
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


func testTime(t *testing.T, c *Corgi) {
    now := time.Now()
    var week string

    switch (now.Weekday()) {

    case time.Sunday:
        week = "Sun"

    case time.Monday:
        week = "Mon"

    case time.Tuesday:
        week = "Tue"

    case time.Wednesday:
        week = "Wed"

    case time.Thursday:
        week = "Thu"

    case time.Friday:
        week = "Fri"

    default:
        week = "Sat"
    }

    expected := fmt.Sprintf("%d%d%d%s%d", now.Year(), int(now.Month()),
                            now.Day(), week, now.Hour())
    data := parse(t, c, "$year$month$day$week$hour")
    if data != expected {
        t.Fatalf("incorrect value, expected \"%s\" but seen \"%s\"", expected,
                 data)
    }

    min, err := strconv.Atoi(parse(t, c, "$minute"))
    if err != nil {
        t.Fatalf("failed to convert $minute to integer: %s", err.Error())
    }

    if math.Abs(float64(min - now.Minute())) > 1 {
        t.Fatalf("incorrect minute: %d", data)
    }

    _, err = strconv.Atoi(parse(t, c, "$second"))
    if err != nil {
        t.Fatalf("failed to convert $second to integer: %s", err.Error())
    }

    // we assume that the second value is correct

    realZone, _ := now.Zone()

    zone := parse(t, c, "${zone}")
    if realZone != zone {
        t.Fatalf("incorrect zone, expected \"%s\" but seen \"%s\"", realZone,
                 zone)
    }

    t1 := time.Now().Unix()
    data = parse(t, c, "${time}")
    t2, err := strconv.ParseInt(data, 10, 64)
    if err != nil {
        t.Fatalf("invalid ${time} value: %s", data)
    }

    if t1 - t2 >= 1 || t2 - t1 >= 1 {
        t.Fatalf("incorrect time, expected \"%d\" but seen \"%d\"", t1, t2)
    }
}


func testENV(t *testing.T, c *Corgi) {
    path := fmt.Sprintf("The PATH is %s", os.Getenv("PATH"))

    data := parse(t, c, "The PATH is $env_PATH")
    if data != path {
        t.Fatalf("incorrect value, expected \"%s\" but seen \"%s\"",
                 path, data)
    }

    path = fmt.Sprintf("The GOPATH is %s", os.Getenv("GOPATH"))

    data = parse(t, c, "The GOPATH is $env_GOPATH")
    if data != path {
        t.Fatalf("incorrect value, expected \"%s\" but seen \"%s\"",
                 path, data)
    }
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
        t.Fatal("failed to create corgi instance failed")
    }

    testPID(t, c)
    testHostname(t, c)
    testTimeLocal(t, c)
    testPWD(t, c)
    testENV(t, c)
    testTime(t, c)
}
