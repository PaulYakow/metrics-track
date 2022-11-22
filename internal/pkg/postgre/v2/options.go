package v2

import "time"

// Option применяет заданную настройку к репозиторию (Postgre).
type Option func(*Postgre)

// MaxOpenConn задаёт максимальное количество подключений к БД
func MaxOpenConn(size int) Option {
	return func(p *Postgre) {
		p.maxOpenConn = size
	}
}

// MaxIdleConn задаёт максимальное количество бездействующих подключений к БД
func MaxIdleConn(size int) Option {
	return func(p *Postgre) {
		p.maxIdleConn = size
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
		p.maxConnLifetime = duration
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
