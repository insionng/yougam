// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package precis

import (
	"unicode"

	"github.com/insionng/yougam/libraries/x/text/runes"
	"github.com/insionng/yougam/libraries/x/text/transform"
	"github.com/insionng/yougam/libraries/x/text/unicode/norm"
	"github.com/insionng/yougam/libraries/x/text/width"
)

var (
	Nickname              Profile = nickname          // Implements the Nickname profile specified in RFC 7700.
	UsernameCaseMapped    Profile = usernamecasemap   // Implements the UsernameCaseMapped profile specified in RFC 7613.
	UsernameCasePreserved Profile = usernamenocasemap // Implements the UsernameCasePreserved profile specified in RFC 7613.
	OpaqueString          Profile = opaquestring      // Implements the OpaqueString profile defined in RFC 7613 for passwords and other secure labels.
)

// TODO: mvl: "Ultimately, I would manually define the structs for the internal
// profiles. This avoid pulling in unneeded tables when they are not used."
var (
	nickname Profile = NewFreeform(
		AdditionalMapping(func() transform.Transformer {
			return &nickAdditionalMapping{}
		}),
		IgnoreCase,
		Norm(norm.NFKC),
		DisallowEmpty,
	)
	usernamecasemap Profile = NewIdentifier(
		AllowWide,
		FoldCase(),
		Norm(norm.NFC),
		// TODO: BIDI rule
	)
	usernamenocasemap Profile = NewIdentifier(
		AllowWide,
		Norm(norm.NFC),
		Width(width.Fold), // TODO: Is this correct?
		// TODO: BIDI rule
	)
	opaquestring Profile = NewFreeform(
		AdditionalMapping(func() transform.Transformer {
			return runes.Map(func(r rune) rune {
				if unicode.Is(unicode.Zs, r) {
					return ' '
				}
				return r
			})
		}),
		Norm(norm.NFC),
		DisallowEmpty,
	)
)
