// operations/operation.go
package operations

type Operation interface {
	Execute() error
}
