package main

type Settings struct {
	Namespace   string   `json: "namespace"`
	UnsafeNames []string `json:"unsafe_names"`
	SafeNames   []string `json:"safe_names"`
}
