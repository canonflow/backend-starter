package security

import "github.com/canonflow/backend-starter/internal/app/security/repository"

func ResolveSecurity() repository.ISecurity {
	return repository.NewAESCBCSecurity()
}
