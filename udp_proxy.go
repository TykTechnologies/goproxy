package main

import (
	"context"
	log "github.com/Sirupsen/logrus"
	"net"
	"time"
)

func handleUDP(source string, targets []string) {
	sourceAddr, err := net.ResolveUDPAddr("udp", source)
	if err != nil {
		log.WithError(err).Fatal("Could not resolve source address:", source)
		return
	}

	var targetAddr []*net.UDPAddr
	for _, v := range targets {
		addr, err := net.ResolveUDPAddr("udp", v)
		if err != nil {
			log.WithError(err).Fatal("Could not resolve target address:", v)
			return
		}
		targetAddr = append(targetAddr, addr)
	}

	sourceConn, err := net.ListenUDP("udp", sourceAddr)
	if err != nil {
		log.WithError(err).Fatal("Could not listen on address:", source)
		return
	}

	defer sourceConn.Close()

	// use only 1
	v := targetAddr[0]
	var targetConn []*net.UDPConn
	conn, err := net.DialUDP("udp", nil, v)
	if err != nil {
		log.WithError(err).Fatal("Could not connect to target address:", v)
		return
	}

	defer conn.Close()
	targetConn = append(targetConn, conn)

	log.Printf(">> Starting udp proxy, Source at %v, Target at %v...", source, targets)

	for {
		b := make([]byte, 10240)
		n, addr, err := sourceConn.ReadFromUDP(b)

		if err != nil {
			log.WithError(err).Error("Could not receive a packet")
			continue
		}

		log.WithField("addr", addr.String()).WithField("bytes", n).Info("Packet received")
		for _, v := range targetConn {
			// don't block
			ctx, _ := context.WithTimeout(context.Background(), time.Second*1)
			//defer cancel()
			go forwardAndWaitForResponse(v, sourceConn, b, n, addr, ctx)
		}
	}
}

func forwardAndWaitForResponse(v, source *net.UDPConn, b []byte, n int, from *net.UDPAddr, ctx context.Context) error {
	if _, err := v.Write(b[0:n]); err != nil {
		log.WithError(err).Warn("Could not forward packet.")
		return err
	}

	for {
		b := make([]byte, 10240)

		select {
		case <-ctx.Done():
			log.Warning(ctx.Err())
			return ctx.Err()
		default:
			n, addr, err := v.ReadFromUDP(b)

			if err != nil {
				log.WithError(err).Error("Could not receive a packet")
				return err
			}

			log.WithField("addr", addr.String()).WithField("bytes", n).Info("Return packet received")

			if _, err := source.WriteTo(b[0:n], from); err != nil {
				log.WithError(err).Warn("Could not return packet.")
				return err
			}

			// leave after response
			return nil
		}

	}

	return nil
}
