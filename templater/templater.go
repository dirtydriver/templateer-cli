package templater

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"text/template"
	"text/template/parse"

	"github.com/Masterminds/sprig/v3"
	"github.com/dirtydriver/templateer-cli/utils"
)

func collectPlaceholders(node parse.Node, placeholders map[string]struct{}) {
	switch n := node.(type) {
	case *parse.ListNode:
		if n == nil {
			return
		}
		for _, child := range n.Nodes {
			collectPlaceholders(child, placeholders)
		}
	case *parse.ActionNode:
		collectPlaceholders(n.Pipe, placeholders)
	case *parse.IfNode:
		collectPlaceholders(n.Pipe, placeholders)
		collectPlaceholders(n.List, placeholders)
		collectPlaceholders(n.ElseList, placeholders)
	case *parse.RangeNode:
		collectPlaceholders(n.Pipe, placeholders)
		collectPlaceholders(n.List, placeholders)
		collectPlaceholders(n.ElseList, placeholders)
	case *parse.WithNode:
		collectPlaceholders(n.Pipe, placeholders)
		collectPlaceholders(n.List, placeholders)
		collectPlaceholders(n.ElseList, placeholders)
	case *parse.TemplateNode:
		// {{ template "name" .Pipe }}
		collectPlaceholders(n.Pipe, placeholders)
	case *parse.PipeNode:
		for _, decl := range n.Decl {
			collectPlaceholders(decl, placeholders)
		}
		for _, cmd := range n.Cmds {
			collectPlaceholders(cmd, placeholders)
		}
	case *parse.CommandNode:
		for _, arg := range n.Args {
			collectPlaceholders(arg, placeholders)
		}
	case *parse.VariableNode:
		// VariableNode represents variables like {{ $name }} or declarations like {{ $name := ... }}.
		// We record the variable name so callers can see which template variables are in use.
		// Note: Ident typically contains a single entry like "$name".
		if len(n.Ident) == 0 {
			return
		}
		placeholder := strings.Join(n.Ident, ".")
		placeholders[placeholder] = struct{}{}
	case *parse.FieldNode:
		// FieldNode represents expressions like {{ .GroupID }}
		// The Ident slice contains the field names, which we join with dots.
		placeholder := strings.Join(n.Ident, ".")
		placeholders[placeholder] = struct{}{}
	case *parse.ChainNode:
		// ChainNode represents chained accesses like {{ .User.Name }}
		var base string
		if fieldNode, ok := n.Node.(*parse.FieldNode); ok {
			base = strings.Join(fieldNode.Ident, ".")
		} else {
			// Fallback: use the node's string representation.
			base = fmt.Sprintf("%v", n.Node)
		}
		if len(n.Field) > 0 {
			base += "." + strings.Join(n.Field, ".")
		}
		placeholders[strings.TrimPrefix(base, ".")] = struct{}{}
	}
}

// CollectParameters analyzes template files and returns a list of unique parameter names used in them.
// It processes templates concurrently for better performance.
func CollectParameters(tempFiles []string) ([]string, error) {

	var (
		wg         sync.WaitGroup
		mu         sync.Mutex
		parameters = make(map[string]struct{})
		errChan    = make(chan error, len(tempFiles))
	)

	for _, file := range tempFiles {
		wg.Add(1)
		go func(file string) {

			defer wg.Done()
			tmpl, err := template.New(filepath.Base(file)).Funcs(sprig.TxtFuncMap()).ParseFiles(file)
			if err != nil {
				errChan <- err
				return
			}
			localParameters := make(map[string]struct{})
			collectPlaceholders(tmpl.Root, localParameters)
			mu.Lock()
			for param := range localParameters {
				parameters[param] = struct{}{}
			}
			mu.Unlock()
		}(file)

	}
	wg.Wait()
	close(errChan)
	var placeholderList []string
	for param := range parameters {

		placeholderList = append(placeholderList, param)
	}

	return utils.RemoveDuplicates(placeholderList), nil

}

// RenderTemplate processes a template file with the given parameters and returns the rendered content.
// It supports all standard Go template functionality plus Sprig template functions (http://masterminds.github.io/sprig/).
// This enables advanced template features like string manipulation, date formatting, math operations, and more.
// Examples of Sprig functions include: upper, lower, title, trim, default, date, repeat, etc.
func RenderTemplate(file string, params map[string]interface{}) (bytes.Buffer, error) {
	// Parse the template file
	tmpl, err := template.New(filepath.Base(file)).Funcs(sprig.TxtFuncMap()).ParseFiles(file)
	if err != nil {
		return bytes.Buffer{}, err // Return the error immediately
	}

	// Create a buffer to store the rendered output
	var output bytes.Buffer

	// Execute the template with the provided parameters
	err = tmpl.Execute(&output, params)
	if err != nil {
		return bytes.Buffer{}, err
	}

	// Return the rendered template as a string
	return output, nil
}

// IsTemplate checks if a file is a template by verifying if it has a .tmpl extension.
func IsTemplate(path string) bool {
	return strings.Contains(filepath.Ext(path), ".tmpl")
}

// WriteTemplate writes the rendered template content to the specified file path.
func WriteTemplate(path string, renderedTemplate *bytes.Buffer) error {

	file, err := os.Create(path)

	if err != nil {
		return err
	}

	defer func() {
		if cerr := file.Close(); cerr != nil {
			if err != nil {
				err = fmt.Errorf("write error: %v; additionally, close error: %v", err, cerr)
			} else {
				err = cerr
			}
		}
	}()

	if _, err := renderedTemplate.WriteTo(file); err != nil {
		return err
	}
	return nil
}
