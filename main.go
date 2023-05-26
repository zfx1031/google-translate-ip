package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"sort"
	"sync"
	"time"
)

func main() {

	type delayList struct {
		ip    string
		delay int64
	}

	
	openFile, err := os.OpenFile("ip.txt", os.O_RDONLY, 0)
	if err != nil && errors.Is(err, os.ErrNotExist) {

		log.Println("当前目录不存在 ip.txt 文件...")

		fmt.Println("按任意键退出...")
		_, err := fmt.Scanln()
		if err != nil {
			return
		}
		return
	}

	
	var ips []string

	scanner := bufio.NewScanner(openFile)

	for scanner.Scan() {
		ips = append(ips, scanner.Text())
	}

	
	var wg sync.WaitGroup

	var delays []delayList

	for _, ip := range ips {

		ip := ip

		wg.Add(1)

		go func() {

			defer wg.Done()
			start := time.Now()
			conn, err := net.DialTimeout("tcp", ip+":80", 5*time.Second)
			if err != nil {
				return
			}

			defer func(conn net.Conn) {
				err := conn.Close()
				if err != nil {
					fmt.Println("连接关闭失败：", err)
					return
				}
			}(conn)

			delay := time.Since(start).Milliseconds()
			delays = append(delays, delayList{
				ip:    ip,
				delay: delay,
			})
		}()
	}

	wg.Wait()

	
	sort.Slice(delays, func(i, j int) bool {
		return delays[i].delay < delays[j].delay
	})

	
	for _, v := range delays {
		fmt.Println(v.ip, "延迟：", v.delay, "毫秒")
	}

	fmt.Println("所有IP测速完毕，按任意键退出...")

	_, err = fmt.Scanln()
	if err != nil {
		return
	}

}
