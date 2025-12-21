package ufo

import (
	"context"
	"testing"

	"github.com/mbakhodurov/examples/week_2/unit_tests/6_unit_test_in_clean_arch/internal/repository/mocks"
	"github.com/stretchr/testify/suite"
)

type ServiceSuite struct {
	suite.Suite

	ctx           context.Context
	ufoRepository *mocks.UFORepository
	service       *UfoService
}

func (s *ServiceSuite) SetupTest() {
	s.ctx = context.Background()

	s.ufoRepository = mocks.NewUFORepository(s.T())

	s.service = NewUfoService(
		s.ufoRepository,
	)
}

func (s *ServiceSuite) TearDownTest() {
}

func TestServiceIntegration(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}
