package users

type UserService struct {
	repository UserRepository
}

func NewService(r UserRepository) *UserService {
	return &UserService{
		repository: r,
	}
}

func (s *UserService) Save(user User) (*User, error) {
	res, err := s.repository.Save(user)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *UserService) GetByUserName(username string) (*User, error) {
	res, err := s.repository.GetByUserName(username)
	if err != nil {
		return nil, err
	}

	return res, nil
}
