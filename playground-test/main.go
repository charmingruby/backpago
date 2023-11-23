package main

import "fmt"

const (
	AwsProvider BucketType = iota
)

type BucketType int

func main() {
	fmt.Printf("AwsProvider: %v\n", AwsProvider)
}
