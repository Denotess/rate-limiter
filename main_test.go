package main

import (
	"testing"
	"time"
)

func TestNewBucketAllowsUpToCapacity(t *testing.T) {
	bucket := NewLeakyBucket(5, 1)

	// First 5 requests should succeed
	for i := 0; i < 5; i++ {
		if !bucket.AllowRequest() {
			t.Errorf("Request %d should be allowed", i)
		}
	}

	// 6th request should fail
	if bucket.AllowRequest() {
		t.Error("Request beyond capacity should be denied")
	}
}

func TestBucketLeaksOverTime(t *testing.T) {
	bucket := NewLeakyBucket(10, 2)
	for i := 0; i < 10; i++ {
		bucket.AllowRequest()
	}

	if bucket.AllowRequest() {
		t.Error("Bucket is full should deny request")
	}
	time.Sleep(1 * time.Second)

	for i := 0; i < 2; i++ {
		if !bucket.AllowRequest() {
			t.Errorf("Request %d should be allowed after leaking", i)
		}
	}
	if bucket.AllowRequest() {
		t.Error("Bucket should be full again after 2 requests")
	}
}
