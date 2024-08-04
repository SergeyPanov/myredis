package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var unknownCommand = errors.New("unknown command")

const separator string = "\r\n"

var invalidArray = errors.New("invalid array")

type cmdType string

const SET cmdType = "SET"
const GET cmdType = "GET"
const PING cmdType = "PING"
const UNKNOWN cmdType = "UNKNOWN"

type Command struct {
	Type  cmdType
	Key   string
	Value string
}

func parseCommand(cmd string) (*Command, error) {
	idxNext := strings.Index(cmd, separator)
	if idxNext == -1 {
		return nil, errors.Join(invalidArray, fmt.Errorf("cmd: %s", cmd))
	}

	arrLen, err := strconv.Atoi(cmd[1:idxNext])
	if err != nil {
		return nil, errors.Join(invalidArray, fmt.Errorf("the array size is invalid; cmd: %s", cmd))
	}

	sepLen := len(separator)
	subArr := cmd[idxNext+sepLen:]
	cmdArr := make([]string, 3)

	for i := 0; i < arrLen; i++ {
		idxNext += sepLen
		idx := strings.Index(subArr, separator)
		elLen, err := strconv.Atoi(subArr[1:idx])
		if err != nil {
			return nil, errors.Join(invalidArray, fmt.Errorf("the array element on index %d is invalid; cmd: %s", idx, cmd))
		}
		val := subArr[idx+sepLen : idx+sepLen+elLen]
		cmdArr[i] = val
		subArr = subArr[idx+2*sepLen+elLen:]
	}

	t, err := convertToCmdType(cmdArr[0])
	if err != nil {
		return nil, err
	}

	return &Command{
		Type:  t,
		Key:   cmdArr[1],
		Value: cmdArr[2],
	}, nil
}

func convertToCmdType(s string) (cmdType, error) {
	switch s {
	case "SET":
		return SET, nil
	case "GET":
		return GET, nil
	case "PING":
		return PING, nil
	default:
		return UNKNOWN, errors.Join(unknownCommand, fmt.Errorf("received: %s", s))
	}
}
