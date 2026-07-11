package pubsub

import "time"

type Config struct {
	App       string // app namespace prefix
	Namespace string

	HealthInterval time.Duration
	SendTimeout    time.Duration
	ChannelSize    int
}

type SubscribeConfig struct {
	topics         []string
	app            string
	namespace      string
	healthInterval time.Duration
	sendTimeout    time.Duration
	channelSize    int
}

type SubscribeOption interface {
	Apply(*SubscribeConfig)
}

type SubscribeOptionFunc func(*SubscribeConfig)

func (f SubscribeOptionFunc) Apply(config *SubscribeConfig) {
	f(config)
}

type PublishConfig struct {
	app       string
	namespace string
}

type PublishOption interface {
	Apply(*PublishConfig)
}

type PublishOptionFunc func(*PublishConfig)

func (f PublishOptionFunc) Apply(config *PublishConfig) {
	f(config)
}

func formatTopic(app, ns, topic string) string {
	return app + ":" + ns + ":" + topic
}

func WithChannelNamespace(value string) SubscribeOption {
	return SubscribeOptionFunc(func(c *SubscribeConfig) {
		c.namespace = value
	})
}

func WithPublishNamespace(value string) PublishOption {
	return PublishOptionFunc(func(c *PublishConfig) {
		c.namespace = value
	})
}
