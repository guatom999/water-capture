package services

import (
	"context"
	"database/sql"
	"fmt"
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
	GetAllLocations(ctx context.Context, limit int) ([]models.LocationWithWaterLevelRes, error)
	GetByLocationID(ctx context.Context, id string) ([]*models.WaterLocationDetailRes, error)
	ScheduleGetWaterLevel(ctx context.Context) (*entities.WaterLevel, error)
	CreateWaterLevel(ctx context.Context, req *models.CreateWaterLevelReq) error

	ScheduleDeleteWaterLevel(ctx context.Context, fileName string, locationID int) error
}

func NewWaterLevelService(repo repositories.WaterLevelRepositoryInterface, baseURL string, cfg *config.Config) WaterLevelServiceInterface {
	return &waterLevelService{
		repo:    repo,
		baseURL: baseURL,
		cfg:     cfg,
	}
}

func (s *waterLevelService) GetAllLocations(ctx context.Context, limit int) ([]models.LocationWithWaterLevelRes, error) {
	locations, err := s.repo.GetAll(ctx, limit)
	if err != nil {
		return nil, err
	}

	locationsRes := make([]models.LocationWithWaterLevelRes, 0)

	// for i, location := range locations {
	// 	if locations[i].Image != nil && *locations[i].Image != "" {
	// 		imageURL := utils.BuildImageURL(s.baseURL, *locations[i].Image)
	// 		locations[i].Image = &imageURL
	// 	}
	// 	locations[i].MeasuredAt = utils.ParseTimeToString(location.MeasuredAt)
	// }
	for _, v := range locations {
		locationsRes = append(locationsRes, models.LocationWithWaterLevelRes{
			LocationID:          v.LocationID,
			LocationName:        v.LocationName,
			LocationDescription: v.LocationDescription,
			Latitude:            v.Latitude,
			Longitude:           v.Longitude,
			IsActive:            v.IsActive,
			BankLevel:           v.BankLevel,
			WaterLevelID:        v.WaterLevelID,
			LevelCm:             v.LevelCm,
			Image: func() *string {
				if v.Image != nil && *v.Image != "" {
					imageURL := utils.BuildImageURL(s.baseURL, *v.Image)
					return &imageURL
				}
				return nil
			}(),
			Danger:     v.Danger,
			IsFlooded:  v.IsFlooded,
			MeasuredAt: utils.ParseTimePtrToString(v.MeasuredAt),
			Note:       v.Note,
		})
	}

	return locationsRes, nil
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

func (s *waterLevelService) ScheduleGetWaterLevel(ctx context.Context) (*entities.WaterLevel, error) {

	fileName := utils.GenerateFileName()

	apiResponse := new(models.ThaiWaterAPIResponse)

	if err := utils.Get("https://api-v3.thaiwater.net/api/v1/thaiwater30/provinces/waterlevel?province_code=13", apiResponse); err != nil {
		return nil, err
	}

	if apiResponse.Result != "OK" {
		return nil, fmt.Errorf("API returned non-OK result: %s", apiResponse.Result)
	}

	stationName := apiResponse.Data[0].Station.TeleStationName.TH

	entity := &entities.WaterLevel{
		LocationID: 28,
		LevelCm:    utils.ConvertStringToFloat64(apiResponse.Data[0].WaterlevelMSL),
		Image:      "",
		Danger: func() string {
			if utils.ConvertStringToFloat64(apiResponse.Data[0].WaterlevelMSL) < 1 {
				return "SAFE"
			} else if utils.ConvertStringToFloat64(apiResponse.Data[0].WaterlevelMSL) < 1.29 {
				return "WATCH"
			} else {
				return "DANGER"
			}
		}(),
		IsFlooded: func() bool {
			if utils.ConvertStringToFloat64(apiResponse.Data[0].WaterlevelMSL) < 1.29 {
				return false
			} else {
				return true
			}
		}(),
		Source:     sql.NullString{String: stationName, Valid: true},
		MeasuredAt: utils.ConvertStringToTime(apiResponse.Data[0].WaterlevelDatetime),
		Note:       "get value of waterLevel from cctv of water",
	}

	if err := s.repo.CreateWaterLevel(ctx, entity); err != nil {
		log.Println("failed to create water level in database:", err)
		filePath := s.cfg.App.UploadDir + "/images/" + fileName
		if deleteErr := utils.DeleteFile(filePath); deleteErr != nil {
			log.Printf("failed to cleanup file %s after DB insert error: %v", fileName, deleteErr)
		}
		return nil, err
	}

	return entity, nil
}

func (s *waterLevelService) ScheduleDeleteWaterLevel(ctx context.Context, fileName string, locationID int) error {

	scheduleAt := time.Now().In(func() *time.Location {
		loc, err := time.LoadLocation("Asia/Bangkok")
		if err != nil {
			log.Printf("failed to load location %v", err.Error())
			return nil
		}
		return loc
	}())

	covertLocationID := int64(28)

	if err := s.repo.MarkForDeletion(ctx, covertLocationID, scheduleAt); err != nil {
		log.Println("failed to mark for deletion", err)
		return err
	}

	results, err := s.repo.GetPendingDeletions(ctx)
	if err != nil {
		log.Println("failed to get pending deletions", err)
		return err
	}

	for _, value := range results {
		if err := utils.DeleteFile(s.cfg.App.UploadDir + "/" + value.Image); err != nil {
			log.Println("failed to delete file", value.Image, err)
		}
		if err := s.repo.HardDelete(ctx, value.ID); err != nil {
			log.Println("failed to delete water level", value.ID, err)
		}
	}

	// if err := utils.DeleteFile(s.cfg.App.UploadDir + "/" + fileName); err != nil {
	// 	log.Println("failed to delete file", err)
	// 	return err
	// }

	// if err := s.repo.DeleteOldestWaterLevels(ctx, locationID, 5); err != nil {
	// 	log.Println("failed to delete water level", err)
	// 	return err
	// }

	return nil
}
