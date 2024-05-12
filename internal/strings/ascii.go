// Copyright 2010 The Go Authors.  All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are
// met:
//
// - Redistributions of source code must retain the above copyright
// notice, this list of conditions and the following disclaimer.
//
// - Redistributions in binary form must reproduce the above
// copyright notice, this list of conditions and the following disclaimer
// in the documentation and/or other materials provided with the
// distribution.
//
// - Neither the name of Google Inc. nor the names of its
// contributors may be used to endorse or promote products derived from
// this software without specific prior written permission.
//
// Original license: https://github.com/golang/protobuf/blob/master/LICENSE
// Original CamelCase func: https://github.com/golang/protobuf/blob/master/protoc-gen-go/generator/generator.go#L2648

package strings

// Is c an ASCII lower-case letter?
func isASCIILower(c byte) bool {
	return 'a' <= c && c <= 'z'
}

// Is c an ASCII digit?
func isASCIIDigit(c byte) bool {
	return '0' <= c && c <= '9'
}

// appendLowercaseSequence appends the lowercase sequence from s that begins at i into t
// returns the new t that contains all the chain of characters that should be lowercase
// and the new index where to start counting from.
func appendLowercaseSequence(s string, i int, t []byte) ([]byte, int) {
	for i+1 < len(s) && isASCIILower(s[i+1]) {
		i++
		t = append(t, s[i])
	}
	return t, i
}
