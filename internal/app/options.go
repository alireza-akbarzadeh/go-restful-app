package app

import "time"

// Option is a functional option for configuring the App
type Option func(*App)

// WithPort sets a custom port
func WithPort(port int) Option {
	return func(a *App) {
		a.config.Port = port
	}
}

// WithJWTSecret sets a custom JWT secret
func WithJWTSecret(secret string) Option {
	return func(a *App) {
		a.config.JWTSecret = secret
	}
}

// WithDatabaseURL sets a custom database URL
func WithDatabaseURL(url string) Option {
	return func(a *App) {
		a.config.DatabaseURL = url
	}
}

// WithIdleTimeout sets a custom idle timeout
func WithIdleTimeout(d time.Duration) Option {
	return func(a *App) {
		a.config.IdleTimeout = d
	}
}

// WithReadTimeout sets a custom read timeout
func WithReadTimeout(d time.Duration) Option {
	return func(a *App) {
		a.config.ReadTimeout = d
	}
}

// WithWriteTimeout sets a custom write timeout
func WithWriteTimeout(d time.Duration) Option {
	return func(a *App) {
		a.config.WriteTimeout = d
	}
}

// WithShutdownTimeout sets a custom shutdown timeout
func WithShutdownTimeout(d time.Duration) Option {
	return func(a *App) {
		a.config.ShutdownTimeout = d
	}
}
