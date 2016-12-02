package crypto

import (
	"fmt"
	"strings"
	"hash"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"errors"
)

const (
	DigestAlgorithmSHA1   DigestAlgorithm = "sha1"
	DigestAlgorithmSHA256 DigestAlgorithm = "sha256"
	DigestAlgorithmSHA512 DigestAlgorithm = "sha512"
)

func CreateHashFromAlgorithm(algorithm DigestAlgorithm) (hash.Hash, error) {
	switch algorithm {
	case DigestAlgorithmSHA1:
		return sha1.New(), nil
	case DigestAlgorithmSHA256:
		return sha256.New(), nil
	case DigestAlgorithmSHA512:
		return sha512.New(), nil
	}

	return nil, errors.New(fmt.Sprintf("Unrecognized digest algorithm: %s", algorithm))
}

func ParseDigestString(digest string) (Digest, error) {
	pieces := strings.SplitN(digest, ":", 2)

	if len(pieces) == 1 {
		// historically digests were only sha1 and did not include a prefix.
		// continue to support that behavior.
		pieces = []string{"sha1", pieces[0]}
	}

	switch pieces[0] {
	case string(DigestAlgorithmSHA1), string(DigestAlgorithmSHA256), string(DigestAlgorithmSHA512):
		return NewDigest(DigestAlgorithm(pieces[0]), pieces[1]), nil
	default:
		return nil, errors.New(fmt.Sprintf("Unrecognized digest algorithm: %s. Supported algorithms: sha1, sha256, sha512", pieces[0]))
	}

	return nil, errors.New(fmt.Sprintf("Parsing digest: %s", digest))
}

func ParseMultipleDigestString(multipleDigest string) (MultipleDigest, error) {
	pieces := strings.Split(multipleDigest, ";")

	digests := []Digest{}

	for _, digest := range pieces {
		parsedDigest, err := ParseDigestString(digest)
		if err != nil {
			return nil, err
		}

		digests = append(digests, parsedDigest)
	}

	return NewMultipleDigest(digests...), nil
}
