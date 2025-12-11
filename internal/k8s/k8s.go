package k8s

// Client is a tiny wrapper to represent a Kubernetes client for tests
type Client struct {
	Config string
}

// InitClient returns a simple Client configured with the provided cfg string.
func InitClient(cfg string) (*Client, error) {
	return &Client{Config: cfg}, nil
}
