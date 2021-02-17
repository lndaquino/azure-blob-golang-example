package errors

import "errors"

var ErrNotFound = errors.New("Not found")

var ErrConflict = errors.New("Conflict")

var ErrInvalidPayload = errors.New("Invalid payload")

var ErrInternalServer = errors.New("Internal Server")

var ErrLoginAlreadyInUse = errors.New("Login already in use")

var ErrLoginOrCNPJCPFAlreadyInUse = errors.New("Login or CNPJ/CPF already in use")

var ErrCNPJCPFAlreadyInUse = errors.New("CNPJ/CPF already in use")

var InvalidCNPJorCPFPair = errors.New("Invalid cnpj/instituicao or cpf/responsavel pair")

var InvalidEntity = errors.New("Invalid entity")

var InvalidEmailOrPassword = errors.New("Invalid email or password")

var ExpiredToken = errors.New("Expired token")

var InvalidClienteID = errors.New("Invalid Cliente ID")

var InvalidPasswordConfirmation = errors.New("Invalid password confirmation")

var ErrNeedToValidateEmail = errors.New("Cliente need to validate email before login")

var ErrNoToken = errors.New("No authorization token provided")

var ErrInvalidToken = errors.New("Invalid token provided")

var ErrInvalidCredentials = errors.New("Invalid credentials")
