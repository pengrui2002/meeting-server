package sfu

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type Client struct {
	baseURL string
	client  *http.Client
}

type CreateRoomRequest struct {
	MeetingID string `json:"meeting_id"`
}

type CreateRoomResponse struct {
	RoomID          string                 `json:"room_id"`
	RtpCapabilities map[string]interface{} `json:"rtp_capabilities"`
}

type CreateTransportRequest struct {
	Type   string `json:"type"`
	PeerID string `json:"peer_id"`
}

type TransportResponse struct {
	TransportID   string `json:"transport_id"`
	ICEParams     any    `json:"ice_params"`
	ICECandidates []any  `json:"ice_candidates"`
	DTLSParams    any    `json:"dtls_params"`
}

type ProduceRequest struct {
	TransportID   string                 `json:"transport_id"`
	Kind         string `json:"kind"`
	RtpParameters map[string]interface{} `json:"rtp_parameters"`
}

type ProduceResponse struct {
	ProducerID string `json:"producer_id"`
}

type ConsumeRequest struct {
	PeerID     string `json:"peer_id"`
	ProducerID string `json:"producer_id"`
}

type ConsumeResponse struct {
	ConsumerID    string                 `json:"consumer_id"`
	ProducerID    string                 `json:"producer_id"`
	PeerID       string                 `json:"peer_id"`
	Kind         string                 `json:"kind"`
	RtpParameters map[string]interface{} `json:"rtp_parameters"`
}

func NewClient(baseURL string) *Client {
	return &Client{
		baseURL: baseURL,
		client:  &http.Client{},
	}
}

func (c *Client) CreateRoom(meetingID string) (*CreateRoomResponse, error) {
	body, _ := json.Marshal(CreateRoomRequest{MeetingID: meetingID})
	resp, err := c.client.Post(c.baseURL+"/api/rooms", "application/json", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result CreateRoomResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) CreateTransport(roomID, transportType, peerID string) (*TransportResponse, error) {
	body, _ := json.Marshal(CreateTransportRequest{Type: transportType, PeerID: peerID})
	url := fmt.Sprintf("%s/api/rooms/%s/transports", c.baseURL, roomID)
	resp, err := c.client.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result TransportResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) CloseTransport(transportID string) error {
	url := fmt.Sprintf("%s/api/transports/%s", c.baseURL, transportID)
	req, _ := http.NewRequest("DELETE", url, nil)
	_, err := c.client.Do(req)
	return err
}

func (c *Client) CloseRoom(roomID string) error {
	url := fmt.Sprintf("%s/api/rooms/%s", c.baseURL, roomID)
	req, _ := http.NewRequest("DELETE", url, nil)
	_, err := c.client.Do(req)
	return err
}