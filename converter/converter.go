package converter

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/russross/blackfriday"
)

// Convert from markdown syntax to confluence wiki syntax
func Convert(markdownText string) string {
	return string(blackfriday.Run(
		[]byte(markdownText),
		blackfriday.WithRenderer(&confluenceRenderer{}),
	))
}

type confluenceRenderer struct {
	listPrefix []string
	buffring   buffering
}

type buffering struct {
	enable bool
	buffer string
}

func (c *confluenceRenderer) RenderNode(w io.Writer, node *blackfriday.Node, entering bool) blackfriday.WalkStatus {
	switch node.Type {
	case blackfriday.Heading:
		fmt.Fprintf(w, "\n")
		if entering {
			fmt.Fprintf(w, "h%d. ", node.HeadingData.Level)
		} else {
			fmt.Fprintf(w, "\n")
		}
	case blackfriday.List:
		if entering {
			prefix := "*"
			if 1 == (node.ListData.ListFlags & blackfriday.ListTypeOrdered) {
				prefix = "#"
			}
			c.listPrefix = append(c.listPrefix, prefix)
		} else {
			c.listPrefix = c.listPrefix[:len(c.listPrefix)-1]
		}
	case blackfriday.Item: // item means list-item
		if entering && len(c.listPrefix) > 0 {
			fmt.Fprintf(w, "%s ", strings.Join(c.listPrefix, ""))
		}
	case blackfriday.Paragraph: // paragraph like a <p>
		if !entering {
			fmt.Fprintf(w, "\n")
		}
	case blackfriday.HorizontalRule:
		fmt.Fprintf(w, "\n----\n\n")
	case blackfriday.Link:
		c.buffring.enable = entering
		if !entering {
			fmt.Fprintf(w, "[%s|%s]", c.buffring.buffer, node.LinkData.Destination)
			c.buffring.buffer = ""
		}
	case blackfriday.Code:
		fmt.Fprintf(w, "{{%s}}", node.Literal)
	case blackfriday.CodeBlock:
		fmt.Fprintf(w, "\n{code}\n%s\n{code}\n\n", strings.TrimSpace(string(node.Literal)))
	case blackfriday.Text:
		if c.buffring.enable {
			c.buffring.buffer += string(node.Literal)
			break
		}
		fmt.Fprintf(w, "%s", node.Literal)

	case blackfriday.Document:
		// do nothing.
	default:
		fmt.Fprintf(os.Stderr, "NodeType = %s is not supported\n", node.Type)
	}
	return blackfriday.GoToNext
}

func (c *confluenceRenderer) RenderHeader(w io.Writer, ast *blackfriday.Node) {
	// TODO
}
func (c *confluenceRenderer) RenderFooter(w io.Writer, ast *blackfriday.Node) {
	// TODO
}
