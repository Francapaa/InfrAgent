package sdk

//this is the folder which the users can install

type ActionFunc func(target string, params map[string]interface{}) error

type AgentSDK struct {
	apiKey     string
	actions    map[string]ActionFunc
	BackendURL string
}
