package metrics

type InternetStatus struct {
	status bool
	err    error
}

func NewInternetStatus(status bool, err error) InternetStatus {
	return InternetStatus{status: status, err: err}
}

func (s InternetStatus) MetricValue() float64 {
	if s.status {
		return 1
	}

	return 0
}

func (s InternetStatus) String() string {
	if s.status {
		return "UP"
	}

	return "DOWN"
}

func (s InternetStatus) IsUp() bool {
	return s.status
}

func (s InternetStatus) Error() error {
	return s.err
}
