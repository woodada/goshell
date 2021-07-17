package goshell

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

func Cd(dirPath string) error {
	return os.Chdir(dirPath)
}

func MustCd(dirPath string) {
	err := Cd(dirPath)
	assert("cd "+dirPath, err)
}

func CurrentDir() (string, error) {
	return os.Getwd()
}

func MustCurrentDir() string {
	dir, err := CurrentDir()
	assert("GetWorkDir", err)
	return dir
}

type Result struct {
	Stdout   []byte
	Stderr   []byte
	ExitCode int
	Success  bool
	// Error    error
}

func (me Result) String() string {
	return fmt.Sprint(me.Success, me.ExitCode, string(me.Stdout), string(me.Stderr))
}

func buildResult(rst *Result, stdout, stderr *bytes.Buffer, state *os.ProcessState) *Result {
	if rst == nil {
		return rst
	}
	if stdout != nil {
		rst.Stdout = stdout.Bytes()
	}
	if stderr != nil {
		rst.Stderr = stderr.Bytes()
	}
	if state != nil {
		rst.ExitCode = state.ExitCode()
		rst.Success = state.Success()
	}
	return rst
}

func Run(arg0 string, argN ...string) (*Result, error) {
	fullPath, err := ResolveArgs0(arg0)
	if err != nil {
		return nil, err
	}
	args := []string{fullPath}
	args = append(args, argN...)
	return _run(args...)
}

func MustRun(arg0 string, argN ...string) string {
	rst, err := Run(arg0, argN...)
	assert(fmt.Sprint(arg0, argN), err)
	if !rst.Success || rst.ExitCode != 0 {
		fmt.Fprintln(os.Stderr, rst.String())
		os.Exit(1)
	}
	return strings.TrimSpace(string(rst.Stdout))
}

func _run(args ...string) (*Result, error) {
	if len(args) <= 0 {
		panic("_runSync:ArgsEmpty")
	}
	cmd := exec.Cmd{
		Path: args[0],
		Args: args,
		Env:  os.Environ(),
	}

	stdin := &bytes.Buffer{}
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	cmd.Stdin = stdin
	if _config.DisableMuti {
		cmd.Stderr = NewMutiWriter(stderr)
		cmd.Stdout = NewMutiWriter(stdout)
	} else {
		cmd.Stderr = NewMutiWriter(os.Stderr, stderr)
		cmd.Stdout = NewMutiWriter(os.Stdout, stdout)
	}

	rst := &Result{}

	fmt.Println(args)
	err := cmd.Start()
	if err != nil {
		return buildResult(rst, stdout, stderr, cmd.ProcessState), err
	}

	err = cmd.Wait()
	if err != nil {
		return buildResult(rst, stdout, stderr, cmd.ProcessState), err
	}

	return buildResult(rst, stdout, stderr, cmd.ProcessState), nil
}

// 按平台实现
// windows下where; mac和linux下用which
func Which(name string) (string, error) {
	if rawpath, ok := _whichMap[name]; ok {
		return rawpath, nil
	}
	arg0 := "which"
	if "windows" == runtime.GOOS {
		arg0 = "where"
	}
	rst, err := _run(arg0, name)
	if err != nil {
		return "", fmt.Errorf("exitcode:%d#err:%v#stderr:%s", rst.ExitCode, err, string(rst.Stderr))
	}
	if rst.ExitCode == 0 {
		arr := strings.Split(strings.TrimSpace(string(rst.Stdout)), "\n")
		if len(arr) <= 0 {
			return "", fmt.Errorf("%s not found", name)
		}
		full := strings.TrimSpace(arr[0])
		SetWhich(name, full)
		return full, nil
	}
	return "", fmt.Errorf("exitcode:%d#err:%v#stderr:%s", rst.ExitCode, err, string(rst.Stderr))
}

func MustWhich(name string) string {
	dir, err := Which(name)
	assert("which "+name, err)
	return dir
}

// 找全路径：返回全路径，或者错误提示
func ResolveArgs0(arg0 string) (string, error) {
	// 全路径
	if filepath.IsAbs(arg0) {
		return arg0, nil
	}
	// 当前路径
	if strings.HasPrefix(arg0, "./") || strings.HasPrefix(arg0, ".\\") {
		return filepath.Abs(arg0)
	}
	// 纯命令
	return Which(arg0)
}

func assert(name string, err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, "cmd:", name, "err：", err)
		os.Exit(1)
	}
}

type MutiWriter struct {
	writers []io.Writer
}

func NewMutiWriter(writes ...io.Writer) *MutiWriter {
	return &MutiWriter{writers: writes}
}

func (this *MutiWriter) Write(data []byte) (int, error) {
	for _, w := range this.writers {
		w.Write(data)
	}
	return len(data), nil
}

var _whichMap = make(map[string]string)

// 设置which的路径，which查找的时候优先用这个
func SetWhich(name, rawpath string) {
	_whichMap[name] = rawpath
}
