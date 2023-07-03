package main

import (
	"fmt"
	"net"
	"time"
)

type Domain struct {
	Name string
	Ping string
}

func ping(domain string, ch chan<- Domain) {
	start := time.Now()
	conn, err := net.Dial("icmp", domain)
	if err != nil {
		ch <- Domain{Name: domain, Ping: fmt.Sprintf("Error resolviendo la dirección IP %s: %s", domain, err.Error())}
		return
	}

	if err != nil {
		ch <- Domain{Name: domain, Ping: fmt.Sprintf("Error al establecer la conexión %s: %s", domain, err.Error())}
		return
	}
	defer conn.Close()
	duration := time.Since(start)
	ch <- Domain{Name: domain, Ping: duration.String()}
}

func orderByPing(domains []Domain) {
	n := len(domains)

	for i := 0; i < n-1; i++ {
		minIdx := i
		for j := i + 1; j < n; j++ {
			// Parsear la duración del ping

			if domains[j].Ping < domains[minIdx].Ping {
				minIdx = j
			}
		}
		// Intercambiar los elementos en las posiciones i y minIdx
		domains[i], domains[minIdx] = domains[minIdx], domains[i]
	}
}

func main() {
	domains := []Domain{
		{Name: "msn.com"},
		{Name: "github.com"},
		{Name: "dailymotion.com"},
		{Name: "google.com"},
		{Name: "ok.ru"},
		{Name: "netflix.com"},
		{Name: "facebook.com"},
		{Name: "reddit.com"},
		{Name: "office.com"},
		{Name: "cnn.com"},
		{Name: "imgur.com"},
		{Name: "whatsapp.com"},
		{Name: "tiktok.com"},
		{Name: "microsoft.com"},
		{Name: "espn.com"},
		{Name: "yahoo.com"},
		{Name: "dropbox.com"},
		{Name: "twitch.tv"},
		{Name: "ebay.com"},
		{Name: "paypal.com"},
		{Name: "youtube.com"},
		{Name: "apple.com"},
		{Name: "wikipedia.org"},
		{Name: "quora.com"},
		{Name: "adobe.com"},
		{Name: "spotify.com"},
		{Name: "xvideos.com"},
		{Name: "pinterest.com"},
		{Name: "bing.com"},
		{Name: "blogspot.com"},
		{Name: "tumblr.com"},
		{Name: "qq.com"},
		{Name: "twitter.com"},
		{Name: "walmart.com"},
		{Name: "stackoverflow.com"},
		{Name: "vk.com"},
		{Name: "baidu.com"},
		{Name: "aliyun.com"},
		{Name: "aliexpress.com"},
		{Name: "yandex.ru"},
		{Name: "nytimes.com"},
		{Name: "wordpress.com"},
		{Name: "tmall.com"},
		{Name: "instagram.com"},
		{Name: "live.com"},
		{Name: "outlook.com"},
		{Name: "amazon.com"},
		{Name: "taobao.com"},
		{Name: "linkedin.com"},
	}

	ch := make(chan Domain)

	for i := range domains {
		go ping(domains[i].Name, ch)
	}

	pingResults := make([]Domain, 0, len(domains))

	for range domains {
		result := <-ch
		pingResults = append(pingResults, result)
	}
	close(ch)
	orderByPing(domains)
	// Mostrar los resultados ordenados
	//fmt.Println("Resultados ordenados:")

	for _, domain := range pingResults {
		fmt.Printf("Dominio: %s - Ping: %s\n", domain.Name, domain.Ping)
	}
}
