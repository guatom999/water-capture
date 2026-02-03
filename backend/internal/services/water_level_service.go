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
			StationID:           v.StationID,
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

	stationID, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}

	waterLevelsRes := make([]*models.WaterLocationDetailRes, 0)

	results, err := s.repo.GetByLocationID(ctx, stationID)
	if err != nil {
		return nil, err
	}

	for _, res := range results {
		waterLevelsRes = append(waterLevelsRes, &models.WaterLocationDetailRes{
			LocationID: res.StationID,
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

	// if err := s.repo.CreateWaterLevel(ctx, &entities.WaterLevel{
	// 	LocationID: req.LocationID,
	// 	LevelCm:    req.LevelCm,
	// 	Image:      req.Image,
	// 	Danger:     req.Danger,
	// 	IsFlooded:  req.IsFlooded,
	// 	MeasuredAt: req.MeasuredAt,
	// 	Note:       req.Note,
	// }); err != nil {
	// 	return err
	// }

	return nil
}

func (s *waterLevelService) ScheduleGetWaterLevel(ctx context.Context) (*entities.WaterLevel, error) {

	fileName := utils.GenerateFileName()

	_ = fileName

	// apiResponses := map[int]models.ThaiWaterAPIResponse{}
	apiResponses := make([]models.ThaiWaterAPIResponse, 0)

	provienceSet := []string{"13"}

	for _, v := range provienceSet {
		apiResponse := models.ThaiWaterAPIResponse{}
		if err := utils.Get("https://api-v3.thaiwater.net/api/v1/thaiwater30/provinces/waterlevel?province_code="+v, &apiResponse); err != nil {
			return nil, err
		}
		apiResponses = append(apiResponses, apiResponse)
	}

	// locations := make([]*entities.Location, 0)
	// for _, v := range apiResponses[0].Data {
	// 	locations = append(locations, &entities.Location{
	// 		StationID:   int64(v.Station.ID),
	// 		Name:        v.Station.TeleStationName.TH,
	// 		Description: v.Station.TeleStationName.TH,
	// 		ProvinceID: func(a any) int {
	// 			provinceCode, _ := strconv.Atoi(v.Geocode.ProvinceCode)
	// 			return provinceCode
	// 		}(v.Geocode.ProvinceCode),
	// 		Latitude:  v.Station.TeleStationLat,
	// 		Longitude: v.Station.TeleStationLong,
	// 		IsActive:  true,
	// 		BankLevel: v.Station.MinBank,
	// 	})
	// }

	// if err := s.repo.CreateStationLocation(ctx, locations); err != nil {
	// 	return nil, err
	// }

	if apiResponses[0].Result != "OK" {
		return nil, fmt.Errorf("API returned non-OK result: %s", apiResponses[0].Result)
	}

	waterLevels := make([]*entities.WaterLevel, 0)

	for _, values := range apiResponses {
		for _, value := range values.Data {
			waterLevels = append(waterLevels, &entities.WaterLevel{
				StationID: int64(value.Station.ID),
				LevelCm:   utils.ConvertStringToFloat64(value.WaterlevelMSL),
				Image:     "",
				Danger: func() string {
					if utils.ConvertStringToFloat64(value.WaterlevelMSL) < value.Station.MinBank-0.2 {
						return "SAFE"
					} else if utils.ConvertStringToFloat64(value.WaterlevelMSL) < value.Station.MinBank {
						return "WATCH"
					} else {
						return "DANGER"
					}
				}(),
				IsFlooded: func() bool {
					if utils.ConvertStringToFloat64(value.WaterlevelMSL) < value.Station.MinBank {
						return false
					} else {
						return true
					}
				}(),
				Source:     sql.NullString{String: value.Station.TeleStationName.TH, Valid: true},
				MeasuredAt: utils.ConvertStringToTime(value.WaterlevelDatetime),
				Note:       "get value of waterLevel from cctv of water",
			})
		}

	}

	if err := s.repo.CreateWaterLevel(ctx, waterLevels); err != nil {
		log.Println("failed to create water level in database:", err)
		filePath := s.cfg.App.UploadDir + "/images/" + fileName
		if deleteErr := utils.DeleteFile(filePath); deleteErr != nil {
			log.Printf("failed to cleanup file %s after DB insert error: %v", fileName, deleteErr)
		}
		return nil, err
	}

	return nil, nil
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
