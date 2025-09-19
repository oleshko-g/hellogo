package main

import (
	"errors"
	"flag"
	"fmt"
	"strconv"
	"strings"
)

type NetAddress struct {
	Host string
	Port int
}

func (na *NetAddress) String() string {
	return fmt.Sprintf("%s:%s", na.Host, strconv.Itoa(na.Port))
}

func (na *NetAddress) Set(flagValue string) error {
	if flagValue == "" {
		return errors.New("empty value")
	}

	splitFlagValue := strings.Split(strings.TrimSpace(flagValue), ":")
	if len(splitFlagValue) == 2 {

		host, err := parseHost(splitFlagValue[0])
		if err != nil {
			return err
		}
		na.Host = host

		port, err := parsePort(splitFlagValue[1])
		if err != nil {
			return err
		}
		na.Port = port
	}

	return nil
}

func parseHost(hostStringValue string) (host string, err error) {
	if hostStringValue == "" {
		return "", errors.New("empty Host")
	}
	// TODO: white space check
	return hostStringValue, nil
}

func parsePort(portStringValue string) (port int, err error) {
	if portStringValue == "" {
		return 0, errors.New("empty Port")
	}
	// TODO: white space check
	portIntValue, err := strconv.Atoi(portStringValue)
	if err != nil {
		return 0, errors.New("port is not an int")
	}
	return portIntValue, nil
}

func main() {
	addr := new(NetAddress)
	_ = flag.Value(addr)
	// проверка реализации
	flag.Var(addr, "addr", "Net address host:port")
	flag.Parse()
	fmt.Println(addr.Host)
	fmt.Println(addr.Port)
}
