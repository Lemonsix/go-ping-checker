package main

import (
	"fmt"
	"net"
	"sort"
	"time"
)

type Domain struct {
	Name string
	Ping string
}

func ping(domain string, ch chan<- Domain) {
	ipAddr, err := net.ResolveIPAddr("ip", domain)
	if err != nil {
		ch <- Domain{Name: domain, Ping: fmt.Sprintf("Error resolviendo la dirección IP para %s: %s", domain, err.Error())}
		return
	}

	start := time.Now()
	conn, err := net.DialIP("ip:icmp", nil, ipAddr)
	if err != nil {
		ch <- Domain{Name: domain, Ping: fmt.Sprintf("Error al establecer la conexión para %s: %s", domain, err.Error())}
		return
	}
	defer conn.Close()

	duration := time.Since(start)
	ch <- Domain{Name: domain, Ping: duration.String()}
}

type ByPing []Domain

func (a ByPing) Len() int      { return len(a) }
func (a ByPing) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByPing) Less(i, j int) bool {
	// Si el ping es "error", el dominio se coloca al final de la lista
	if a[i].Ping == "error" {
		return false
	}
	if a[j].Ping == "error" {
		return true
	}

	// Convertir los valores de ping a duraciones
	pingDuration1, err1 := time.ParseDuration(a[i].Ping)
	pingDuration2, err2 := time.ParseDuration(a[j].Ping)

	// Si hay un error al convertir las duraciones, colocar el dominio al final de la lista
	if err1 != nil {
		return false
	}
	if err2 != nil {
		return true
	}

	// Comparar las duraciones de ping
	return pingDuration1 < pingDuration2
}

func main() {
	domains := []string{
		"msn.com",
		"github.com",
		"dailymotion.com",
		"google.com",
		"ok.ru",
		"netflix.com",
		"facebook.com",
		"reddit.com",
		"office.com",
		"cnn.com",
		"imgur.com",
		"whatsapp.com",
		"tiktok.com",
		"microsoft.com",
		"espn.com",
		"yahoo.com",
		"dropbox.com",
		"twitch.tv",
		"ebay.com",
		"paypal.com",
		"youtube.com",
		"apple.com",
		"wikipedia.org",
		"quora.com",
		"adobe.com",
		"spotify.com",
		"xvideos.com",
		"pinterest.com",
		"bing.com",
		"blogspot.com",
		"tumblr.com",
		"qq.com",
		"twitter.com",
		"walmart.com",
		"stackoverflow.com",
		"vk.com",
		"baidu.com",
		"aliyun.com",
		"aliexpress.com",
		"yandex.ru",
		"nytimes.com",
		"wordpress.com",
		"tmall.com",
		"instagram.com",
		"live.com",
		"outlook.com",
		"amazon.com",
		"taobao.com",
		"linkedin.com",
	}

	ch := make(chan Domain)

	for _, domain := range domains {
		go ping(domain, ch)
	}

	pingResults := make([]Domain, 0, len(domains))

	for range domains {
		result := <-ch
		pingResults = append(pingResults, result)
	}

	close(ch)

	// Ordenar los resultados por ping
	sort.Sort(ByPing(pingResults))

	// Mostrar los resultados ordenados
	fmt.Println("Resultados ordenados:")
	for _, result := range pingResults {
		fmt.Printf("Dominio: %s - Ping: %s\n", result.Name, result.Ping)
	}
}
