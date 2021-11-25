package util

import (
	"fmt"

	"github.com/gookit/color"
)

// CliError 输出错误
func CliError(msg string) {
	fmt.Println(color.FgLightRed.Render(msg))
}

// CliInfo 输出信息
func CliInfo(msg string) {
	fmt.Println(color.FgLightWhite.Render(msg))
}
