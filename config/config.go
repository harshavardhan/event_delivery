package config

import "time"

// list of possible destinations
var Destinations = []string{"A", "B", "C", "D"}

var Delta = 20 * time.Second // delta for exponential backoff
