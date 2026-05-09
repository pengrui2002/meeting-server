package media

import (
	"errors"
)

// SFUServer provides WebRTC SFU functionality for video conferencing
// Note: This is a placeholder implementation. Full SFU support requires:
// - Node.js mediasoup server (https://mediasoup.org)
// - Or pion/ion-sfu which needs network access to install
type SFUServer struct {
	// sfu *sfu.SFU // TODO: Uncomment when ion-sfu is available
}

// NewSFUServer creates a new SFU server instance
// Returns an error if SFU library is not available
func NewSFUServer() (*SFUServer, error) {
	// TODO: Uncomment when ion-sfu is available
	// config := sfu.Config{}
	// s, err := sfu.New(config)
	// if err != nil {
	//     return nil, err
	// }
	// return &SFUServer{sfu: s}, nil

	return nil, errors.New("SFU not implemented: requires Node.js mediasoup or network access to install pion/ion-sfu")
}

// Close shuts down the SFU server
func (s *SFUServer) Close() error {
	// TODO: Uncomment when ion-sfu is available
	// return s.sfu.Close()

	return errors.New("SFU not implemented")
}
