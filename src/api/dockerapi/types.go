package dockerapi

import (
	"archive/tar"
	"bytes"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type Config struct {
	Image   string
	Command []string
	// k: 宿主端口, v: 容器端口
	Ports   map[string]string
	CPU     int64
	Memory  int64
	GPU     int
	Volumes []Volume
	EnvVars []string

	// k: 文件路径, v: 文件内容
	ConfigData map[string]string
}

func (c Config) HasConfigData(s string) bool {
	_, ok := c.ConfigData[s]
	return ok
}

func (c Config) Tar(name, workspace string) (*bytes.Buffer, error) {
	err := c.saveConfigToLocal(name, workspace)
	if err != nil {
		return nil, err
	}

	return c.tarDirectory(name, workspace)
}

func (c Config) saveConfigToLocal(name, workspace string) (err error) {
	dir := filepath.Join(workspace, name)
	err = os.MkdirAll(dir, 0755)
	if err != nil {
		err = errors.Wrap(err, "os.MkdirAll")
		return
	}

	for k, v := range c.ConfigData {
		var f *os.File
		dataFilePath := filepath.Join(dir, k)
		f, err = os.OpenFile(dataFilePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0777)
		if err != nil {
			err = fmt.Errorf("%s openFile err: %s", dataFilePath, err.Error())
			return
		}
		_, err = f.Write([]byte(v))
		if err != nil {
			err = fmt.Errorf("%s write err: %s", dataFilePath, err.Error())
			return
		}

		f.Close()
	}

	return
}

func (c Config) tarDirectory(name, workspace string) (*bytes.Buffer, error) {
	source := filepath.Join(workspace, name)
	buffer := new(bytes.Buffer)
	tw := tar.NewWriter(buffer)
	defer tw.Close()

	// 递归地添加目录中的文件和子目录到 tar 归档
	err := filepath.Walk(source, func(file string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 创建 tar 头部
		header, err := tar.FileInfoHeader(fi, file)
		if err != nil {
			return err
		}

		// 必要时调整头部信息
		header.Name = strings.TrimPrefix(file, source)

		// 写入头部信息
		if err := tw.WriteHeader(header); err != nil {
			return err
		}

		// 如果不是普通文件则不继续
		if !fi.Mode().IsRegular() {
			return nil
		}

		// 打开文件
		f, err := os.Open(file)
		if err != nil {
			return err
		}
		defer f.Close()

		// 将文件内容复制到 tar 归档
		if _, err := io.Copy(tw, f); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// 确保 tar 归档写入完成
	if err := tw.Close(); err != nil {
		return nil, err
	}

	return buffer, nil
}

type Volume struct {
	// 宿主路径
	Key string
	// 容器路径
	Value string
}
