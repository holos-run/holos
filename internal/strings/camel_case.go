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

// CamelCase is a special case of PascalCase with the difference that
// CamelCase will return the first value in the string as a lowercase character.
// In short, _my_field_name_2 becomes xMyFieldName_2.
func CamelCase(s string) string {
	if s == "" {
		return ""
	}
	t := make([]byte, 0, 32)
	i := 0
	if s[0] == '_' {
		// Need a capital letter; drop the '_'.
		t = append(t, 'x')
		i++
	}
	// If the first letter is a lowercase, we keep it as is.
	if len(t) == 0 && isASCIILower(s[i]) {
		t = append(t, s[i])
		t, i = appendLowercaseSequence(s, i, t)
		i++
	}
	return string(append(t, lookupAndReplacePascalCaseWords(s, i)...))
}
