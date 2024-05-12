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

// PascalCase returns the PascalCased name.
// If there is an interior underscore followed by a lower case letter,
// drop the underscore and convert the letter to upper case.
// There is a remote possibility of this rewrite causing a name collision,
// but it's so remote we're prepared to pretend it's nonexistent - since the
// C++ generator lowercases names, it's extremely unlikely to have two fields
// with different capitalizations.
// In short, _my_field_name_2 becomes XMyFieldName_2.
func PascalCase(s string) string {
	if s == "" {
		return ""
	}
	t := make([]byte, 0, 32)
	i := 0
	if s[0] == '_' {
		// Need a capital letter; drop the '_'.
		t = append(t, 'X')
		i++
	}
	return string(append(t, lookupAndReplacePascalCaseWords(s, i)...))
}

// lookupAndReplacePascalCaseWords lookups for words in the string starting in position i
// and replaces the snake case format to PascalCase.
// Invariant: if the next letter is lower case, it must be converted
// to upper case.
// That is, we process a word at a time, where words are marked by _ or
// upper case letter. Digits are treated as words.
func lookupAndReplacePascalCaseWords(s string, i int) []byte {
	t := make([]byte, 0, 32)
	for ; i < len(s); i++ {
		c := s[i]
		if c == '_' && i+1 < len(s) && isASCIILower(s[i+1]) {
			continue // Skip the underscore in s.
		}
		if isASCIIDigit(c) {
			t = append(t, c)
			continue
		}
		// Assume we have a letter now - if not, it's a bogus identifier.
		// The next word is a sequence of characters that must start upper case.
		if isASCIILower(c) {
			c ^= ' ' // Make it a capital letter.
		}
		t = append(t, c) // Guaranteed not lower case.
		// Accept lower case sequence that follows.
		t, i = appendLowercaseSequence(s, i, t)
	}
	return t
}
