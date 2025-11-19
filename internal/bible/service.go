package bible

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/hackermanpeter/fx-learning/internal/config"
	"github.com/hackermanpeter/fx-learning/internal/database"
	"github.com/hackermanpeter/fx-learning/internal/models"
	"go.uber.org/zap"
)

type BibleService struct {
	logger *zap.Logger
	cfg    *config.Config
	db     *database.Database
}

func NewBibleService(logger *zap.Logger, cfg *config.Config, db *database.Database) *BibleService {
	return &BibleService{
		logger,
		cfg,
		db,
	}
}

func (b *BibleService) GetRandomVerse() fiber.Handler {
	return func(c *fiber.Ctx) error {
		a := fiber.AcquireAgent()
		req := a.Request()
		req.Header.SetMethod(fiber.MethodGet)
		req.SetRequestURI(fmt.Sprintf("%v/data/kjv/random", b.cfg.Bible.URL))

		if err := a.Parse(); err != nil {
			b.logger.Sugar().Errorw("🚀 ~ bible ~ b.GetRandomVerse ~ unable to initialize request agent", "err", err)
			return c.Status(fiber.StatusServiceUnavailable).JSON(models.NewErrorResponse("service currently unavailable"))
		}

		var data models.BibleAPIResponse
		if _, _, errs := a.Struct(&data); len(errs) != 0 {
			b.logger.Sugar().Errorw("🚀 ~ bible ~ b.GetRandomVerse ~ unable to reach bible api", "errs", errs)
			return c.Status(fiber.StatusServiceUnavailable).JSON(models.NewErrorResponse("service currently unavailable"))
		}

		response := models.RandomVerseResponseDTO{
			Translation: data.Translation.Name,
			Book:        data.RandomVerse.Book,
			Chapter:     data.RandomVerse.Chapter,
			Verse:       data.RandomVerse.Verse,
			Text:        data.RandomVerse.Text,
		}

		return c.Status(fiber.StatusOK).JSON(models.NewSuccessResponse(response, "successfully fetched random bible verse"))
	}
}
