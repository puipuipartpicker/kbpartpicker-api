package di

import iDI "github.com/puipuipartpicker/kbpartpicker/api/pkg/di"


func (s *Server) setupRoutes() {
	logger := iDI.GetContextLogger()
	s.server.Use(s.handleError)

	s.server.Use(logger.Access)

	s.server.Get("/", s.healthCheck)

	{
		v1 := s.server.Group("/v1")

		// s.installBot(v1, &honda.Honda{})
		// s.installBot(v1, &curves.Curves{})
		// s.installBot(v1, &lacoco.Lacoco{})
		// s.installBot(v1, &james.James{})
		// s.installBot(v1, &rolandbl.Rolandbl{})
		// s.installBot(v1, &agaskinwoman.AgaskinWoman{})
		// s.installBot(v1, &epireserve.EpiReserve{})
		// s.installBot(v1, &lava.Lava{})
		// s.installBot(v1, &bbt.Bbt{})
		// s.installBot(v1, &sasala.Sasala{})
		// s.installBot(v1, &kidsduo.KidsDuo{})
		// s.installBot(v1, &freya.Freya{})
		// s.installBot(v1, &workout.Workout{})
		// s.installBot(v1, &kireimo.Kireimo{})
		// s.installBot(v1, &chickengym.Chickengym{})

		// v1.Use("/swagger", filesystem.New(filesystem.Config{Root: docs.SwaggerAPI()}))
	}
}