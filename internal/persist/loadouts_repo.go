// File: internal/persist/loadouts_repo.go
package persist

import (
	"context"
	"database/sql"
)

type Loadout struct {
	UserID        int64
	ElementID     int64
	ArmyAbility1  int64
	ArmyAbility2  int64
	ArmyAbility3  int64
	ArmyAbility4  int64
	AbilityPawn   int64
	AbilityKnight int64
	AbilityBishop int64
	AbilityRook   int64
	AbilityQueen  int64
	AbilityKing   int64
	Item1         int64
	Item2         int64
	Item3         int64
	Item4         int64
	UpdatedAt     int64
}

type LoadoutsRepo struct {
	db      *sql.DB
	dialect Dialect
}

func NewLoadoutsRepo(db *sql.DB, dialect Dialect) *LoadoutsRepo {
	return &LoadoutsRepo{
		db:      db,
		dialect: dialect,
	}
}

func (r *LoadoutsRepo) Get(ctx context.Context, userID int64) (Loadout, error) {
	if r.db == nil {
		return Loadout{}, ErrNilDB
	}
	row := r.db.QueryRowContext(ctx, selectLoadoutByUserID(r.dialect), userID)
	var loadout Loadout
	if err := row.Scan(
		&loadout.UserID,
		&loadout.ElementID,
		&loadout.ArmyAbility1,
		&loadout.ArmyAbility2,
		&loadout.ArmyAbility3,
		&loadout.ArmyAbility4,
		&loadout.AbilityPawn,
		&loadout.AbilityKnight,
		&loadout.AbilityBishop,
		&loadout.AbilityRook,
		&loadout.AbilityQueen,
		&loadout.AbilityKing,
		&loadout.Item1,
		&loadout.Item2,
		&loadout.Item3,
		&loadout.Item4,
		&loadout.UpdatedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			return Loadout{}, ErrNotFound
		}
		return Loadout{}, err
	}
	return loadout, nil
}

func (r *LoadoutsRepo) Upsert(ctx context.Context, loadout Loadout) error {
	if r.db == nil {
		return ErrNilDB
	}
	_, err := r.db.ExecContext(ctx, upsertLoadout(r.dialect),
		loadout.UserID,
		loadout.ElementID,
		loadout.ArmyAbility1,
		loadout.ArmyAbility2,
		loadout.ArmyAbility3,
		loadout.ArmyAbility4,
		loadout.AbilityPawn,
		loadout.AbilityKnight,
		loadout.AbilityBishop,
		loadout.AbilityRook,
		loadout.AbilityQueen,
		loadout.AbilityKing,
		loadout.Item1,
		loadout.Item2,
		loadout.Item3,
		loadout.Item4,
		loadout.UpdatedAt,
	)
	return err
}

func selectLoadoutByUserID(dialect Dialect) string {
	if dialect == DialectPostgres {
		return `SELECT user_id, element_id, army_ability_1, army_ability_2, army_ability_3, army_ability_4, ability_pawn, ability_knight, ability_bishop, ability_rook, ability_queen, ability_king, item_1, item_2, item_3, item_4, updated_at FROM army_loadouts WHERE user_id = $1`
	}
	return `SELECT user_id, element_id, army_ability_1, army_ability_2, army_ability_3, army_ability_4, ability_pawn, ability_knight, ability_bishop, ability_rook, ability_queen, ability_king, item_1, item_2, item_3, item_4, updated_at FROM army_loadouts WHERE user_id = ?`
}

func upsertLoadout(dialect Dialect) string {
	if dialect == DialectPostgres {
		return `INSERT INTO army_loadouts (user_id, element_id, army_ability_1, army_ability_2, army_ability_3, army_ability_4, ability_pawn, ability_knight, ability_bishop, ability_rook, ability_queen, ability_king, item_1, item_2, item_3, item_4, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17) ON CONFLICT (user_id) DO UPDATE SET element_id = EXCLUDED.element_id, army_ability_1 = EXCLUDED.army_ability_1, army_ability_2 = EXCLUDED.army_ability_2, army_ability_3 = EXCLUDED.army_ability_3, army_ability_4 = EXCLUDED.army_ability_4, ability_pawn = EXCLUDED.ability_pawn, ability_knight = EXCLUDED.ability_knight, ability_bishop = EXCLUDED.ability_bishop, ability_rook = EXCLUDED.ability_rook, ability_queen = EXCLUDED.ability_queen, ability_king = EXCLUDED.ability_king, item_1 = EXCLUDED.item_1, item_2 = EXCLUDED.item_2, item_3 = EXCLUDED.item_3, item_4 = EXCLUDED.item_4, updated_at = EXCLUDED.updated_at`
	}
	return `INSERT INTO army_loadouts (user_id, element_id, army_ability_1, army_ability_2, army_ability_3, army_ability_4, ability_pawn, ability_knight, ability_bishop, ability_rook, ability_queen, ability_king, item_1, item_2, item_3, item_4, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?) ON CONFLICT(user_id) DO UPDATE SET element_id = excluded.element_id, army_ability_1 = excluded.army_ability_1, army_ability_2 = excluded.army_ability_2, army_ability_3 = excluded.army_ability_3, army_ability_4 = excluded.army_ability_4, ability_pawn = excluded.ability_pawn, ability_knight = excluded.ability_knight, ability_bishop = excluded.ability_bishop, ability_rook = excluded.ability_rook, ability_queen = excluded.ability_queen, ability_king = excluded.ability_king, item_1 = excluded.item_1, item_2 = excluded.item_2, item_3 = excluded.item_3, item_4 = excluded.item_4, updated_at = excluded.updated_at`
}
