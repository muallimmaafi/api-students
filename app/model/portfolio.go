package model

import "time"

type Portfolio struct {
	StudentID    int           `bson:"student_id" json:"student_id"`
	Headline     string        `bson:"headline" json:"headline"`
	Bio          string        `bson:"bio" json:"bio"`
	Projects     []Project     `bson:"projects" json:"projects"`
	Certificates []Certificate `bson:"certificates" json:"certificates"`
}

type Project struct {
	ID           string   `bson:"id" json:"id"`
	Title        string   `bson:"title" json:"title"`
	Technologies []string `bson:"technologies" json:"technologies"`
	Link         string   `bson:"link" json:"link"`
	Year         int      `bson:"year" json:"year"`
}

type Certificate struct {
	Name   string    `bson:"name" json:"name"`
	Issuer string    `bson:"issuer" json:"issuer"`
	Date   time.Time `bson:"date" json:"date"`
}
