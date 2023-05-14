package structs

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"torch/src/utils"
)

var colorCodeRegex = regexp.MustCompile(`(?:§|&)([a-fA-F0-9k-oK-OrR]|#[a-fA-F0-9]{6})([^§&]*)`)
var cleanRegex = regexp.MustCompile(`[§&][a-fA-F0-9k-oK-OrR]|[§&]#[a-fA-F0-9]{6}`)

type ParsedText struct {
	Raw   string `json:"raw"`
	Clean string `json:"clean"`
	Html  string `json:"html"`
	Json  string `json:"json"`
}

type JsonSegment struct {
	Text   string   `json:"text"`
	Styles []string `json:"styles"`
}

func Parse(object interface{}) (res *ParsedText) {
	switch v := object.(type) {
	case string:
		json := Json(v)
		return &ParsedText{v, Clean(v), strings.ReplaceAll(Html(json), "{{newline}}", "<br />"), json}
	case map[string]interface{}:
		raw := ParseTextObject(v)
		json := Json(raw)
		return &ParsedText{raw, Clean(raw), strings.ReplaceAll(Html(json), "{{newline}}", "<br />"), json}
	default:
		return nil
	}
}

func ParseTextObject(object map[string]interface{}) (result string) {
	{
		if v, ok := object["color"].(string); ok {
			color, ok := utils.ParseColor(v)
			if ok {
				result += color.ToRaw()
			}
		}
	}

	if v, ok := object["bold"]; ok && parseBool(v) {
		result += "\u00A7l"
	}

	if v, ok := object["italic"]; ok && parseBool(v) {
		result += "\u00A7o"
	}

	if v, ok := object["underlined"]; ok && parseBool(v) {
		result += "\u00A7n"
	}

	if v, ok := object["strikethrough"]; ok && parseBool(v) {
		result += "\u00A7m"
	}

	if v, ok := object["obfuscated"]; ok && parseBool(v) {
		result += "\u00A7k"
	}

	if text, ok := object["text"].(string); ok {
		result += text
	}

	if extra, ok := object["extra"].([]interface{}); ok {
		for _, v := range extra {
			if v2, ok := v.(map[string]interface{}); ok {
				result += ParseTextObject(v2)
			}
		}
	}

	return

}

func Clean(str string) (clean string) {
	return cleanRegex.ReplaceAllString(str, "")
}

func Json(str string) string {

	str = strings.ReplaceAll(str, " ", "{{space}}")
	str = strings.ReplaceAll(str, "\n", "{{newline}}")
	str = "§r" + str + "§r"

	fmt.Println(str)

	matches := colorCodeRegex.FindAllStringSubmatch(str, -1)

	segments := make(map[string]JsonSegment)
	styles := []string{}
	segmentIndex := 1

	for _, match := range matches {
		code := strings.ToLower(match[1])
		stylesCopy := make([]string, len(styles))
		copy(stylesCopy, styles)

		if len(code) == 1 {
			switch code {
			case "k":
				stylesCopy = append(stylesCopy, "obfuscated")
			case "l":
				stylesCopy = append(stylesCopy, "bold")
			case "m":
				stylesCopy = append(stylesCopy, "strikethrough")
			case "n":
				stylesCopy = append(stylesCopy, "underline")
			case "o":
				stylesCopy = append(stylesCopy, "italic")
			case "r":
				stylesCopy = []string{}
			default:
				if color, ok := utils.ParseColor(code); ok {
					stylesCopy = []string{fmt.Sprintf("color=%s", color.ToHex())}
				}
			}
		} else if strings.HasPrefix(code, "#") {
			stylesCopy = []string{fmt.Sprintf("color=%s", code)}
		}
		styles = stylesCopy

		text := strings.TrimSpace(match[2])
		if text != "" {
			segments[fmt.Sprint(segmentIndex)] = JsonSegment{
				Text:   text,
				Styles: styles,
			}
			segmentIndex++
		}
	}

	jsonOutput, _ := json.Marshal(segments)
	return strings.ReplaceAll(strings.ReplaceAll(string(jsonOutput), "{{newline}}", "\n"), "{{space}}", " ")
}

func Html(str string) string {

	str = strings.ReplaceAll(str, "\n", "{{newline}}")

	html := strings.Builder{}
	html.WriteString("<span>")

	segments := make(map[string]JsonSegment)
	err := json.Unmarshal([]byte(str), &segments)
	if err != nil {
		return ""
	}

	htmlOutput := strings.Builder{}
	for i := 1; i <= len(segments); i++ {
		index := fmt.Sprint(i)
		segment := segments[index]
		styleString := ""
		for _, style := range segment.Styles {
			if strings.HasPrefix(style, "color=") {
				color := strings.TrimPrefix(style, "color=")
				styleString += fmt.Sprintf("color: %s; ", color)
			} else {
				switch style {
				case "bold":
					styleString += "font-weight: bold; "
				case "strikethrough":
					styleString += "text-decoration: line-through; "
				case "underline":
					styleString += "text-decoration: underline; "
				case "italic":
					styleString += "font-style: italic; "
				}
			}
		}
		textString := strings.ReplaceAll(segment.Text, " ", "&nbsp;")
		htmlOutput.WriteString(fmt.Sprintf("<span style=\"%s\">%s</span>", styleString, textString))
	}

	html.WriteString(htmlOutput.String())
	html.WriteString("</span>")

	return html.String()
}

func parseBool(value interface{}) bool {
	switch v := value.(type) {
	case bool:
		return v
	case string:
		return strings.ToLower(v) == "true"
	case int, uint, int8, uint8, int16, uint16, int32, uint32, int64, uint64:
		return v == 1
	default:
		return false
	}
}
