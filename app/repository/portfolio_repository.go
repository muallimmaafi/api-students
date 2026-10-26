package repository

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

const PortfolioCollectionName = "student_portfolios"

type PortfolioRepository struct {
	collection *mongo.Collection
}

func NewPortfolioRepository(db *mongo.Database) *PortfolioRepository {
	return &PortfolioRepository{
		collection: db.Collection(PortfolioCollectionName),
	}
}

// EnsureIndexes membuat indeks yang dibutuhkan collection portofolio.
func (r *PortfolioRepository) EnsureIndexes(ctx context.Context) error {
	_, err := r.collection.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{
			Keys: bson.D{{Key: "student_id", Value: 1}},
			Options: options.Index().
				SetUnique(true).
				SetName("uniq_portfolio_student_id"),
		},
		{
			Keys: bson.D{{Key: "projects.technologies", Value: 1}},
			Options: options.Index().
				SetName("idx_portfolio_project_technologies"),
		},
	})
	if err != nil {
		return fmt.Errorf("membuat indeks portofolio: %w", err)
	}

	return nil
}

// Ping memastikan collection dapat diakses melalui koneksi MongoDB.
func (r *PortfolioRepository) Ping(ctx context.Context) error {
	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	return r.collection.Database().Client().Ping(pingCtx, nil)
}
