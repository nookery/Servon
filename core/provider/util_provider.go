package provider

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// UtilProvider 工具提供者
type UtilProvider struct {
	RootCmd *cobra.Command
}

func NewUtilProvider() UtilProvider {
	return UtilProvider{
		RootCmd: &cobra.Command{},
	}
}

// PrintCommandErrorAndExit 打印命令错误并退出，用于打印命令错误信息
func (p *UtilProvider) PrintCommandErrorAndExit(err error) error {
	errorStyle := color.New(color.FgRed, color.Bold)
	errorStyle.Print("❌ Error: ")
	errorStyle.Println(err)
	os.Exit(1)
	return err
}

// PrintCommandSuccess 打印命令成功，用于打印命令成功信息
func (p *UtilProvider) PrintCommandSuccess(msg string) {
	color.New(color.FgGreen).Println(msg)
}

// PrintStep 打印步骤，用于打印步骤信息
func (p *UtilProvider) PrintStep(msg string) {
	color.New(color.FgBlue).Println("🔍 " + msg)
}

// PrintStepSuccess 打印步骤成功，用于打印步骤成功信息
func (p *UtilProvider) PrintStepSuccess(msg string) {
	color.New(color.FgGreen).Println("🎉 " + msg)
}

// PrintStepFinish 打印步骤完成，用于打印步骤完成信息
func (p *UtilProvider) PrintStepFinish(msg string) {
	color.New(color.FgGreen).Println("✅ " + msg)
}

// PrintStepError 打印步骤错误，用于打印步骤错误信息
func (p *UtilProvider) PrintStepError(msg string) {
	color.New(color.FgRed).Println("❌ " + msg)
}

// PrintSuccess 打印成功信息
func (p *UtilProvider) PrintSuccess(msg string) {
	color.New(color.FgGreen).Println("✅ " + msg)
}

// PrintError 打印错误信息
func (p *UtilProvider) PrintError(msg string) {
	color.New(color.FgRed).Println("❌ " + msg)
}

// PrintAndReturnError 打印错误信息并返回错误
func (p *UtilProvider) PrintAndReturnError(errMsg string) error {
	p.PrintError(errMsg)
	return fmt.Errorf("❌ %s", errMsg)
}
