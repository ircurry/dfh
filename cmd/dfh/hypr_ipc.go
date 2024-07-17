package main

import (
	"fmt"
	"net"
	"os"
)

func listenEvents(con *net.Conn, lim int, ch chan string) {
    b := make([]byte, 1)
    str := ""

    for i := 0; (i < lim) || (lim <= 0); i++{
		for {
			(*con).Read(b)
			if string(b[:]) == "\n" {
				break
			}
			str = str + string(b[:])
		}
		ch <- str
		str = ""
    }

    close(ch)
    return
}

func printEvents(ch chan string) {
    for {
		msg, ok := <-ch
		if !ok {
			break
		}
		fmt.Println(msg)
    }
    return
}

func hyprMessage(con net.Conn, msg string) {
    _, err := con.Write([]byte(msg + "\n"))
    if err != nil {
		panic(err)
    }
}

func getHyprDir() string {
    var xdgRunDir, his string
    var ok bool
    xdgRunDir, ok = os.LookupEnv("XDG_RUNTIME_DIR")
    if !ok {
		panic("XDG_RUNTIME_DIR not found")
    }
    
    his, ok = os.LookupEnv("HYPRLAND_INSTANCE_SIGNATURE")
    if !ok {
		panic("HYPRLAND_INSTANCE_SIGNATURE not found")
    }
    
    return xdgRunDir + "/hypr/" + his
}

func getHyprCtlSocket() (net.Conn) {
    hyprDir := getHyprDir()
    sock1, err := net.Dial("unix", hyprDir + "/.socket.sock")
    if err != nil {
		panic(err)
    }
    return sock1
}

func getEventSocket() (net.Conn) {
    hyprDir := getHyprDir()
    sock2, err := net.Dial("unix", hyprDir + "/.socket2.sock")
    if err != nil {
		panic(err)
    }
    return sock2
}
