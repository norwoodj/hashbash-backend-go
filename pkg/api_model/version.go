package api_model

type VersionInfo struct {
	BuildTimestamp string `json:"build_timestamp"`
	GitRevision    string `json:"git_revision"`
	Version        string `json:"version"`
}
