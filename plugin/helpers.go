package secretsengine

import "fmt"

// copied from vault/login_mfa.go in main vault repo
// uuidRegex crafts a regex for use in URL paths, somewhat similar to framework.GenericNameRegex, but only accepting
// UUIDs, and only lowercase UUIDs at that.
// It is currently exclusively used for the method_id parameter for MFA methods.
// Think twice before making use of it in any other context, as restricting the valid input in the URL regex results in
// an "unsupported path" error, given input which does not match the regex, which is a fairly unclear way to report an
// invalid parameter value, unless the person seeing the error has an excellent understanding of Vault URL routing.
func uuidRegex(name string) string {
	return fmt.Sprintf("(?P<%s>[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12})", name)
}
