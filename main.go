package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"gopkg.in/ini.v1"
)

var (
	CurrentDir, _      = os.Getwd() // 当前所在目录
	dirPassword        string
	filePassword       string
	dirComSuffix       string
	fileComSuffix      string
	dirDeleteOriginal  bool // 是否删除原文件夹
	fileDeleteOriginal bool // 是否删除原文件
	excludeFileSuffix  string
)

func init() {
	cfg, err := ini.Load(filepath.Join(CurrentDir, "config.ini"))
	if err != nil {
		log.Printf("初始化失败: %v", err)
		return
	}
	dirPassword = cfg.Section("compression").Key("dirPassword").String()
	filePassword = cfg.Section("compression").Key("filePassword").String()
	dirComSuffix = cfg.Section("compression").Key("dirComSuffix").String()
	fileComSuffix = cfg.Section("compression").Key("fileComSuffix").String()
	dirDeleteOriginal, _ = cfg.Section("compression").Key("dirDeleteOriginal").Bool()
	fileDeleteOriginal, _ = cfg.Section("compression").Key("fileDeleteOriginal").Bool()
	excludeFileSuffix = cfg.Section("compression").Key("excludeFileSuffix").String()

}
func main() {
	fmt.Print("此程序可以压缩当前文件夹内的所有文件夹或文件\n类型：1.压缩文件夹，2.压缩文件\n请输入压缩类型：")
	var compressionType int
	fmt.Scanln(&compressionType)
	if compressionType != 1 && compressionType != 2 {
		fmt.Println("压缩类型输入错误")
		return
	}
	if compressionType == 1 {
		compressionDir()
		return
	}
	if compressionType == 2 {
		compressionFile()
		return
	}
}

func compressionDir() {
	// 获取当前目录下的所有子文件夹
	subFolders, err := os.ReadDir(CurrentDir)
	if err != nil {
		fmt.Println("读取当前目录失败:", err)
		return
	}
	// 遍历子文件夹，并进行压缩
	for _, subFolder := range subFolders {
		if subFolder.IsDir() {
			// 子文件夹的路径
			subFolderPath := filepath.Join(CurrentDir, subFolder.Name())
			// 压缩文件保存的路径和文件名
			zipFilePath := fmt.Sprintf("%s."+dirComSuffix, subFolderPath)
			// 构造7z命令行参数
			args := []string{"a", "-tzip", "-p" + dirPassword, zipFilePath, subFolderPath}
			// 创建一个Cmd对象
			cmd := exec.Command("7z", args...)
			// 设置输出和错误输出
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			// 执行命令
			err := cmd.Run()
			if err != nil {
				fmt.Printf("压缩文件夹 %s 失败: %s\n", subFolder.Name(), err)
				continue
			}
			fmt.Printf("文件夹 %s 已成功压缩并设置密码：%s\n", subFolder.Name(), zipFilePath)
			// 根据开关设置是否删除原文件夹
			if dirDeleteOriginal {
				err = os.RemoveAll(subFolderPath)
				if err != nil {
					fmt.Printf("删除文件夹 %s 失败: %s\n", subFolder.Name(), err)
					continue
				}
				fmt.Printf("文件夹 %s 已成功删除\n", subFolder.Name())
			}
		}
	}
	main()
}
func compressionFile() {
	// 获取待压缩文件夹下的所有文件和子文件夹
	files, err := os.ReadDir(CurrentDir)
	if err != nil {
		fmt.Println("读取文件夹失败:", err)
		return
	}
	excludeFileSuffixList := strings.Split(excludeFileSuffix, ",")
	// 遍历文件和文件夹，并进行压缩
	for _, file := range files {
		filePath := filepath.Join(CurrentDir, file.Name())
		if !file.IsDir() && !contains(excludeFileSuffixList, filepath.Ext(filePath)) {
			// 压缩文件保存的路径和文件名
			zipFilePath := fmt.Sprintf("%s."+fileComSuffix, strings.TrimSuffix(filePath, filepath.Ext(filePath)))
			// 构造7z命令行参数
			args := []string{"a", "-tzip", "-p" + filePassword, zipFilePath, filePath}
			// 创建一个Cmd对象
			cmd := exec.Command("7z", args...)
			// 设置输出和错误输出
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			// 执行命令
			err := cmd.Run()
			if err != nil {
				fmt.Printf("压缩文件 %s 失败: %s\n", file.Name(), err)
				continue
			}
			fmt.Printf("文件 %s 已成功压缩并设置密码：%s\n", file.Name(), zipFilePath)
			if fileDeleteOriginal {
				// 根据开关状态，决定是否删除原文件
				err = os.Remove(filePath)
				if err != nil {
					fmt.Printf("删除原文件 %s 失败: %s\n", file.Name(), err)
				} else {
					fmt.Printf("原文件 %s 已删除\n", file.Name())
				}
			}
		}
	}
	main()
}

// 辅助函数，用于判断切片中是否包含某个值
func contains(slice []string, value string) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}
