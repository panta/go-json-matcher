# JSON Matcher

[![tag](https://img.shields.io/github/tag/panta/go-json-matcher.svg)](https://github.com/panta/go-json-matcher/releases)
[![Go Reference](https://pkg.go.dev/badge/github.com/panta/go-json-matcher.svg)](https://pkg.go.dev/github.com/panta/go-json-matcher)
![Build Status](https://github.com/panta/go-json-matcher/actions/workflows/test.yml/badge.svg)
[![Go report](https://goreportcard.com/badge/github.com/panta/go-json-matcher)](https://goreportcard.com/report/github.com/panta/go-json-matcher)

JSON Matcher is a Go (_golang_) library to verify conformance of JSON objects to a
desired structure, according to provided patterns.

This is especially useful when writing unit-/integration- tests where exact comparisons
won't be viable (because some parts are particularly dynamic, think about
current timestamps, JWT tokens, UUIDs, ...).

For example, suppose you have a JSON API response like the following:

```json
{
  "id": "adb43c69-f8d9-4108-a2da-d740a2a800ec",
  "title": "A short article.",
  "body": "Lorem ipsum dolor sit amet, consectetur adipiscing elit, ...",
  "publish": true,
  "type": "articles",
  "created": "2022-07-20T09:56:29.000Z",
  "updated": "2022-07-20T10:12:47.000Z",
  "section_id": 42,
  "tags": [ "society", "essays", "history" ]
}
```

you could check the above structure with (error checks omitted here for brevity, you should check errors in actual code):

```go
matches, _ := matcher.JSONMatches(responseString, `{
  "id": "#uuid",
  "title": "#string",
  "body": "#string",
  "publish": "#boolean",
  "type": "articles",
  "created": "#datetime",
  "updated": "#datetime",
  "section_id": "#number",
  "tags": [ "#array-of", "#string" ],
  "error": "#notpresent"
}`)
if !matches {
    // ...
}
```

## Install

Using JSON Matcher is easy. Use `go get` to install the latest version of the library:

```shell
go get github.com/panta/go-json-matcher@latest
```

then import the library in you application:

```go
import "github.com/panta/go-json-matcher"
```

## Usage

The function `JSONMatches()` checks that the JSON string provided with the first
argument satisfies the pattern specified with the second argument.
The pattern can be a valid literal value (in that case an exact match will be required),
a special marker beginning with a `#` character as described below, or any combination
of these via arrays and objects.

### Supported markers

Marker | Description
------ | -----------
`#ignore` | Ignore the value or field
`#null` | Requires that the value is `null` (the element must be present though) 
`#notnull` | Requires that the value is not `null`
`#present` | Requires that the value is present (but it may be `null`)
`#notpresent` | Requires that the value is NOT present (not even `null`)
`#array` | Requires the value to be an array
`#object` | Requires the value to be an object
`#boolean` | Requires the value to be a boolean (either `true` or `false`)
`#number` | Requires the value to be a number
`#string` | Requires the value to be a string
`#uuid` | Requires the value to be a string conforming to a UUID
`#uuid-v4` | Requires the value to be a string conforming to a V4 UUID according to [RFC4122](https://datatracker.ietf.org/doc/html/rfc4122)
`#date` | Requires the value to be a string representing a valid ISO8601 date (format _YYYY-MM-DD_)
`#datetime` | Requires the value to be a string representing a valid RFC3339 / ISO8601 datetime
`#regex RE` | Requires the value to be a string matching the regular expression provided in `RE`

## License

Copyright (C) 2022 Marco Pantaleoni.

	Licensed under the Apache License, Version 2.0 (the "License");
	you may not use this software except in compliance with the License.
	You may obtain a copy of the License at
	
	       http://www.apache.org/licenses/LICENSE-2.0
	
	Unless required by applicable law or agreed to in writing, software
	distributed under the License is distributed on an "AS IS" BASIS,
	WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
	See the License for the specific language governing permissions and
	limitations under the License.

See the full license in [LICENSE](https://github.com/panta/go-json-matcher/blob/main/LICENSE) file.

## Acknowledgements

This library has been inspired by [orangain/json-fuzzy-match](https://github.com/orangain/json-fuzzy-match).
