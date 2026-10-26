package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	"api-students/app/model"
	"github.com/google/uuid"
)

const PortfolioCollectionName = "student_portfolios"

var (
	ErrPortfolioNotFound = errors.New("portofolio tidak ditemukan")
	ErrProjectLimit      = errors.New("batas 20 proyek tercapai")
)

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

func (r *PortfolioRepository) FindByStudentID(
	ctx context.Context,
	studentID int,
) (model.Portfolio, error) {
	var portfolio model.Portfolio

	err := r.collection.FindOne(ctx, bson.D{
		{Key: "student_id", Value: studentID},
	}).Decode(&portfolio)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return model.Portfolio{}, ErrNotFound
		}
		return model.Portfolio{}, fmt.Errorf("mengambil portofolio: %w", err)
	}

	return portfolio, nil
}

func (r *PortfolioRepository) Replace(
	ctx context.Context,
	portfolio model.Portfolio,
) error {
	opts := options.Replace().SetUpsert(true)

	_, err := r.collection.ReplaceOne(
		ctx,
		bson.D{{Key: "student_id", Value: portfolio.StudentID}},
		portfolio,
		opts,
	)
	if err != nil {
		return fmt.Errorf("menyimpan portofolio: %w", err)
	}

	return nil
}

// AddProject menambahkan proyek secara atomik jika jumlahnya belum 20.

func (r *PortfolioRepository) AddProject(
	ctx context.Context,
	studentID int,
	project model.Project,
) (string, error) {
	project.ID = uuid.NewString()
	projectID := project.ID

	filter := bson.D{
		{Key: "student_id", Value: studentID},
		{Key: "$expr", Value: bson.D{
			{Key: "$lt", Value: bson.A{
				bson.D{{Key: "$size", Value: bson.D{
					{Key: "$ifNull", Value: bson.A{"$projects", bson.A{}}},
				}}},
				20,
			}},
		}},
	}

	update := bson.D{
		{Key: "$push", Value: bson.D{
			{Key: "projects", Value: project},
		}},
	}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return "", fmt.Errorf("menambahkan proyek: %w", err)
	}

	if result.MatchedCount == 0 {
		// Bedakan portofolio yang tidak ditemukan dengan batas proyek.
		count, err := r.collection.CountDocuments(ctx, bson.D{
			{Key: "student_id", Value: studentID},
		})
		if err != nil {
			return "", fmt.Errorf("memeriksa portofolio: %w", err)
		}

		if count == 0 {
			return "", ErrPortfolioNotFound
		}

		return "", ErrProjectLimit
	}

	return projectID, nil
}

func (r *PortfolioRepository) DeleteProject(
	ctx context.Context,
	studentID int,
	projectID string,
) error {
	filter := bson.D{
		{Key: "student_id", Value: studentID},
	}

	update := bson.D{
		{Key: "$pull", Value: bson.D{
			{Key: "projects", Value: bson.D{
				{Key: "id", Value: projectID},
			}},
		}},
	}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("menghapus proyek: %w", err)
	}
	if result.MatchedCount == 0 || result.ModifiedCount == 0 {
		return ErrNotFound
	}

	return nil
}

type TechnologyStat struct {
	Technology string `bson:"_id" json:"technology"`
	Count      int    `bson:"count" json:"count"`
}

// TechnologyStats menghitung 10 teknologi yang paling sering digunakan.
func (r *PortfolioRepository) TechnologyStats(
	ctx context.Context,
) ([]TechnologyStat, error) {
	pipeline := mongo.Pipeline{
		bson.D{{Key: "$unwind", Value: "$projects"}},
		bson.D{{Key: "$unwind", Value: "$projects.technologies"}},
		bson.D{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$projects.technologies"},
			{Key: "count", Value: bson.D{
				{Key: "$sum", Value: 1},
			}},
		}}},
		bson.D{{Key: "$sort", Value: bson.D{
			{Key: "count", Value: -1},
			{Key: "_id", Value: 1},
		}}},
		bson.D{{Key: "$limit", Value: 10}},
	}

	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, fmt.Errorf("menghitung statistik teknologi: %w", err)
	}
	defer cursor.Close(ctx)

	stats := make([]TechnologyStat, 0)
	if err := cursor.All(ctx, &stats); err != nil {
		return nil, fmt.Errorf("membaca statistik teknologi: %w", err)
	}

	return stats, nil
}
