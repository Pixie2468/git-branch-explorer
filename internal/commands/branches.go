package commands

func (s *Service) Switch(branchName string) error {
	_, err := s.git(
		"switch",
		"--",
		branchName,
	)
	return err
}

func (s *Service) List() ([]byte, error) {
	return s.git(
		"branch",
		"--format=%(refname:short)",
	)
}

func (s *Service) GetUpstream() ([]byte, error) {
	return s.git(
		"for-each-ref",
		"--format={\"branch\":\"%(refname:short)\",\"track\":\"%(upstream:track)\"}",
		"refs/heads",
	)
}
