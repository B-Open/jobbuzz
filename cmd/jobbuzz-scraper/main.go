package main

import (
	"fmt"
	"os"

	"github.com/b-open/jobbuzz/internal/config"
	"github.com/b-open/jobbuzz/pkg/scraper"
	"github.com/b-open/jobbuzz/pkg/service"
	"github.com/mattn/go-isatty"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	// set pretty print if using terminal
	isTerm := isatty.IsTerminal(os.Stdout.Fd())
	if isTerm {
		log.Logger = log.Output(
			zerolog.ConsoleWriter{
				Out:     os.Stderr,
				NoColor: false,
			},
		)
	}

	configuration, err := config.LoadConfig("../../")

	if err != nil {
		log.Fatal().Err(err).Msg("Fail to load db config")
	}

	db, err := configuration.GetDb()

	if err != nil {
		log.Fatal().Err(err).Msg("Fail to get db connection")
	}

	service := service.Service{DB: db}

	// Scrape JobCenter

	log.Info().Msg("Fetching jobs from JobCenter")
	jobCentreScraper := scraper.NewJobCentreScraper()
	jobs, companies, err := jobCentreScraper.ScrapeJobs()
	if err != nil {
		log.Error().Err(err).Msg("Fail to scrape jobs from jobcenter")
	} else {
		log.Info().Msgf("Found %d jobs", len(jobs))

		err := service.CreateJobsAndCompanies(jobs, companies)
		if err != nil {
			log.Error().Err(err).Msg("Failed to add jobs and companies to database")
		}
	}

	// Scrape Bruneida

	log.Info().Msg("Fetching jobs from Bruneida")
	bruneidaScraper := scraper.NewBruneidaScraper()
	jobs, err = bruneidaScraper.ScrapeJobs()
	if err != nil {
		log.Error().Err(err).Msg("Fail to scrape jobs from Bruneida")
	} else {
		log.Info().Msgf("Found %d jobs", len(jobs))

		_, err = service.CreateJobs(service.DB, jobs)
		if err != nil {
			log.Error().Err(err).Str("job", fmt.Sprintf("%+v", jobs))
		}
	}
}
