/**
 * @Author:      leafney
 * @GitHub:      https://github.com/leafney
 * @Project:     whisky
 * @Date:        2024-07-06 11:27
 * @Description:
 */

package utils

import "os/exec"

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
	cmd := exec.Command("/bin/sh", "-c", shellStr)
	cmd.Args = append(cmd.Args, args...)
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(output), nil
}
