package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
)

type Args struct {
	programName string
	startPage   int
	endPage     int
	srcFile     string
	pageLength  int
	pageType    bool // true： '\f'
	desProgram  string
}

func main() {
	var myArgs Args

	setUsage(&myArgs)
	argsProcess(&myArgs)
	fileProcess(&myArgs)
}

// 对参数进行逻辑处理
func argsProcess(myArgs *Args) {
	var proName = myArgs.programName
	// 必须输入起始页和终止页
	if myArgs.startPage < 0 || myArgs.endPage < 0 {
		processError(proName, fmt.Sprintf("not enough arguments, -s and -e is indispensable."))
	}
	// 起始页要小于等于终止页
	if myArgs.startPage > myArgs.endPage {
		processError(proName, fmt.Sprintf("STARTPAGE(%d) should not be greater than ENDPAGE(%d).", myArgs.startPage, myArgs.endPage))
	}

	myArgs.srcFile = flag.Arg(0)

	var defaulLength int = 72
	if myArgs.pageType == true {
		// 保证-f和-l的互斥
		if myArgs.pageLength != -1 {
			processError(proName, "-f and -l=PAGELENGTH are mutually-exclusive.")
		}
	} else {
		// 为每一页的长度赋缺省值
		if myArgs.pageLength < 1 {
			myArgs.pageLength = defaulLength
		}
	}

}

// 文件处理
func fileProcess(myArgs *Args) {
	// 若用户规定了输入文件
	if flag.NArg() == 1 {
		myArgs.srcFile = flag.Arg(0)
	}

	if myArgs.srcFile == "" {
		// 若为空字符串，表示输入为标准输入
		inputReader := bufio.NewReader(os.Stdin)

		if myArgs.pageType == true {
			readByPage(inputReader, myArgs)
		} else {
			readByLine(inputReader, myArgs)
		}
	} else {
		inputFile, err := os.Open(myArgs.srcFile)
		check(err)
		inputReader := bufio.NewReader(inputFile)
		defer inputFile.Close()

		if myArgs.pageType == true {
			readByPage(inputReader, myArgs)
		} else {
			readByLine(inputReader, myArgs)
		}
	}
}

func readByPage(inputReader *bufio.Reader, myArgs *Args) {
	// 记录当前页数
	pageCount := 1
	// 读取所有页
	for {
		page, err := inputReader.ReadString('\f')
		check(err)
		// 当页数在所要选取的范围时
		if pageCount >= myArgs.startPage && pageCount <= myArgs.endPage {
			// 若输出为标准输出
			if myArgs.desProgram == "" {
				fmt.Printf(page)
			} else {
				// 打开./go的输入管道，将该程序输出传输到管道
				cmd := exec.Command("./out")          // 创建命令"./out"
				echoInPipe, err := cmd.StdinPipe()    // 打开./out的标准输入管道
				check(err)                            // 错误检测
				echoInPipe.Write([]byte(page + "\n")) // 向管道中写入文本
				echoInPipe.Close()                    // 关闭管道
				cmd.Stdout = os.Stdout                // ./out将会输出到屏幕
				cmd.Run()                             // 运行./out命令
			}
		}
		if err == io.EOF {
			break
		}
		pageCount++
	}
	// 当起始页大于总页数输出为空
	if myArgs.startPage > pageCount {
		fmt.Printf("Warning:\n\tSTARTPAGE(%d) is greater than number of total pages(%d).\noutput will be empty.\n", myArgs.startPage, pageCount)
	}
	// 当终止页大于总页数
	if myArgs.endPage > pageCount {
		fmt.Printf("Warning:\n\tENDPAGE(%d) is greater than number of total pages(%d).\nthere will be less output than expected.\n", myArgs.endPage, pageCount)
	}
}

func readByLine(inputReader *bufio.Reader, myArgs *Args) {
	lineCount := 1
	for {
		line, err := inputReader.ReadString('\n')
		check(err)

		if lineCount > myArgs.pageLength*(myArgs.startPage-1) && lineCount <= myArgs.pageLength*myArgs.endPage {
			if myArgs.desProgram == "" {
				fmt.Printf(line)
			} else {
				cmd := exec.Command("./out")
				echoInPipe, err := cmd.StdinPipe()
				check(err)
				echoInPipe.Write([]byte(line))
				echoInPipe.Close()
				cmd.Stdout = os.Stdout
				cmd.Run()
			}
		}
		if err == io.EOF {
			break
		}
		lineCount++
	}
	if myArgs.startPage > lineCount/myArgs.pageLength+1 {
		fmt.Printf("Warning:\n\tSTARTPAGE(%d) is greater than number of total pages(%d).\noutput will be empty.\n", myArgs.startPage, lineCount/myArgs.pageLength+1)
	}
	if myArgs.endPage > lineCount/myArgs.pageLength+1 {
		fmt.Printf("Warning:\n\tENDPAGE(%d) is greater than number of total pages(%d).\nthere will be less output than expected.\n", myArgs.endPage, lineCount/myArgs.pageLength+1)
	}
}

func check(err error) {
	if err != nil && err != io.EOF {
		panic(err)
	}
}

func processError(name string, errorStr string) {
	fmt.Fprintf(os.Stderr, "%s: %s\n", name, errorStr)
	flag.Usage()
	os.Exit(1)
}
