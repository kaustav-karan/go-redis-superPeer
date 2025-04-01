package models

import "time"

type Peer struct {
	TrackUri string `json:"trackUri" redis:"trackUri"`
}

type TrackMetadata struct {
	PublisherName string `json:"publisherName" redis:"publisherName"`
	Size          int    `json:"size" redis:"size"`
	PeerAvailable bool   `json:"peerAvailable" redis:"peerAvailable"`
	PeerList      []Peer `json:"peerList" redis:"peerList"`
}

type TrackLog struct {
	TrackId      string       `json:"trackId" redis:"trackId"`
	TrackMetadata TrackMetadata `json:"trackMetadata" redis:"trackMetadata"`
	TimeStamp   time.Time    `json:"timeStamp" redis:"timeStamp"`
}