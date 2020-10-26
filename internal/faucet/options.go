package faucet

const (
	DefaultAppCli        = "gaiad"
	DefaultKeyName       = "faucet"
	DefaultDenom         = "uatom"
	DefaultCreditAmount  = 10000000
	DefaultMaximumCredit = 100000000
)

func defaultOptions() *Options {
	return &Options{
		AppCli:       DefaultAppCli,
		KeyName:      DefaultKeyName,
		Denom:        DefaultDenom,
		CreditAmount: DefaultCreditAmount,
		MaxCredit:    DefaultMaximumCredit,
	}
}

type Options struct {
	AppCli          string
	KeyringPassword string
	KeyName         string
	KeyMnemonic     string
	Denom           string
	CreditAmount    uint64
	MaxCredit       uint64
}

type Option func(*Options)

func CliName(s string) Option {
	return func(opts *Options) {
		opts.AppCli = s
	}
}

func KeyringPassword(s string) Option {
	return func(opts *Options) {
		opts.KeyringPassword = s
	}
}

func KeyName(s string) Option {
	return func(opts *Options) {
		opts.KeyName = s
	}
}

func WithMnemonic(s string) Option {
	return func(opts *Options) {
		opts.KeyMnemonic = s
	}
}

func Denom(s string) Option {
	return func(opts *Options) {
		opts.Denom = s
	}
}

func CreditAmount(v uint64) Option {
	return func(opts *Options) {
		opts.CreditAmount = v
	}
}

func MaxCredit(v uint64) Option {
	return func(opts *Options) {
		opts.MaxCredit = v
	}
}
