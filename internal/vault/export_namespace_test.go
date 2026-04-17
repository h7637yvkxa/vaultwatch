package vault

// NewNamespaceCheckerFromParts exposes the constructor for external tests.
func NewNamespaceCheckerFromParts(address, token string) *NamespaceChecker {
	return NewNamespaceChecker(address, token, nil)
}
