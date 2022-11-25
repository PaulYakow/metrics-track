package v1

import "time"

// Option применяет заданную настройку к репозиторию (Postgre).
type Option func(*Postgre)

// MaxPoolSize задаёт максимальный размер пула подключений.
func MaxPoolSize(size int) Option {
	return func(p *Postgre) {
		p.maxPollSize = size
	}
}

// MaxConnIdleTime задаёт время, после которого бездействующее соединение будет закрыто.
func MaxConnIdleTime(duration time.Duration) Option {
	return func(p *Postgre) {
		p.maxConnIdleTime = duration
	}
}

// MaxConnLifeTime задаёт время с момента создания, после которого соединение будет закрыто.
func MaxConnLifeTime(duration time.Duration) Option {
	return func(p *Postgre) {
		p.maxConnLifeTime = duration
	}
}

// ConnAttempts задаёт количество попыток подключения.
func ConnAttempts(attempts int) Option {
	return func(p *Postgre) {
		p.connAttempts = attempts
	}
}

// ConnTimeout задаёт таймаут между попытками подключения.
func ConnTimeout(timeout time.Duration) Option {
	return func(p *Postgre) {
		p.connTimeout = timeout
	}
}
