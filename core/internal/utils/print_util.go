package utils

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"
	"sort"
	"strings"

	"github.com/fatih/color"
)

var DefaultPrinter = newPrinter()

type Printer struct {
	Color   *color.Color
	LogChan chan string // 添加日志通道
	enabled bool        // 是否启用通道
}

func newPrinter() *Printer {
	p := &Printer{
		Color:   color.New(color.FgCyan),
		LogChan: make(chan string, 100), // 使用带缓冲的通道，避免阻塞
		enabled: true,
	}

	// go func() {
	// 	for msg := range p.LogChan {
	// 		fmt.Println("=====:" + msg)
	// 	}
	// }()

	return p
}

type LogType struct {
	Name   string
	Color  *color.Color
	Symbol string
}
type LocationType string

var (
	LogTypeInfo LogType = LogType{
		Name:   "info",
		Color:  color.New(color.FgCyan),
		Symbol: "🍋",
	}
	LogTypeError LogType = LogType{
		Name:   "error",
		Color:  color.New(color.FgRed),
		Symbol: "❌",
	}
	LogTypeWarn LogType = LogType{
		Name:   "warn",
		Color:  color.New(color.FgYellow),
		Symbol: "🚨",
	}
	LogTypeSuccess LogType = LogType{
		Name:   "success",
		Color:  color.New(color.FgGreen),
		Symbol: "✅",
	}
	LogTypeDebug LogType = LogType{
		Name:   "debug",
		Color:  color.New(color.FgBlue),
		Symbol: "🔍",
	}
	LogTypeCommand LogType = LogType{
		Name:   "command",
		Color:  color.New(color.FgMagenta),
		Symbol: "📺",
	}
	LogTypeTitle LogType = LogType{
		Name:   "title",
		Color:  color.New(color.FgGreen),
		Symbol: "",
	}
	LogTypeCommandOutput LogType = LogType{
		Name:   "command_output",
		Color:  color.New(color.FgMagenta),
		Symbol: "",
	}
	LogTypeRaw LogType = LogType{
		Name:   "raw",
		Color:  color.New(color.FgWhite),
		Symbol: "",
	}
	LogTypeAlert LogType = LogType{
		Name:   "alert",
		Color:  color.New(color.FgRed),
		Symbol: "🐸",
	}
)

const (
	LocationTypeLong  LocationType = "long"
	LocationTypeShort LocationType = "short"
	LocationTypeNone  LocationType = "none"
)

// Define a structured log message
type LogMessage struct {
	Type    string `json:"type"`    // The log type (info, error, warn, etc.)
	Symbol  string `json:"symbol"`  // The emoji symbol
	Message string `json:"message"` // The actual message
}

// EnableLogging 启用日志通道
func (p *Printer) EnableLogging(enable bool) {
	p.enabled = enable
}

// Close 关闭日志通道
func (p *Printer) Close() {
	if p.LogChan != nil {
		close(p.LogChan)
	}
}

// sendToChannel 发送日志到通道
func (p *Printer) sendToChannel(message string, logType LogType) {
	if p.enabled && p.LogChan != nil {
		logMsg := LogMessage{
			Type:    logType.Name,
			Symbol:  logType.Symbol,
			Message: message,
		}

		jsonMsg, err := json.Marshal(logMsg)
		if err != nil {
			return
		}

		select {
		case p.LogChan <- string(jsonMsg):
			// Successfully sent
		default:
			// Channel is full, discard log
		}
	}
}

// print handles all printing operations
func (p *Printer) print(level LogType, message string, locationType LocationType, sendToChannel bool) {
	var color = level.Color
	var callerInfo string

	_, thisFile, _, _ := runtime.Caller(0)
	skip := 2
	_, callerFile, callerLine, _ := runtime.Caller(skip)

	// 持续往上追溯直到找到非当前文件的调用者
	for callerFile == thisFile {
		skip++
		_, callerFile, callerLine, _ = runtime.Caller(skip)
	}

	// 生成调用者信息
	if locationType == LocationTypeLong {
		callerInfo = fmt.Sprintf("[%s:%d]", callerFile, callerLine)
	} else if locationType == LocationTypeShort {
		shortFile := callerFile[strings.LastIndex(callerFile, "/")+1:]
		callerInfo = fmt.Sprintf("%s\n📃 位置: %s:%d", message, shortFile, callerLine)
	} else {
		callerInfo = ""
	}

	// 如果是打包后的软件，则不打印位置
	if DefaultDevUtil.IsDev() {
		// 在开发环境中运行
	} else {
		// 在打包环境中运行
		callerInfo = ""
	}

	// emoji
	emoji := level.Symbol

	// Generate the complete message with newline
	completeMessage := emoji + " " + callerInfo + " " + message + "\n"

	// Print in a single operation
	color.Print(completeMessage)

	if sendToChannel {
		p.sendToChannel(message, level)
	}
}

// PrintCyan 打印青色信息
func (p *Printer) PrintCyan(format string, args ...interface{}) {
	p.print(LogTypeInfo, fmt.Sprintf(format, args...), LocationTypeNone, true)
}

func (p *Printer) PrintGreen(format string, args ...interface{}) {
	p.print(LogTypeSuccess, fmt.Sprintf(format, args...), LocationTypeNone, true)
}

// PrintRed 打印红色信息
func (p *Printer) PrintRed(format string, args ...interface{}) {
	p.print(LogTypeRaw, fmt.Sprintf(format, args...), LocationTypeNone, true)
}

// PrintWhite 打印白色信息
func (p *Printer) PrintWhite(format string, args ...interface{}) {
	p.print(LogTypeDebug, fmt.Sprintf(format, args...), LocationTypeNone, true)
}

func (p *Printer) PrintYellow(format string, args ...interface{}) {
	p.print(LogTypeWarn, fmt.Sprintf(format, args...), LocationTypeNone, true)
}

func (p *Printer) PrintAndReturnError(errMsg string) error {
	p.print(LogTypeError, fmt.Sprintf("错误: %s", errMsg), LocationTypeNone, true)
	return fmt.Errorf("%s", errMsg)
}

// PrintInfo 打印信息
func (p *Printer) PrintInfo(info string) {
	p.print(LogTypeInfo, info, LocationTypeLong, true)
}

// PrintInfof 打印信息
func (p *Printer) PrintInfof(format string, args ...interface{}) {
	p.PrintInfo(fmt.Sprintf(format, args...))
}

// PrintInfofWithoutLocation 打印信息
func (p *Printer) PrintInfofWithoutLocation(format string, args ...interface{}) {
	p.print(LogTypeInfo, fmt.Sprintf(format, args...), LocationTypeNone, true)
}

// PrintInfofWithoutLocationAndEmoji 打印信息
func (p *Printer) PrintInfofWithoutLocationAndEmoji(format string, args ...interface{}) {
	p.print(LogTypeRaw, fmt.Sprintf(format, args...), LocationTypeNone, true)
}

// PrintTitle 打印标题
func (p *Printer) PrintTitle(title string) {
	p.print(LogTypeTitle, title, LocationTypeNone, true)
}

// PrintLn 打印换行
func (p *Printer) PrintLn() {
	fmt.Println()
}

// Printf 打印格式化信息
func (p *Printer) Printf(format string, args ...interface{}) {
	p.print(LogTypeDebug, fmt.Sprintf(format, args...), LocationTypeNone, true)
}

// PrintError 打印错误信息
func (p *Printer) PrintError(err error) {
	p.PrintErrorMessage(err.Error())
}

// PrintErrorf 打印错误信息
func (p *Printer) PrintErrorf(format string, args ...interface{}) {
	p.PrintErrorMessage(fmt.Sprintf(format, args...))
}

// PrintAndReturnErrorf 打印错误信息并返回错误
func (p *Printer) PrintAndReturnErrorf(format string, args ...interface{}) error {
	p.print(LogTypeError, fmt.Sprintf("错误: %s", fmt.Sprintf(format, args...)), LocationTypeNone, true)
	return fmt.Errorf("%s", fmt.Sprintf(format, args...))
}

// PrintErrorMessage 打印错误信息
func (p *Printer) PrintErrorMessage(message string) {
	_, thisFile, _, _ := runtime.Caller(0) // 获取当前文件路径
	_, file, line, _ := runtime.Caller(1)
	// 如果错误来自当前文件，则往上找一级调用者
	if file == thisFile {
		_, file, line, _ = runtime.Caller(2)
	}

	p.PrintLn()
	p.PrintRed("❌ 错误: %s", strings.TrimSpace(message))
	p.PrintRed("📃 位置: %s:%d", file, line)
	p.PrintLn()

	p.sendToChannel(fmt.Sprintf("错误: %s\n位置: %s:%d", message, file, line), LogTypeError)
}

// PrintList 打印列表
func (p *Printer) PrintList(list []string, title string) {
	p.Color.Println()
	p.Color.Println(title)
	if len(list) == 0 {
		p.Color.Println("  暂无数据")
		return
	}
	for _, item := range list {
		p.Color.Printf("  ▶️  %s\n", item)
	}
	fmt.Println()
}

// PrintSuccess 打印成功信息
func (p *Printer) PrintSuccess(format string) {
	p.print(LogTypeSuccess, format, LocationTypeLong, true)
}

// PrintSuccessf 打印成功信息
func (p *Printer) PrintSuccessf(format string, args ...interface{}) {
	p.print(LogTypeSuccess, fmt.Sprintf(format, args...), LocationTypeNone, true)
}

// PrintWarn 打印警告信息
func (p *Printer) PrintWarn(format string) {
	p.print(LogTypeWarn, format, LocationTypeNone, true)
}

// PrintWarnf 打印警告信息
func (p *Printer) PrintWarnf(format string, args ...interface{}) {
	p.print(LogTypeWarn, fmt.Sprintf(format, args...), LocationTypeNone, true)
}

// PrintAlert 打印提示信息
func (p *Printer) PrintAlert(format string) {
	p.print(LogTypeAlert, format, LocationTypeNone, true)
}

// PrintCommand 打印命令信息
func (p *Printer) PrintCommand(command string) {
	p.print(LogTypeCommand, command, LocationTypeLong, true)
}

// PrintCommandOutput 打印命令输出
func (p *Printer) PrintCommandOutput(output string) {
	p.print(LogTypeCommandOutput, output, LocationTypeNone, true)
}

func (p *Printer) RunShell(command string, args ...string) error {
	if len(args) == 0 {
		return fmt.Errorf("command is required")
	}

	// 去除command开头的空格
	command = strings.TrimSpace(command)

	p.PrintCommand(fmt.Sprintf("%s %s", command, JoinArgs(args)))

	execCmd := exec.Command(command, args...)

	// 创建管道用于捕获输出
	stdoutPipe, err := execCmd.StdoutPipe()
	if err != nil {
		return err
	}
	stderrPipe, err := execCmd.StderrPipe()
	if err != nil {
		return err
	}

	// 启动命令
	if err := execCmd.Start(); err != nil {
		return err
	}

	// 处理标准输出
	go p.processOutput(stdoutPipe, "stdout")

	// 处理标准错误
	go p.processOutput(stderrPipe, "stderr")

	// 等待命令完成
	return execCmd.Wait()
}

// RunShellWithSudo 运行命令并使用 sudo
func (p *Printer) RunShellWithSudo(command string, args ...string) error {
	sudoPrefix := p.GetSudoPrefix()

	if sudoPrefix == "" {
		return p.RunShell(command, args...)
	}

	return p.RunShell(sudoPrefix+command, args...)
}

// processOutput 处理输出流
func (p *Printer) processOutput(pipe io.ReadCloser, source string) {
	scanner := bufio.NewScanner(pipe)
	for scanner.Scan() {
		line := scanner.Text()
		// 打印到控制台并发送到日志通道
		fmt.Println(line)
		p.sendToChannel(line, LogTypeCommandOutput)
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("读取%s错误: %v\n", source, err)
	}
}

// RunShellWithOutput 运行命令并返回输出
func (p *Printer) RunShellWithOutput(command string, args ...string) (string, error) {
	if len(args) == 0 {
		return "", fmt.Errorf("command is required")
	}

	p.PrintCommand(fmt.Sprintf("%s %s", command, JoinArgs(args)))

	execCmd := exec.Command(command, args...)

	output, err := execCmd.CombinedOutput()

	return string(output), err
}

// RunShellInFolder 在指定文件夹中运行命令
func (p *Printer) RunShellInFolder(folder string, command string, args ...string) error {
	p.PrintInfo(fmt.Sprintf("切换到文件夹: %s", folder))

	// 切换到指定文件夹
	err := os.Chdir(folder)
	if err != nil {
		return err
	}

	return p.RunShell(command, args...)
}

// GetSudoPrefix 获取 sudo 前缀
func (p *Printer) GetSudoPrefix() string {
	if os.Geteuid() != 0 {
		return "sudo "
	}
	return ""
}

// PrintKeyValue 打印键值对
func (p *Printer) PrintKeyValue(key string, value string) {
	p.print(LogTypeDebug, fmt.Sprintf("%s: %s", key, value), LocationTypeNone, true)
}

// PrintKeyValues 打印键值对列表
func (p *Printer) PrintKeyValues(keyValues map[string]string) {
	// 找出最长的 key 长度
	maxKeyLength := 0
	for key := range keyValues {
		if len(key) > maxKeyLength {
			maxKeyLength = len(key)
		}
	}

	// 按照 key 排序，保证输出顺序一致
	keys := make([]string, 0, len(keyValues))
	for key := range keyValues {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	// 输出对齐的键值对
	for _, key := range keys {
		// 使用空格补充 key 到最大长度
		paddedKey := fmt.Sprintf("%-*s", maxKeyLength, key)
		p.PrintKeyValue(paddedKey, keyValues[key])
	}
}

// PrintEmojiForBool 打印布尔值的emoji
func (p *Printer) PrintEmojiForBool(value bool) {
	if value {
		p.PrintSuccess("✅")
	} else {
		p.PrintError(fmt.Errorf("❌"))
	}
}
