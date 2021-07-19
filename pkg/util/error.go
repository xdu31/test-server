package util

import (
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ApiError struct {
	s *status.Status
}

func (e *ApiError) Error() string { return e.s.Message() }

func AlreadyExistsErr(format string, args ...interface{}) error {
	return &ApiError{s: status.Newf(codes.AlreadyExists, format, args...)}
}

func InvalidArgumentErr(format string, args ...interface{}) error {
	return &ApiError{s: status.Newf(codes.InvalidArgument, format, args...)}
}

func FailedPreconditionErr(format string, args ...interface{}) error {
	return &ApiError{s: status.Newf(codes.FailedPrecondition, format, args...)}
}

func NotFoundErr(format string, args ...interface{}) error {
	return &ApiError{s: status.Newf(codes.NotFound, format, args...)}
}

func InternalErr(format string, args ...interface{}) error {
	return &ApiError{s: status.Newf(codes.Internal, format, args...)}
}

func PermissionDeniedErr(format string, args ...interface{}) error {
	return &ApiError{s: status.Newf(codes.PermissionDenied, format, args...)}
}

func DuplicateNameErr(name interface{}) error {
	return &ApiError{s: status.Newf(codes.AlreadyExists, "Cannot use duplicate name %q", name)}
}

func UnimplementedErr(format string, args ...interface{}) error {
	return &ApiError{s: status.Newf(codes.PermissionDenied, format, args...)}
}

func UnauthenticatedErr(format string, args ...interface{}) error {
	return &ApiError{s: status.Newf(codes.Unauthenticated, format, args...)}
}

func IsNotFound(e error) bool {
	apiErr, ok := e.(*ApiError)
	return ok && apiErr.s.Code() == codes.NotFound
}

func ErrStatus(e error) error {
	var err error
	if apiErr, ok := e.(*ApiError); ok {
		err = apiErr.s.Err()
	} else {
		err = status.Error(codes.Internal, "Internal Server Error")
		log.Errorf("Internal Server Error: %v", e)
	}

	return err
}
