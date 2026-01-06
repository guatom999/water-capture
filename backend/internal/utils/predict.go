package utils

import (
	"context"
	"fmt"
	"log"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/guatom999/self-boardcast/internal/models"
)

func PredictWaterLevel(imageProcessingDir string, fileName string) (*models.PredictWater, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// fileName := GenerateFileName()

	pythonScript := filepath.Join(imageProcessingDir, "create_waterlevel_file.py")

	cmd := exec.CommandContext(ctx, "python", pythonScript, fileName)
	cmd.Dir = imageProcessingDir

	framOutput, err := cmd.CombinedOutput()

	if ctx.Err() != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return nil, fmt.Errorf("failed to create water level file")
		}
		return nil, fmt.Errorf("operation cancelled: %v", ctx.Err())
	}

	if err != nil {
		log.Printf("Python script error output: %s\n", string(framOutput))
		return nil, fmt.Errorf("failed to run Python script: %v", err)
	}

	time.Sleep(5 * time.Second)

	predictScriptPath := filepath.Join(imageProcessingDir, "predict.py")

	fmt.Println("create file name is", fileName)

	cmd = exec.CommandContext(ctx, "python", predictScriptPath, fileName)
	cmd.Dir = imageProcessingDir
	output, err := cmd.CombinedOutput()

	if ctx.Err() != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return nil, fmt.Errorf("failed to predict water level")
		}
		return nil, fmt.Errorf("operation cancelled: %v", ctx.Err())
	}

	if err != nil {
		return nil, fmt.Errorf("failed to predict water level: %v, output: %s", err, string(output))
	}

	result, err := strconv.ParseFloat(strings.TrimSpace(string(output)), 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse result: %v, output: %s", err, string(output))
	}

	return &models.PredictWater{
		FileName:   fileName,
		WaterLevel: result,
	}, nil

}

func PredictBatch(folderPath string) error {

	return nil
}
