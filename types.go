package main

// Settings struct are the allowed types to use in the policy
type Settings struct {
	// Namespace is a string which select the destination NS to work with
	Namespace string `json:"namespace"`
	// UnsafeNames blacklists the names
	UnsafeNames []string `json:"unsafe_names"`
	// SafeNames whitelists the names
	SafeNames []string `json:"safe_names"`
}
