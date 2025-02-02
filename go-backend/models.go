package main

import "time"

type Song struct {
	MD5Hash   string
	Title     string
	Artist    string
	Creator   string
	Folder    string
	File      string
	Audio     string
	Mapper    string
	TotalTime time.Duration
	Url       string
	Image     string
}

type CollectionPreview struct {
	name  string
	image string
	index int
	items int
}

type Collection struct {
	name  string
	items int
	Songs []Song
}

type ActiveSearch struct {
	artist string
	Songs  []Song
}
