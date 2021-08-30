/*
Copyright 2021 The Vitess Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package vttls

import (
	"crypto/x509"
	"crypto/x509/pkix"
	"errors"
	"fmt"
	"io/ioutil"
	"time"
)

type verifyPeerCertificateFunc func([][]byte, [][]*x509.Certificate) error

func revocationCheck(cert *x509.Certificate, crl *pkix.CertificateList) (revoked bool, err error) {
	if crl.HasExpired(time.Now()) {
		return false, errors.New("CRL has expired")
	}

	for _, revoked := range crl.TBSCertList.RevokedCertificates {
		if cert.SerialNumber.Cmp(revoked.SerialNumber) == 0 {
			return true, nil
		}
	}
	return false, nil
}

func verifyPeerCertificateAgainstCRL(crl string) (verifyPeerCertificateFunc, error) {
	body, err := ioutil.ReadFile(crl)
	if err != nil {
		return nil, err
	}

	parsedCRL, err := x509.ParseCRL(body)
	if err != nil {
		return nil, err
	}

	return func(_ [][]byte, verifiedChains [][]*x509.Certificate) error {
		for _, chain := range verifiedChains {
			for _, cert := range chain {
				revoked, err := revocationCheck(cert, parsedCRL)
				if err != nil { // The CRL needs to be updated...
					// TODO: should we be returning errors or soft-failing?
					return errors.New("Error checking CRL")
				} else if revoked {
					return fmt.Errorf("Certificate revoked: CommonName=%v", cert.Subject.CommonName)
				}
			}
		}
		return nil
	}, nil
}
