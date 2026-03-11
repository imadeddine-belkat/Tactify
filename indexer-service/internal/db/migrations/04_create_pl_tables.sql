\connect fpl;

-- ==========================================
-- 1. BASE ENTITIES (code-keyed layer)
-- ==========================================

CREATE TABLE IF NOT EXISTS element_types (
                                             id                      INTEGER NOT NULL,
                                             plural_name             VARCHAR(255) NOT NULL,
                                             plural_name_short       VARCHAR(50) NOT NULL,
                                             singular_name           VARCHAR(255) NOT NULL,
                                             singular_name_short     VARCHAR(50) NOT NULL,
                                             squad_select            INTEGER NOT NULL,
                                             squad_min_select        INTEGER,
                                             squad_max_select        INTEGER,
                                             squad_min_play          INTEGER NOT NULL,
                                             squad_max_play          INTEGER NOT NULL,
                                             ui_shirt_specific       BOOLEAN NOT NULL DEFAULT FALSE,
                                             sub_positions_locked    INTEGER[] NOT NULL DEFAULT '{}',
                                             PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS seasons (
                                       code         INTEGER NOT NULL,
                                       name         VARCHAR(255) NOT NULL,
                                       is_current   BOOLEAN NOT NULL,
                                       PRIMARY KEY (code)
);

CREATE TABLE IF NOT EXISTS teams (
                                     code       INTEGER NOT NULL,
                                     name       VARCHAR(255),
                                     short_name VARCHAR(50),
                                     PRIMARY KEY (code)
);

CREATE TABLE IF NOT EXISTS players (
                                       code            INTEGER NOT NULL,
                                       first_name      VARCHAR(255),
                                       second_name     VARCHAR(255),
                                       web_name        VARCHAR(255),
                                       birth_date      DATE,
                                       opta_code       VARCHAR(50),
                                       PRIMARY KEY (code)
);

CREATE TABLE IF NOT EXISTS fixtures (
                                        code INTEGER NOT NULL,
                                        PRIMARY KEY (code)
);

-- ==========================================
-- BRIDGE TABLES (code <-> id per season)
-- ==========================================

CREATE TABLE IF NOT EXISTS team_seasons (
                                            team_code   INTEGER NOT NULL,
                                            season_code INTEGER NOT NULL,
                                            team_id     INTEGER NOT NULL,
                                            PRIMARY KEY (team_code, season_code),
                                            FOREIGN KEY (team_code)   REFERENCES teams(code)   ON DELETE CASCADE ON UPDATE CASCADE,
                                            FOREIGN KEY (season_code) REFERENCES seasons(code)  ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE UNIQUE INDEX IF NOT EXISTS uq_team_seasons_id
    ON team_seasons (team_id, season_code);

CREATE TABLE IF NOT EXISTS player_seasons (
                                              player_code INTEGER NOT NULL,
                                              season_code INTEGER NOT NULL,
                                              player_id   INTEGER NOT NULL,
                                              PRIMARY KEY (player_code, season_code),
                                              FOREIGN KEY (player_code) REFERENCES players(code) ON DELETE CASCADE ON UPDATE CASCADE,
                                              FOREIGN KEY (season_code) REFERENCES seasons(code)  ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE UNIQUE INDEX IF NOT EXISTS uq_player_seasons_id
    ON player_seasons (player_id, season_code);

CREATE TABLE IF NOT EXISTS fixture_seasons (
                                               season_code  INTEGER NOT NULL,
                                               fixture_code INTEGER NOT NULL,
                                               fixture_id   INTEGER NOT NULL,
                                               PRIMARY KEY (season_code, fixture_code),
                                               FOREIGN KEY (fixture_code) REFERENCES fixtures(code) ON DELETE CASCADE ON UPDATE CASCADE,
                                               FOREIGN KEY (season_code) REFERENCES seasons(code)   ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE UNIQUE INDEX IF NOT EXISTS uq_fixture_seasons_id
    ON fixture_seasons (fixture_id, season_code);

CREATE INDEX IF NOT EXISTS idx_fixture_seasons_code ON fixture_seasons (fixture_code);

-- ==========================================
-- 2. FIXTURES
-- ==========================================

CREATE TABLE IF NOT EXISTS fixture_info (
                                            season_code       INTEGER NOT NULL,
                                            fixture_id        INTEGER NOT NULL,
                                            event             INTEGER,
                                            team_h_id         INTEGER,
                                            team_a_id         INTEGER,
                                            team_h_score      INTEGER,
                                            team_a_score      INTEGER,
                                            kickoff_time      TIMESTAMP,
                                            finished          BOOLEAN,
                                            team_h_difficulty INTEGER,
                                            team_a_difficulty INTEGER,
                                            PRIMARY KEY (season_code, fixture_id),
                                            FOREIGN KEY (fixture_id, season_code) REFERENCES fixture_seasons(fixture_id, season_code) ON DELETE CASCADE ON UPDATE CASCADE,
                                            FOREIGN KEY (team_h_id, season_code)  REFERENCES team_seasons(team_id, season_code)       ON DELETE CASCADE ON UPDATE CASCADE,
                                            FOREIGN KEY (team_a_id, season_code)  REFERENCES team_seasons(team_id, season_code)       ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_fixture_info_event    ON fixture_info (season_code, event);
CREATE INDEX IF NOT EXISTS idx_fixture_info_team_h   ON fixture_info (team_h_id, season_code);
CREATE INDEX IF NOT EXISTS idx_fixture_info_team_a   ON fixture_info (team_a_id, season_code);
CREATE INDEX IF NOT EXISTS idx_fixture_info_finished ON fixture_info (season_code, finished);
CREATE INDEX IF NOT EXISTS idx_fixture_info_kickoff  ON fixture_info (season_code, kickoff_time);

CREATE TABLE IF NOT EXISTS fixture_stats_general (
                                                     fixture_code            INTEGER NOT NULL,
                                                     team_code               INTEGER NOT NULL,
                                                     possession_percentage   REAL,
                                                     touches                 REAL,
                                                     touches_in_opp_box      REAL,
                                                     total_distance          REAL,
                                                     duel_won                REAL,
                                                     duel_lost               REAL,
                                                     won_contest             REAL,
                                                     total_contest           REAL,
                                                     challenge_lost          REAL,
                                                     dispossessed            REAL,
                                                     unsuccessful_touch      REAL,
                                                     overrun                 REAL,
                                                     poss_lost_all           REAL,
                                                     poss_lost_ctrl          REAL,
                                                     fk_foul_won             REAL,
                                                     fk_foul_lost            REAL,
                                                     attempted_tackle_foul   REAL,
                                                     yellow_card             REAL,
                                                     total_yel_card          REAL,
                                                     red_card                REAL,
                                                     total_red_card          REAL,
                                                     subs_made               REAL,
                                                     subs_goals              REAL,
                                                     goals                   REAL,
                                                     goals_conceded          REAL,
                                                     goals_conceded_ibox     REAL,
                                                     goals_conceded_obox     REAL,
                                                     own_goals               REAL,
                                                     winning_goal            REAL,
                                                     goal_assist             REAL,
                                                     goal_assist_openplay    REAL,
                                                     goal_assist_intentional REAL,
                                                     fastest_player_speed    REAL,
                                                     fastest_player_id       TEXT,
                                                     PRIMARY KEY (fixture_code, team_code),
                                                     FOREIGN KEY (fixture_code) REFERENCES fixtures(code) ON DELETE CASCADE ON UPDATE CASCADE,
                                                     FOREIGN KEY (team_code)    REFERENCES teams(code)    ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_fsg_team ON fixture_stats_general (team_code);

CREATE TABLE IF NOT EXISTS fixture_stats_attacking (
                                                       fixture_code                      INTEGER NOT NULL,
                                                       team_code                         INTEGER NOT NULL,
                                                       expected_goals                    REAL,
                                                       expected_goals_on_target          REAL,
                                                       expected_goals_on_target_conceded REAL,
                                                       expected_assists                  REAL,
                                                       total_scoring_att                 REAL,
                                                       ontarget_scoring_att              REAL,
                                                       blocked_scoring_att               REAL,
                                                       shot_off_target                   REAL,
                                                       post_scoring_att                  REAL,
                                                       hit_woodwork                      REAL,
                                                       att_post_high                     REAL,
                                                       attempts_ibox                     REAL,
                                                       attempts_obox                     REAL,
                                                       att_ibox_goal                     REAL,
                                                       att_ibox_target                   REAL,
                                                       att_ibox_miss                     REAL,
                                                       att_ibox_post                     REAL,
                                                       att_ibox_blocked                  REAL,
                                                       att_obox_goal                     REAL,
                                                       att_obox_target                   REAL,
                                                       att_obox_miss                     REAL,
                                                       att_obox_blocked                  REAL,
                                                       att_obx_centre                    REAL,
                                                       att_obp_goal                      REAL,
                                                       att_bx_left                       REAL,
                                                       att_bx_right                      REAL,
                                                       att_bx_centre                     REAL,
                                                       att_openplay                      REAL,
                                                       att_setpiece                      REAL,
                                                       att_corner                        REAL,
                                                       att_fastbreak                     REAL,
                                                       att_freekick_goal                 REAL,
                                                       att_pen_goal                      REAL,
                                                       goals_openplay                    REAL,
                                                       goal_fastbreak                    REAL,
                                                       big_chance_created                REAL,
                                                       big_chance_scored                 REAL,
                                                       big_chance_missed                 REAL,
                                                       total_att_assist                  REAL,
                                                       ontarget_att_assist               REAL,
                                                       offtarget_att_assist              REAL,
                                                       att_assist_openplay               REAL,
                                                       att_assist_setplay                REAL,
                                                       shot_fastbreak                    REAL,
                                                       total_fastbreak                   REAL,
                                                       pen_area_entries                  REAL,
                                                       final_third_entries               REAL,
                                                       put_through                       REAL,
                                                       successful_put_through            REAL,
                                                       defender_goals                    REAL,
                                                       midfielder_goals                  REAL,
                                                       keeper_goals                      REAL,
                                                       PRIMARY KEY (fixture_code, team_code),
                                                       FOREIGN KEY (fixture_code) REFERENCES fixtures(code) ON DELETE CASCADE ON UPDATE CASCADE,
                                                       FOREIGN KEY (team_code)    REFERENCES teams(code)    ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_fsa_team ON fixture_stats_attacking (team_code);

CREATE TABLE IF NOT EXISTS fixture_stats_shooting (
                                                      fixture_code         INTEGER NOT NULL,
                                                      team_code            INTEGER NOT NULL,
                                                      att_rf_total         REAL,
                                                      att_rf_target        REAL,
                                                      att_rf_goal          REAL,
                                                      att_lf_total         REAL,
                                                      att_lf_goal          REAL,
                                                      att_hd_total         REAL,
                                                      att_hd_target        REAL,
                                                      att_hd_miss          REAL,
                                                      att_hd_goal          REAL,
                                                      att_goal_low_left    REAL,
                                                      att_goal_low_right   REAL,
                                                      att_goal_low_centre  REAL,
                                                      att_goal_high_left   REAL,
                                                      att_goal_high_right  REAL,
                                                      att_goal_high_centre REAL,
                                                      att_miss_left        REAL,
                                                      att_miss_right       REAL,
                                                      att_miss_high        REAL,
                                                      att_miss_high_left   REAL,
                                                      att_miss_high_right  REAL,
                                                      att_sv_low_centre    REAL,
                                                      att_sv_low_right     REAL,
                                                      att_sv_high_centre   REAL,
                                                      PRIMARY KEY (fixture_code, team_code),
                                                      FOREIGN KEY (fixture_code) REFERENCES fixtures(code) ON DELETE CASCADE ON UPDATE CASCADE,
                                                      FOREIGN KEY (team_code)    REFERENCES teams(code)    ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_fss_team ON fixture_stats_shooting (team_code);

CREATE TABLE IF NOT EXISTS fixture_stats_passing (
                                                     fixture_code                   INTEGER NOT NULL,
                                                     team_code                      INTEGER NOT NULL,
                                                     total_pass                     REAL,
                                                     accurate_pass                  REAL,
                                                     open_play_pass                 REAL,
                                                     successful_open_play_pass      REAL,
                                                     fwd_pass                       REAL,
                                                     backward_pass                  REAL,
                                                     passes_left                    REAL,
                                                     passes_right                   REAL,
                                                     leftside_pass                  REAL,
                                                     rightside_pass                 REAL,
                                                     total_fwd_zone_pass            REAL,
                                                     accurate_fwd_zone_pass         REAL,
                                                     total_back_zone_pass           REAL,
                                                     accurate_back_zone_pass        REAL,
                                                     total_final_third_passes       REAL,
                                                     successful_final_third_passes  REAL,
                                                     total_long_balls               REAL,
                                                     accurate_long_balls            REAL,
                                                     total_through_ball             REAL,
                                                     accurate_through_ball          REAL,
                                                     total_chipped_pass             REAL,
                                                     accurate_chipped_pass          REAL,
                                                     total_cross                    REAL,
                                                     accurate_cross                 REAL,
                                                     total_cross_nocorner           REAL,
                                                     accurate_cross_nocorner        REAL,
                                                     crosses_18yard                 REAL,
                                                     crosses_18yard_plus            REAL,
                                                     total_layoffs                  REAL,
                                                     accurate_layoffs               REAL,
                                                     total_launches                 REAL,
                                                     accurate_launches              REAL,
                                                     total_flick_on                 REAL,
                                                     accurate_flick_on              REAL,
                                                     total_pull_back                REAL,
                                                     long_pass_own_to_opp           REAL,
                                                     long_pass_own_to_opp_success   REAL,
                                                     freekick_cross                 REAL,
                                                     accurate_freekick_cross        REAL,
                                                     blocked_pass                   REAL,
                                                     PRIMARY KEY (fixture_code, team_code),
                                                     FOREIGN KEY (fixture_code) REFERENCES fixtures(code) ON DELETE CASCADE ON UPDATE CASCADE,
                                                     FOREIGN KEY (team_code)    REFERENCES teams(code)    ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_fsp_team ON fixture_stats_passing (team_code);

CREATE TABLE IF NOT EXISTS fixture_stats_defending (
                                                       fixture_code              INTEGER NOT NULL,
                                                       team_code                 INTEGER NOT NULL,
                                                       total_tackle              REAL,
                                                       won_tackle                REAL,
                                                       total_clearance           REAL,
                                                       effective_clearance       REAL,
                                                       head_clearance            REAL,
                                                       effective_head_clearance  REAL,
                                                       clearance_off_line        REAL,
                                                       interception              REAL,
                                                       interception_won          REAL,
                                                       interceptions_ibox        REAL,
                                                       ball_recovery             REAL,
                                                       blocked_cross             REAL,
                                                       effective_blocked_cross   REAL,
                                                       outfielder_block          REAL,
                                                       six_yard_block            REAL,
                                                       aerial_won                REAL,
                                                       aerial_lost               REAL,
                                                       poss_won_def_3rd          REAL,
                                                       poss_won_mid_3rd          REAL,
                                                       poss_won_att_3rd          REAL,
                                                       fouled_final_third        REAL,
                                                       shield_ball_oop           REAL,
                                                       attempts_conceded_ibox    REAL,
                                                       attempts_conceded_obox    REAL,
                                                       error_lead_to_goal        REAL,
                                                       PRIMARY KEY (fixture_code, team_code),
                                                       FOREIGN KEY (fixture_code) REFERENCES fixtures(code) ON DELETE CASCADE ON UPDATE CASCADE,
                                                       FOREIGN KEY (team_code)    REFERENCES teams(code)    ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_fsd_team ON fixture_stats_defending (team_code);

CREATE TABLE IF NOT EXISTS fixture_stats_goalkeeping (
                                                         fixture_code             INTEGER NOT NULL,
                                                         team_code                INTEGER NOT NULL,
                                                         saves                    REAL,
                                                         saved_ibox               REAL,
                                                         saved_obox               REAL,
                                                         diving_save              REAL,
                                                         punches                  REAL,
                                                         good_high_claim          REAL,
                                                         total_high_claim         REAL,
                                                         total_keeper_sweeper     REAL,
                                                         accurate_keeper_sweeper  REAL,
                                                         keeper_throws            REAL,
                                                         accurate_keeper_throws   REAL,
                                                         goal_kicks               REAL,
                                                         accurate_goal_kicks      REAL,
                                                         PRIMARY KEY (fixture_code, team_code),
                                                         FOREIGN KEY (fixture_code) REFERENCES fixtures(code) ON DELETE CASCADE ON UPDATE CASCADE,
                                                         FOREIGN KEY (team_code)    REFERENCES teams(code)    ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_fsgk_team ON fixture_stats_goalkeeping (team_code);

CREATE TABLE IF NOT EXISTS fixture_stats_set_pieces (
                                                        fixture_code              INTEGER NOT NULL,
                                                        team_code                 INTEGER NOT NULL,
                                                        corner_taken              REAL,
                                                        won_corners               REAL,
                                                        lost_corners              REAL,
                                                        total_corners_intobox     REAL,
                                                        accurate_corners_intobox  REAL,
                                                        total_throws              REAL,
                                                        accurate_throws           REAL,
                                                        total_offside             REAL,
                                                        PRIMARY KEY (fixture_code, team_code),
                                                        FOREIGN KEY (fixture_code) REFERENCES fixtures(code) ON DELETE CASCADE ON UPDATE CASCADE,
                                                        FOREIGN KEY (team_code)    REFERENCES teams(code)    ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_fssp_team ON fixture_stats_set_pieces (team_code);

-- ==========================================
-- 3. TEAM STATS (code-keyed, from PL API)
-- ==========================================

CREATE TABLE IF NOT EXISTS team_standings (
                                              team_code       INTEGER NOT NULL,
                                              season_code     INTEGER NOT NULL,
                                              scope           VARCHAR(10) NOT NULL CHECK (scope IN ('overall', 'home', 'away')),
                                              position        INTEGER,
                                              played          INTEGER,
                                              won             INTEGER,
                                              drawn           INTEGER,
                                              lost            INTEGER,
                                              goals_for       INTEGER,
                                              goals_against   INTEGER,
                                              goal_difference INTEGER GENERATED ALWAYS AS (goals_for - goals_against) STORED,
                                              points          INTEGER,
                                              PRIMARY KEY (team_code, season_code, scope),
                                              FOREIGN KEY (team_code)   REFERENCES teams(code)   ON DELETE CASCADE ON UPDATE CASCADE,
                                              FOREIGN KEY (season_code) REFERENCES seasons(code)  ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_standings_season_scope ON team_standings (season_code, scope, position);

CREATE TABLE IF NOT EXISTS team_stats_general (
                                                  team_code                             INTEGER NOT NULL,
                                                  season_code                           INTEGER NOT NULL,
                                                  games_played                          REAL,
                                                  goals                                 REAL,
                                                  home_goals                            REAL,
                                                  away_goals                            REAL,
                                                  goals_conceded                        REAL,
                                                  goals_conceded_ibox                   REAL,
                                                  goals_conceded_obox                   REAL,
                                                  own_goals_conceded                    REAL,
                                                  clean_sheets                          REAL,
                                                  expected_goals                        REAL,
                                                  expected_goals_on_target              REAL,
                                                  expected_goals_on_target_conceded     REAL,
                                                  expected_goals_freekick               REAL,
                                                  expected_assists                      REAL,
                                                  goal_assists                          REAL,
                                                  goal_conversion                       REAL,
                                                  possession_percentage                 REAL,
                                                  points_gained_from_losing_positions   REAL,
                                                  points_dropped_from_winning_positions REAL,
                                                  total_fouls_won                       REAL,
                                                  total_fouls_conceded                  REAL,
                                                  foul_attempted_tackle                 REAL,
                                                  foul_won_penalty                      REAL,
                                                  penalties_conceded                    REAL,
                                                  yellow_cards                          REAL,
                                                  total_red_cards                       REAL,
                                                  straight_red_cards                    REAL,
                                                  handballs_conceded                    REAL,
                                                  offsides                              REAL,
                                                  overruns                              REAL,
                                                  total_losses_of_possession            REAL,
                                                  recoveries                            REAL,
                                                  PRIMARY KEY (team_code, season_code),
                                                  FOREIGN KEY (team_code)   REFERENCES teams(code)   ON DELETE CASCADE ON UPDATE CASCADE,
                                                  FOREIGN KEY (season_code) REFERENCES seasons(code)  ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE IF NOT EXISTS team_stats_attacking (
                                                    team_code                     INTEGER NOT NULL,
                                                    season_code                   INTEGER NOT NULL,
                                                    total_shots                   REAL,
                                                    shots_on_target_inc_goals     REAL,
                                                    shots_off_target_inc_woodwork REAL,
                                                    blocked_shots                 REAL,
                                                    hit_woodwork                  REAL,
                                                    shooting_accuracy             REAL,
                                                    ibox_target                   REAL,
                                                    ibox_blocked                  REAL,
                                                    obox_target                   REAL,
                                                    obox_blocked                  REAL,
                                                    touches_in_opp_box            REAL,
                                                    right_foot_goals              REAL,
                                                    left_foot_goals               REAL,
                                                    headed_goals                  REAL,
                                                    penalty_goals                 REAL,
                                                    penalty_goals_conceded        REAL,
                                                    set_pieces_goals              REAL,
                                                    attempts_from_set_pieces      REAL,
                                                    freekick_total                REAL,
                                                    key_passes_attempt_assists    REAL,
                                                    successful_dribbles           REAL,
                                                    unsuccessful_dribbles         REAL,
                                                    PRIMARY KEY (team_code, season_code),
                                                    FOREIGN KEY (team_code)   REFERENCES teams(code)   ON DELETE CASCADE ON UPDATE CASCADE,
                                                    FOREIGN KEY (season_code) REFERENCES seasons(code)  ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE IF NOT EXISTS team_stats_passing (
                                                  team_code                            INTEGER NOT NULL,
                                                  season_code                          INTEGER NOT NULL,
                                                  total_passes                         REAL,
                                                  passing_accuracy                     REAL,
                                                  open_play_passes                     REAL,
                                                  successful_open_play_passes          REAL,
                                                  successful_short_passes              REAL,
                                                  unsuccessful_short_passes            REAL,
                                                  successful_long_passes               REAL,
                                                  unsuccessful_long_passes             REAL,
                                                  successful_passes_own_half           REAL,
                                                  unsuccessful_passes_own_half         REAL,
                                                  successful_passes_opp_half           REAL,
                                                  unsuccessful_passes_opp_half         REAL,
                                                  passing_percent_opp_half             REAL,
                                                  successful_layoffs                   REAL,
                                                  unsuccessful_layoffs                 REAL,
                                                  successful_launches                  REAL,
                                                  unsuccessful_launches                REAL,
                                                  successful_crosses_open_play         REAL,
                                                  unsuccessful_crosses_open_play       REAL,
                                                  successful_crosses_and_corners       REAL,
                                                  unsuccessful_crosses_and_corners     REAL,
                                                  crossing_accuracy                    REAL,
                                                  putthrough_blocked_distribution      REAL,
                                                  putthrough_blocked_distribution_won  REAL,
                                                  PRIMARY KEY (team_code, season_code),
                                                  FOREIGN KEY (team_code)   REFERENCES teams(code)   ON DELETE CASCADE ON UPDATE CASCADE,
                                                  FOREIGN KEY (season_code) REFERENCES seasons(code)  ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE IF NOT EXISTS team_stats_defending (
                                                    team_code                INTEGER NOT NULL,
                                                    season_code              INTEGER NOT NULL,
                                                    total_clearances         REAL,
                                                    clearances_off_the_line  REAL,
                                                    interceptions            REAL,
                                                    blocks                   REAL,
                                                    last_player_tackle       REAL,
                                                    tackles_won              REAL,
                                                    tackles_lost             REAL,
                                                    tackle_success           REAL,
                                                    times_tackled            REAL,
                                                    duels                    REAL,
                                                    duels_won                REAL,
                                                    duels_lost               REAL,
                                                    ground_duels             REAL,
                                                    ground_duels_won         REAL,
                                                    ground_duels_lost        REAL,
                                                    aerial_duels             REAL,
                                                    aerial_duels_won         REAL,
                                                    aerial_duels_lost        REAL,
                                                    total_shots_conceded     REAL,
                                                    shots_conceded_ibox      REAL,
                                                    shots_conceded_obox      REAL,
                                                    PRIMARY KEY (team_code, season_code),
                                                    FOREIGN KEY (team_code)   REFERENCES teams(code)   ON DELETE CASCADE ON UPDATE CASCADE,
                                                    FOREIGN KEY (season_code) REFERENCES seasons(code)  ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE IF NOT EXISTS team_stats_goalkeeping (
                                                      team_code                    INTEGER NOT NULL,
                                                      season_code                  INTEGER NOT NULL,
                                                      catches                      REAL,
                                                      drops                        REAL,
                                                      goal_kicks                   REAL,
                                                      gk_successful_distribution   REAL,
                                                      gk_unsuccessful_distribution REAL,
                                                      PRIMARY KEY (team_code, season_code),
                                                      FOREIGN KEY (team_code)   REFERENCES teams(code)   ON DELETE CASCADE ON UPDATE CASCADE,
                                                      FOREIGN KEY (season_code) REFERENCES seasons(code)  ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE IF NOT EXISTS team_stats_set_pieces (
                                                     team_code                        INTEGER NOT NULL,
                                                     season_code                      INTEGER NOT NULL,
                                                     corners_taken_incl_short_corners REAL,
                                                     corners_won                      REAL,
                                                     successful_corners_into_box      REAL,
                                                     unsuccessful_corners_into_box    REAL,
                                                     throw_ins_to_own_player          REAL,
                                                     throw_ins_to_opp_player          REAL,
                                                     PRIMARY KEY (team_code, season_code),
                                                     FOREIGN KEY (team_code)   REFERENCES teams(code)   ON DELETE CASCADE ON UPDATE CASCADE,
                                                     FOREIGN KEY (season_code) REFERENCES seasons(code)  ON DELETE CASCADE ON UPDATE CASCADE
);

-- ==========================================
-- 4. PLAYER STATS (code-keyed, from PL API)
-- ==========================================

CREATE TABLE IF NOT EXISTS player_stats_general (
                                                    player_code                           INTEGER NOT NULL,
                                                    season_code                           INTEGER NOT NULL,
                                                    appearances                           REAL,
                                                    games_played                          REAL,
                                                    starts                                REAL,
                                                    time_played                           REAL,
                                                    substitute_on                         REAL,
                                                    substitute_off                        REAL,
                                                    goals                                 REAL,
                                                    home_goals                            REAL,
                                                    away_goals                            REAL,
                                                    goals_conceded                        REAL,
                                                    goals_conceded_inside_box             REAL,
                                                    goals_conceded_outside_box            REAL,
                                                    own_goal_scored                       REAL,
                                                    clean_sheets                          REAL,
                                                    expected_goals                        REAL,
                                                    expected_goals_on_target              REAL,
                                                    expected_goals_on_target_conceded     REAL,
                                                    expected_goals_freekick               REAL,
                                                    expected_assists                      REAL,
                                                    goal_assists                          REAL,
                                                    second_goal_assists                   REAL,
                                                    total_fouls_won                       REAL,
                                                    total_fouls_conceded                  REAL,
                                                    foul_attempted_tackle                 REAL,
                                                    foul_won_penalty                      REAL,
                                                    penalties_conceded                    REAL,
                                                    yellow_cards                          REAL,
                                                    total_red_cards                       REAL,
                                                    straight_red_cards                    REAL,
                                                    red_cards_2nd_yellow                  REAL,
                                                    handballs_conceded                    REAL,
                                                    offsides                              REAL,
                                                    overruns                              REAL,
                                                    total_losses_of_possession            REAL,
                                                    recoveries                            REAL,
                                                    penalty_goals_conceded                REAL,
                                                    touches                               REAL,
                                                    PRIMARY KEY (player_code, season_code),
                                                    FOREIGN KEY (player_code) REFERENCES players(code) ON DELETE CASCADE ON UPDATE CASCADE,
                                                    FOREIGN KEY (season_code) REFERENCES seasons(code)  ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_psg_season ON player_stats_general (season_code);

CREATE TABLE IF NOT EXISTS player_stats_attacking (
                                                      player_code                           INTEGER NOT NULL,
                                                      season_code                           INTEGER NOT NULL,
                                                      total_shots                           REAL,
                                                      shots_on_target_inc_goals             REAL,
                                                      shots_off_target_inc_woodwork         REAL,
                                                      blocked_shots                         REAL,
                                                      hit_woodwork                          REAL,
                                                      ibox_target                           REAL,
                                                      ibox_blocked                          REAL,
                                                      obox_target                           REAL,
                                                      obox_blocked                          REAL,
                                                      total_touches_in_opposition_box       REAL,
                                                      right_foot_goals                      REAL,
                                                      left_foot_goals                       REAL,
                                                      headed_goals                          REAL,
                                                      penalty_goals                         REAL,
                                                      set_pieces_goals                      REAL,
                                                      set_piece_goals                       REAL,
                                                      assists_intentional                   REAL,
                                                      attempts_from_set_pieces              REAL,
                                                      freekick_total                        REAL,
                                                      key_passes_attempt_assists            REAL,
                                                      successful_dribbles                   REAL,
                                                      unsuccessful_dribbles                 REAL,
                                                      winning_goal                          REAL,
                                                      goals_from_inside_box                 REAL,
                                                      goals_from_outside_box                REAL,
                                                      penalties_taken                       REAL,
                                                      other_goals                           REAL,
                                                      PRIMARY KEY (player_code, season_code),
                                                      FOREIGN KEY (player_code) REFERENCES players(code) ON DELETE CASCADE ON UPDATE CASCADE,
                                                      FOREIGN KEY (season_code) REFERENCES seasons(code)  ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_psa_season ON player_stats_attacking (season_code);

CREATE TABLE IF NOT EXISTS player_stats_passing (
                                                    player_code                                             INTEGER NOT NULL,
                                                    season_code                                             INTEGER NOT NULL,
                                                    total_passes                                            REAL,
                                                    open_play_passes                                        REAL,
                                                    successful_open_play_passes                             REAL,
                                                    successful_short_passes                                 REAL,
                                                    unsuccessful_short_passes                               REAL,
                                                    successful_long_passes                                  REAL,
                                                    unsuccessful_long_passes                                REAL,
                                                    successful_passes_own_half                              REAL,
                                                    unsuccessful_passes_own_half                            REAL,
                                                    successful_passes_opposition_half                       REAL,
                                                    successful_layoffs                                      REAL,
                                                    unsuccessful_layoffs                                    REAL,
                                                    successful_launches                                     REAL,
                                                    unsuccessful_launches                                   REAL,
                                                    successful_crosses_open_play                            REAL,
                                                    unsuccessful_crosses_open_play                          REAL,
                                                    successful_crosses_and_corners                          REAL,
                                                    unsuccessful_crosses_and_corners                        REAL,
                                                    putthrough_blocked_distribution                         REAL,
                                                    putthrough_blocked_distribution_won                     REAL,
                                                    through_balls                                           REAL,
                                                    forward_passes                                          REAL,
                                                    backward_passes                                         REAL,
                                                    leftside_passes                                         REAL,
                                                    rightside_passes                                        REAL,
                                                    total_successful_passes_excl_crosses_and_corners        REAL,
                                                    total_unsuccessful_passes_excl_crosses_and_corners      REAL,
                                                    PRIMARY KEY (player_code, season_code),
                                                    FOREIGN KEY (player_code) REFERENCES players(code) ON DELETE CASCADE ON UPDATE CASCADE,
                                                    FOREIGN KEY (season_code) REFERENCES seasons(code)  ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_psp_season ON player_stats_passing (season_code);

CREATE TABLE IF NOT EXISTS player_stats_defending (
                                                      player_code               INTEGER NOT NULL,
                                                      season_code               INTEGER NOT NULL,
                                                      total_clearances          REAL,
                                                      clearances_off_the_line   REAL,
                                                      interceptions             REAL,
                                                      blocks                    REAL,
                                                      last_player_tackle        REAL,
                                                      total_tackles             REAL,
                                                      tackles_won               REAL,
                                                      tackles_lost              REAL,
                                                      times_tackled             REAL,
                                                      duels                     REAL,
                                                      duels_won                 REAL,
                                                      duels_lost                REAL,
                                                      ground_duels              REAL,
                                                      ground_duels_won          REAL,
                                                      ground_duels_lost         REAL,
                                                      aerial_duels              REAL,
                                                      aerial_duels_won          REAL,
                                                      aerial_duels_lost         REAL,
                                                      fifty_fifty               REAL,
                                                      successful_fifty_fifty    REAL,
                                                      PRIMARY KEY (player_code, season_code),
                                                      FOREIGN KEY (player_code) REFERENCES players(code) ON DELETE CASCADE ON UPDATE CASCADE,
                                                      FOREIGN KEY (season_code) REFERENCES seasons(code)  ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_psd_season ON player_stats_defending (season_code);

CREATE TABLE IF NOT EXISTS player_stats_goalkeeping (
                                                        player_code                  INTEGER NOT NULL,
                                                        season_code                  INTEGER NOT NULL,
                                                        catches                      REAL,
                                                        drops                        REAL,
                                                        goal_kicks                   REAL,
                                                        gk_successful_distribution   REAL,
                                                        gk_unsuccessful_distribution REAL,
                                                        saves_made                   REAL,
                                                        saves_made_caught            REAL,
                                                        saves_made_parried           REAL,
                                                        saves_made_from_inside_box   REAL,
                                                        saves_made_from_outside_box  REAL,
                                                        saves_from_penalty           REAL,
                                                        penalties_saved              REAL,
                                                        penalties_faced              REAL,
                                                        punches                      REAL,
                                                        crosses_not_claimed          REAL,
                                                        goalkeeper_smother           REAL,
                                                        PRIMARY KEY (player_code, season_code),
                                                        FOREIGN KEY (player_code) REFERENCES players(code) ON DELETE CASCADE ON UPDATE CASCADE,
                                                        FOREIGN KEY (season_code) REFERENCES seasons(code)  ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_psgk_season ON player_stats_goalkeeping (season_code);

CREATE TABLE IF NOT EXISTS player_stats_set_pieces (
                                                       player_code                      INTEGER NOT NULL,
                                                       season_code                      INTEGER NOT NULL,
                                                       corners_taken_incl_short_corners REAL,
                                                       corners_won                      REAL,
                                                       successful_corners_into_box      REAL,
                                                       unsuccessful_corners_into_box    REAL,
                                                       throw_ins_to_own_player          REAL,
                                                       throw_ins_to_opposition_player   REAL,
                                                       PRIMARY KEY (player_code, season_code),
                                                       FOREIGN KEY (player_code) REFERENCES players(code) ON DELETE CASCADE ON UPDATE CASCADE,
                                                       FOREIGN KEY (season_code) REFERENCES seasons(code)  ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_pssp_season ON player_stats_set_pieces (season_code);

-- ====================================================================
-- FPL-SPECIFIC TABLES
-- ====================================================================

CREATE TABLE IF NOT EXISTS scoring_rules (
                                             stat            VARCHAR(255) NOT NULL,
                                             element_type_id INTEGER NOT NULL,
                                             points          INTEGER NOT NULL DEFAULT 0,
                                             PRIMARY KEY (stat, element_type_id),
                                             FOREIGN KEY (element_type_id) REFERENCES element_types(id) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_scoring_rules_type ON scoring_rules (element_type_id);

CREATE TABLE IF NOT EXISTS teams_strengths (
                                               team_code              INTEGER NOT NULL,
                                               season_code            INTEGER NOT NULL,
                                               strength               INTEGER,
                                               strength_overall_home  INTEGER,
                                               strength_overall_away  INTEGER,
                                               strength_attack_home   INTEGER,
                                               strength_attack_away   INTEGER,
                                               strength_defence_home  INTEGER,
                                               strength_defence_away  INTEGER,
                                               unavailable            BOOLEAN NOT NULL DEFAULT FALSE,
                                               updated_at             TIMESTAMPTZ NOT NULL DEFAULT NOW(),
                                               PRIMARY KEY (team_code, season_code),
                                               FOREIGN KEY (team_code)   REFERENCES teams(code)   ON DELETE CASCADE ON UPDATE CASCADE,
                                               FOREIGN KEY (season_code) REFERENCES seasons(code)  ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE IF NOT EXISTS player_bootstrap (
                                                season_code    INTEGER NOT NULL,
                                                player_code    INTEGER NOT NULL,
                                                team_code      INTEGER NOT NULL,
                                                player_id      INTEGER NOT NULL,
                                                element_type_id INTEGER,
                                                status         VARCHAR(50) NOT NULL,
                                                now_cost       INTEGER NOT NULL,
                                                photo          VARCHAR(255),
                                                squad_number   INTEGER,
                                                can_transact   BOOLEAN,
                                                can_select     BOOLEAN,
                                                in_dreamteam   BOOLEAN,
                                                dreamteam_count INTEGER,
                                                special        BOOLEAN,
                                                removed        BOOLEAN,
                                                unavailable    BOOLEAN,
                                                updated_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
                                                PRIMARY KEY (player_code, season_code),
                                                FOREIGN KEY (season_code) REFERENCES seasons(code) ON DELETE CASCADE ON UPDATE CASCADE,
                                                FOREIGN KEY (team_code)   REFERENCES teams(code)   ON DELETE CASCADE ON UPDATE CASCADE,
                                                FOREIGN KEY (player_code) REFERENCES players(code) ON DELETE CASCADE ON UPDATE CASCADE,
                                                FOREIGN KEY (element_type_id) REFERENCES element_types(id) ON DELETE SET NULL ON UPDATE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_bootstrap_team      ON player_bootstrap (team_code, season_code);
CREATE INDEX IF NOT EXISTS idx_bootstrap_player_id ON player_bootstrap (player_id, season_code);
CREATE INDEX IF NOT EXISTS idx_bootstrap_status    ON player_bootstrap (season_code, status);

CREATE TABLE IF NOT EXISTS player_stats_fixture (
                                                    season_code                     INTEGER NOT NULL,
                                                    player_code                     INTEGER NOT NULL,
                                                    team_code                       INTEGER NOT NULL,
                                                    fixture_id                      INTEGER NOT NULL,
                                                    total_points                    INTEGER,
                                                    minutes                         INTEGER,
                                                    goals_scored                    INTEGER,
                                                    assists                         INTEGER,
                                                    clean_sheets                    INTEGER,
                                                    goals_conceded                  INTEGER,
                                                    own_goals                       INTEGER,
                                                    penalties_saved                 INTEGER,
                                                    penalties_missed                INTEGER,
                                                    yellow_cards                    INTEGER,
                                                    red_cards                       INTEGER,
                                                    saves                           INTEGER,
                                                    bonus                           INTEGER,
                                                    bps                             INTEGER,
                                                    influence                       DECIMAL,
                                                    creativity                      DECIMAL,
                                                    threat                          DECIMAL,
                                                    ict_index                       DECIMAL,
                                                    clearances_blocks_interceptions INTEGER,
                                                    recoveries                      INTEGER,
                                                    tackles                         INTEGER,
                                                    defensive_contribution          INTEGER,
                                                    expected_goals                  DECIMAL,
                                                    expected_assists                DECIMAL,
                                                    expected_goal_involvements      DECIMAL,
                                                    expected_goals_conceded         DECIMAL,
                                                    value                           INTEGER,
                                                    transfers_balance               INTEGER,
                                                    selected                        INTEGER,
                                                    transfers_in                    INTEGER,
                                                    transfers_out                   INTEGER,
                                                    updated_at                      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
                                                    PRIMARY KEY (player_code, season_code, fixture_id),
                                                    FOREIGN KEY (player_code, season_code) REFERENCES player_bootstrap(player_code, season_code) ON DELETE CASCADE ON UPDATE CASCADE,
                                                    FOREIGN KEY (fixture_id, season_code) REFERENCES fixture_seasons(fixture_id, season_code) ON DELETE CASCADE ON UPDATE CASCADE,
                                                    FOREIGN KEY (team_code) REFERENCES teams(code) ON DELETE CASCADE ON UPDATE CASCADE,
                                                    FOREIGN KEY (season_code) REFERENCES seasons(code) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_psf_fixture      ON player_stats_fixture (fixture_id, season_code);
CREATE INDEX IF NOT EXISTS idx_psf_team         ON player_stats_fixture (team_code, season_code);
CREATE INDEX IF NOT EXISTS idx_psf_team_fixture ON player_stats_fixture (team_code, season_code, fixture_id);

-- ==========================================
-- SEED DATA
-- ==========================================

INSERT INTO element_types VALUES
                              (1, 'Goalkeepers', 'GKP', 'Goalkeeper', 'GKP', 2, NULL, NULL, 1, 1, TRUE, '{12}'),
                              (2, 'Defenders',   'DEF', 'Defender',   'DEF', 5, NULL, NULL, 3, 5, FALSE, '{}'),
                              (3, 'Midfielders', 'MID', 'Midfielder', 'MID', 5, NULL, NULL, 2, 5, FALSE, '{}'),
                              (4, 'Forwards',    'FWD', 'Forward',    'FWD', 3, NULL, NULL, 1, 3, FALSE, '{}');

INSERT INTO seasons VALUES
                        (2025, '2025/26', TRUE),
                        (2024, '2024/25', FALSE),
                        (2023, '2023/24', FALSE),
                        (2022, '2022/23', FALSE),
                        (2021, '2021/22', FALSE),
                        (2020, '2020/21', FALSE),
                        (2019, '2019/20', FALSE),
                        (2018, '2018/19', FALSE),
                        (2017, '2017/18', FALSE),
                        (2016, '2016/17', FALSE),
                        (2015, '2015/16', FALSE),
                        (2014, '2014/15', FALSE),
                        (2013, '2013/14', FALSE),
                        (2012, '2012/13', FALSE),
                        (2011, '2011/12', FALSE),
                        (2010, '2010/11', FALSE),
                        (2009, '2009/10', FALSE),
                        (2008, '2008/09', FALSE),
                        (2007, '2007/08', FALSE),
                        (2006, '2006/07', FALSE);

INSERT INTO scoring_rules (stat, element_type_id, points) VALUES
                                                              ('long_play',        1, 2), ('long_play',        2, 2), ('long_play',        3, 2), ('long_play',        4, 2),
                                                              ('short_play',       1, 1), ('short_play',       2, 1), ('short_play',       3, 1), ('short_play',       4, 1),
                                                              ('goals_scored',     1, 6), ('goals_scored',     2, 6), ('goals_scored',     3, 5), ('goals_scored',     4, 4),
                                                              ('assists',          1, 3), ('assists',          2, 3), ('assists',          3, 3), ('assists',          4, 3),
                                                              ('clean_sheets',     1, 4), ('clean_sheets',     2, 4), ('clean_sheets',     3, 1), ('clean_sheets',     4, 0),
                                                              ('goals_conceded',   1,-1), ('goals_conceded',   2,-1), ('goals_conceded',   3, 0), ('goals_conceded',   4, 0),
                                                              ('saves',            1, 1), ('saves',            2, 0), ('saves',            3, 0), ('saves',            4, 0),
                                                              ('penalties_saved',  1, 5), ('penalties_saved',  2, 5), ('penalties_saved',  3, 5), ('penalties_saved',  4, 5),
                                                              ('penalties_missed', 1,-2), ('penalties_missed', 2,-2), ('penalties_missed', 3,-2), ('penalties_missed', 4,-2),
                                                              ('yellow_cards',     1,-1), ('yellow_cards',     2,-1), ('yellow_cards',     3,-1), ('yellow_cards',     4,-1),
                                                              ('red_cards',        1,-3), ('red_cards',        2,-3), ('red_cards',        3,-3), ('red_cards',        4,-3),
                                                              ('own_goals',        1,-2), ('own_goals',        2,-2), ('own_goals',        3,-2), ('own_goals',        4,-2),
                                                              ('bonus',            1, 1), ('bonus',            2, 1), ('bonus',            3, 1), ('bonus',            4, 1);