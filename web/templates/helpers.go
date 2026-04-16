package templates

import (
	"fmt"

	"github.com/a-h/templ"
)

// StyleAttr returns a templ.Attributes map with a style key
func StyleAttr(style string) templ.Attributes {
	return templ.Attributes{"style": style}
}

// PillStyle returns style for a pill component
func PillStyle(bg, color string) templ.Attributes {
	return templ.Attributes{"style": fmt.Sprintf("background:%s; color:%s", bg, color)}
}

// DotStyle returns style for a dot element
func DotStyle(bg string) templ.Attributes {
	return templ.Attributes{"style": fmt.Sprintf("background:%s", bg)}
}

// WidthStyle returns style for progress bar width
func WidthStyle(pct int) templ.Attributes {
	return templ.Attributes{"style": fmt.Sprintf("width:%d%%", pct)}
}

// WidthBgStyle returns style for progress bar width + custom background
func WidthBgStyle(pct int, bg string) templ.Attributes {
	return templ.Attributes{"style": fmt.Sprintf("width:%d%%; background:%s", pct, bg)}
}

// ColorStyle returns a style with just color
func ColorStyle(color string) templ.Attributes {
	return templ.Attributes{"style": fmt.Sprintf("color:%s", color)}
}

// HeightStyle returns style for height percentage
func HeightStyle(height string) templ.Attributes {
	return templ.Attributes{"style": fmt.Sprintf("height:%s", height)}
}

// BorderLeftStyle returns border-left style
func BorderLeftStyle(color string) templ.Attributes {
	return templ.Attributes{"style": fmt.Sprintf("border-left: 3px solid %s", color)}
}
