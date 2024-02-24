package feed

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func clearImage(filePath string) error {
	absolutePath, err := filepath.Abs(filePath)
	if err != nil {
		return fmt.Errorf("Error converting to absolute path: %w", err)
	}

	if _, err := os.Stat(absolutePath); err == nil {
		if err := os.Remove(absolutePath); err != nil {
			return fmt.Errorf("Error removing file: %w", err)
		} else {
			slog.Info(fmt.Sprintf("%s removed successfully", absolutePath))
		}
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("Error checking file existence: %w", err)
	}

	return nil
}

func getUserIdFromLocals(c *fiber.Ctx) (primitive.ObjectID, error) {
	userIdInterface := c.Locals("user_id")
	if userIdInterface != nil {
		userId, ok := userIdInterface.(string)
		if !ok {
			return primitive.ObjectID{}, fmt.Errorf("type assertion failed")
		}
		userIdObj, err := primitive.ObjectIDFromHex(userId)
		if err != nil {
			return primitive.ObjectID{}, err
		}
		return userIdObj, nil
	}
	return primitive.NewObjectID(), fmt.Errorf("user not found")
}
