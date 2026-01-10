-- FPL Service Database Schema
\connect fpl;

-- ==========================================
-- 1. BASE ENTITIES
-- ==========================================

-- Teams
CREATE TABLE IF NOT EXISTS teams (
                                     team_id INTEGER NOT NULL,
                                     season_id INTEGER NOT NULL,
                                     team_code INTEGER,
                                     name VARCHAR(255),
                                     short_name VARCHAR(50),
                                     strength INTEGER,
                                     form DECIMAL,
                                     position INTEGER,
                                     points INTEGER,
                                     played INTEGER,
                                     win INTEGER,
                                     draw INTEGER,
                                     loss INTEGER,
                                     team_division INTEGER,
                                     unavailable BOOLEAN,
                                     pulse_id INTEGER,
                                     strength_overall_home INTEGER,
                                     strength_overall_away INTEGER,
                                     strength_attack_home INTEGER,
                                     strength_attack_away INTEGER,
                                     strength_defence_home INTEGER,
                                     strength_defence_away INTEGER,
                                     updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                     PRIMARY KEY (team_id, season_id)
);

-- TeamPlayers (Bootstrap info)
CREATE TABLE IF NOT EXISTS players (
                                       player_id INTEGER NOT NULL,
                                       season_id INTEGER NOT NULL,
                                       player_code INTEGER,
                                       first_name VARCHAR(255),
                                       second_name VARCHAR(255),
                                       web_name VARCHAR(255),
                                       team_id INTEGER,
                                       team_code INTEGER,
                                       element_type_id INTEGER,
                                       status VARCHAR(50),
                                       photo VARCHAR(255),
                                       squad_number INTEGER,
                                       birth_date DATE,
                                       team_join_date DATE,
                                       region INTEGER,
                                       opta_code VARCHAR(50),
                                       can_transact BOOLEAN,
                                       can_select BOOLEAN,
                                       in_dreamteam BOOLEAN,
                                       dreamteam_count INTEGER,
                                       special BOOLEAN,
                                       removed BOOLEAN,
                                       unavailable BOOLEAN,
                                       updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

                                       PRIMARY KEY (player_id, season_id),

    -- Relationships
                                       CONSTRAINT fk_players_team FOREIGN KEY (team_id, season_id) REFERENCES teams(team_id, season_id)
);

-- ==========================================
-- 2. PLAYER STATS & DETAILS
--    (Added season_id to these tables to allow Foreign Keys to players)
-- ==========================================

-- Player Costs
CREATE TABLE IF NOT EXISTS player_costs (
                                            player_id INTEGER NOT NULL,
                                            season_id INTEGER NOT NULL, -- Added for Relation
                                            now_cost INTEGER,
                                            cost_change_event INTEGER,
                                            cost_change_event_fall INTEGER,
                                            cost_change_start INTEGER,
                                            cost_change_start_fall INTEGER,
                                            updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

                                            PRIMARY KEY (player_id, season_id),
                                            CONSTRAINT fk_costs_player FOREIGN KEY (player_id, season_id) REFERENCES players(player_id, season_id)
);

-- Player Season Stats
CREATE TABLE IF NOT EXISTS player_season_stats (
                                                   player_id INTEGER NOT NULL,
                                                   season_id INTEGER NOT NULL, -- Added for Relation
                                                   dreamteam_count INTEGER,
                                                   total_points INTEGER,
                                                   event_points INTEGER,
                                                   points_per_game DECIMAL,
                                                   form DECIMAL,
                                                   selected_by_percent DECIMAL,
                                                   value_form DECIMAL,
                                                   value_season DECIMAL,
                                                   minutes INTEGER,
                                                   goals_scored INTEGER,
                                                   assists INTEGER,
                                                   clean_sheets INTEGER,
                                                   goals_conceded INTEGER,
                                                   own_goals INTEGER,
                                                   penalties_saved INTEGER,
                                                   penalties_missed INTEGER,
                                                   yellow_cards INTEGER,
                                                   red_cards INTEGER,
                                                   saves INTEGER,
                                                   bonus INTEGER,
                                                   bps INTEGER,
                                                   starts INTEGER,
                                                   clearances_blocks_interceptions INTEGER,
                                                   recoveries INTEGER,
                                                   tackles INTEGER,
                                                   defensive_contribution INTEGER,
                                                   updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

                                                   PRIMARY KEY (player_id, season_id),
                                                   CONSTRAINT fk_season_stats_player FOREIGN KEY (player_id, season_id) REFERENCES players(player_id, season_id)
);

-- Player ICT Stats
CREATE TABLE IF NOT EXISTS player_ict_stats (
                                                player_id INTEGER NOT NULL,
                                                season_id INTEGER NOT NULL, -- Added for Relation
                                                influence DECIMAL,
                                                creativity DECIMAL,
                                                threat DECIMAL,
                                                ict_index DECIMAL,
                                                influence_rank INTEGER,
                                                influence_rank_type INTEGER,
                                                creativity_rank INTEGER,
                                                creativity_rank_type INTEGER,
                                                threat_rank INTEGER,
                                                threat_rank_type INTEGER,
                                                ict_index_rank INTEGER,
                                                ict_index_rank_type INTEGER,
                                                updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

                                                PRIMARY KEY (player_id, season_id),
                                                CONSTRAINT fk_ict_stats_player FOREIGN KEY (player_id, season_id) REFERENCES players(player_id, season_id)
);

-- Player Expected Stats
CREATE TABLE IF NOT EXISTS player_expected_stats (
                                                     player_id INTEGER NOT NULL,
                                                     season_id INTEGER NOT NULL, -- Added for Relation
                                                     expected_goals DECIMAL,
                                                     expected_assists DECIMAL,
                                                     expected_goal_involvements DECIMAL,
                                                     expected_goals_conceded DECIMAL,
                                                     expected_goals_per_90 DECIMAL,
                                                     expected_assists_per_90 DECIMAL,
                                                     expected_goal_involvements_per_90 DECIMAL,
                                                     expected_goals_conceded_per_90 DECIMAL,
                                                     saves_per_90 DECIMAL,
                                                     goals_conceded_per_90 DECIMAL,
                                                     starts_per_90 DECIMAL,
                                                     clean_sheets_per_90 DECIMAL,
                                                     defensive_contribution_per_90 DECIMAL,
                                                     updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

                                                     PRIMARY KEY (player_id, season_id),
                                                     CONSTRAINT fk_expected_stats_player FOREIGN KEY (player_id, season_id) REFERENCES players(player_id, season_id)
);

-- Player Rankings
CREATE TABLE IF NOT EXISTS player_rankings (
                                               player_id INTEGER NOT NULL,
                                               season_id INTEGER NOT NULL, -- Added for Relation
                                               now_cost_rank INTEGER,
                                               now_cost_rank_type INTEGER,
                                               form_rank INTEGER,
                                               form_rank_type INTEGER,
                                               points_per_game_rank INTEGER,
                                               points_per_game_rank_type INTEGER,
                                               selected_rank INTEGER,
                                               selected_rank_type INTEGER,
                                               updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

                                               PRIMARY KEY (player_id, season_id),
                                               CONSTRAINT fk_rankings_player FOREIGN KEY (player_id, season_id) REFERENCES players(player_id, season_id)
);

-- ==========================================
-- 3. FIXTURES
-- ==========================================

-- Fixtures
CREATE TABLE IF NOT EXISTS fixtures (
                                        fixture_id INTEGER NOT NULL,
                                        season_id INTEGER NOT NULL,
                                        fixture_code INTEGER,
                                        event INTEGER,
                                        team_h INTEGER,
                                        team_a INTEGER,
                                        kickoff_time TIMESTAMP,
                                        team_h_score INTEGER,
                                        team_a_score INTEGER,
                                        finished BOOLEAN,
                                        minutes INTEGER,
                                        provisional_start_time BOOLEAN,
                                        team_h_difficulty INTEGER,
                                        team_a_difficulty INTEGER,
                                        pulse_id INTEGER,
                                        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

                                        PRIMARY KEY (fixture_id, season_id),

    -- Relationships (Composite keys for teams)
                                        CONSTRAINT fk_fixture_home_team FOREIGN KEY (team_h, season_id) REFERENCES teams(team_id, season_id),
                                        CONSTRAINT fk_fixture_away_team FOREIGN KEY (team_a, season_id) REFERENCES teams(team_id, season_id)
);

-- Fixture Stats
CREATE TABLE IF NOT EXISTS fixture_stats (
                                             fixture_id INTEGER NOT NULL,
                                             player_id INTEGER NOT NULL,
                                             identifier VARCHAR(50) NOT NULL,
                                             season_id INTEGER NOT NULL,
                                             value INTEGER,
                                             updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

                                             PRIMARY KEY (fixture_id, player_id, identifier, season_id),

                                             CONSTRAINT fk_fixture_stats_fixture FOREIGN KEY (fixture_id, season_id) REFERENCES fixtures(fixture_id, season_id),
                                             CONSTRAINT fk_fixture_stats_player FOREIGN KEY (player_id, season_id) REFERENCES players(player_id, season_id)
);

-- ==========================================
-- 4. HISTORY & LOGS
-- ==========================================

-- Player Gameweek Stats (History)
CREATE TABLE IF NOT EXISTS player_gameweek_stats (
                                                     player_id INTEGER NOT NULL,
                                                     fixture_id INTEGER NOT NULL,
                                                     season_id INTEGER NOT NULL,
                                                     event INTEGER,
                                                     opponent_team_id INTEGER,
                                                     kickoff_time TIMESTAMP,
                                                     was_home BOOLEAN,
                                                     team_h_score INTEGER,
                                                     team_a_score INTEGER,
                                                     minutes INTEGER,
                                                     goals_scored INTEGER,
                                                     assists INTEGER,
                                                     clean_sheets INTEGER,
                                                     goals_conceded INTEGER,
                                                     own_goals INTEGER,
                                                     penalties_saved INTEGER,
                                                     penalties_missed INTEGER,
                                                     yellow_cards INTEGER,
                                                     red_cards INTEGER,
                                                     saves INTEGER,
                                                     bonus INTEGER,
                                                     bps INTEGER,
                                                     starts INTEGER,
                                                     clearances_blocks_interceptions INTEGER,
                                                     recoveries INTEGER,
                                                     tackles INTEGER,
                                                     defensive_contribution INTEGER,
                                                     influence DECIMAL,
                                                     creativity DECIMAL,
                                                     threat DECIMAL,
                                                     ict_index DECIMAL,
                                                     expected_goals DECIMAL,
                                                     expected_assists DECIMAL,
                                                     expected_goal_involvements DECIMAL,
                                                     expected_goals_conceded DECIMAL,
                                                     total_points INTEGER,
                                                     value INTEGER,
                                                     transfers_balance INTEGER,
                                                     selected INTEGER,
                                                     transfers_in INTEGER,
                                                     transfers_out INTEGER,
                                                     modified BOOLEAN,
                                                     updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

                                                     PRIMARY KEY (player_id, fixture_id, season_id),

                                                     CONSTRAINT fk_gw_stats_player FOREIGN KEY (player_id, season_id) REFERENCES players(player_id, season_id),
                                                     CONSTRAINT fk_gw_stats_fixture FOREIGN KEY (fixture_id, season_id) REFERENCES fixtures(fixture_id, season_id),
                                                     CONSTRAINT fk_gw_stats_opp_team FOREIGN KEY (opponent_team_id, season_id) REFERENCES teams(team_id, season_id)
);

-- Player Gameweek Explain
CREATE TABLE IF NOT EXISTS player_gameweek_explain (
                                                       player_id INTEGER NOT NULL,
                                                       fixture_id INTEGER NOT NULL,
                                                       season_id INTEGER NOT NULL,
                                                       identifier VARCHAR(50) NOT NULL,
                                                       event INTEGER,
                                                       points INTEGER,
                                                       value INTEGER,
                                                       points_modification INTEGER,
                                                       updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

                                                       PRIMARY KEY (player_id, season_id, fixture_id, identifier),

                                                       CONSTRAINT fk_gw_explain_player FOREIGN KEY (player_id, season_id) REFERENCES players(player_id, season_id),
                                                       CONSTRAINT fk_gw_explain_fixture FOREIGN KEY (fixture_id, season_id) REFERENCES fixtures(fixture_id, season_id)
);

-- Player Past Seasons (Note: Links via player_code conceptually, no strict FK to players due to multi-season redundancy)
CREATE TABLE IF NOT EXISTS player_past_seasons (
                                                   player_code INTEGER NOT NULL,
                                                   season_id INTEGER NOT NULL,
                                                   season_name VARCHAR(50),
                                                   start_cost INTEGER,
                                                   end_cost INTEGER,
                                                   total_points INTEGER,
                                                   minutes INTEGER,
                                                   goals_scored INTEGER,
                                                   assists INTEGER,
                                                   clean_sheets INTEGER,
                                                   goals_conceded INTEGER,
                                                   own_goals INTEGER,
                                                   penalties_saved INTEGER,
                                                   penalties_missed INTEGER,
                                                   yellow_cards INTEGER,
                                                   red_cards INTEGER,
                                                   saves INTEGER,
                                                   bonus INTEGER,
                                                   bps INTEGER,
                                                   starts INTEGER,
                                                   clearances_blocks_interceptions INTEGER,
                                                   recoveries INTEGER,
                                                   tackles INTEGER,
                                                   defensive_contribution INTEGER,
                                                   influence DECIMAL,
                                                   creativity DECIMAL,
                                                   threat DECIMAL,
                                                   ict_index DECIMAL,
                                                   expected_goals DECIMAL,
                                                   expected_assists DECIMAL,
                                                   expected_goal_involvements DECIMAL,
                                                   expected_goals_conceded DECIMAL,
                                                   updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                                   PRIMARY KEY (player_code, season_id)
);

-- ==========================================
-- 5. MANAGERS
-- ==========================================

-- Managers
CREATE TABLE IF NOT EXISTS managers (
                                        manager_id INTEGER NOT NULL,
                                        season_id INTEGER NOT NULL,
                                        manager_name VARCHAR(255),
                                        player_first_name VARCHAR(255),
                                        player_last_name VARCHAR(255),
                                        player_region_id INTEGER,
                                        player_region_name VARCHAR(255),
                                        player_region_iso_code_short VARCHAR(10),
                                        player_region_iso_code_long VARCHAR(10),
                                        favourite_team_id INTEGER,
                                        joined_time TIMESTAMP,
                                        started_event INTEGER,
                                        years_active INTEGER,
                                        summary_overall_points INTEGER,
                                        summary_overall_rank INTEGER,
                                        summary_points INTEGER,
                                        summary_rank INTEGER,
                                        current_event INTEGER,
                                        name_change_blocked BOOLEAN,
                                        last_deadline_bank INTEGER,
                                        last_deadline_value INTEGER,
                                        last_deadline_total_transfers INTEGER,
                                        club_badge_src VARCHAR(255),
                                        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

                                        PRIMARY KEY (manager_id, season_id),

                                        CONSTRAINT fk_manager_fav_team FOREIGN KEY (favourite_team_id, season_id) REFERENCES teams(team_id, season_id)
);

-- Manager Picks
CREATE TABLE IF NOT EXISTS manager_picks (
                                             manager_id INTEGER NOT NULL,
                                             season_id INTEGER NOT NULL,
                                             event INTEGER NOT NULL,
                                             player_id INTEGER NOT NULL,
                                             is_captain BOOLEAN,
                                             is_vice_captain BOOLEAN,
                                             multiplier INTEGER,
                                             position INTEGER,
                                             element_type INTEGER,
                                             updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

                                             PRIMARY KEY (manager_id, season_id, event, player_id),

                                             CONSTRAINT fk_picks_manager FOREIGN KEY (manager_id, season_id) REFERENCES managers(manager_id, season_id),
                                             CONSTRAINT fk_picks_player FOREIGN KEY (player_id, season_id) REFERENCES players(player_id, season_id)
);

-- Manager Automatic Subs
CREATE TABLE IF NOT EXISTS manager_automatic_subs (
                                                      manager_id INTEGER NOT NULL,
                                                      season_id INTEGER NOT NULL,
                                                      event INTEGER NOT NULL,
                                                      player_out_id INTEGER NOT NULL,
                                                      player_in_id INTEGER NOT NULL,
                                                      updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

                                                      PRIMARY KEY (manager_id, season_id, event, player_out_id, player_in_id),

                                                      CONSTRAINT fk_subs_manager FOREIGN KEY (manager_id, season_id) REFERENCES managers(manager_id, season_id),
                                                      CONSTRAINT fk_subs_player_in FOREIGN KEY (player_in_id, season_id) REFERENCES players(player_id, season_id),
                                                      CONSTRAINT fk_subs_player_out FOREIGN KEY (player_out_id, season_id) REFERENCES players(player_id, season_id)
);

-- Manager Transfers
CREATE TABLE IF NOT EXISTS manager_transfers (
                                                 manager_id INTEGER NOT NULL,
                                                 season_id INTEGER NOT NULL,
                                                 event INTEGER NOT NULL,
                                                 player_in_id INTEGER NOT NULL,
                                                 player_out_id INTEGER NOT NULL,
                                                 player_in_cost INTEGER,
                                                 player_out_cost INTEGER,
                                                 updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

                                                 PRIMARY KEY (manager_id, season_id, event, player_in_id, player_out_id),

                                                 CONSTRAINT fk_transfers_manager FOREIGN KEY (manager_id, season_id) REFERENCES managers(manager_id, season_id),
                                                 CONSTRAINT fk_transfers_player_in FOREIGN KEY (player_in_id, season_id) REFERENCES players(player_id, season_id),
                                                 CONSTRAINT fk_transfers_player_out FOREIGN KEY (player_out_id, season_id) REFERENCES players(player_id, season_id)
);

-- Manager Gameweek History
CREATE TABLE IF NOT EXISTS manager_gameweek_history (
                                                        manager_id INTEGER NOT NULL,
                                                        season_id INTEGER NOT NULL,
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
                                                        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

                                                        PRIMARY KEY (manager_id, season_id, event),

                                                        CONSTRAINT fk_gw_history_manager FOREIGN KEY (manager_id, season_id) REFERENCES managers(manager_id, season_id)
);

-- Manager Season History (Past)
CREATE TABLE IF NOT EXISTS manager_season_history (
                                                      manager_id INTEGER NOT NULL,
                                                      season_id INTEGER NOT NULL,
                                                      season_name VARCHAR(50),
                                                      total_points INTEGER,
                                                      rank INTEGER,
                                                      updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

                                                      PRIMARY KEY (manager_id, season_id)
    -- No strict FK to managers(manager_id) because this table may contain history for seasons not present in the main 'managers' table
);

-- Manager Chips
CREATE TABLE IF NOT EXISTS manager_chips (
                                             manager_id INTEGER NOT NULL,
                                             season_id INTEGER NOT NULL,
                                             event INTEGER NOT NULL,
                                             chip_name VARCHAR(50) NOT NULL,
                                             time TIMESTAMP,
                                             updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

                                             PRIMARY KEY (manager_id, season_id, event, chip_name),

                                             CONSTRAINT fk_chips_manager FOREIGN KEY (manager_id, season_id) REFERENCES managers(manager_id, season_id)
);