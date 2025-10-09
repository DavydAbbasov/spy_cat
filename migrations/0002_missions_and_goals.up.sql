CREATE TABLE IF NOT EXISTS missions (
  id          BIGSERIAL PRIMARY KEY,
  title       TEXT        NOT NULL,
  description TEXT        DEFAULT '',
  status      TEXT        NOT NULL DEFAULT 'planned',
  cat_id      BIGINT      NULL,
  created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at  TIMESTAMPTZ NOT NULL DEFAULT now(),

  CONSTRAINT chk_mission_status CHECK (status IN ('planned','active','completed')),

  CONSTRAINT fk_missions_cat
    FOREIGN KEY (cat_id) REFERENCES cats(id) ON DELETE SET NULL
);
CREATE TABLE IF NOT EXISTS mission_goals (
  id          BIGSERIAL PRIMARY KEY,
  mission_id  BIGINT      NOT NULL REFERENCES missions(id) ON DELETE CASCADE,
  name        TEXT        NOT NULL,
  country     TEXT        NOT NULL,
  notes       TEXT        DEFAULT '',
  status      TEXT        NOT NULL DEFAULT 'todo',
  created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at  TIMESTAMPTZ NOT NULL DEFAULT now(),

  CONSTRAINT chk_goal_status CHECK (status IN ('todo','done'))
);
CREATE INDEX IF NOT EXISTS idx_missions_status_created_at ON missions(status, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_missions_cat_id           ON missions(cat_id);
CREATE INDEX IF NOT EXISTS idx_mission_goals_mid         ON mission_goals(mission_id);
CREATE INDEX IF NOT EXISTS idx_mission_goals_status      ON mission_goals(status);