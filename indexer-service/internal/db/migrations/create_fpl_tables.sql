-- =====================================================
-- Season-Proof FPL Database Schema
-- Using surrogate keys with season isolation
-- =====================================================

-- Drop existing tables (in reverse order of dependencies)
DROP TABLE IF EXISTS player_past_seasons CASCADE;
DROP TABLE IF EXISTS player_gameweek_explain CASCADE;
DROP TABLE IF EXISTS player_gameweek_stats CASCADE;
DROP TABLE IF EXISTS automatic_subs CASCADE;
DROP TABLE IF EXISTS chip_usage CASCADE;
DROP TABLE IF EXISTS manager_picks CASCADE;
DROP TABLE IF EXISTS manager_season_history CASCADE;
DROP TABLE IF EXISTS manager_gameweek_history CASCADE;
DROP TABLE IF EXISTS managers CASCADE;
DROP TABLE IF EXISTS fixture_stats CASCADE;
DROP TABLE IF EXISTS fixtures CASCADE;
DROP TABLE IF EXISTS player_rankings CASCADE;
DROP TABLE IF EXISTS player_expected_stats CASCADE;
DROP TABLE IF EXISTS player_ict_stats CASCADE;
DROP TABLE IF EXISTS player_season_stats CASCADE;
DROP TABLE IF EXISTS player_costs CASCADE;
DROP TABLE IF EXISTS players CASCADE;
DROP TABLE IF EXISTS teams CASCADE;
DROP TABLE IF EXISTS seasons CASCADE;

-- =====================================================
-- Core Reference Tables
-- =====================================================

CREATE TABLE seasons (
                         season_id SERIAL PRIMARY KEY,
                         season_name VARCHAR(20) UNIQUE NOT NULL,  -- e.g., '2024/25'
                         start_date DATE NOT NULL,
                         end_date DATE,
                         is_current BOOLEAN DEFAULT FALSE,
                         created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
                         updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

COMMENT ON TABLE seasons IS 'Premier League seasons';
COMMENT ON COLUMN seasons.is_current IS 'Only one season should have is_current = true';

-- =====================================================
-- Teams
-- =====================================================

CREATE TABLE teams (
                       id SERIAL PRIMARY KEY,
                       season_id INTEGER NOT NULL REFERENCES seasons(season_id) ON DELETE CASCADE,
                       fpl_team_id INTEGER NOT NULL,
                       team_code INTEGER NOT NULL,
                       name VARCHAR(100) NOT NULL,
                       short_name VARCHAR(50) NOT NULL,
                       strength INTEGER,
                       form VARCHAR(10),
                       position INTEGER,
                       points INTEGER,
                       played INTEGER,
                       win INTEGER,
                       draw INTEGER,
                       loss INTEGER,
                       team_division INTEGER,
                       unavailable BOOLEAN DEFAULT FALSE,
                       pulse_id INTEGER,
                       strength_overall_home INTEGER,
                       strength_overall_away INTEGER,
                       strength_attack_home INTEGER,
                       strength_attack_away INTEGER,
                       strength_defence_home INTEGER,
                       strength_defence_away INTEGER,
                       created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
                       updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
                       UNIQUE(fpl_team_id, season_id),
                       UNIQUE(team_code, season_id)
);

COMMENT ON TABLE teams IS 'Premier League teams by season';
COMMENT ON COLUMN teams.fpl_team_id IS 'FPL API team ID (resets each season)';
COMMENT ON COLUMN teams.team_code IS 'FPL API team code (resets each season)';

CREATE INDEX idx_teams_season ON teams(season_id);
CREATE INDEX idx_teams_fpl_id ON teams(fpl_team_id, season_id);

-- =====================================================
-- Players
-- =====================================================

CREATE TABLE players (
                         id SERIAL PRIMARY KEY,
                         season_id INTEGER NOT NULL REFERENCES seasons(season_id) ON DELETE CASCADE,
                         fpl_player_id INTEGER NOT NULL,
                         player_code INTEGER NOT NULL,
                         first_name VARCHAR(100),
                         second_name VARCHAR(100),
                         web_name VARCHAR(50),
                         team_id INTEGER REFERENCES teams(id) ON DELETE SET NULL,
                         element_type_id INTEGER NOT NULL,
                         status VARCHAR(20),
                         photo VARCHAR(255),
                         squad_number INTEGER,
                         birth_date DATE,
                         team_join_date DATE,
                         region INTEGER,
                         opta_code VARCHAR(20),
                         can_transact BOOLEAN DEFAULT TRUE,
                         can_select BOOLEAN DEFAULT TRUE,
                         in_dreamteam BOOLEAN DEFAULT FALSE,
                         special BOOLEAN DEFAULT FALSE,
                         removed BOOLEAN DEFAULT FALSE,
                         unavailable BOOLEAN DEFAULT FALSE,
                         dreamteam_count INTEGER,
                         created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
                         updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
                         UNIQUE(fpl_player_id, season_id),
                         UNIQUE(player_code, season_id)
);

COMMENT ON TABLE players IS 'Player master data by season';
COMMENT ON COLUMN players.fpl_player_id IS 'FPL API player ID (resets each season)';
COMMENT ON COLUMN players.player_code IS 'FPL API player code (may persist across seasons)';
COMMENT ON COLUMN players.element_type_id IS '1=GKP, 2=DEF, 3=MID, 4=FWD';

CREATE INDEX idx_players_season ON players(season_id);
CREATE INDEX idx_players_team ON players(team_id);
CREATE INDEX idx_players_element_type ON players(element_type_id);
CREATE INDEX idx_players_fpl_id ON players(fpl_player_id, season_id);

-- =====================================================
-- Player Stats Tables
-- =====================================================

CREATE TABLE player_costs (
                              id SERIAL PRIMARY KEY,
                              player_id INTEGER REFERENCES players(id) ON DELETE CASCADE,
                              now_cost INTEGER,
                              cost_change_event INTEGER,
                              cost_change_event_fall INTEGER,
                              cost_change_start INTEGER,
                              cost_change_start_fall INTEGER,
                              updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

COMMENT ON TABLE player_costs IS 'Player cost and value information';
COMMENT ON COLUMN player_costs.now_cost IS 'Cost in tenths of £m (e.g., 105 = £10.5m)';

CREATE TABLE player_season_stats (
                                     id SERIAL PRIMARY KEY,
                                     player_id INTEGER REFERENCES players(id) ON DELETE CASCADE,
                                     dreamteam_count INTEGER,
                                     total_points INTEGER,
                                     event_points INTEGER,
                                     points_per_game NUMERIC(10, 2),
                                     form NUMERIC(10, 2),
                                     selected_by_percent NUMERIC(10, 2),
                                     value_form NUMERIC(10, 2),
                                     value_season NUMERIC(10, 2),
                                     minutes INTEGER DEFAULT 0,
                                     goals_scored INTEGER DEFAULT 0,
                                     assists INTEGER DEFAULT 0,
                                     clean_sheets INTEGER DEFAULT 0,
                                     goals_conceded INTEGER DEFAULT 0,
                                     own_goals INTEGER DEFAULT 0,
                                     penalties_saved INTEGER DEFAULT 0,
                                     penalties_missed INTEGER DEFAULT 0,
                                     yellow_cards INTEGER DEFAULT 0,
                                     red_cards INTEGER DEFAULT 0,
                                     saves INTEGER DEFAULT 0,
                                     bonus INTEGER DEFAULT 0,
                                     bps INTEGER DEFAULT 0,
                                     starts INTEGER DEFAULT 0,
                                     clearances_blocks_interceptions INTEGER DEFAULT 0,
                                     recoveries INTEGER DEFAULT 0,
                                     tackles INTEGER DEFAULT 0,
                                     defensive_contribution INTEGER DEFAULT 0,
                                     updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

COMMENT ON TABLE player_season_stats IS 'Player cumulative season statistics';

CREATE TABLE player_ict_stats (
                                  id SERIAL PRIMARY KEY,
                                  player_id INTEGER REFERENCES players(id) ON DELETE CASCADE,
                                  influence NUMERIC(10, 2),
                                  creativity NUMERIC(10, 2),
                                  threat NUMERIC(10, 2),
                                  ict_index NUMERIC(10, 2),
                                  influence_rank INTEGER,
                                  influence_rank_type INTEGER,
                                  creativity_rank INTEGER,
                                  creativity_rank_type INTEGER,
                                  threat_rank INTEGER,
                                  threat_rank_type INTEGER,
                                  ict_index_rank INTEGER,
                                  ict_index_rank_type INTEGER,
                                  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

COMMENT ON TABLE player_ict_stats IS 'Player ICT (Influence, Creativity, Threat) metrics';

CREATE TABLE player_expected_stats (
                                       id SERIAL PRIMARY KEY,
                                       player_id INTEGER REFERENCES players(id) ON DELETE CASCADE,
                                       expected_goals NUMERIC(10, 2),
                                       expected_assists NUMERIC(10, 2),
                                       expected_goal_involvements NUMERIC(10, 2),
                                       expected_goals_conceded NUMERIC(10, 2),
                                       expected_goals_per_90 NUMERIC(10, 2),
                                       expected_assists_per_90 NUMERIC(10, 2),
                                       expected_goal_involvements_per_90 NUMERIC(10, 2),
                                       expected_goals_conceded_per_90 NUMERIC(10, 2),
                                       saves_per_90 NUMERIC(10, 2),
                                       goals_conceded_per_90 NUMERIC(10, 2),
                                       starts_per_90 NUMERIC(10, 2),
                                       clean_sheets_per_90 NUMERIC(10, 2),
                                       defensive_contribution_per_90 NUMERIC(10, 2),
                                       updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

COMMENT ON TABLE player_expected_stats IS 'Player expected goals and per-90 statistics';

CREATE TABLE player_rankings (
                                 id SERIAL PRIMARY KEY,
                                 player_id INTEGER REFERENCES players(id) ON DELETE CASCADE,
                                 now_cost_rank INTEGER,
                                 now_cost_rank_type INTEGER,
                                 form_rank INTEGER,
                                 form_rank_type INTEGER,
                                 points_per_game_rank INTEGER,
                                 points_per_game_rank_type INTEGER,
                                 selected_rank INTEGER,
                                 selected_rank_type INTEGER,
                                 updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

COMMENT ON TABLE player_rankings IS 'Player rankings across various metrics';

-- =====================================================
-- Fixtures
-- =====================================================

CREATE TABLE fixtures (
                          id SERIAL PRIMARY KEY,
                          season_id INTEGER NOT NULL REFERENCES seasons(season_id) ON DELETE CASCADE,
                          fpl_fixture_id INTEGER NOT NULL,
                          fixture_code INTEGER NOT NULL,
                          event INTEGER NOT NULL,
                          team_h INTEGER REFERENCES teams(id) ON DELETE SET NULL,
                          team_a INTEGER REFERENCES teams(id) ON DELETE SET NULL,
                          team_h_score INTEGER,
                          team_a_score INTEGER,
                          team_h_difficulty INTEGER,
                          team_a_difficulty INTEGER,
                          kickoff_time TIMESTAMP WITH TIME ZONE,
                          minutes INTEGER DEFAULT 0,
                          started BOOLEAN DEFAULT FALSE,
                          finished BOOLEAN DEFAULT FALSE,
                          finished_provisional BOOLEAN DEFAULT FALSE,
                          provisional_start_time BOOLEAN DEFAULT FALSE,
                          pulse_id INTEGER,
                          created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
                          updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
                          UNIQUE(fpl_fixture_id, season_id),
                          UNIQUE(fixture_code, season_id)
);

COMMENT ON TABLE fixtures IS 'Match fixtures by season';
COMMENT ON COLUMN fixtures.fpl_fixture_id IS 'FPL API fixture ID (resets each season)';
COMMENT ON COLUMN fixtures.event IS 'Gameweek number';

CREATE INDEX idx_fixtures_season ON fixtures(season_id);
CREATE INDEX idx_fixtures_event ON fixtures(season_id, event);
CREATE INDEX idx_fixtures_teams ON fixtures(team_h, team_a);
CREATE INDEX idx_fixtures_fpl_id ON fixtures(fpl_fixture_id, season_id);

CREATE TABLE fixture_stats (
                               id SERIAL PRIMARY KEY,
                               fixture_id INTEGER NOT NULL REFERENCES fixtures(id) ON DELETE CASCADE,
                               identifier VARCHAR(50) NOT NULL,
                               player_id INTEGER REFERENCES players(id) ON DELETE CASCADE,
                               value INTEGER,
                               created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
                               updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                               UNIQUE(fixture_id, identifier, player_id)
);

COMMENT ON TABLE fixture_stats IS 'Detailed fixture statistics by player';

CREATE INDEX idx_fixture_stats_fixture ON fixture_stats(fixture_id);
CREATE INDEX idx_fixture_stats_player ON fixture_stats(player_id);

-- =====================================================
-- Managers
-- =====================================================

CREATE TABLE managers (
                          id SERIAL PRIMARY KEY,
                          fpl_manager_id INTEGER NOT NULL UNIQUE,
                          manager_name VARCHAR(100),
                          player_first_name VARCHAR(100),
                          player_last_name VARCHAR(100),
                          player_region_id INTEGER,
                          player_region_name VARCHAR(100),
                          player_region_iso_code_short VARCHAR(10),
                          player_region_iso_code_long VARCHAR(10),
                          favourite_team INTEGER REFERENCES teams(id) ON DELETE SET NULL,
                          joined_time TIMESTAMP WITH TIME ZONE,
                          started_event INTEGER,
                          years_active INTEGER,
                          summary_overall_points INTEGER,
                          summary_overall_rank INTEGER,
                          summary_event_points INTEGER,
                          summary_event_rank INTEGER,
                          current_event INTEGER,
                          name_change_blocked BOOLEAN DEFAULT FALSE,
                          last_deadline_bank INTEGER,
                          last_deadline_value INTEGER,
                          last_deadline_total_transfers INTEGER,
                          club_badge_src VARCHAR(255),
                          created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
                          updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

COMMENT ON TABLE managers IS 'FPL manager/entry information';
COMMENT ON COLUMN managers.fpl_manager_id IS 'FPL API manager/entry ID (persists across seasons)';

CREATE INDEX idx_managers_fpl_id ON managers(fpl_manager_id);

CREATE TABLE manager_gameweek_history (
                                          id SERIAL PRIMARY KEY,
                                          manager_id INTEGER NOT NULL REFERENCES managers(id) ON DELETE CASCADE,
                                          season_id INTEGER NOT NULL REFERENCES seasons(season_id) ON DELETE CASCADE,
                                          event INTEGER NOT NULL,
                                          points INTEGER,
                                          total_points INTEGER,
                                          rank INTEGER,
                                          rank_sort INTEGER,
                                          overall_rank INTEGER,
                                          percentile_rank INTEGER,
                                          bank INTEGER,
                                          value INTEGER,
                                          event_transfers INTEGER,
                                          event_transfers_cost INTEGER,
                                          points_on_bench INTEGER,
                                          created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
                                          UNIQUE(manager_id, season_id, event)
);

COMMENT ON TABLE manager_gameweek_history IS 'Manager performance per gameweek';
COMMENT ON COLUMN manager_gameweek_history.value IS 'Team value in tenths of £m';

CREATE INDEX idx_manager_gameweek_season ON manager_gameweek_history(season_id);
CREATE INDEX idx_manager_gameweek_event ON manager_gameweek_history(season_id, event);
CREATE INDEX idx_manager_gameweek_manager ON manager_gameweek_history(manager_id);

CREATE TABLE manager_season_history (
                                        id SERIAL PRIMARY KEY,
                                        manager_id INTEGER NOT NULL REFERENCES managers(id) ON DELETE CASCADE,
                                        season_name VARCHAR(20) NOT NULL,
                                        total_points INTEGER,
                                        rank INTEGER,
                                        created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
                                        UNIQUE(manager_id, season_name)
);

COMMENT ON TABLE manager_season_history IS 'Manager performance per season';

CREATE INDEX idx_manager_season_manager ON manager_season_history(manager_id);

CREATE TABLE manager_picks (
                               id SERIAL PRIMARY KEY,
                               manager_id INTEGER NOT NULL REFERENCES managers(id) ON DELETE CASCADE,
                               season_id INTEGER NOT NULL REFERENCES seasons(season_id) ON DELETE CASCADE,
                               event INTEGER NOT NULL,
                               player_id INTEGER NOT NULL REFERENCES players(id) ON DELETE CASCADE,
                               position INTEGER NOT NULL,
                               multiplier INTEGER NOT NULL,
                               is_captain BOOLEAN DEFAULT FALSE,
                               is_vice_captain BOOLEAN DEFAULT FALSE,
                               element_type INTEGER,
                               created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
                               UNIQUE(manager_id, season_id, event, player_id)
);

COMMENT ON TABLE manager_picks IS 'Manager team selection per gameweek';
COMMENT ON COLUMN manager_picks.position IS 'Position 1-15 (1-11 starting, 12-15 bench)';
COMMENT ON COLUMN manager_picks.multiplier IS '0=bench, 1=playing, 2=captain, 3=triple captain';

CREATE INDEX idx_manager_picks_season ON manager_picks(season_id);
CREATE INDEX idx_manager_picks_event ON manager_picks(season_id, event);
CREATE INDEX idx_manager_picks_manager ON manager_picks(manager_id);

CREATE TABLE automatic_subs (
                                id SERIAL PRIMARY KEY,
                                manager_id INTEGER NOT NULL REFERENCES managers(id) ON DELETE CASCADE,
                                season_id INTEGER NOT NULL REFERENCES seasons(season_id) ON DELETE CASCADE,
                                event INTEGER NOT NULL,
                                element_in INTEGER NOT NULL REFERENCES players(id) ON DELETE CASCADE,
                                element_out INTEGER NOT NULL REFERENCES players(id) ON DELETE CASCADE,
                                created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

COMMENT ON TABLE automatic_subs IS 'Automatic substitutions made';

CREATE INDEX idx_automatic_subs_manager_event ON automatic_subs(manager_id, season_id, event);

CREATE TABLE chip_usage (
                            id SERIAL PRIMARY KEY,
                            manager_id INTEGER NOT NULL REFERENCES managers(id) ON DELETE CASCADE,
                            season_id INTEGER NOT NULL REFERENCES seasons(season_id) ON DELETE CASCADE,
                            chip_name VARCHAR(50) NOT NULL,
                            event INTEGER NOT NULL,
                            time TIMESTAMP WITH TIME ZONE,
                            created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
                            UNIQUE(manager_id, season_id, chip_name, event)
);

COMMENT ON TABLE chip_usage IS 'Manager chip usage (wildcard, free hit, etc.)';

CREATE INDEX idx_chip_usage_season ON chip_usage(season_id);
CREATE INDEX idx_chip_usage_event ON chip_usage(season_id, event);

-- =====================================================
-- Player Gameweek Stats
-- =====================================================

CREATE TABLE player_gameweek_stats (
                                       id SERIAL PRIMARY KEY,
                                       player_id INTEGER NOT NULL REFERENCES players(id) ON DELETE CASCADE,
                                       fixture_id INTEGER NOT NULL REFERENCES fixtures(id) ON DELETE CASCADE,
                                       season_id INTEGER NOT NULL REFERENCES seasons(season_id) ON DELETE CASCADE,
                                       event INTEGER NOT NULL,
                                       opponent_team_id INTEGER REFERENCES teams(id) ON DELETE SET NULL,
                                       kickoff_time TIMESTAMP WITH TIME ZONE,
                                       was_home BOOLEAN,
                                       team_h_score INTEGER,
                                       team_a_score INTEGER,
                                       minutes INTEGER DEFAULT 0,
                                       goals_scored INTEGER DEFAULT 0,
                                       assists INTEGER DEFAULT 0,
                                       clean_sheets INTEGER DEFAULT 0,
                                       goals_conceded INTEGER DEFAULT 0,
                                       own_goals INTEGER DEFAULT 0,
                                       penalties_saved INTEGER DEFAULT 0,
                                       penalties_missed INTEGER DEFAULT 0,
                                       yellow_cards INTEGER DEFAULT 0,
                                       red_cards INTEGER DEFAULT 0,
                                       saves INTEGER DEFAULT 0,
                                       bonus INTEGER DEFAULT 0,
                                       bps INTEGER DEFAULT 0,
                                       starts INTEGER DEFAULT 0,
                                       clearances_blocks_interceptions INTEGER DEFAULT 0,
                                       recoveries INTEGER DEFAULT 0,
                                       tackles INTEGER DEFAULT 0,
                                       defensive_contribution INTEGER DEFAULT 0,
                                       influence NUMERIC(10, 2),
                                       creativity NUMERIC(10, 2),
                                       threat NUMERIC(10, 2),
                                       ict_index NUMERIC(10, 2),
                                       expected_goals NUMERIC(10, 2),
                                       expected_assists NUMERIC(10, 2),
                                       expected_goal_involvements NUMERIC(10, 2),
                                       expected_goals_conceded NUMERIC(10, 2),
                                       total_points INTEGER,
                                       value INTEGER,
                                       transfers_balance INTEGER,
                                       selected INTEGER,
                                       transfers_in INTEGER,
                                       transfers_out INTEGER,
                                       modified BOOLEAN DEFAULT FALSE,
                                       created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
                                       updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                       UNIQUE(player_id, fixture_id)
);

COMMENT ON TABLE player_gameweek_stats IS 'Player performance per gameweek/fixture';

CREATE INDEX idx_player_gameweek_season ON player_gameweek_stats(season_id);
CREATE INDEX idx_player_gameweek_event ON player_gameweek_stats(season_id, event);
CREATE INDEX idx_player_gameweek_player ON player_gameweek_stats(player_id);
CREATE INDEX idx_player_gameweek_fixture ON player_gameweek_stats(fixture_id);

CREATE TABLE player_gameweek_explain (
                                         id SERIAL PRIMARY KEY,
                                         player_id INTEGER NOT NULL REFERENCES players(id) ON DELETE CASCADE,
                                         fixture_id INTEGER NOT NULL REFERENCES fixtures(id) ON DELETE CASCADE,
                                         season_id INTEGER NOT NULL REFERENCES seasons(season_id) ON DELETE CASCADE,
                                         event INTEGER NOT NULL,
                                         identifier VARCHAR(50) NOT NULL,
                                         points INTEGER,
                                         value INTEGER,
                                         points_modification INTEGER,
                                         created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

COMMENT ON TABLE player_gameweek_explain IS 'Detailed points breakdown per fixture';

CREATE INDEX idx_player_gameweek_explain_player_fixture ON player_gameweek_explain(player_id, fixture_id);
CREATE INDEX idx_player_gameweek_explain_season ON player_gameweek_explain(season_id);

-- =====================================================
-- Historical Data
-- =====================================================

CREATE TABLE player_past_seasons (
                                     id SERIAL PRIMARY KEY,
                                     player_code INTEGER NOT NULL,
                                     season_name VARCHAR(20) NOT NULL,
                                     start_cost INTEGER,
                                     end_cost INTEGER,
                                     total_points INTEGER,
                                     minutes INTEGER DEFAULT 0,
                                     goals_scored INTEGER DEFAULT 0,
                                     assists INTEGER DEFAULT 0,
                                     clean_sheets INTEGER DEFAULT 0,
                                     goals_conceded INTEGER DEFAULT 0,
                                     own_goals INTEGER DEFAULT 0,
                                     penalties_saved INTEGER DEFAULT 0,
                                     penalties_missed INTEGER DEFAULT 0,
                                     yellow_cards INTEGER DEFAULT 0,
                                     red_cards INTEGER DEFAULT 0,
                                     saves INTEGER DEFAULT 0,
                                     bonus INTEGER DEFAULT 0,
                                     bps INTEGER DEFAULT 0,
                                     starts INTEGER DEFAULT 0,
                                     clearances_blocks_interceptions INTEGER DEFAULT 0,
                                     recoveries INTEGER DEFAULT 0,
                                     tackles INTEGER DEFAULT 0,
                                     defensive_contribution INTEGER DEFAULT 0,
                                     influence NUMERIC(10, 2),
                                     creativity NUMERIC(10, 2),
                                     threat NUMERIC(10, 2),
                                     ict_index NUMERIC(10, 2),
                                     expected_goals NUMERIC(10, 2),
                                     expected_assists NUMERIC(10, 2),
                                     expected_goal_involvements NUMERIC(10, 2),
                                     expected_goals_conceded NUMERIC(10, 2),
                                     created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
                                     UNIQUE(player_code, season_name)
);

COMMENT ON TABLE player_past_seasons IS 'Historical player data from past seasons';

CREATE INDEX idx_player_past_seasons_code ON player_past_seasons(player_code);
CREATE INDEX idx_player_past_seasons_season ON player_past_seasons(season_name);