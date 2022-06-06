package main

import (
	"bufio"
	"log"
	"net"
	"strconv"
	"strings"
)

var debug bool

// mock 下 redis server
func main() {
	debug = true

	listener, err := net.Listen("tcp", "127.0.0.1:6389")
	if err != nil {
		log.Panicf("error: %s", err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("failed to accept: %s", err)
			continue
		}

		// 对每一个接受到的请求，都起一个goroutine处理
		go handleConn(conn)
	}
}

/*
+ 开头，代表这是一个字符串，simple string，也就是非二进制安全的字符串。然后以 \r\n 结尾，比如如果服务端返回 OK，那么实际返回的内容是 +OK\r\n。
$ 开头，代表这是一个字符串，但是是二进制安全的字符串。当然，也可以传输简单的字符串，比如上面的OK，会被传输为 $2\r\nOK\r\n，可以看出来，和简单字符串的不同之处在于，最前面告诉了我们内容到底有多长。有一个特例，那就是NULL，表示为 $-1\r\n。
- 开头，代表这是一个错误，比如 -Error message\r\n，实际上要显示的错误就是 Error message，也就是说，中间的部分就是错误信息。
: 开头，代表这是一个数字。比如 :10000\r\n 代表10000，而 :0\r\n 就是0。
* 开头，代表这是一个数组，比如由foo和bar两个字符串组成的数组，就应该返回为："*2\r\n$3\r\nfoo\r\n$3\r\nbar\r\n" 。*2 代表这个数组有两个元素，后续的内容其实就是上面几种内容的组合。有两个特例：长度为0的数组表示为 *0\r\n，空数组表示为 *-1\r\n。
*/
func handleConn(conn net.Conn) {
	defer conn.Close()

	buf := bufio.NewReader(conn)
	sb := strings.Builder{}

	// eg. "*3\r\n$3\r\nset\r\n$1\r\na\r\n$1\r\nb\r\n"
	for {
		// 读取第一个\r\n结尾的
		bs, err := buf.ReadBytes('\n')
		if err != nil {
			log.Printf("err: %s", err)
			break
		}

		// 如果不是数组，我们就直接panic了
		if bs[0] != '*' {
			log.Panicf("bad bs: %s", bs)
		}

		// 数组里有多少个元素，我们要解析出来，然后读取
		length, err := strconv.ParseUint(string(bs[1:len(bs)-2]), 10, 64)
		if err != nil {
			log.Panicf("bad length: %s", err)
		}
		if debug {
			log.Printf("length: %d", length)
		}

		// 把最开始读到的命令写进去
		sb.Write(bs)
		var i uint64 = 0
		for ; i < length; i++ {
			bs, err = buf.ReadBytes('\n')
			if err != nil {
				log.Printf("err: %s", err)
				break
			}

			if bs[0] == '$' {
				// 如果是复杂字符串，那么就有两个\r\n
				sb.Write(bs)
				sz, _ := strconv.ParseUint(string(bs[1:len(bs)-2]), 10, 64)
				bs, _ = buf.ReadBytes('\n')

				if debug {
					// eg.
					//  i: 0 size: 3 str: set
					//  i: 1 size: 1 str: a
					//  i: 2 size: 1 str: b
					log.Printf("i: %d size: %d str: %s", i, sz, bs[:sz])
				}
			}
			sb.Write(bs)
		}

		// 打印出来
		log.Printf("content: %#v", sb.String())
		conn.Write([]byte("+OK\r\n"))

		// 重置
		sb.Reset()
	}
}
