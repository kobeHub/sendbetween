package io

import "bytes"

var MSG_DELIM string = `\q`

var ScanMsg = SplitAt(MSG_DELIM)

func SplitAt(delim string) func(data []byte, atEOF bool) (advance int, token []byte, err error) {
	splitByte := []byte(delim)
	splitLen := len(splitByte)

	return func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		dataLen := len(data)

		// return nothing if at end
		if atEOF && dataLen == 0 {
			return 0, nil, nil
		}

		if i := bytes.Index(data, splitByte); i > 0 {
			return i + splitLen, data[0:i], nil
		}

		// return data at end
		if atEOF {
			return dataLen, data, nil
		}

		return 0, nil, nil
	}
}
