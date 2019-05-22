### useage

    Usage: COMMAND args1 args2 args3
    args1 is email address
    args2 is the mesages's title
    args3 is messages's content

### conf

配置文件和程序放到一个目录下。

     [main]
     #smtp服务器地址
     smtp_server = smtp.163.com
     #服务端口
     port = 465
     #email 地址
     email = yang@163.com
     #授权码
     passworld = yang

### 编译

    centos
    yum install golang -y
    cd
    mkdir go{pkg,bin,src}
    cd go/src
    go get https://github.com/yangyouwei/gomail.git
    cd github.com/yangyouwei/gomail
    go build gomail.go