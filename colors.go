package main

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
)

func colorFuncFromStr(s string) (colorFunc, error) {
	switch strings.ToLower(s) {
	case "black":
		return color.BlackString, nil
	case "red":
		return color.RedString, nil
	case "green":
		return color.GreenString, nil
	case "yellow":
		return color.YellowString, nil
	case "blue":
		return color.BlueString, nil
	case "magenta":
		return color.MagentaString, nil
	case "cyan":
		return color.CyanString, nil
	case "white":
		return color.WhiteString, nil
	case "hiblack":
		return color.HiBlackString, nil
	case "hired":
		return color.HiRedString, nil
	case "higreen":
		return color.HiGreenString, nil
	case "hiyellow":
		return color.HiYellowString, nil
	case "hiblue":
		return color.HiBlueString, nil
	case "himagenta":
		return color.HiMagentaString, nil
	case "hicyan":
		return color.HiCyanString, nil
	case "hiwhite":
		return color.HiWhiteString, nil
	default:
		return nil, fmt.Errorf("invalid color %q", s)
	}
}
