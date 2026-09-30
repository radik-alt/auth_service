package middleware

import (
	"context"
	"net/http"

	"google.golang.org/grpc/metadata"
)

// UserIDMetadataKey — под этим ключом сервисы получают id пользователя, проверенный gateway.
const UserIDMetadataKey = "x-user-id"

func UserIDMetadata(_ context.Context, r *http.Request) metadata.MD {
	userID, ok := UserIDFromContext(r.Context())
	if !ok {
		return nil
	}
	return metadata.Pairs(UserIDMetadataKey, userID)
}

// BlockIncomingHeaders — для runtime.WithIncomingHeaderMatcher:
// не пробрасывает в сервисы metadata (headers)
func BlockIncomingHeaders(string) (string, bool) {
	return "", false
}
