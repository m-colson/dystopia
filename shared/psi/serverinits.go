package psi

import "crypto/tls"

func Specific[T Server, E any](cb func(Server, E) error) func(m T, e E) error {
	return func(m T, e E) error {
		return cb(m, e)
	}
}

func SpecificTLS[T TLSServer, E any](cb func(TLSServer, E) error) func(m T, e E) error {
	return func(m T, e E) error {
		return cb(m, e)
	}
}

func AddTLSCert(certPath string, keyPath string) func(*tls.Config) error {
	return func(cfg *tls.Config) error {
		cert, err := tls.LoadX509KeyPair(certPath, keyPath)
		if err != nil {
			return err
		}

		cfg.Certificates = append(cfg.Certificates, cert)
		return nil
	}
}
