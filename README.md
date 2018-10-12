#CLI命令行使用程序开发基础
##设计说明
参考了以下博客  
[https://www.ibm.com/developerworks/cn/linux/shell/clutil/index.html](https://www.ibm.com/developerworks/cn/linux/shell/clutil/index.html)   
引用如下包   
```
import (
    "bufio"
    "flag"
    "fmt"
    "io"
    "os"
    "os/exec"
    "github.com/spf13/pflag"
)
```
结构体   
```
type selpg_args struct {
    start_page  int
    end_page    int
    in_filename []string
    page_len    int
    page_type   byte
    print_dest  string
}
```
函数   
```
func main() //程序入口
func parseArgs(args *selpg_args) //用pflag解析命令
func processArgs(args selpg_args) //检查参数
func processInput(args selpg_args) //处理参数和输出
```
##Usage
```
[hechx6@centos-manager bin]$ selpg -h
Usage of selpg:
  -e, --end int            End page number (default 1)
      --f                  Page type
  -l, --pagelength int     Line number of each page (default 72)
  -d, --printdest string   Output to destination pipe
  -s, --start int          Start page number (default 1)
pflag: help requested
```
example
```
selpg -s 10 -e 20 test.txt
selpg -s 10 -e 20 -l 10 test.txt
```

##测试
测试文件：   
test.txt   
```
line 1
line 2
line 3
line 4
line 5
line 6
line 7
line 8
line 9
line 10
line 11
line 12
line 13
line 14
line 15
line 16
line 17
line 18
line 19
line 20
line 21
line 22
line 23
line 24
line 25
line 26
line 27
line 28
line 29
line 30
line 31
line 32
line 33
line 34
line 35
line 36
```
1.从标准输入选定特定页写至标准输出(即屏幕)   
```
[hechx6@centos-manager bin]$ selpg -s 1 -e 1
aabbc
aabbc
bbccd
bbccd
efegdg
efegdg
erg
erg
^Z
[3]+  Stopped                 selpg -s 1 -e 1
```
2.命令将把“input_file”特定页写至标准输出（也就是屏幕）  
```
[hechx6@centos-manager bin]$ selpg -s 1 -e 1 -l 2 test.txt
line 1
line 2
``` 
3.重定向输出
```
[hechx6@centos-manager bin]$ selpg -s 1 -e 3 -l 5 test.txt > test
[hechx6@centos-manager bin]$ cat test
line 1
line 2
line 3
line 4
line 5
line 6
line 7
line 8
line 9
line 10
line 11
line 12
line 13
line 14
line 15
```
4.当页码超出范围错误提示
```
[hechx6@centos-manager bin]$ selpg -s 3 -e 50 -l 2 test.txt
line 5
line 6
line 7
line 8
line 9
line 10
line 11
line 12
line 13
line 14
line 15
line 16
line 17
line 18
line 19
line 20
line 21
line 22
line 23
line 24
line 25
line 26
line 27
line 28
line 29
line 30
line 31
line 32
line 33
line 34
line 35
line 36
End page (50) is greater than total pages (18)
```
