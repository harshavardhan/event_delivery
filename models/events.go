package models

type Event struct {
	UserID  string `json:"userId"`
	Payload string
}

type EventResponse struct {
	Msg string `json:"msg"`
}

type EventMetadata struct {
	// duplicating Event fields to avoid doing embedded structs flatten before write to redis
	UserID        string `redis:"userId"`
	Payload       string `redis:"payload"`
	Timestamp     int64  `redis:"timestamp"`     // (to be generally set by the client) set to server time if not sent (unix epoch nanoseconds)
	ExecTimestamp int64  `redis:"execTimestamp"` // initially equal to timestamp, on retries or waits set to next exec timestamp
	RetryCount    int    `redis:"retryCount"`
}
