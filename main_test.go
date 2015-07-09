package main

import (
	"fmt"
	"testing"
	"text/template"

	"github.com/trayio/bunny/nodes"
)

type fake struct {
	content string
}

func (f *fake) Write(p []byte) (int, error) {
	// string concatenation
	f.content = fmt.Sprintf("%s%s", f.content, p)
	return len(p), nil
}

func TestTemplate(t *testing.T) {
	expected := `[{rabbit, [{cluster_nodes, {['rabbit@rabbit1', 'rabbit@rabbit2'], disc}}]}].`
	f := &fake{}

	tmpl := template.Must(template.New("config").Parse(configTemplate))

	rabbits := []*nodes.Node{
		{Host: "rabbit1"},
		{Host: "rabbit2"},
	}
	if err := tmpl.Execute(f, rabbits); err != nil {
		t.Errorf("Failed to parse execute template\n")
	}

	if expected != f.content {
		t.Errorf("Unexpected content while parsing template:\nexpected: %s\nreceived: %s\n", expected, f.content)
	}

}
