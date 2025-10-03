-- ============================================
-- REFERENCE TABLES
-- ============================================

CREATE TABLE teams (
                       team_id INTEGER PRIMARY KEY,
                       team_code INTEGER UNIQUE NOT NULL,
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
                       unavailable BOOLEAN DEFAULT false,
                       pulse_id INTEGER,
                       strength_overall_home INTEGER,
                       strength_overall_away INTEGER,
                       strength_attack_home INTEGER,
                       strength_attack_away INTEGER,
                       strength_defence_home INTEGER,
                       strength_defence_away INTEGER,
                       created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
                       updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE players (
                         player_id INTEGER PRIMARY KEY,
                         player_code INTEGER UNIQUE NOT NULL,
                         first_name VARCHAR(100),
                         second_name VARCHAR(100),
                         web_name VARCHAR(50),
                         team_id INTEGER REFERENCES teams(team_id),
                         team_code INTEGER,
                         element_type_id INTEGER NOT NULL,
                         status VARCHAR(20),
                         photo VARCHAR(255),
                         squad_number INTEGER,
                         birth_date DATE,
                         team_join_date DATE,
                         region INTEGER,
                         opta_code VARCHAR(20),
                         can_transact BOOLEAN DEFAULT true,
                         can_select BOOLEAN DEFAULT true,
                         in_dreamteam BOOLEAN DEFAULT false,
                         special BOOLEAN DEFAULT false,
                         removed BOOLEAN DEFAULT false,
                         unavailable BOOLEAN DEFAULT false,
                         created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
                         updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE player_costs (
                              player_id INTEGER PRIMARY KEY REFERENCES players(player_id),
                              now_cost INTEGER,
                              cost_change_event INTEGER,
                              cost_change_event_fall INTEGER,
                              cost_change_start INTEGER,
                              cost_change_start_fall INTEGER,
                              updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE player_season_stats (
                                     player_id INTEGER PRIMARY KEY REFERENCES players(player_id),
                                     dreamteam_count INTEGER,
                                     total_points INTEGER,
                                     event_points INTEGER,
                                     points_per_game NUMERIC(10,2),
                                     form NUMERIC(10,2),
                                     selected_by_percent NUMERIC(10,2),
                                     value_form NUMERIC(10,2),
                                     value_season NUMERIC(10,2),
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
                                     updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE player_ict_stats (
                                  player_id INTEGER PRIMARY KEY REFERENCES players(player_id),
                                  influence NUMERIC(10,2),
                                  creativity NUMERIC(10,2),
                                  threat NUMERIC(10,2),
                                  ict_index NUMERIC(10,2),
                                  influence_rank INTEGER,
                                  influence_rank_type INTEGER,
                                  creativity_rank INTEGER,
                                  creativity_rank_type INTEGER,
                                  threat_rank INTEGER,
                                  threat_rank_type INTEGER,
                                  ict_index_rank INTEGER,
                                  ict_index_rank_type INTEGER,
                                  updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE player_expected_stats (
                                       player_id INTEGER PRIMARY KEY REFERENCES players(player_id),
                                       expected_goals NUMERIC(10,2),
                                       expected_assists NUMERIC(10,2),
                                       expected_goal_involvements NUMERIC(10,2),
                                       expected_goals_conceded NUMERIC(10,2),
                                       expected_goals_per_90 NUMERIC(10,2),
                                       expected_assists_per_90 NUMERIC(10,2),
                                       expected_goal_involvements_per_90 NUMERIC(10,2),
                                       expected_goals_conceded_per_90 NUMERIC(10,2),
                                       saves_per_90 NUMERIC(10,2),
                                       goals_conceded_per_90 NUMERIC(10,2),
                                       starts_per_90 NUMERIC(10,2),
                                       clean_sheets_per_90 NUMERIC(10,2),
                                       defensive_contribution_per_90 NUMERIC(10,2),
                                       updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE player_rankings (
                                 player_id INTEGER PRIMARY KEY REFERENCES players(player_id),
                                 now_cost_rank INTEGER,
                                 now_cost_rank_type INTEGER,
                                 form_rank INTEGER,
                                 form_rank_type INTEGER,
                                 points_per_game_rank INTEGER,
                                 points_per_game_rank_type INTEGER,
                                 selected_rank INTEGER,
                                 selected_rank_type INTEGER,
                                 updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE fixtures (
                          fixture_id INTEGER PRIMARY KEY,
                          fixture_code INTEGER UNIQUE NOT NULL,
                          event INTEGER NOT NULL,
                          team_h INTEGER REFERENCES teams(team_id),
                          team_a INTEGER REFERENCES teams(team_id),
                          team_h_score INTEGER,
                          team_a_score INTEGER,
                          team_h_difficulty INTEGER,
                          team_a_difficulty INTEGER,
                          kickoff_time TIMESTAMPTZ,
                          minutes INTEGER DEFAULT 0,
                          started BOOLEAN DEFAULT false,
                          finished BOOLEAN DEFAULT false,
                          finished_provisional BOOLEAN DEFAULT false,
                          provisional_start_time BOOLEAN DEFAULT false,
                          pulse_id INTEGER,
                          created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
                          updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE fixture_stats (
                               id SERIAL PRIMARY KEY,
                               fixture_id INTEGER REFERENCES fixtures(fixture_id),
                               identifier VARCHAR(50) NOT NULL,
                               player_id INTEGER REFERENCES players(player_id),
                               value INTEGER,
                               is_home BOOLEAN NOT NULL,
                               created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
                               CONSTRAINT unique_fixture_stat UNIQUE(fixture_id, identifier, player_id)
);

-- ============================================
-- MANAGER TABLES
-- ============================================

CREATE TABLE managers (
                          manager_id INTEGER PRIMARY KEY,
                          manager_name VARCHAR(100),
                          player_first_name VARCHAR(100),
                          player_last_name VARCHAR(100),
                          player_region_id INTEGER,
                          player_region_name VARCHAR(100),
                          player_region_iso_code_short VARCHAR(10),
                          player_region_iso_code_long VARCHAR(10),
                          favourite_team INTEGER REFERENCES teams(team_id),
                          joined_time TIMESTAMPTZ,
                          started_event INTEGER,
                          years_active INTEGER,
                          summary_overall_points INTEGER,
                          summary_overall_rank INTEGER,
                          summary_event_points INTEGER,
                          summary_event_rank INTEGER,
                          current_event INTEGER,
                          name_change_blocked BOOLEAN DEFAULT false,
                          last_deadline_bank INTEGER,
                          last_deadline_value INTEGER,
                          last_deadline_total_transfers INTEGER,
                          club_badge_src VARCHAR(255),
                          created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
                          updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE manager_gameweek_history (
                                          id SERIAL PRIMARY KEY,
                                          manager_id INTEGER REFERENCES managers(manager_id),
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
                                          created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
                                          CONSTRAINT unique_manager_event UNIQUE(manager_id, event)
);

CREATE TABLE manager_season_history (
                                        id SERIAL PRIMARY KEY,
                                        manager_id INTEGER REFERENCES managers(manager_id),
                                        season_name VARCHAR(20) NOT NULL,
                                        total_points INTEGER,
                                        rank INTEGER,
                                        created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
                                        CONSTRAINT unique_manager_season UNIQUE(manager_id, season_name)
);

CREATE TABLE manager_picks (
                               id SERIAL PRIMARY KEY,
                               manager_id INTEGER REFERENCES managers(manager_id),
                               event INTEGER NOT NULL,
                               player_id INTEGER REFERENCES players(player_id),
                               position INTEGER NOT NULL,
                               multiplier INTEGER NOT NULL,
                               is_captain BOOLEAN DEFAULT false,
                               is_vice_captain BOOLEAN DEFAULT false,
                               element_type INTEGER,
                               created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
                               CONSTRAINT unique_manager_event_player UNIQUE(manager_id, event, player_id)
);

CREATE TABLE automatic_subs (
                                id SERIAL PRIMARY KEY,
                                manager_id INTEGER REFERENCES managers(manager_id),
                                event INTEGER NOT NULL,
                                element_in INTEGER REFERENCES players(player_id),
                                element_out INTEGER REFERENCES players(player_id),
                                created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE chip_usage (
                            id SERIAL PRIMARY KEY,
                            manager_id INTEGER REFERENCES managers(manager_id),
                            chip_name VARCHAR(50) NOT NULL,
                            event INTEGER NOT NULL,
                            time TIMESTAMPTZ,
                            created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
                            CONSTRAINT unique_manager_chip_event UNIQUE(manager_id, chip_name, event)
);

-- ============================================
-- GAMEWEEK/EVENT TABLES
-- ============================================

CREATE TABLE player_gameweek_stats (
                                       id SERIAL PRIMARY KEY,
                                       player_id INTEGER REFERENCES players(player_id),
                                       fixture_id INTEGER REFERENCES fixtures(fixture_id),
                                       event INTEGER NOT NULL,
                                       opponent_team_id INTEGER REFERENCES teams(team_id),
                                       kickoff_time TIMESTAMPTZ,
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
                                       influence NUMERIC(10,2),
                                       creativity NUMERIC(10,2),
                                       threat NUMERIC(10,2),
                                       ict_index NUMERIC(10,2),
                                       expected_goals NUMERIC(10,2),
                                       expected_assists NUMERIC(10,2),
                                       expected_goal_involvements NUMERIC(10,2),
                                       expected_goals_conceded NUMERIC(10,2),
                                       total_points INTEGER,
                                       value INTEGER,
                                       transfers_balance INTEGER,
                                       selected INTEGER,
                                       transfers_in INTEGER,
                                       transfers_out INTEGER,
                                       in_dreamteam BOOLEAN DEFAULT false,
                                       modified BOOLEAN DEFAULT false,
                                       created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
                                       CONSTRAINT unique_player_fixture UNIQUE(player_id, fixture_id)
);

CREATE TABLE player_gameweek_explain (
                                         id SERIAL PRIMARY KEY,
                                         player_id INTEGER REFERENCES players(player_id),
                                         fixture_id INTEGER REFERENCES fixtures(fixture_id),
                                         event INTEGER NOT NULL,
                                         identifier VARCHAR(50) NOT NULL,
                                         points INTEGER,
                                         value INTEGER,
                                         points_modification INTEGER,
                                         created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

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
                                     influence NUMERIC(10,2),
                                     creativity NUMERIC(10,2),
                                     threat NUMERIC(10,2),
                                     ict_index NUMERIC(10,2),
                                     expected_goals NUMERIC(10,2),
                                     expected_assists NUMERIC(10,2),
                                     expected_goal_involvements NUMERIC(10,2),
                                     expected_goals_conceded NUMERIC(10,2),
                                     created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
                                     CONSTRAINT unique_player_season UNIQUE(player_code, season_name)
);

-- ============================================
-- INDEXES
-- ============================================

CREATE INDEX idx_players_team_id ON players(team_id);
CREATE INDEX idx_players_element_type ON players(element_type_id);
CREATE INDEX idx_fixtures_event ON fixtures(event);
CREATE INDEX idx_fixtures_teams ON fixtures(team_h, team_a);
CREATE INDEX idx_fixture_stats_fixture ON fixture_stats(fixture_id);
CREATE INDEX idx_fixture_stats_player ON fixture_stats(player_id);
CREATE INDEX idx_manager_gameweek_history_event ON manager_gameweek_history(event);
CREATE INDEX idx_manager_picks_event ON manager_picks(event);
CREATE INDEX idx_player_gameweek_stats_event ON player_gameweek_stats(event);
CREATE INDEX idx_player_gameweek_stats_player ON player_gameweek_stats(player_id);
CREATE INDEX idx_player_gameweek_explain_player_fixture ON player_gameweek_explain(player_id, fixture_id);
CREATE INDEX idx_chip_usage_event ON chip_usage(event);
CREATE INDEX idx_automatic_subs_manager_event ON automatic_subs(manager_id, event);

-- ============================================
-- COMMENTS
-- ============================================

COMMENT ON TABLE teams IS 'Premier League teams';
COMMENT ON TABLE players IS 'Player master data';
COMMENT ON TABLE player_costs IS 'Player cost and value information';
COMMENT ON TABLE player_season_stats IS 'Player cumulative season statistics';
COMMENT ON TABLE player_ict_stats IS 'Player ICT (Influence, Creativity, Threat) metrics';
COMMENT ON TABLE player_expected_stats IS 'Player expected goals and per-90 statistics';
COMMENT ON TABLE player_rankings IS 'Player rankings across various metrics';
COMMENT ON TABLE fixtures IS 'Match fixtures';
COMMENT ON TABLE fixture_stats IS 'Detailed fixture statistics by player';
COMMENT ON TABLE managers IS 'FPL manager/entry information';
COMMENT ON TABLE manager_gameweek_history IS 'Manager performance per gameweek';
COMMENT ON TABLE manager_season_history IS 'Manager performance per season';
COMMENT ON TABLE manager_picks IS 'Manager team selection per gameweek';
COMMENT ON TABLE automatic_subs IS 'Automatic substitutions made';
COMMENT ON TABLE chip_usage IS 'Manager chip usage (wildcard, free hit, etc.)';
COMMENT ON TABLE player_gameweek_stats IS 'Player performance per gameweek/fixture';
COMMENT ON TABLE player_gameweek_explain IS 'Detailed points breakdown per fixture';
COMMENT ON TABLE player_past_seasons IS 'Historical player data from past seasons';

COMMENT ON COLUMN players.element_type_id IS '1=GKP, 2=DEF, 3=MID, 4=FWD';
COMMENT ON COLUMN player_costs.now_cost IS 'Cost in tenths of £m (e.g., 105 = £10.5m)';
COMMENT ON COLUMN manager_gameweek_history.value IS 'Team value in tenths of £m';
COMMENT ON COLUMN manager_picks.multiplier IS '0=bench, 1=playing, 2=captain, 3=triple captain';
COMMENT ON COLUMN manager_picks.position IS 'Position 1-15 (1-11 starting, 12-15 bench)';