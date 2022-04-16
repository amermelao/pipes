package conduit

type logTestBench interface {
	Fatal(args ...any)
	Fatalf(format string, args ...any)
	Error(args ...any)
	Errorf(format string, args ...any)
	Log(args ...any)
	Logf(format string, args ...any)
}
