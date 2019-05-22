package converter

import (
	"fmt"
	"io"
	"strings"

	"github.com/russross/blackfriday/v2"
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
	case blackfriday.HorizontalRule:
		fmt.Fprintf(w, "----\n")
	case blackfriday.Link:
		c.buffring.enable = entering
		if !entering {
			fmt.Fprintf(w, "[%s|%s]", c.buffring.buffer, node.LinkData.Destination)
			c.buffring.buffer = ""
		}
	case blackfriday.Text:
		if c.buffring.enable {
			c.buffring.buffer += string(node.Literal)
			break
		}

		if len(c.listPrefix) > 0 {
			fmt.Fprintf(w, "%s ", strings.Join(c.listPrefix, ""))
		}
		fmt.Fprintf(w, "%s\n", node.Literal)
	}
	return blackfriday.GoToNext
}

func (c *confluenceRenderer) RenderHeader(w io.Writer, ast *blackfriday.Node) {
	// TODO
}
func (c *confluenceRenderer) RenderFooter(w io.Writer, ast *blackfriday.Node) {
	// TODO
}
