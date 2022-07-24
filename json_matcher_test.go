package matcher_test

import (
	"testing"

	matcher "github.com/panta/go-json-matcher"
)

func TestJsonStringMatches(t *testing.T) {
	type args struct {
		j     string
		jSpec string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{name: "bool-true-eq", args: args{
			j:     `true`,
			jSpec: `true`,
		}, want: true},
		{name: "bool-false-eq", args: args{
			j:     `false`,
			jSpec: `false`,
		}, want: true},
		{name: "bool-ne", args: args{
			j:     `true`,
			jSpec: `false`,
		}, want: false},
		{name: "bool-true-spec", args: args{
			j:     `true`,
			jSpec: `"#boolean"`,
		}, want: true},
		{name: "bool-false-spec", args: args{
			j:     `true`,
			jSpec: `"#boolean"`,
		}, want: true},
		{name: "bool-spec-wrongtype", args: args{
			j:     `"hello"`,
			jSpec: `"#boolean"`,
		}, want: false},
		{name: "bool-alias-spec", args: args{
			j:     `true`,
			jSpec: `"#bool"`,
		}, want: true},
		{name: "int-eq", args: args{
			j:     `123`,
			jSpec: `123`,
		}, want: true},
		{name: "int-ne", args: args{
			j:     `123`,
			jSpec: `124`,
		}, want: false},
		{name: "int-spec", args: args{
			j:     `123`,
			jSpec: `"#number"`,
		}, want: true},
		{name: "float-eq", args: args{
			j:     `123.52`,
			jSpec: `123.52`,
		}, want: true},
		{name: "float-ne", args: args{
			j:     `123.52`,
			jSpec: `124.52`,
		}, want: false},
		{name: "float-spec", args: args{
			j:     `123.52`,
			jSpec: `"#number"`,
		}, want: true},
		{name: "number-spec-wrongtype", args: args{
			j:     `"hello"`,
			jSpec: `"#number"`,
		}, want: false},
		{name: "string-eq", args: args{
			j:     `"the quick brown fox"`,
			jSpec: `"the quick brown fox"`,
		}, want: true},
		{name: "string-zero-eq", args: args{
			j:     `""`,
			jSpec: `""`,
		}, want: true},
		{name: "string-ne", args: args{
			j:     `"the quick brown fox"`,
			jSpec: `"the slow brown fox"`,
		}, want: false},
		{name: "string-spec", args: args{
			j:     `"the quick brown fox"`,
			jSpec: `"#string"`,
		}, want: true},
		{name: "string-zero-spec", args: args{
			j:     `""`,
			jSpec: `"#string"`,
		}, want: true},
		{name: "string-spec-wrongtype", args: args{
			j:     `123.52`,
			jSpec: `"#string"`,
		}, want: false},
		{name: "array-eq", args: args{
			j:     `[ true, 42, 5.52, "hello" ]`,
			jSpec: `[ true, 42, 5.52, "hello" ]`,
		}, want: true},
		{name: "array-zero-eq", args: args{
			j:     `[]`,
			jSpec: `[]`,
		}, want: true},
		{name: "array-ne", args: args{
			j:     `[ true, 42, 5.52, "hello" ]`,
			jSpec: `[ true, 42, 0.52, "hello" ]`,
		}, want: false},
		{name: "array-with-notarray", args: args{
			j:     `[ true, 42, 5.52, "hello" ]`,
			jSpec: `"hello"`,
		}, want: false},
		{name: "first-array-superset", args: args{
			j:     `[ true, 42, 5.52, "hello", "uh?" ]`,
			jSpec: `[ true, 42, 5.52, "hello" ]`,
		}, want: false},
		{name: "first-array-subset", args: args{
			j:     `[ true, 42, 5.52 ]`,
			jSpec: `[ true, 42, 5.52, "hello" ]`,
		}, want: false},
		{name: "array-with-bad-pattern", args: args{
			j:     `[ true, 42, 5.52, "hello" ]`,
			jSpec: `[ true, 42, 5.52, "#regex *+" ]`,
		}, want: false, wantErr: true},
		{name: "array-spec", args: args{
			j:     `[ true, 42, 5.52, "hello" ]`,
			jSpec: `"#array"`,
		}, want: true},
		{name: "array-zero-spec", args: args{
			j:     `[]`,
			jSpec: `"#array"`,
		}, want: true},
		{name: "array-spec-wrongtype", args: args{
			j:     `5`,
			jSpec: `"#array"`,
		}, want: false},
		{name: "array-of-spec", args: args{
			j:     `[ 12, 42, 5.52, 0, 7 ]`,
			jSpec: `[ "#array-of", "#number" ]`,
		}, want: true},
		{name: "array-of-spec-extra-arg", args: args{
			j:     `[ 12, 42, 5.52, 0, 7 ]`,
			jSpec: `[ "#array-of", "#number", "uh?" ]`,
		}, want: false},
		{name: "array-of-obj-spec", args: args{
			j:     `[ { "id": 1, "name": "joe" }, { "id": 1, "name": "jack" } ]`,
			jSpec: `[ "#array-of", { "id": "#number", "name": "#string" } ]`,
		}, want: true},
		{name: "array-of-spec-fail", args: args{
			j:     `[ 12, 42, 5.52, 0, 7 ]`,
			jSpec: `[ "#array-of", "#boolean" ]`,
		}, want: false},
		{name: "object-eq", args: args{
			j:     `{ "foo": "bar", "i": 42, "f": 5.52, "ok": true }`,
			jSpec: `{ "foo": "bar", "i": 42, "f": 5.52, "ok": true }`,
		}, want: true},
		{name: "object-ne", args: args{
			j:     `{ "foo": "bar", "i": 42, "f": 5.52, "ok": true }`,
			jSpec: `{ "foo": "bar", "i": 52, "f": 5.52, "ok": true }`,
		}, want: false},
		{name: "object-with-notobject", args: args{
			j:     `{ "foo": "bar", "i": 42, "f": 5.52, "ok": true }`,
			jSpec: `[ 2, 3 ]`,
		}, want: false},
		{name: "first-object-superset", args: args{
			j:     `{ "foo": "bar", "i": 42, "f": 5.52, "ok": true }`,
			jSpec: `{ "foo": "bar", "i": 42, "f": 5.52 }`,
		}, want: true},
		{name: "first-object-subset", args: args{
			j:     `{ "foo": "bar", "i": 42, "f": 5.52 }`,
			jSpec: `{ "foo": "bar", "i": 42, "f": 5.52, "ok": true }`,
		}, want: false},
		{name: "object-with-bad-pattern", args: args{
			j:     `{ "foo": "bar", "i": 42, "f": 5.52, "ok": true, "name": "hello" }`,
			jSpec: `{ "foo": "bar", "i": 52, "f": 5.52, "ok": true, "name": "#regex *+" }`,
		}, want: false, wantErr: true},
		{name: "object-zero-eq", args: args{
			j:     `{}`,
			jSpec: `{}`,
		}, want: true},
		{name: "object-spec", args: args{
			j:     `{ "foo": "bar", "i": 42, "f": 5.52, "ok": true }`,
			jSpec: `"#object"`,
		}, want: true},
		{name: "object-zero-spec", args: args{
			j:     `{}`,
			jSpec: `"#object"`,
		}, want: true},
		{name: "object-spec-wrongtype", args: args{
			j:     `5`,
			jSpec: `"#object"`,
		}, want: false},
		{name: "object-of-spec", args: args{
			j: `{ "foo": "bar", "id": 2, "ok": false, "els": [ "foo", "bar", "xyz" ],
"objs": [ { "id": 1, "name": "joe" }, { "id": 1, "name": "jack" } ] }`,
			jSpec: `{ "foo": "#string", "id": "#number", "ok": "#boolean", "els": "#array",
"objs": [ "#array-of", { "id": "#number", "name": "#string" } ] }`,
		}, want: true},
		{name: "ignore-spec-number", args: args{
			j:     `15`,
			jSpec: `"#ignore"`,
		}, want: true},
		{name: "ignore-spec-string", args: args{
			j:     `"hello"`,
			jSpec: `"#ignore"`,
		}, want: true},
		{name: "ignore-spec-null", args: args{
			j:     `null`,
			jSpec: `"#ignore"`,
		}, want: true},
		{name: "ignore-spec-empty-array", args: args{
			j:     `[]`,
			jSpec: `"#ignore"`,
		}, want: true},
		{name: "ignore-spec-object-key-present", args: args{
			j:     `{ "key": "score", "value": 5 }`,
			jSpec: `{ "key": "#ignore", "value": 5 }`,
		}, want: true},
		{name: "ignore-spec-object-key-absent", args: args{
			j:     `{ "value": 5 }`,
			jSpec: `{ "key": "#ignore", "value": 5 }`,
		}, want: true},
		{name: "ignore-spec-object-key-present", args: args{
			j:     `{ "key": "score", "value": 5 }`,
			jSpec: `{ "key": "#ignore" }`,
		}, want: true},
		{name: "null-eq", args: args{
			j:     `null`,
			jSpec: `null`,
		}, want: true},
		{name: "null-ne", args: args{
			j:     `null`,
			jSpec: `false`,
		}, want: false},
		{name: "null-spec", args: args{
			j:     `null`,
			jSpec: `"#null"`,
		}, want: true},
		{name: "notnull-spec", args: args{
			j:     `5`,
			jSpec: `"#notnull"`,
		}, want: true},
		{name: "notnull-spec-fail", args: args{
			j:     `null`,
			jSpec: `"#notnull"`,
		}, want: false},
		{name: "present-spec", args: args{
			j:     `5`,
			jSpec: `"#present"`,
		}, want: true},
		{name: "present-null-spec", args: args{
			j:     `null`,
			jSpec: `"#present"`,
		}, want: true},
		{name: "object-key-present-spec", args: args{
			j:     `{ "foo": "bar", "id": 2, "ok": false }`,
			jSpec: `{ "foo": "#present", "id": "#number", "ok": "#boolean" }`,
		}, want: true},
		{name: "object-key-nullvalue-present-spec", args: args{
			j:     `{ "foo": null, "id": 2, "ok": false }`,
			jSpec: `{ "foo": "#present", "id": "#number", "ok": "#boolean" }`,
		}, want: true},
		{name: "object-key-present-spec-fail", args: args{
			j:     `{ "id": 2, "ok": false }`,
			jSpec: `{ "foo": "#present", "id": "#number", "ok": "#boolean" }`,
		}, want: false},
		{name: "notpresent-spec", args: args{
			j:     `5`,
			jSpec: `"#notpresent"`,
		}, want: false},
		{name: "notpresent-null-spec", args: args{
			j:     `null`,
			jSpec: `"#notpresent"`,
		}, want: false},
		{name: "object-key-notpresent-spec-fail", args: args{
			j:     `{ "foo": "bar", "id": 2, "ok": false }`,
			jSpec: `{ "foo": "#notpresent", "id": "#number", "ok": "#boolean" }`,
		}, want: false},
		{name: "object-key-nullvalue-notpresent-spec-fail", args: args{
			j:     `{ "foo": null, "id": 2, "ok": false }`,
			jSpec: `{ "foo": "#notpresent", "id": "#number", "ok": "#boolean" }`,
		}, want: false},
		{name: "object-key-notpresent-spec-ok", args: args{
			j:     `{ "id": 2, "ok": false }`,
			jSpec: `{ "foo": "#notpresent", "id": "#number", "ok": "#boolean" }`,
		}, want: true},
		{name: "date-spec", args: args{
			j:     `"2012-09-27"`,
			jSpec: `"#date"`,
		}, want: true},
		{name: "date-spec-fail", args: args{
			j:     `2012`,
			jSpec: `"#date"`,
		}, want: false},
		{name: "datetime-spec", args: args{
			j:     `"2012-09-27T13:42:24+02:00"`,
			jSpec: `"#datetime"`,
		}, want: true},
		{name: "datetime-spec-fail", args: args{
			j:     `"2012-09-27"`,
			jSpec: `"#datetime"`,
		}, want: false},
		{name: "datetime-spec-fail-2", args: args{
			j:     `2012`,
			jSpec: `"#date"`,
		}, want: false},
		{name: "uuid-spec-v4", args: args{
			j:     `"a5bf6b35-61b2-4187-8396-463a3d6c742b"`,
			jSpec: `"#uuid"`,
		}, want: true},
		{name: "uuid-spec-v1", args: args{
			j:     `"f183ee98-07a3-11ed-861d-0242ac120002"`,
			jSpec: `"#uuid"`,
		}, want: true},
		{name: "uuid-spec-fail", args: args{
			j:     `"gosh-a30ae5f3-818a-4b29-8815-320f8561021a"`,
			jSpec: `"#uuid"`,
		}, want: false},
		{name: "uuid-spec-wrongtype", args: args{
			j:     `123.52`,
			jSpec: `"#uuid"`,
		}, want: false},
		{name: "uuid-v4-spec-v4", args: args{
			j:     `"a5bf6b35-61b2-4187-8396-463a3d6c742b"`,
			jSpec: `"#uuid-v4"`,
		}, want: true},
		{name: "uuid-v4-spec-fail-v1", args: args{
			j:     `"f183ee98-07a3-11ed-861d-0242ac120002"`,
			jSpec: `"#uuid-v4"`,
		}, want: false},
		{name: "uuid-v4-spec-fail", args: args{
			j:     `"a5bf6b35-61b2-3187-8396-463a3d6c742b"`,
			jSpec: `"#uuid-v4"`,
		}, want: false},
		{name: "uuid-v4-spec-fail-2", args: args{
			j:     `"gosh-a30ae5f3-818a-4b29-8815-320f8561021a"`,
			jSpec: `"#uuid-v4"`,
		}, want: false},
		{name: "uuid-v4-spec-wrongtype", args: args{
			j:     `123.52`,
			jSpec: `"#uuid-v4"`,
		}, want: false},
		{name: "regex-spec", args: args{
			j:     `"This is fun"`,
			jSpec: `"#regex ^This is [a-z]{3}$"`,
		}, want: true},
		{name: "regex-spec-fail", args: args{
			j:     `"This is f4n"`,
			jSpec: `"#regex ^This is [a-z]{3}$"`,
		}, want: false},
		{name: "regex-spec-fail-2", args: args{
			j:     `"This is fun."`,
			jSpec: `"#regex ^This is [a-z]{3}$"`,
		}, want: false},
		{name: "regex-spec-invalid", args: args{
			j:     `"This is fun"`,
			jSpec: `"#regex +*{3a"`,
		}, want: false, wantErr: true},
		{name: "regex-spec-wrongtype", args: args{
			j:     `42`,
			jSpec: `"#regex ^This is [a-z]{3}$"`,
		}, want: false},
		{name: "readme-example", args: args{
			j: `{
  "id": "adb43c69-f8d9-4108-a2da-d740a2a800ec",
  "title": "A short article.",
  "body": "Lorem ipsum dolor sit amet, consectetur adipiscing elit, ...",
  "publish": true,
  "type": "articles",
  "created": "2022-07-20T09:56:29.000Z",
  "updated": "2022-07-20T10:12:47.000Z",
  "section_id": 42,
  "tags": [ "society", "essays", "history" ]
}`,
			jSpec: `{
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
}`,
		}, want: true},
		{name: "bad-json-first-argument", args: args{
			j:     `[A237`,
			jSpec: `#ignore`,
		}, want: false, wantErr: true},
		{name: "bad-json-second-argument", args: args{
			j:     `15`,
			jSpec: `[A237`,
		}, want: false, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := matcher.JSONStringMatches(tt.args.j, tt.args.jSpec)
			if (err != nil) != tt.wantErr {
				t.Errorf("JSONStringMatches() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("JSONStringMatches() got = %v, want %v", got, tt.want)
			}
		})
	}
}
