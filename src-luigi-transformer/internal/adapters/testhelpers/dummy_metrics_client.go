package testhelpers

type DummyMetricsClient struct{}

func (DummyMetricsClient) Histogram(_ string, _ float64, _ []string) {}
func (DummyMetricsClient) Count(_ string, _ int64, _ []string)       {}
