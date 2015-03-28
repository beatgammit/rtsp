package rtsp

import (
	"bufio"
	"errors"
	"io"
	"strconv"
	"strings"
)

type SessionSection struct {
	Version               int
	Originator            string
	SessionName           string
	SessionInformation    string
	URI                   string
	Email                 string
	Phone                 string
	ConnectionInformation string
	BandwidthInformation  string
}

func ParseSdp(r io.Reader) (SessionSection, error) {
	var packet SessionSection
	s := bufio.NewScanner(r)
	for s.Scan() {
		parts := strings.SplitN(s.Text(), "=", 2)
		if len(parts) == 2 {
			if len(parts[0]) != 1 {
				return packet, errors.New("SDP only allows 1-character variables")
			}

			switch parts[0] {
			// version
			case "v":
				ver, err := strconv.Atoi(parts[1])
				if err != nil {
					return packet, err
				}
				packet.Version = ver
			// owner/creator and session identifier
			case "o":
				// o=<username> <session id> <version> <network type> <address type> <address>
				// TODO: parse this
				packet.Originator = parts[1]
			// session name
			case "s":
				packet.SessionName = parts[1]
			// session information
			case "i":
				packet.SessionInformation = parts[1]
			// URI of description
			case "u":
				packet.URI = parts[1]
			// email address
			case "e":
				packet.Email = parts[1]
			// phone number
			case "p":
				packet.Phone = parts[1]
			// connection information - not required if included in all media
			case "c":
				// TODO: parse this
				packet.ConnectionInformation = parts[1]
			// bandwidth information
			case "b":
				// TODO: parse this
				packet.BandwidthInformation = parts[1]
			}
		}
	}
	return packet, nil
}
