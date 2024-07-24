package ad

import (
	"ads-service/internal/models"
)

func (s *service) Create(ad *models.Ad) (int, error) {
	int, err := s.store.Ad.Create(ad)
	if err != nil {
		s.logger.Error("unable to create ad", "err", err.Error())
		return 0, err
	}
	
	return int, nil
}
