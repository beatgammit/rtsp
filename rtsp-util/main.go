package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"

	"rtsp"
)

func init() {
	flag.Parse()
}

const sampleRequest = `OPTIONS rtsp://example.com/media.mp4 RTSP/1.0
CSeq: 1
Require: implicit-play
Proxy-Require: gzipped-messages

`

const sampleResponse = `RTSP/1.0 200 OK
CSeq: 1
Public: DESCRIBE, SETUP, TEARDOWN, PLAY, PAUSE

`

func main() {
	if len(flag.Args()) >= 1 {
		rtspUrl := flag.Args()[0]

		sess := rtsp.NewSession()
		res, err := sess.Options(rtspUrl)
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Println("Options:")
		fmt.Println(res)

		res, err = sess.Describe(rtspUrl)
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Println("Describe:")
		fmt.Println(res)

		p, err := rtsp.ParseSdp(&io.LimitedReader{R: res.Body, N: res.ContentLength})
		if err != nil {
			log.Fatalln(err)
		}
		log.Printf("%+v", p)

		rtpPort, rtcpPort := 8000, 8001
		res, err = sess.Setup(rtspUrl, fmt.Sprintf("RTP/AVP;unicast;client_port=%d-%d", rtpPort, rtcpPort))
		if err != nil {
			log.Fatalln(err)
		}
		log.Println(res)

		res, err = sess.Play(rtspUrl, res.Header.Get("Session"))
		if err != nil {
			log.Fatalln(err)
		}
		log.Println(res)
	} else {
		r, err := rtsp.ReadRequest(bytes.NewBufferString(sampleRequest))
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(r)
		}

		res, err := rtsp.ReadResponse(bytes.NewBufferString(sampleResponse))
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(res)
		}
	}
}
