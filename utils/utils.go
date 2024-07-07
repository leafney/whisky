/**
 * @Author:      leafney
 * @GitHub:      https://github.com/leafney
 * @Project:     whisky
 * @Date:        2024-07-06 11:27
 * @Description:
 */

package utils

import (
	"bytes"
	"fmt"
	"github.com/klauspost/compress/gzip"
	"github.com/leafney/rose"
	"github.com/leafney/whisky/global/vars"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func RunScript(script string, args ...string) (string, error) {
	cmd := exec.Command("/bin/sh", script)
	cmd.Args = append(cmd.Args, args...)
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(output), nil
}

func RunBash(shellStr string, args ...string) (string, error) {
	// 这种方式经测试后发现无法正确传入参数
	//cmd := exec.Command("/bin/sh", "-c", shellStr)
	//cmd.Args = append(cmd.Args, args...)

	// 后来改用这种方式，经测试，这种方式也不对
	//command := fmt.Sprintf("%s %s", shellStr, strings.Join(args, " "))
	//cmd := exec.Command("/bin/sh", "-c", command)

	// 第三种方式
	//cmd := exec.Command("/bin/sh", "-c", shellStr)
	//cmd.Stdin = strings.NewReader(strings.Join(args, " "))

	// 第四种方法
	args = append([]string{"-c", shellStr}, args...)
	cmd := exec.Command("/bin/sh", args...)

	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(output), nil
}

func RunBashFile(shellFilePath string, args ...string) (string, error) {
	command := fmt.Sprintf("%s %s", shellFilePath, strings.Join(args, " "))
	cmd := exec.Command("/bin/sh", "-c", command)
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(output), nil
}

// LoadByteBashFile 载入 shell 脚本文件
func LoadByteBashFile(shellName []byte) (string, error) {
	// 先判断本地是否存在该文件
	fileName := fmt.Sprintf("%s.sky", rose.Md5HashBuf(shellName))

	filePath := filepath.Join(vars.ShellTempDir, fileName)

	// 如果文件存在，则直接返回，否则通过解压后得到
	if rose.FIsExist(filePath) {
		return filePath, nil
	}

	// 创建文件保存目录
	if err := rose.DEnsurePathExist(filePath); err != nil {
		return "", err
	}

	// 解压缩
	reader, err := gzip.NewReader(bytes.NewReader(shellName))
	if err != nil {
		return "", err
	}
	defer reader.Close()

	//	将解压缩后的内容写入本地文件
	targetFile, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer targetFile.Close()

	if _, err := io.Copy(targetFile, reader); err != nil {
		return "", err
	}

	// 为脚本文件添加执行权限
	if err := rose.ExecCmdRun("chmod", "+x", filePath); err != nil {
		return "", err
	}

	// 返回文件路径
	return filePath, nil
}

func DeleteFilesByExtension(dirPath string, ext string) error {
	// 读取目录下的所有文件和文件夹
	files, err := os.ReadDir(dirPath)
	if err != nil {
		return fmt.Errorf("读取目录失败: %w", err)
	}

	// 遍历所有文件
	for _, file := range files {
		// 获取文件路径
		filePath := filepath.Join(dirPath, file.Name())

		// 检查文件是否为指定扩展名
		if filepath.Ext(filePath) == ext {
			// 删除文件
			err := os.Remove(filePath)
			if err != nil {
				return fmt.Errorf("删除文件 %s 失败: %w", filePath, err)
			}
			fmt.Printf("已删除文件: %s\n", filePath)
		}
	}

	return nil
}
