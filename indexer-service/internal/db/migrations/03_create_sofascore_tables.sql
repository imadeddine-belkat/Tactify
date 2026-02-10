-- Sofascore Service Database Schema
\connect sofascore;

-- ==========================================
-- 1. BASE TABLES
-- ==========================================

-- Leagues
CREATE TABLE IF NOT EXISTS leagues (
                                       league_id INTEGER PRIMARY KEY,
                                       name VARCHAR(255),
                                       country VARCHAR(255),
                                       updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Seasons
CREATE TABLE IF NOT EXISTS seasons (
                                       season_id INTEGER PRIMARY KEY,
                                       league_id INTEGER NOT NULL,
                                       name VARCHAR(255),
                                       year VARCHAR(20),
                                       is_current BOOLEAN,
                                       updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                       CONSTRAINT fk_seasons_league FOREIGN KEY (league_id) REFERENCES leagues(league_id)
);

-- Teams
CREATE TABLE IF NOT EXISTS teams (
                                     team_id INTEGER NOT NULL,
                                     league_id INTEGER NOT NULL,
                                     name VARCHAR(255),
                                     primary_color VARCHAR(50),
                                     secondary_color VARCHAR(50),
                                     PRIMARY KEY (team_id, league_id),
                                     CONSTRAINT fk_teams_league FOREIGN KEY (league_id) REFERENCES leagues(league_id)
);

CREATE TABLE IF NOT EXISTS players (
                                    player_id INTEGER NOT NULL,
                                    season_id INTEGER NOT NULL,
                                    team_id INTEGER NOT NULL,
                                    league_id INTEGER NOT NULL,
                                    player_name varchar(100) not null,
                                    player_short_name varchar(50),
                                    position varchar(50),
                                    height int,
                                    preferred_foot varchar(10),
                                    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                    PRIMARY KEY (player_id, season_id, team_id, league_id),
                                    CONSTRAINT fk_players_season FOREIGN KEY (season_id) REFERENCES seasons(season_id),
                                    CONSTRAINT fk_players_team FOREIGN KEY (team_id, league_id) REFERENCES teams(team_id, league_id)
);

-- Matches
CREATE TABLE IF NOT EXISTS matches (
                                       match_id INTEGER NOT NULL,
                                       season_id INTEGER NOT NULL,
                                       league_id INTEGER NOT NULL,
                                       home_team_id INTEGER,
                                       away_team_id INTEGER,
                                       home_team_name VARCHAR(255),
                                       away_team_name VARCHAR(255),
                                       start_time TIMESTAMP,
                                       round VARCHAR(50),
                                       status VARCHAR(50),
                                       status_description VARCHAR(255),
                                       updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

                                       PRIMARY KEY (match_id, season_id, league_id),
    -- Constraint to ensure match_id is unique globally, allowing FK references from stats tables
                                       CONSTRAINT uq_matches_match_id UNIQUE (match_id),

    -- Foreign Keys
                                       CONSTRAINT fk_matches_season FOREIGN KEY (season_id) REFERENCES seasons(season_id),
                                       CONSTRAINT fk_matches_league FOREIGN KEY (league_id) REFERENCES leagues(league_id),
    -- Link teams using composite key (team + league)
                                       CONSTRAINT fk_matches_home_team FOREIGN KEY (home_team_id, league_id) REFERENCES teams(team_id, league_id),
                                       CONSTRAINT fk_matches_away_team FOREIGN KEY (away_team_id, league_id) REFERENCES teams(team_id, league_id)
);

-- Team Overall Stats
CREATE TABLE IF NOT EXISTS team_overall_stats (
    -- Primary Keys
                                                  team_id INTEGER NOT NULL,
                                                  league_id INTEGER NOT NULL,
                                                  season_id INTEGER NOT NULL,

    -- Offensive Stats
                                                  goals_scored INTEGER,
                                                  goals_conceded INTEGER,
                                                  own_goals INTEGER,
                                                  assists INTEGER,
                                                  shots INTEGER,
                                                  shots_on_target INTEGER,
                                                  shots_off_target INTEGER,
                                                  penalty_goals INTEGER,
                                                  penalties_taken INTEGER,
                                                  free_kick_goals INTEGER,
                                                  free_kick_shots INTEGER,

    -- Positional Goals
                                                  goals_from_inside_box INTEGER,
                                                  goals_from_outside_box INTEGER,
                                                  headed_goals INTEGER,
                                                  left_foot_goals INTEGER,
                                                  right_foot_goals INTEGER,

    -- Chances
                                                  big_chances INTEGER,
                                                  big_chances_created INTEGER,
                                                  big_chances_missed INTEGER,

    -- Shooting Details
                                                  shots_from_inside_box INTEGER,
                                                  shots_from_outside_box INTEGER,
                                                  blocked_scoring_attempt INTEGER,
                                                  hit_woodwork INTEGER,

    -- Dribbling
                                                  successful_dribbles INTEGER,
                                                  dribble_attempts INTEGER,

    -- Set Pieces
                                                  corners INTEGER,
                                                  free_kicks INTEGER,
                                                  throw_ins INTEGER,
                                                  goal_kicks INTEGER,

    -- Fast Breaks
                                                  fast_breaks INTEGER,
                                                  fast_break_goals INTEGER,
                                                  fast_break_shots INTEGER,

    -- Possession & Passing
                                                  average_ball_possession DECIMAL,
                                                  total_passes INTEGER,
                                                  accurate_passes INTEGER,
                                                  accurate_passes_percentage DECIMAL,

    -- Passing by Zone
                                                  total_own_half_passes INTEGER,
                                                  accurate_own_half_passes INTEGER,
                                                  accurate_own_half_passes_percentage DECIMAL,
                                                  total_opposition_half_passes INTEGER,
                                                  accurate_opposition_half_passes INTEGER,
                                                  accurate_opposition_half_passes_percentage DECIMAL,

    -- Long Balls & Crosses
                                                  total_long_balls INTEGER,
                                                  accurate_long_balls INTEGER,
                                                  accurate_long_balls_percentage DECIMAL,
                                                  total_crosses INTEGER,
                                                  accurate_crosses INTEGER,
                                                  accurate_crosses_percentage DECIMAL,

    -- Defensive Stats
                                                  clean_sheets INTEGER,
                                                  tackles INTEGER,
                                                  interceptions INTEGER,
                                                  saves INTEGER,
                                                  clearances INTEGER,
                                                  clearances_off_line INTEGER,
                                                  last_man_tackles INTEGER,
                                                  ball_recovery INTEGER,

    -- Errors
                                                  errors_leading_to_goal INTEGER,
                                                  errors_leading_to_shot INTEGER,

    -- Penalties
                                                  penalties_commited INTEGER,
                                                  penalty_goals_conceded INTEGER,

    -- Duels
                                                  total_duels INTEGER,
                                                  duels_won INTEGER,
                                                  duels_won_percentage DECIMAL,
                                                  total_ground_duels INTEGER,
                                                  ground_duels_won INTEGER,
                                                  ground_duels_won_percentage DECIMAL,
                                                  total_aerial_duels INTEGER,
                                                  aerial_duels_won INTEGER,
                                                  aerial_duels_won_percentage DECIMAL,

    -- Discipline
                                                  possession_lost INTEGER,
                                                  offsides INTEGER,
                                                  fouls INTEGER,
                                                  yellow_cards INTEGER,
                                                  yellow_red_cards INTEGER,
                                                  red_cards INTEGER,

    -- Performance
                                                  avg_rating DECIMAL,
                                                  matches INTEGER,
                                                  awarded_matches INTEGER,

    -- Stats Against
                                                  accurate_final_third_passes_against INTEGER,
                                                  accurate_opposition_half_passes_against INTEGER,
                                                  accurate_own_half_passes_against INTEGER,
                                                  accurate_passes_against INTEGER,
                                                  big_chances_against INTEGER,
                                                  big_chances_created_against INTEGER,
                                                  big_chances_missed_against INTEGER,
                                                  clearances_against INTEGER,
                                                  corners_against INTEGER,
                                                  crosses_successful_against INTEGER,
                                                  crosses_total_against INTEGER,
                                                  dribble_attempts_total_against INTEGER,
                                                  dribble_attempts_won_against INTEGER,
                                                  errors_leading_to_goal_against INTEGER,
                                                  errors_leading_to_shot_against INTEGER,
                                                  hit_woodwork_against INTEGER,
                                                  interceptions_against INTEGER,
                                                  key_passes_against INTEGER,
                                                  long_balls_successful_against INTEGER,
                                                  long_balls_total_against INTEGER,
                                                  offsides_against INTEGER,
                                                  red_cards_against INTEGER,
                                                  shots_against INTEGER,
                                                  shots_blocked_against INTEGER,
                                                  shots_from_inside_box_against INTEGER,
                                                  shots_from_outside_box_against INTEGER,
                                                  shots_off_target_against INTEGER,
                                                  shots_on_target_against INTEGER,
                                                  blocked_scoring_attempt_against INTEGER,
                                                  tackles_against INTEGER,
                                                  total_final_third_passes_against INTEGER,
                                                  opposition_half_passes_total_against INTEGER,
                                                  own_half_passes_total_against INTEGER,
                                                  total_passes_against INTEGER,
                                                  yellow_cards_against INTEGER,

                                                  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                                  PRIMARY KEY (team_id, league_id, season_id),

    -- Foreign Keys
                                                  CONSTRAINT fk_ovr_stats_team FOREIGN KEY (team_id, league_id) REFERENCES teams(team_id, league_id),
                                                  CONSTRAINT fk_ovr_stats_season FOREIGN KEY (season_id) REFERENCES seasons(season_id)
);

-- ==========================================
-- 2. DYNAMIC MATCH STATS TABLES
-- ==========================================

-- Match Overview
CREATE TABLE IF NOT EXISTS match_overview (
                                              match_id INTEGER NOT NULL,
                                              period VARCHAR(20) NOT NULL,
                                              home_team_id INTEGER,
                                              away_team_id INTEGER,
                                              updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

                                              ball_possession_home DECIMAL,
                                              ball_possession_away DECIMAL,
                                              expected_goals_home DECIMAL,
                                              expected_goals_away DECIMAL,
                                              big_chances_home INTEGER,
                                              big_chances_away INTEGER,
                                              total_shots_home INTEGER,
                                              total_shots_away INTEGER,
                                              goalkeeper_saves_home INTEGER,
                                              goalkeeper_saves_away INTEGER,
                                              corner_kicks_home INTEGER,
                                              corner_kicks_away INTEGER,
                                              fouls_home INTEGER,
                                              fouls_away INTEGER,
                                              passes_home INTEGER,
                                              passes_away INTEGER,
                                              tackles_home INTEGER,
                                              tackles_away INTEGER,
                                              free_kicks_home INTEGER,
                                              free_kicks_away INTEGER,
                                              yellow_cards_home INTEGER,
                                              yellow_cards_away INTEGER,

                                              PRIMARY KEY (match_id, period),
                                              CONSTRAINT fk_overview_match FOREIGN KEY (match_id) REFERENCES matches(match_id)
);

-- Match Shots
CREATE TABLE IF NOT EXISTS match_shots (
                                           match_id INTEGER NOT NULL,
                                           period VARCHAR(20) NOT NULL,
                                           home_team_id INTEGER,
                                           away_team_id INTEGER,
                                           updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

                                           total_shots_home INTEGER,
                                           total_shots_away INTEGER,
                                           shots_on_target_home INTEGER,
                                           shots_on_target_away INTEGER,
                                           hit_woodwork_home INTEGER,
                                           hit_woodwork_away INTEGER,
                                           shots_off_target_home INTEGER,
                                           shots_off_target_away INTEGER,
                                           blocked_shots_home INTEGER,
                                           blocked_shots_away INTEGER,
                                           shots_inside_box_home INTEGER,
                                           shots_inside_box_away INTEGER,
                                           shots_outside_box_home INTEGER,
                                           shots_outside_box_away INTEGER,

                                           PRIMARY KEY (match_id, period),
                                           CONSTRAINT fk_shots_match FOREIGN KEY (match_id) REFERENCES matches(match_id)
);

-- Match Attack
CREATE TABLE IF NOT EXISTS match_attack (
                                            match_id INTEGER NOT NULL,
                                            period VARCHAR(20) NOT NULL,
                                            home_team_id INTEGER,
                                            away_team_id INTEGER,
                                            updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

                                            big_chances_scored_home INTEGER,
                                            big_chances_scored_away INTEGER,
                                            big_chances_missed_home INTEGER,
                                            big_chances_missed_away INTEGER,
                                            through_balls_home INTEGER,
                                            through_balls_away INTEGER,
                                            touches_in_penalty_area_home INTEGER,
                                            touches_in_penalty_area_away INTEGER,
                                            fouled_in_final_third_home INTEGER,
                                            fouled_in_final_third_away INTEGER,
                                            offsides_home INTEGER,
                                            offsides_away INTEGER,

                                            PRIMARY KEY (match_id, period),
                                            CONSTRAINT fk_attack_match FOREIGN KEY (match_id) REFERENCES matches(match_id)
);

-- Match Passes
CREATE TABLE IF NOT EXISTS match_passes (
                                            match_id INTEGER NOT NULL,
                                            period VARCHAR(20) NOT NULL,
                                            home_team_id INTEGER,
                                            away_team_id INTEGER,
                                            updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

                                            accurate_passes_home INTEGER,
                                            accurate_passes_away INTEGER,
                                            throw_ins_home INTEGER,
                                            throw_ins_away INTEGER,
                                            final_third_entries_home INTEGER,
                                            final_third_entries_away INTEGER,
                                            final_third_phase_home INTEGER,
                                            final_third_phase_away INTEGER,
                                            final_third_phase_total_home INTEGER,
                                            final_third_phase_total_away INTEGER,
                                            long_balls_home INTEGER,
                                            long_balls_away INTEGER,
                                            long_balls_total_home INTEGER,
                                            long_balls_total_away INTEGER,
                                            crosses_home INTEGER,
                                            crosses_away INTEGER,
                                            crosses_total_home INTEGER,
                                            crosses_total_away INTEGER,

                                            PRIMARY KEY (match_id, period),
                                            CONSTRAINT fk_passes_match FOREIGN KEY (match_id) REFERENCES matches(match_id)
);

-- Match Duels
CREATE TABLE IF NOT EXISTS match_duels (
                                           match_id INTEGER NOT NULL,
                                           period VARCHAR(20) NOT NULL,
                                           home_team_id INTEGER,
                                           away_team_id INTEGER,
                                           updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

                                           duels_won_percent_home DECIMAL,
                                           duels_won_percent_away DECIMAL,
                                           dispossessed_home INTEGER,
                                           dispossessed_away INTEGER,
                                           ground_duels_home INTEGER,
                                           ground_duels_away INTEGER,
                                           ground_duels_total_home INTEGER,
                                           ground_duels_total_away INTEGER,
                                           aerial_duels_home INTEGER,
                                           aerial_duels_away INTEGER,
                                           aerial_duels_total_home INTEGER,
                                           aerial_duels_total_away INTEGER,
                                           dribbles_home INTEGER,
                                           dribbles_away INTEGER,
                                           dribbles_total_home INTEGER,
                                           dribbles_total_away INTEGER,

                                           PRIMARY KEY (match_id, period),
                                           CONSTRAINT fk_duels_match FOREIGN KEY (match_id) REFERENCES matches(match_id)
);

-- Match Defending
CREATE TABLE IF NOT EXISTS match_defending (
                                               match_id INTEGER NOT NULL,
                                               period VARCHAR(20) NOT NULL,
                                               home_team_id INTEGER,
                                               away_team_id INTEGER,
                                               updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

                                               tackles_won_home INTEGER,
                                               tackles_won_away INTEGER,
                                               tackles_won_total_home INTEGER,
                                               tackles_won_total_away INTEGER,
                                               total_tackles_home INTEGER,
                                               total_tackles_away INTEGER,
                                               interceptions_home INTEGER,
                                               interceptions_away INTEGER,
                                               recoveries_home INTEGER,
                                               recoveries_away INTEGER,
                                               clearances_home INTEGER,
                                               clearances_away INTEGER,
                                               errors_lead_to_shot_home INTEGER,
                                               errors_lead_to_shot_away INTEGER,

                                               PRIMARY KEY (match_id, period),
                                               CONSTRAINT fk_defending_match FOREIGN KEY (match_id) REFERENCES matches(match_id)
);

-- Match Goalkeeping
CREATE TABLE IF NOT EXISTS match_goalkeeping (
                                                 match_id INTEGER NOT NULL,
                                                 period VARCHAR(20) NOT NULL,
                                                 home_team_id INTEGER,
                                                 away_team_id INTEGER,
                                                 updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

                                                 total_saves_home INTEGER,
                                                 total_saves_away INTEGER,
                                                 goals_prevented_home DECIMAL,
                                                 goals_prevented_away DECIMAL,
                                                 big_saves_home INTEGER,
                                                 big_saves_away INTEGER,
                                                 high_claims_home INTEGER,
                                                 high_claims_away INTEGER,
                                                 punches_home INTEGER,
                                                 punches_away INTEGER,
                                                 goal_kicks_home INTEGER,
                                                 goal_kicks_away INTEGER,

                                                 PRIMARY KEY (match_id, period),
                                                 CONSTRAINT fk_goalkeeping_match FOREIGN KEY (match_id) REFERENCES matches(match_id)
);