package services

import (
	"context"
	"database/sql"
	"log"
	"strconv"
	"time"

	"github.com/guatom999/self-boardcast/internal/config"
	"github.com/guatom999/self-boardcast/internal/entities"
	"github.com/guatom999/self-boardcast/internal/models"
	"github.com/guatom999/self-boardcast/internal/repositories"
	"github.com/guatom999/self-boardcast/internal/utils"
)

// WaterLevelService handles business logic
type waterLevelService struct {
	repo    repositories.WaterLevelRepositoryInterface
	baseURL string
	cfg     *config.Config
}

type WaterLevelServiceInterface interface {
	// ProcessImage(ctx context.Context, imageURL string) (*models.WaterLevel, error)
	GetAllLocations(ctx context.Context, limit int) ([]models.LocationWithWaterLevel, error)
	GetByLocationID(ctx context.Context, id string) ([]*models.WaterLocationDetailRes, error)
	ScheduleGetWaterLevel(ctx context.Context) (string, int, error)
	ScheduleDeleteWaterLevel(ctx context.Context, fileName string, locationID int) error
	CreateWaterLevel(ctx context.Context, req *models.CreateWaterLevelReq) error
}

func NewWaterLevelService(repo repositories.WaterLevelRepositoryInterface, baseURL string, cfg *config.Config) WaterLevelServiceInterface {
	return &waterLevelService{
		repo:    repo,
		baseURL: baseURL,
		cfg:     cfg,
	}
}

func (s *waterLevelService) GetAllLocations(ctx context.Context, limit int) ([]models.LocationWithWaterLevel, error) {
	locations, err := s.repo.GetAll(ctx, limit)
	if err != nil {
		return nil, err
	}

	for i := range locations {
		if locations[i].Image != nil && *locations[i].Image != "" {
			imageURL := utils.BuildImageURL(s.baseURL, *locations[i].Image)
			locations[i].Image = &imageURL
		}
	}

	return locations, nil
}

func (s *waterLevelService) GetByLocationID(ctx context.Context, id string) ([]*models.WaterLocationDetailRes, error) {

	locationID, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}

	waterLevelsRes := make([]*models.WaterLocationDetailRes, 0)

	results, err := s.repo.GetByLocationID(ctx, locationID)
	if err != nil {
		return nil, err
	}

	for _, res := range results {
		waterLevelsRes = append(waterLevelsRes, &models.WaterLocationDetailRes{
			LocationID: res.LocationID,
			LevelCm:    res.LevelCm,
			Image:      res.Image,
			Danger:     res.Danger,
			IsFlooded:  res.IsFlooded,
			Source:     res.Source,
			MeasuredAt: res.MeasuredAt,
			Note:       res.Note,
		})
	}
	return waterLevelsRes, nil
}

func (s *waterLevelService) CreateWaterLevel(ctx context.Context, req *models.CreateWaterLevelReq) error {

	if err := s.repo.CreateWaterLevel(ctx, &entities.WaterLevel{
		LocationID: req.LocationID,
		LevelCm:    req.LevelCm,
		Image:      req.Image,
		Danger:     req.Danger,
		IsFlooded:  req.IsFlooded,
		MeasuredAt: req.MeasuredAt,
		Note:       req.Note,
	}); err != nil {
		return err
	}

	return nil
}

func (s *waterLevelService) ScheduleGetWaterLevel(ctx context.Context) (string, int, error) {

	fileName := utils.GenerateFileName()

	result, err := utils.PredictWaterLevel(s.cfg.App.ImageProcessingDir, fileName)
	if err != nil {
		log.Println(err)
		return "", 0, err
	}

	entity := &entities.WaterLevel{
		LocationID: 28,
		LevelCm:    result.WaterLevel * 100,
		Image:      result.FileName,
		Danger: func(waterLevel float64) string {
			if waterLevel < 1 {
				return "SAFE"
			} else if waterLevel < 2 {
				return "WATCH"
			} else {
				return "DANGER"
			}
		}(result.WaterLevel),
		IsFlooded:  false,
		Source:     sql.NullString{String: "sensor-1"},
		MeasuredAt: time.Now(),
		Note:       "get value of waterLevel from cctv of water",
	}

	if err := s.repo.CreateWaterLevel(ctx, entity); err != nil {
		log.Println("failed to create water level", err)
		return "", 0, err
	}

	return fileName, int(entity.LocationID), nil
}

func (s *waterLevelService) ScheduleDeleteWaterLevel(ctx context.Context, fileName string, locationID int) error {

	if err := utils.DeleteFile(s.cfg.App.UploadDir + "/" + fileName); err != nil {
		log.Println("failed to delete file", err)
		return err
	}

	if err := s.repo.DeleteOldestWaterLevels(ctx, locationID, 5); err != nil {
		log.Println("failed to delete water level", err)
		return err
	}

	return nil
}
