package model

// TODO: Server settings:
// - password policy
// - hash algorithm with parameters
// - pepper
type PasswordHashType string

var DefaultPasswordHashParams = PasswordHashParams{
	Type:       PasswordHashArgon2i,
	SaltLength: 32,
	Argon:      &DefaultPasswordHashArgonParams,
}

var DefaultPasswordHashArgonParams = PasswordHashArgonParams{
	Memory:      64 * 1024, // 64MB
	Iterations:  3,
	Parallelism: 2,
	KeyLength:   32,
}

const (
	PasswordHashBcrypt  PasswordHashType = "bcrypt"
	PasswordHashArgon2i PasswordHashType = "argon2i"
)

type PasswordHashParams struct {
	Type       PasswordHashType
	SaltLength uint32
	Argon      *PasswordHashArgonParams
	Bcrypt     *PasswordHashBcryptParams
}

// https://www.alexedwards.net/blog/how-to-hash-and-verify-passwords-with-argon2-in-go
// The Argon2 algorithm accepts a number of configurable parameters
// Memory — The amount of memory used by the algorithm (in kibibytes).
// Iterations — The number of iterations (or passes) over the memory.
// Parallelism — The number of threads (or lanes) used by the algorithm.
// Salt length — Length of the random salt. 16 bytes is recommended for password hashing.
// Key length — Length of the generated key (or password hash). 16 bytes or more is recommended.
//
// The recommended process for choosing the parameters can be paraphrased as follows:
// 1. Set the parallelism and memory parameters to the largest amount you are willing to afford, bearing in mind that you probably don't want to max these out completely unless your machine is dedicated to password hashing.
// 2. Increase the number of iterations until you reach your maximum runtime limit (for example, 500ms).
// 3. If you're already exceeding the your maximum runtime limit with the number of iterations = 1, then you should reduce the memory parameter.
// More details here: https://datatracker.ietf.org/doc/html/draft-irtf-cfrg-argon2-04#section-4
type PasswordHashArgonParams struct {
	Memory      uint32
	Iterations  uint32
	Parallelism uint8
	KeyLength   uint32
}

type PasswordHashBcryptParams struct {
	Cost int
}
