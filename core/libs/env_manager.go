package libs

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type EnvManager struct {
	envExample []byte
	env        []byte
}

func NewEnvManager() *EnvManager {
	return &EnvManager{}
}

func (e *EnvManager) LoadEnvExample() {
	envExample, err := os.ReadFile(".env.example")
	if err != nil {
		log.Fatal("Error reading .env.example file")
	}
	e.envExample = envExample
}

func (e *EnvManager) LoadEnv() {
	env, err := os.ReadFile(".env")
	if err != nil {
		log.Fatal("Error reading .env file")
	}
	e.env = env

	godotenv.Load()
}

func (e *EnvManager) CompareEnv() error {
	// 读取并解析两个文件的键
	envKeys, err := e.getEnvKeys(".env")
	if err != nil {
		return fmt.Errorf("Error reading .env file: %v", err)
	}

	exampleKeys, err := e.getEnvKeys(".env.example")
	if err != nil {
		return fmt.Errorf("Error reading .env.example file: %v", err)
	}

	// 检查不一致
	var missingKeys []string
	var extraKeys []string

	// 检查 .env 中是否缺少 .env.example 中的键
	for key := range exampleKeys {
		if !envKeys[key] {
			missingKeys = append(missingKeys, key)
		}
	}

	// 检查 .env 中是否有多余的键
	for key := range envKeys {
		if !exampleKeys[key] {
			extraKeys = append(extraKeys, key)
		}
	}

	// 如果有任何不一致，返回错误
	if len(missingKeys) > 0 || len(extraKeys) > 0 {
		var errMsg strings.Builder
		if len(missingKeys) > 0 {
			errMsg.WriteString(fmt.Sprintf(".env 文件缺少以下键: %v\n", missingKeys))
		}
		if len(extraKeys) > 0 {
			errMsg.WriteString(fmt.Sprintf(".env 文件包含多余的键: %v\n", extraKeys))
		}
		return fmt.Errorf(errMsg.String())
	}

	return nil
}

// getEnvKeys 读取env文件中的所有键
func (e *EnvManager) getEnvKeys(filename string) (map[string]bool, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	keys := make(map[string]bool)
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		// 跳过空行和注释
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// 分割键值对
		parts := strings.SplitN(line, "=", 2)
		if len(parts) >= 1 {
			key := strings.TrimSpace(parts[0])
			keys[key] = true
		}
	}

	return keys, scanner.Err()
}

func (e *EnvManager) GetEnv(key string) string {
	return os.Getenv(key)
}
