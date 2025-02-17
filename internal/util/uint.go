package util

import (
	"strconv"

	"github.com/abhakash/gosaavn/internal/logging"
)

func StringToUInt(value string, bitsize int) (any, error) {
	parsedInt, err := strconv.ParseUint(value, 10, bitsize)
	if err != nil {
		logging.Log.Warn("Failed to parse {} into Int {}", value, err)
		return nil, err
	}
	switch bitsize {
	case 16:
		return uint16(parsedInt), nil
	case 32:
		return uint32(parsedInt), nil
	default:
		return parsedInt, nil
	}
}
