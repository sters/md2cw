package converter

import (
	"strings"
	"testing"
)

func TestConvert(t *testing.T) {
	tests := []struct {
		name     string
		markdown string
		want     string
	}{
		// headers
		{"h1", "# foo", "h1. foo"},
		{"h2", "## foo", "h2. foo"},
		{"h3", "### foo", "h3. foo"},
		{"h4", "#### foo", "h4. foo"},
		{"h5", "##### foo", "h5. foo"},
		{"h6", "###### foo", "h6. foo"},

		// list
		{"list-asterisk", "* foo", "* foo"},
		{"list-hypfon", "- foo", "* foo"},
		{"list-plus", "+ foo", "* foo"},
		{"sub-list-2spaces", "- foo\n  - foo", "* foo\n** foo"},
		// 2spaces has parser's bug, maybe...
		{"deep-list-2spaces", "* foo\n  * foo\n    * foo\n  * foo\n* foo", "* foo\n** foo\n** foo\n** foo\n* foo"},
		// 4spaces works fine.
		{"sub-list-4spaces", "- foo\n    - foo", "* foo\n** foo"},
		{"deep-list-4spaces", "* foo\n    * foo\n        * foo\n    * foo\n* foo", "* foo\n** foo\n*** foo\n** foo\n* foo"},

		// number-list
		{"numberlist", "1. foo", "# foo"},
		{"numberlist items", "1. foo\n1. foo", "# foo\n# foo"},
		{"numberlist items number", "1. foo\n2. foo", "# foo\n# foo"},
		{"sub-numberlist-2spaces", "1. foo\n  1. foo", "# foo\n## foo"},
		// 2spaces has parser's bug, maybe...
		{"deep-numberlist-2spaces", "1. foo\n  1. foo\n    1. foo\n  1. foo\n1. foo", "# foo\n## foo\n## foo\n## foo\n# foo"},
		// 4spaces works fine.
		{"sub-numberlist-4spaces", "1. foo\n    1. foo", "# foo\n## foo"},
		{"deep-numberlist-4spaces", "1. foo\n    1. foo\n        1. foo\n    1. foo\n1. foo", "# foo\n## foo\n### foo\n## foo\n# foo"},

		// mixed list
		{"list-number", "- foo\n    1. foo\n- foo", "* foo\n*# foo\n* foo"},
		{"number-list", "1. foo\n    - foo\n1. foo", "# foo\n#* foo\n# foo"},

		// hr
		{"hr", "foo\n\n----\n\nfoo", "foo\n\n----\n\nfoo"},

		// link
		{"link", "[foo](https://github.com/)", "[foo|https://github.com/]"},

		// list with link
		{"list with link", "- [foo](https://github.com/)\n- [foo](https://github.com/)", "* [foo|https://github.com/]\n* [foo|https://github.com/]"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := Convert(test.markdown); strings.TrimSpace(test.want) != strings.TrimSpace(got) {
				t.Errorf(
					"\nConvert wants:\n-----\n%s\n-----\n\nConvert got:\n-----\n%s\n-----\n",
					strings.TrimSpace(test.want),
					strings.TrimSpace(got),
				)
			}
		})
	}
}
