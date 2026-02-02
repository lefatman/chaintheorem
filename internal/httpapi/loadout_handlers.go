// File: internal/httpapi/loadout_handlers.go
package httpapi

import (
	"errors"
	"net/http"
	"strings"

	"example.com/mvp-repo/internal/auth"
	"example.com/mvp-repo/internal/persist"
)

type LoadoutInput struct {
	ElementID     int64 `json:"element_id"`
	ArmyAbility1  int64 `json:"army_ability_1"`
	ArmyAbility2  int64 `json:"army_ability_2"`
	ArmyAbility3  int64 `json:"army_ability_3"`
	ArmyAbility4  int64 `json:"army_ability_4"`
	AbilityPawn   int64 `json:"ability_pawn"`
	AbilityKnight int64 `json:"ability_knight"`
	AbilityBishop int64 `json:"ability_bishop"`
	AbilityRook   int64 `json:"ability_rook"`
	AbilityQueen  int64 `json:"ability_queen"`
	AbilityKing   int64 `json:"ability_king"`
	Item1         int64 `json:"item_1"`
	Item2         int64 `json:"item_2"`
	Item3         int64 `json:"item_3"`
	Item4         int64 `json:"item_4"`
}

type loadoutResponse struct {
	UserID        int64 `json:"user_id"`
	ElementID     int64 `json:"element_id"`
	ArmyAbility1  int64 `json:"army_ability_1"`
	ArmyAbility2  int64 `json:"army_ability_2"`
	ArmyAbility3  int64 `json:"army_ability_3"`
	ArmyAbility4  int64 `json:"army_ability_4"`
	AbilityPawn   int64 `json:"ability_pawn"`
	AbilityKnight int64 `json:"ability_knight"`
	AbilityBishop int64 `json:"ability_bishop"`
	AbilityRook   int64 `json:"ability_rook"`
	AbilityQueen  int64 `json:"ability_queen"`
	AbilityKing   int64 `json:"ability_king"`
	Item1         int64 `json:"item_1"`
	Item2         int64 `json:"item_2"`
	Item3         int64 `json:"item_3"`
	Item4         int64 `json:"item_4"`
	UpdatedAt     int64 `json:"updated_at"`
}

func (s *Server) handleLoadout(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		s.handleLoadoutGet(w, r)
	case http.MethodPost:
		s.handleLoadoutPost(w, r)
	default:
		writeError(w, http.StatusMethodNotAllowed, "method_not_allowed", "method not allowed")
	}
}

func (s *Server) handleLoadoutGet(w http.ResponseWriter, r *http.Request) {
	token := bearerToken(r)
	if token == "" {
		writeError(w, http.StatusUnauthorized, "missing_token", "authorization token required")
		return
	}
	session, err := s.auth.ValidateToken(r.Context(), token)
	if err != nil {
		writeAuthError(w, err)
		return
	}
	loadout, err := s.loadouts.Get(r.Context(), session.UserID)
	if err != nil {
		if errors.Is(err, persist.ErrNotFound) {
			writeError(w, http.StatusNotFound, "loadout_not_found", "loadout not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "server_error", "unable to load loadout")
		return
	}
	writeData(w, http.StatusOK, toLoadoutResponse(loadout))
}

func (s *Server) handleLoadoutPost(w http.ResponseWriter, r *http.Request) {
	token := bearerToken(r)
	if token == "" {
		writeError(w, http.StatusUnauthorized, "missing_token", "authorization token required")
		return
	}
	session, err := s.auth.ValidateToken(r.Context(), token)
	if err != nil {
		writeAuthError(w, err)
		return
	}
	var req LoadoutInput
	if err := s.decodeJSON(w, r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_json", "invalid json payload")
		return
	}
	loadout, err := s.loadouts.Update(r.Context(), session.UserID, req)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "server_error", "unable to update loadout")
		return
	}
	writeData(w, http.StatusOK, toLoadoutResponse(loadout))
}

func writeAuthError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, auth.ErrInvalidToken), errors.Is(err, auth.ErrTokenExpired):
		writeError(w, http.StatusUnauthorized, "invalid_token", "invalid or expired token")
	default:
		writeError(w, http.StatusInternalServerError, "server_error", "unable to validate token")
	}
}

func bearerToken(r *http.Request) string {
	header := strings.TrimSpace(r.Header.Get("Authorization"))
	if header == "" {
		return ""
	}
	const prefix = "Bearer "
	if len(header) < len(prefix) || !strings.EqualFold(header[:len(prefix)], prefix) {
		return ""
	}
	return strings.TrimSpace(header[len(prefix):])
}

func toLoadoutResponse(loadout persist.Loadout) loadoutResponse {
	return loadoutResponse{
		UserID:        loadout.UserID,
		ElementID:     loadout.ElementID,
		ArmyAbility1:  loadout.ArmyAbility1,
		ArmyAbility2:  loadout.ArmyAbility2,
		ArmyAbility3:  loadout.ArmyAbility3,
		ArmyAbility4:  loadout.ArmyAbility4,
		AbilityPawn:   loadout.AbilityPawn,
		AbilityKnight: loadout.AbilityKnight,
		AbilityBishop: loadout.AbilityBishop,
		AbilityRook:   loadout.AbilityRook,
		AbilityQueen:  loadout.AbilityQueen,
		AbilityKing:   loadout.AbilityKing,
		Item1:         loadout.Item1,
		Item2:         loadout.Item2,
		Item3:         loadout.Item3,
		Item4:         loadout.Item4,
		UpdatedAt:     loadout.UpdatedAt,
	}
}
