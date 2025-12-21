package ufo

import "context"

func (s *UfoService) DeleteSight(ctx context.Context, uuids string) error {
	if err := s.uforepo.DeleteByUUID(ctx, uuids); err != nil {
		return err
	}
	return nil
}
