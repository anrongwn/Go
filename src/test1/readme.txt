https://blog.csdn.net/yxf771hotmail/article/details/88233857

https://www.cnblogs.com/zhangweizhong/category/1275863.html
https://studygolang.com/pkgdoc
https://books.studygolang.com/The-Golang-Standard-Library-by-Example/chapter01/01.1.html

//Go ϵ�н̳̣�Golang tutorial series��
https://studygolang.com/subject/2



https://studygolang.com/articles/12315

https://www.jb51.net/article/126998.htm

https://www.cnblogs.com/justdoyou/p/9853520.html


//go tuntap
https://github.com/songgao/water

cd /d/myProgramming/Go/src/golang.org/x
git clone https://github.com/golang/sys

gogm 


//手动安装go相关插件
3.安装失败插件处理
在你的$GOPATH/src/golang.org/x目录下从github上下载插件源码,如果没有该文件夹则自己创建即可.执行命令如下:

git clone https://github.com/golang/tools.git tools
git clone https://github.com/golang/lint.git lint
go install golang.org/x/lint/golint

//更新工具
go install github.com/ramya-rao-a/go-outline
go install github.com/acroca/go-symbols
go install golang.org/x/tools/cmd/guru
go install golang.org/x/tools/cmd/gorename
go install github.com/josharian/impl
go install github.com/rogpeppe/godef
go install github.com/sqs/goreturns
go install github.com/golang/lint/golint
go install github.com/cweill/gotests/gotests
go install github.com/ramya-rao-a/go-outline
go install github.com/acroca/go-symbols
go install golang.org/x/tools/cmd/guru
go install golang.org/x/tools/cmd/gorename
go install github.com/josharian/impl
go install github.com/rogpeppe/godef
go install github.com/sqs/goreturns
go install github.com/cweill/gotests/gotests