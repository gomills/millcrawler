package urlextraction

import (
	"context"
	"slices"
	"strings"

	"github.com/gomills/gofocusedcrawler/internal/config"
	"github.com/gomills/gofocusedcrawler/pkg/queue"
	"github.com/gomills/gofocusedcrawler/pkg/urlvalidator"
	tree_sitter "github.com/tree-sitter/go-tree-sitter"
	tree_sitter_javascript "github.com/tree-sitter/tree-sitter-javascript/bindings/go"
)

// relevantParentTypes lists JavaScript AST node types (names according to Tree-Sitter) relevant for URL extraction. This is completely heuristic
var relevantParentTypes = []string{
	"call_expression", "arguments", "import_statement",
	"pair", "binary_expression", "assignment_pattern",
	"variable_declarator", "assignment_expression",
}

// TraverseJs parses JavaScript code, extracts possible URLs, validates them, and adds them to the queue.
func TraverseJs(ctx context.Context, config *config.Config, code []byte, qp *queue.Queue, domain string, registeredDomain string) {

	foundUrls := []string{}
	parser := tree_sitter.NewParser()
	parser.SetLanguage(tree_sitter.NewLanguage(tree_sitter_javascript.Language()))
	defer parser.Close() // you must always call Close on an object that allocates memory from C

	tree := parser.Parse(code, nil)
	defer tree.Close()
	if tree == nil {
		return
	}

	root := tree.RootNode()
	cursor := root.Walk()

	traverseLoop(cursor, &foundUrls, code)

	for _, parsedUrl := range foundUrls {
		parsedUrl, err := urlvalidator.ValidateStringForUrl(config, parsedUrl, domain, registeredDomain)
		if err == nil {
			qp.AddUrl(ctx, parsedUrl)
		}
	}

}

// traverseLoop walks the AST and processes each node for potential URLs.
func traverseLoop(cursor *tree_sitter.TreeCursor, foundUrls *[]string, code []byte) {

	for {

		if cursor.GotoFirstChild() {
			processNode(cursor.Node(), foundUrls, code)
			continue
		}

		for !cursor.GotoNextSibling() {
			if !cursor.GotoParent() {
				return
			}
		}

		processNode(cursor.Node(), foundUrls, code)

	}
}

// processNode dispatches node processing based on node type (string, template_string, comment).
func processNode(node *tree_sitter.Node, foundUrls *[]string, code []byte) {

	nodeType := node.Kind()

	switch nodeType {

	case "string", "template_string":
		processStringNode(node, foundUrls, code)

	case "comment":
		processCommentNode(node, foundUrls, code)
	}

}

// processStringNode extracts string values from relevant AST nodes and adds them to foundUrls.
func processStringNode(node *tree_sitter.Node, foundUrls *[]string, code []byte) {

	parent := node.Parent()
	if parent == nil {
		return
	}

	// Make sure the string is inside an interesting block for us
	for !slices.Contains(relevantParentTypes, parent.Kind()) {
		parent = parent.Parent()
		if parent == nil {
			return
		}
	}

	node_string := strings.Trim(node.Utf8Text(code), `"\()'`)
	node_string = strings.Trim(node.Utf8Text(code), "`")
	if node_string == "" {
		return
	}

	*foundUrls = append(*foundUrls, node_string)

}

// processCommentNode extracts absolute URLs from comment nodes and adds them to foundUrls.
func processCommentNode(node *tree_sitter.Node, foundUrls *[]string, code []byte) {

	if matches := AbsoluteUrlPattern.FindAllString(node.Utf8Text(code), -1); matches != nil {

		for _, match := range matches {
			if match != "" {
				*foundUrls = append(*foundUrls, match)
			}

		}
	}

}
