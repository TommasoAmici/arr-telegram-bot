package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	imdb "github.com/tommasoamici/arr-telegram-bot/pkg/imdb"
	tvdb "github.com/tommasoamici/arr-telegram-bot/pkg/tdvb"
	"golift.io/starr"
	"golift.io/starr/radarr"
	"golift.io/starr/sonarr"
)

type Config struct {
	TelegramToken          string
	RadarrHost             string
	RadarrToken            string
	RadarrRoot             string
	RadarrQualityProfileID int64
	SonarrHost             string
	SonarrToken            string
	SonarrRoot             string
	SonarrQualityProfileID int64
}

var cfg Config

func initConfig() {
	slog.Info("reading configuration from environment")
	cfg = Config{
		TelegramToken: os.Getenv("TELEGRAM_TOKEN"),
		RadarrHost:    os.Getenv("RADARR_HOST"),
		RadarrToken:   os.Getenv("RADARR_TOKEN"),
		RadarrRoot:    os.Getenv("RADARR_ROOT"),
		SonarrHost:    os.Getenv("SONARR_HOST"),
		SonarrToken:   os.Getenv("SONARR_TOKEN"),
		SonarrRoot:    os.Getenv("SONARR_ROOT"),
	}
	radarrQP, err := strconv.ParseInt(os.Getenv("RADARR_QUALITY_PROFILE_ID"), 10, 64)
	if err != nil {
		slog.Error("RADARR_QUALITY_PROFILE_ID must be an integer")
		os.Exit(1)
	}
	cfg.RadarrQualityProfileID = radarrQP
	sonarrQP, err := strconv.ParseInt(os.Getenv("SONARR_QUALITY_PROFILE_ID"), 10, 64)
	if err != nil {
		slog.Error("SONARR_QUALITY_PROFILE_ID must be an integer")
		os.Exit(1)
	}
	cfg.SonarrQualityProfileID = sonarrQP
}

func addMovieToRadarr(imdbID string) error {
	slog.Info("adding movie to radarr", "imdb_id", imdbID)

	c := starr.New(cfg.RadarrToken, cfg.RadarrHost, 5*time.Second)
	r := radarr.New(c)
	mov, err := r.LookupIMDB(imdbID)
	if err != nil {
		slog.Error("failed to lookup movie in IMDB")
		return err
	}
	addOptions := radarr.AddMovieOptions{
		SearchForMovie: true,
		Monitor:        "movieOnly",
	}
	movie := radarr.AddMovieInput{
		Monitored:           true,
		MinimumAvailability: radarr.AvailabilityReleased,
		AddOptions:          &addOptions,
		QualityProfileID:    cfg.RadarrQualityProfileID,
		RootFolderPath:      cfg.RadarrRoot,
		TmdbID:              mov.TmdbID,
	}
	_, err = r.AddMovie(&movie)
	return err
}

func addSeriesToSonarr(imdbID string) error {
	slog.Info("adding series to sonarr", "imdb_id", imdbID)

	tvdb, err := tvdb.LookupFromIMDBId(imdbID)
	if err != nil {
		return err
	}

	c := starr.New(cfg.SonarrToken, cfg.SonarrHost, 5*time.Second)
	s := sonarr.New(c)

	addOptions := sonarr.AddSeriesOptions{
		SearchForMissingEpisodes: true,
		Monitor:                  sonarr.MonitorPilot,
	}
	toAdd := sonarr.AddSeriesInput{
		Title:            "tbd",
		Monitored:        true,
		SeasonFolder:     true,
		TvdbID:           tvdb.Series.SeriesID,
		QualityProfileID: cfg.SonarrQualityProfileID,
		RootFolderPath:   cfg.SonarrRoot,
		AddOptions:       &addOptions,
	}
	_, err = s.AddSeries(&toAdd)
	return err
}

func handler(ctx context.Context, b *bot.Bot, update *models.Update) {
	for _, e := range update.Message.Entities {
		var url string
		if e.Type == models.MessageEntityTypeTextLink {
			url = e.URL
		}
		if e.Type == models.MessageEntityTypeURL {
			url = update.Message.Text[e.Offset : e.Offset+e.Length]
		}
		if url == "" {
			continue
		}
		slog.Info("url found", "url", url)

		imdb, err := imdb.LookupIMDB(url)
		if err != nil {
			slog.Error("failed to lookup url in imdb", "err", err)
		}
		if imdb.IsMovie {
			err = addMovieToRadarr(imdb.ID)
			if err != nil {
				slog.Error("failed to add movie to radarr", "err", err)
			}
		} else {
			err = addSeriesToSonarr(imdb.ID)
			if err != nil {
				slog.Error("failed to add series to sonarr", "err", err)
			}
		}
	}
}

func main() {
	initConfig()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	opts := []bot.Option{
		bot.WithDefaultHandler(handler),
	}

	b, err := bot.New(cfg.TelegramToken, opts...)
	if err != nil {
		panic(err)
	}

	slog.Info("starting bot")
	b.Start(ctx)
}
