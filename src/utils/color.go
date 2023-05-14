package utils

type Color rune

var (
	Black        Color = '0'
	DarkBlue     Color = '1'
	DarkGreen    Color = '2'
	DarkAqua     Color = '3'
	DarkRed      Color = '4'
	DarkPurple   Color = '5'
	Gold         Color = '6'
	Gray         Color = '7'
	DarkGray     Color = '8'
	Blue         Color = '9'
	Green        Color = 'a'
	Aqua         Color = 'b'
	Red          Color = 'c'
	LightPurple  Color = 'd'
	Yellow       Color = 'e'
	White        Color = 'f'
	MinecoinGold Color = 'g'
)

type Formatting rune

var (
	Obfuscated    Formatting = 'k'
	Bold          Formatting = 'l'
	Strikethrough Formatting = 'm'
	Underline     Formatting = 'n'
	Italic        Formatting = 'o'
	Reset         Formatting = 'r'
)

func ParseFormatting(value interface{}) (Formatting, bool) {
	switch value {
	case "k", "K", Obfuscated:
		return Obfuscated, true
	case "l", "L", Bold:
		return Bold, true
	case "m", "M", Strikethrough:
		return Strikethrough, true
	case "n", "N", Underline:
		return Underline, true
	case "o", "O", Italic:
		return Italic, true
	case "r", "R", Reset:
		return Reset, true
	default:
		return Reset, false
	}
}

func (f Formatting) ToName() string {
	switch f {
	case Obfuscated:
		return "Obfuscated"
	case Bold:
		return "Bold"
	case Strikethrough:
		return "Strikethrough"
	case Underline:
		return "Underline"
	case Italic:
		return "Italic"
	case Reset:
		return "Reset"
	default:
		return ""
	}
}

func (f Formatting) ToRaw() string {
	return "\u00A7" + string(f)
}

func ParseColor(value interface{}) (Color, bool) {
	switch value {
	case "0", "black", Black:
		return Black, true
	case "1", "dark_blue", DarkBlue:
		return DarkBlue, true
	case "2", "dark_green", DarkGreen:
		return DarkGreen, true
	case "3", "dark_aqua", DarkAqua:
		return DarkAqua, true
	case "4", "dark_red", DarkRed:
		return DarkRed, true
	case "5", "dark_purple", DarkPurple:
		return DarkPurple, true
	case "6", "gold", Gold:
		return Gold, true
	case "7", "gray", Gray:
		return Gray, true
	case "8", "dark_gray", DarkGray:
		return DarkGray, true
	case "9", "blue", Blue:
		return Blue, true
	case "a", "green", Green:
		return Green, true
	case "b", "aqua", Aqua:
		return Aqua, true
	case "c", "red", Red:
		return Red, true
	case "d", "light_purple", LightPurple:
		return LightPurple, true
	case "e", "yellow", Yellow:
		return Yellow, true
	case "f", "white", White:
		return White, true
	case "g", "minecoin_gold", MinecoinGold:
		return MinecoinGold, true
	default:
		return White, false
	}
}

// ToRaw returns the encoded Minecraft formatting of the color (§ + code)
func (c Color) ToRaw() string {
	return "\u00A7" + string(c)
}

// ToHex returns the hex string of the color prefixed with a # symbol
func (c Color) ToHex() string {
	switch c {
	case Black:
		return "#000000"
	case DarkBlue:
		return "#0000aa"
	case DarkGreen:
		return "#00aa00"
	case DarkAqua:
		return "#00aaaa"
	case DarkRed:
		return "#aa0000"
	case DarkPurple:
		return "#aa00aa"
	case Gold:
		return "#ffaa00"
	case Gray:
		return "#aaaaaa"
	case DarkGray:
		return "#555555"
	case Blue:
		return "#5555ff"
	case Green:
		return "#55ff55"
	case Aqua:
		return "#55ffff"
	case Red:
		return "#ff5555"
	case LightPurple:
		return "#ff55ff"
	case Yellow:
		return "#ffff55"
	case White:
		return "#ffffff"
	case MinecoinGold:
		return "#ddd605"
	default:
		return "#ffffff"
	}
}
