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

	conn, err := net.Dial("ip4:icmp", domain)
	if err != nil {
		ch <- Domain{Name: domain, Ping: fmt.Sprintf("Error al crear la conexión ICMP: %s", err.Error())}
		return
	}
	defer conn.Close()

	// Configurar timeout de lectura
	conn.SetReadDeadline(time.Now().Add(2 * time.Second))

	// Crear un mensaje ICMP Echo Request
	icmpMsg := make([]byte, 8)
	icmpMsg[0] = 8 // ICMP Echo Request Type
	icmpMsg[1] = 0 // ICMP Echo Request Code
	icmpMsg[2] = 0 // ICMP Echo Request Checksum (zeroed for now)
	icmpMsg[3] = 0 // ICMP Echo Request Checksum (zeroed for now)
	icmpMsg[4] = 0 // ICMP Echo Request Identifier (zeroed for now)
	icmpMsg[5] = 0 // ICMP Echo Request Identifier (zeroed for now)
	icmpMsg[6] = 0 // ICMP Echo Request Sequence Number (zeroed for now)
	icmpMsg[7] = 0 // ICMP Echo Request Sequence Number (zeroed for now)

	// Calcular el checksum para el mensaje ICMP
	checksum := checkSum(icmpMsg)
	icmpMsg[2] = byte(checksum >> 8)
	icmpMsg[3] = byte(checksum)

	// Enviar el mensaje ICMP al dominio
	_, err = conn.Write(icmpMsg)
	if err != nil {
		ch <- Domain{Name: domain, Ping: fmt.Sprintf("Error al enviar el mensaje ICMP: %s", err.Error())}
		return
	}

	// Esperar la respuesta ICMP
	reply := make([]byte, 1500)
	_, err = conn.Read(reply)
	if err != nil {
		if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
			ch <- Domain{Name: domain, Ping: "Tiempo de espera agotado o superior a 2000ms"}
		} else {
			ch <- Domain{Name: domain, Ping: fmt.Sprintf("Error al leer la respuesta ICMP: %s", err.Error())}
		}
		return
	}

	// Calcular la duración de latencia
	duration := time.Since(start)
	ch <- Domain{Name: domain, Ping: duration.String()}
}
func checkSum(msg []byte) uint16 {
	sum := 0
	for i := 0; i < len(msg)-1; i += 2 {
		sum += int(msg[i])*256 + int(msg[i+1])
	}
	if len(msg)%2 == 1 {
		sum += int(msg[len(msg)-1]) * 256
	}
	sum = (sum >> 16) + (sum & 0xffff)
	sum += sum >> 16
	return uint16(^sum)
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
		{Name: "reduno.com.ar"},
		{Name: "radios.lu17.com"},
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
