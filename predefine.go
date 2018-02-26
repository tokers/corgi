/* Copyright (C) Alex Zhang */

package corgi

import (
    "os"
    "time"
    "strconv"
)


var predefineVariables []*Variable = []*Variable {
    &Variable {
        Name  : "hostname",
        Get   : predefineVariableHostname,
        Flags : VARIABLE_CHANGEABLE,
    },

    &Variable {
        Name  : "time_local",
        Get   : predefineVariableTimeLocal,
        Flags : VARIABLE_CHANGEABLE,
    },

    &Variable {
        Name  : "pid",
        Get   : predefineVariablePID,
        Flags : VARIABLE_CHANGEABLE,
    },

    &Variable {
        Name  : "pwd",
        Get   : predefineVariablePWD,
        Flags : VARIABLE_CHANGEABLE,
    },

    &Variable {
        Name  : "year",
        Get   : predefineVariableTime,
        Flags : VARIABLE_CHANGEABLE,
    },

    &Variable {
        Name  : "month",
        Get   : predefineVariableTime,
        Flags : VARIABLE_CHANGEABLE,
    },

    &Variable {
        Name  : "week",
        Get   : predefineVariableTime,
        Flags : VARIABLE_CHANGEABLE,
    },

    &Variable {
        Name  : "day",
        Get   : predefineVariableTime,
        Flags : VARIABLE_CHANGEABLE,
    },

    &Variable {
        Name  : "hour",
        Get   : predefineVariableTime,
        Flags : VARIABLE_CHANGEABLE,
    },

    &Variable {
        Name  : "minute",
        Get   : predefineVariableTime,
        Flags : VARIABLE_CHANGEABLE,
    },

    &Variable {
        Name  : "zone",
        Get   : predefineVariableTime,
        Flags : VARIABLE_CHANGEABLE,
    },

    &Variable {
        Name  : "second",
        Get   : predefineVariableTime,
        Flags : VARIABLE_CHANGEABLE,
    },

    &Variable {
        Name  : "env_",
        Get   : predefineVariableENV,
        Flags : VARIABLE_UNKNOWN,
    },
}


func predefineVariableTime(value *VariableValue, _ interface{}, component string) error {
    now := time.Now()

    value.NotFound = false
    value.Cacheable = false

    if component == "year" {
        value.Value = strconv.Itoa(now.Year())
        return nil
    }

    if component == "month" {
        value.Value = strconv.Itoa(int(now.Month()))
        return nil
    }

    if component == "week" {
        week := now.Weekday()
        switch (week) {

        case time.Sunday:
            value.Value = "Sun"

        case time.Monday:
            value.Value = "Mon"

        case time.Tuesday:
            value.Value = "Tue"

        case time.Wednesday:
            value.Value = "Wed"

        case time.Thursday:
            value.Value = "Thur"

        case time.Friday:
            value.Value = "Fri"

        default:
            value.Value = "Sat"
        }

        return nil
    }

    if component == "day" {
        value.Value = strconv.Itoa(now.Day())
        return nil
    }

    if component == "hour" {
        value.Value = strconv.Itoa(now.Hour())
        return nil
    }

    if component == "minute" {
        value.Value = strconv.Itoa(now.Minute())
        return nil
    }

    if component == "second" {
        value.Value = strconv.Itoa(now.Second())
        return nil
    }

    if component == "zone" {
        value.Value, _ = now.Zone()
        return nil
    }

    value.NotFound = true
    value.Cacheable = true

    return nil
}


func predefineVariableHostname(value *VariableValue, _ interface{}, _ string) error {
    if name, err := os.Hostname(); err != nil {
        value.NotFound = true
        value.Cacheable = false
        return err

    } else {
        value.Value = name
        value.NotFound = false
        value.Cacheable = true
    }

    return nil
}


func predefineVariableTimeLocal(value *VariableValue, _ interface{}, _ string) error {
    value.Value = time.Now().Format("02/Jan/2006:15:04:05 -0700")
    value.Cacheable = false
    value.NotFound = false

    return nil
}


func predefineVariablePID(value *VariableValue, _ interface{}, _ string) error {
    pid := os.Getpid()
    value.Value = strconv.Itoa(pid)

    value.NotFound = false
    value.Cacheable = true

    return nil
}


func predefineVariablePWD(value *VariableValue, _ interface{}, _ string) error {
    value.Cacheable = false

    if dir, err := os.Getwd(); err != nil {
        value.NotFound = true
        return err

    } else {
        value.Value = dir
    }

    value.NotFound = false

    return nil
}


func predefineVariableENV(value *VariableValue, _ interface{}, key string) error {
    val := os.Getenv(key)

    value.Cacheable = false

    if val == "" {
        value.NotFound = true

    } else {
        value.NotFound = false
        value.Value = val
    }

    return nil
}
