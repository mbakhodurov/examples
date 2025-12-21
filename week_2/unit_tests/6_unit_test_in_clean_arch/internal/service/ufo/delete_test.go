package ufo

import "github.com/brianvoe/gofakeit/v7"

func (s *ServiceSuite) TestDeleteSuccess() {
	uuid := gofakeit.UUID()

	s.ufoRepository.On("DeleteByUUID", s.ctx, uuid).Return(nil)

	err := s.service.DeleteSight(s.ctx, uuid)
	s.NoError(err)
}

func (s *ServiceSuite) TestDeleteRepoError() {
	var (
		repoErr = gofakeit.Error()
		uuid    = gofakeit.UUID()
	)

	s.ufoRepository.On("DeleteByUUID", s.ctx, uuid).Return(repoErr)

	err := s.service.DeleteSight(s.ctx, uuid)
	s.Error(err)
	s.ErrorIs(err, repoErr)
}
