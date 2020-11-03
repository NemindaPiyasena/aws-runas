package config

type chainLoader struct {
	loaders []Loader
}

// NewChainLoader returns a Loader which will resolve configuration and credentials according to the provided array
// of Loaders.  The loaders are consulted sequentially, first to last.
func NewChainLoader(chain []Loader) *chainLoader {
	return &chainLoader{loaders: chain}
}

// Config will build an AwsConfig object using values looked up via the array of Loaders given to the constructor.
// If an error occurs, the next loader in the chain is consulted until the end of the array.  As such, this method will
// never return an error, but is required to satisfy the Loader interface.
//
// Values retrieved via the various loaders are merged using the AwsConfig.MergeIn() method
func (l *chainLoader) Config(profile string, sources ...interface{}) (*AwsConfig, error) {
	c := new(AwsConfig)

	for _, ldr := range l.loaders {
		cf, err := ldr.Config(profile, sources...)
		if err != nil {
			logger.Debugf("error loading configuration: %v", err)
			continue
		}
		c.MergeIn(cf)
	}

	return c, nil
}

// Credentials will build an AwsCredentials object using values looked up via the array of Loaders given to the constructor.
// If an error occurs, the next loader in the chain is consulted until the end of the array.  As such, this method will
// never return an error, but is required to satisfy the Loader interface.
//
// Values retrieved via the various loaders are merged using the AwsCredentials.MergeIn() method
func (l *chainLoader) Credentials(profile string, sources ...interface{}) (*AwsCredentials, error) {
	c := new(AwsCredentials)

	for _, ldr := range l.loaders {
		cr, err := ldr.Credentials(profile, sources...)
		if err != nil {
			logger.Debugf("error loading credentials: %v", err)
			continue
		}
		c.MergeIn(cr)
	}

	return c, nil
}
