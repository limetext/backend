package sublime

import "github.com/limetext/lime-backend/lib/parser"

type Syntax struct {
	l *Language
}

func newSyntax(path string) (*Syntax, error) {
	l, err := Provider.LanguageFromFile(path)
	if err != nil {
		return nil, err
	}
	return &Syntax{l: l}, nil
}

func (s *Syntax) Parser(data string) (parser.Parser, error) {
	l := s.l.copy()
	return &LanguageParser{l: l, data: []rune(data)}, nil
}

func (s *Syntax) Name() string {
	return s.l.Name
}

func (s *Syntax) FileTypes() []string {
	return s.l.FileTypes
}
