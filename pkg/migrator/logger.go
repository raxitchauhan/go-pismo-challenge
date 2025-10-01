package migrator

import "github.com/rs/zerolog/log"

type logger struct{}

func (l *logger) Fatal(v ...any) {
	l.Fatalf("%v", v)
}

func (l *logger) Fatalf(format string, v ...any) {
	log.Fatal().Msgf(format, v...)
}

func (l *logger) Print(v ...any) {
	l.Println(v...)
}

func (l *logger) Println(v ...any) {
	for _, vv := range v {
		l.Printf("%v", vv)
	}
}

func (l *logger) Printf(format string, v ...any) {
	log.Debug().Msgf(format, v...)
}
