SET foreign_key_checks = 0;

TRUNCATE user;
INSERT INTO user (id, lower_name, name, full_name, email, keep_email_private, email_notifications_preference, passwd, passwd_hash_algo, must_change_password, login_type, login_source, login_name, type, location, website, rands, salt, language, description, created_unix, updated_unix, last_login_unix, last_repo_visibility, max_repo_creation, is_active, is_admin, is_restricted, allow_git_hook, allow_import_local, allow_create_organization, prohibit_login, avatar, avatar_email, use_custom_avatar, num_followers, num_following, num_stars, num_repos, num_teams, num_members, visibility, repo_admin_change_team_access, diff_view_style, theme, keep_activity_private) VALUES (1,'terraform','terraform','','terraform@example.com',0,'enabled','58c34406ea25d3662ac0336579642900635b741723349be758ad6994a96e87a4fbe064eb0f829fc4498bcc1fe057d0003d55','pbkdf2',0,0,0,'',0,'','','9d3b614be9727bd47914e582bd26b821','74ec824754a7a8f31c8421bb9270e974','en-US','',1651421906,1651422082,1651422082,0,-1,1,1,0,0,0,0,0,'','terraform@example.com',0,0,0,0,3,0,0,0,0,'','auto',0);

TRUNCATE oauth2_application;
INSERT INTO oauth2_application (id, uid, name, client_id, client_secret, redirect_uris, created_unix, updated_unix) VALUES (1,1,'drone','9819efc6-6716-4568-90bc-4f28757f6721','$2a$10$kanEdHCvIdRjw3eHSzHmyOCq06YZGlPA9BgiZd6hg9QetxsduBSQO','[\"http://localhost:8000/login\"]',1651422130,1651422130);

TRUNCATE oauth2_authorization_code;

TRUNCATE oauth2_grant;
INSERT INTO oauth2_grant (id, user_id, application_id, counter, scope, nonce, created_unix, updated_unix) VALUES (1,1,1,0,'','',1651422201,1651422201);

SET foreign_key_checks = 1;