package locator

import "context"

// Locator tracks player session locations across nodes.
type Locator interface {
	// ID returns the locator implementation identifier.
	ID() string

	// Gate returns the gate node for a user ID.
	Gate(ctx context.Context, uid int64) (string, error)

	// BindGate associates a user ID with a gate node.
	BindGate(ctx context.Context, uid int64, node string) error

	// UnBindGate removes user ID's gate node association.
	UnBindGate(ctx context.Context, uid int64, node string) error

	// Game returns the game node for a user ID.
	Game(ctx context.Context, uid int64) (string, error)

	// BindGame associates a user ID with a game node.
	BindGame(ctx context.Context, uid int64, node string) error

	// UnBindGame removes user ID's game node association.
	UnBindGame(ctx context.Context, uid int64, node string) error
}
