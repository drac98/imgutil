package index

import (
	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/types"
)

type Option func(*Options) error

type Options struct {
	keychain          authn.Keychain
	xdgPath, repoName string
	insecure          bool
	format            types.MediaType
}

func (o *Options) Keychain() authn.Keychain {
	return o.keychain
}

func (o *Options) XDGRuntimePath() string {
	return o.xdgPath
}

func (o *Options) RepoName() string {
	return o.repoName
}

func (o *Options) Insecure() bool {
	return o.insecure
}

func (o *Options) Format() types.MediaType {
	return o.format
}

func WithKeychain(keychain authn.Keychain) Option {
	return func(o *Options) error {
		o.keychain = keychain
		return nil
	}
}

func WithXDGRuntimePath(xdgPath string) Option {
	return func(o *Options) error {
		o.xdgPath = xdgPath
		return nil
	}
}

func WithRepoName(repoName string) Option {
	return func(o *Options) error {
		if o.insecure {
			_, err := name.ParseReference(repoName, name.Insecure, name.WeakValidation)
			if err != nil {
				return err
			}
		} else {
			_, err := name.ParseReference(repoName, name.WeakValidation)
			if err != nil {
				return err
			}
		}
		o.repoName = repoName
		return nil
	}
}

func WithInsecure(insecure bool) Option {
	return func(o *Options) error {
		o.insecure = insecure
		return nil
	}
}

func WithFormat(format types.MediaType) Option {
	return func(o *Options) error {
		o.format = format
		return nil
	}
}
