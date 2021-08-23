package v1

type TypeMeta struct {
	Kind       string `json:"kind,omitempty"`
	APIVersion string `json:"apiVersion,omitempty"`
	Success    bool   `json:"success"`
	TraceId    string `json:"traceId,omitempty"`
	Host       string `json:"host,omitempty"`
}

type Pages struct {
	Current  int   `json:"current,omitempty"`
	PageSize int   `json:"pageSize,omitempty"`
	Total    int64 `json:"total,omitempty"`
	Success  bool  `json:"success"`
}
