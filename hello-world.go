package main

import (
	"fmt"
	"net"
	"os"
	"sync"
)

func listenEvents(con *net.Conn, lim int, ch chan string, wg *sync.WaitGroup) {
    defer wg.Done()
    b := make([]byte, 1)
    c := 0
    
    for {
	wg.Add(1)
	if c >= lim {
	    wg.Done()
	    break
	}
	(*con).Read(b)
	str := string(b[:])
	if str == "\n" {
	    c++
	}
	ch <- str
	wg.Done()
    }

    close(ch)
    return
}

func printStream(ch chan string, wg *sync.WaitGroup) {
    defer wg.Done()
    for {
	wg.Add(1)
	msg, ok := <-ch
	if !ok {
	    wg.Done()
	    break
	}
	fmt.Print(msg)
	wg.Done()
    }
    return
}

func printMessages(conn *net.Conn, lim int, ch chan string, wg *sync.WaitGroup) {
    wg.Add(2)
    defer wg.Done()
    go printStream(ch, wg)
    go listenEvents(conn, lim, ch, wg)
    return
}

// func hyprMessage()

func main() {
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
    
    hyprDir := xdgRunDir + "/hypr/" + his
    // sock1, err := net.Dial("unix", hyprDir + "/.socket.sock")
    // if err != nil {
    // 	panic(err)
    // }
    sock2, err := net.Dial("unix", hyprDir + "/.socket2.sock")
    if err != nil {
	panic(err)
    }

    // _, err = sock1.Write([]byte("dispatch exec alacritty\n"))
    // if err != nil {
    // 	panic(err)
    // }

    var wg sync.WaitGroup
    wg.Add(1)
    ch := make(chan string, 20)
    printMessages(&sock2, 20, ch, &wg)
    wg.Wait()
    return
}
