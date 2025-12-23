package bootstrap

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWithDisabledResource(t *testing.T) {
	t.Run("disables rabbitmq resource", func(t *testing.T) {
		opts := &BootstrapOptions{}
		option := WithDisabledResource("rabbitmq")
		option(opts)

		assert.True(t, opts.IsResourceDisabled("rabbitmq"))
		assert.False(t, opts.IsResourceDisabled("s3"))
	})

	t.Run("disables s3 resource", func(t *testing.T) {
		opts := &BootstrapOptions{}
		option := WithDisabledResource("s3")
		option(opts)

		assert.True(t, opts.IsResourceDisabled("s3"))
		assert.False(t, opts.IsResourceDisabled("rabbitmq"))
	})

	t.Run("disables multiple resources", func(t *testing.T) {
		opts := &BootstrapOptions{}
		WithDisabledResource("rabbitmq")(opts)
		WithDisabledResource("s3")(opts)

		assert.True(t, opts.IsResourceDisabled("rabbitmq"))
		assert.True(t, opts.IsResourceDisabled("s3"))
	})
}

func TestIsResourceDisabled(t *testing.T) {
	t.Run("returns false for nil options", func(t *testing.T) {
		var opts *BootstrapOptions
		assert.False(t, opts.IsResourceDisabled("rabbitmq"))
	})

	t.Run("returns false for nil DisabledResources map", func(t *testing.T) {
		opts := &BootstrapOptions{}
		assert.False(t, opts.IsResourceDisabled("rabbitmq"))
	})

	t.Run("returns false for non-disabled resource", func(t *testing.T) {
		opts := &BootstrapOptions{
			DisabledResources: map[string]bool{
				"rabbitmq": true,
			},
		}
		assert.False(t, opts.IsResourceDisabled("s3"))
	})

	t.Run("returns true for disabled resource", func(t *testing.T) {
		opts := &BootstrapOptions{
			DisabledResources: map[string]bool{
				"rabbitmq": true,
			},
		}
		assert.True(t, opts.IsResourceDisabled("rabbitmq"))
	})
}

func TestWithOptionalResource(t *testing.T) {
	t.Run("marks resource as optional", func(t *testing.T) {
		opts := &BootstrapOptions{}
		option := WithOptionalResource("rabbitmq")
		option(opts)

		assert.True(t, opts.OptionalResources["rabbitmq"])
	})

	t.Run("marks multiple resources as optional", func(t *testing.T) {
		opts := &BootstrapOptions{}
		WithOptionalResource("rabbitmq")(opts)
		WithOptionalResource("s3")(opts)

		assert.True(t, opts.OptionalResources["rabbitmq"])
		assert.True(t, opts.OptionalResources["s3"])
	})
}

func TestDefaultResourceConfig(t *testing.T) {
	config := DefaultResourceConfig()

	t.Run("logger is required and enabled", func(t *testing.T) {
		assert.False(t, config["logger"].Optional)
		assert.True(t, config["logger"].Enabled)
	})

	t.Run("postgresql is required and enabled", func(t *testing.T) {
		assert.False(t, config["postgresql"].Optional)
		assert.True(t, config["postgresql"].Enabled)
	})

	t.Run("mongodb is required and enabled", func(t *testing.T) {
		assert.False(t, config["mongodb"].Optional)
		assert.True(t, config["mongodb"].Enabled)
	})

	t.Run("rabbitmq is optional and enabled", func(t *testing.T) {
		assert.True(t, config["rabbitmq"].Optional)
		assert.True(t, config["rabbitmq"].Enabled)
	})

	t.Run("s3 is optional and enabled", func(t *testing.T) {
		assert.True(t, config["s3"].Optional)
		assert.True(t, config["s3"].Enabled)
	})
}
