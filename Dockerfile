
FROM golang:1.10.7

#作者
MAINTAINER bobo "bobo.huang@ucloud.cn"
 
#创建工作目录
RUN mkdir -p /go/src/Demo
#进入工作目录
WORKDIR /go/src/Demo
 
#将当前目录下的所有文件复制到指定位置
COPY . /go/src/Demo
 
 
#端口
EXPOSE 8000
 
#运行
ENTRYPOINT ["./Demo"]
