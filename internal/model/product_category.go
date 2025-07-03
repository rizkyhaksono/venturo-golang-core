package model

type ProductCategory uint8 // Use uint8 since it's a small, non-negative number

const (
	_ ProductCategory = iota // Start with 0, but we'll ignore it
	Goods
	Service
	Subscription
)
