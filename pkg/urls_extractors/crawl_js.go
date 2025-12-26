package urls_extractors

import (
	"context"
	"slices"
	"strings"

	"github.com/gomills/gofocusedcrawler/internal/config"
	"github.com/gomills/gofocusedcrawler/pkg/queue"
	"github.com/gomills/gofocusedcrawler/pkg/url_validator"
	tree_sitter "github.com/tree-sitter/go-tree-sitter"
	tree_sitter_javascript "github.com/tree-sitter/tree-sitter-javascript/bindings/go"
)

// relevantParentTypes lists JavaScript AST node types (names according to Tree-Sitter) relevant for URL extraction. This is completely heuristic
var relevantParentTypes = []string{
	"call_expression", "arguments", "import_statement",
	"pair", "binary_expression", "assignment_pattern",
	"variable_declarator", "assignment_expression",
}

// CrawlJs parses JavaScript code, extracts possible URLs, validates them, and adds them to the queue.
func CrawlJs(ctx context.Context, config *config.Config, code []byte, qp *queue.Queue, domain string, registeredDomain string) {

	foundUrls := []string{}

	// instantiate new Tree-Sitter parser with .js as language
	parser := tree_sitter.NewParser()
	parser.SetLanguage(tree_sitter.NewLanguage(tree_sitter_javascript.Language()))
	defer parser.Close() // must always call Close on an object that allocates memory from C

	// parse the .js file and retrieve the AST tree
	tree := parser.Parse(code, nil)
	defer tree.Close()
	if tree == nil {
		return
	}

	// retrieve the root node and the cursor for iteration
	root := tree.RootNode()
	cursor := root.Walk()

	// move to this function to do the whole AST iteration (not recursive)
	traverseLoop(cursor, &foundUrls, code)

	for _, parsedUrl := range foundUrls {
		parsedUrl, err := url_validator.ValidateStringForUrl(config, parsedUrl, domain, registeredDomain)
		if err == nil {
			qp.AddUrl(ctx, parsedUrl)
		}
	}

}

// traverseLoop walks the AST from bottom left to right until root and processes each node for potential URLs.
// It's NOT recursive.
func traverseLoop(cursor *tree_sitter.TreeCursor, foundUrls *[]string, code []byte) {

	for {

		// try to go to child. If success, process it for urls in processNode()
		if cursor.GotoFirstChild() {
			processNode(cursor.Node(), foundUrls, code)
			continue
		}

		// if there was no child move to the next sibling.
		// if no sibling, go to parent until there's a sibling.
		for !cursor.GotoNextSibling() {

			// if neither sibling nor parent, we're in root again, break.
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

	// get the parent
	parent := node.Parent()
	if parent == nil {
		return
	}

	// and make sure the string is inside an interesting block. Keep walking upwards in AST
	// for nested code scenarios.
	for !slices.Contains(relevantParentTypes, parent.Kind()) {
		parent = parent.Parent()
		if parent == nil {
			return
		}
	}

	// trim all surrounding delimitator elements found in both string and template_string parents.
	node_string := strings.Trim(node.Utf8Text(code), `"\()'`)
	node_string = strings.Trim(node.Utf8Text(code), "`")

	if node_string == "" {
		return
	}

	// grow the found urls slice
	*foundUrls = append(*foundUrls, node_string)

}

// processCommentNode regexes for absolute URLs in comment nodes and adds them to foundUrls.
func processCommentNode(node *tree_sitter.Node, foundUrls *[]string, code []byte) {

	if matches := AbsoluteUrlPattern.FindAllString(node.Utf8Text(code), -1); matches != nil {

		for _, match := range matches {
			if match != "" {
				*foundUrls = append(*foundUrls, match)
			}

		}
	}

}
