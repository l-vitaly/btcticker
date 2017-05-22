package btcticker

type FetcherCfg struct {
	Name      string
	Exchanges []string
}

type Config struct {
	Fetchers []FetcherCfg
}
