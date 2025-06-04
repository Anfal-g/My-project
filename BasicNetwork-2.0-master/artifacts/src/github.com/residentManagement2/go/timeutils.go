package main

import "time"

// GetUniformTimestamp ensures both chaincodes use same time format
func GetUniformTimestamp() int64 {
    return time.Now().UTC().UnixMilli() // Milliseconds since epoch in UTC
}

// IsWithinWindow checks if current time is within range
func IsWithinWindow(start, end int64) bool {
    now := GetUniformTimestamp()
    return now >= start && now <= end
}