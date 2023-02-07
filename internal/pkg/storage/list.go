package storage

func (s *DiskStorage) List() ([]ImageRow, error) {
	result := make([]ImageRow, 0, len(s.files))
	for _, v := range s.files {
		result = append(result, v)
	}
	return result, nil

}
