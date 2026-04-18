package ufo

import "context"

func (s *service) Delete(ctx context.Context, uuid string) error {
	if err := s.ufoRepo.Delete(ctx, uuid); err != nil {
		return err
	}

	_ = s.ufoCache.Delete(ctx, uuid)

	return nil
}
