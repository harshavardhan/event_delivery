package config

import "time"

// list of possible destinations
var Destinations = []string{"A", "B", "C", "D"}

var Delta = 200 * time.Millisecond // delta for exponential backoff
