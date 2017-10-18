# 第六周作业:selgo


这次作业要求用刚配置好的golang来完成，开发Linux命令行实用程序中的selpg，总体而言逻辑上不是很难，只不过需要阅读的文档和资料都挺多的，其中也对管道什么的纠结了一会。主要参考文献:
[开发Linux命令行实用程序][1]


  [1]: https://www.ibm.com/developerworks/cn/linux/shell/clutil/index.html
# 具体实现内容
**演示如何用Go语言编写与 cat、ls、pr 和 mv 等标准命令类似的 Linux 命令行实用程序。**我选择了一个名为 selpg 的实用程序，这个名称代表 SELect PaGes。selpg允许用户指定从输入文本抽取的页的范围，这些输入文本可以来自文件/标准输入/另一个进程。页的范围由起始页和终止页决定。**在管道中，输入输出和错误流重定向的情况下也可使用该工具**。

# 代码具体使用过程

    $ selpg -s1 -e1 read_in.txt

该命令将把“`input_file`”的第 1 页写至标准输出（也就是屏幕），因为这里没有重定向或管道。

    $ selpg -s1 -e1 < read_in.txt

该命令与示例 1 所做的工作相同，但在本例中，selpg 读取标准输入，而标准输入已被 shell／内核重定向为来自“`read_in.txt`”而不是显式命名的文件名参数。输入的第 1 页被写至屏幕。

    $ ./in | ./selpg | selpg -s10 -e50

将第 10 页到第 50 页写至 selpg 的标准输出（屏幕）。

    $ selpg -s10 -e50 read_in.txt >read_out.txt

selpg 将第 10 页到第 50 页写至标准输出；标准输出被 shell／内核重定向至“output_file”。
1

    $ selpg -s10 -e50 read_in.txt 2>read_error.txt

selpg 将第 10 页到第 50 页写至标准输出（屏幕）；所有的错误消息被 shell／内核重定向至“`read_error.txt`”。请注意：在“2”和“>”之间不能有空格；这是 shell 语法的一部分（请参阅“man bash”或“man sh”）。

    $ selpg -s10 -e50 read_in.txt >read_out.txt 2>read_error.txt

selpg 将第 10 页到第 50 页写至标准输出，标准输出被重定向至“`read_out.txt`”；selpg 写至标准错误的所有内容都被重定向至“`read_error.txt`”。当“`read_in.txt`”很大时可使用这种调用；您不会想坐在那里等着 selpg 完成工作，并且您希望对输出和错误都进行保存。

    $ selpg -s10 -e20 read_in.txt >read_out.txt 2>/dev/null

selpg 将第 10 页到第 50 页写至标准输出，标准输出被重定向至“`read_out.txt`”；selpg 写至标准错误的所有内容都被重定向至 /dev/null（空设备），这意味着错误消息被丢弃了。设备文件 /dev/null 废弃所有写至它的输出，当从该设备文件读取时，会立即返回 EOF。

    $ selpg -s10 -e50  read_in.txt>/dev/null

selpg 将第 10 页到第 50 页写至标准输出，标准输出被丢弃；错误消息在屏幕出现。这可作为测试 selpg 的用途，此时您也许只想（对一些测试情况）检查错误消息，而不想看到正常输出。

    $ selpg -s10 -e50  read_in.txt|./read_out

selpg 的标准输出透明地被 shell／内核重定向，成为“`./read_out`”的标准输入，第 10 页到第 50 页被写至该标准输入。“`./read_out`”的示例可以是 lp，它使输出在系统缺省打印机上打印。“`./read_out`”的示例也可以 wc，它会显示选定范围的页中包含的行数、字数和字符数。“`./read_out`”可以是任何其它能从其标准输入读取的命令。错误消息仍在屏幕显示。