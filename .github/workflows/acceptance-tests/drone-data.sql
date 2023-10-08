SET foreign_key_checks = 0;

TRUNCATE users;
INSERT INTO users (user_id, user_login, user_email, user_admin, user_machine, user_active, user_avatar, user_syncing, user_synced, user_created, user_updated, user_last_login, user_oauth_token, user_oauth_refresh, user_oauth_expiry, user_hash) VALUES (1,'terraform','terraform@example.com',1,0,1,'https://secure.gravatar.com/avatar/aad3747997ae1a5cc149ef80ba32a78b?d=identicon',0,1651434882,1651419573,1651419573,1651435613,_binary 'eyJhbGciOiJSUzI1NiIsImtpZCI6ImlmRkw3ZnNGQ2pQLUxnR1lSX0FDd2xFLXpxVDhGbE5Lc085S3NlcFBsRnMiLCJ0eXAiOiJKV1QifQ.eyJnbnQiOjEsInR0IjowLCJleHAiOjE2NTE0MzkyMTMsImlhdCI6MTY1MTQzNTYxM30.bzm2vjjsE-vTqPQyOq2V_Eh_Cab23rtzeVHgxaWEIdO-skmLbEig_rFilo3g-tUsLStHyqCjM2lx6TXXH2WkrmeZVwtdO6aP2Bo4Op0qEeCvEOBVCYOIcVd-UCXzicSLi5oiX2SvbdZH_p9tMDiWFtQxuPClhqENoznGYoRVBMMPgD64hMj4J6706wserjyQdD6PuyoQAahdOnCzNkplhSmBDSMgUCeqgVMp55OqMKcjKjRu8MfpwmM2T01WvKrEBb7E5NpWc2nxVM109yh_e71Benac-ixGahkdpM-ho3IMuWuQxkO-19KYAJIBf0j-Sw311O1JCurFgWmYvCSYWcfE-AvMNokqKhC4KhElSNybo9hR-M4EyWuOxuiWGtUcx-x7BciTYGEDOp57m9GcLeOb7YKRmXYZusWF7p0k5eBunuJ9J2kfOEpz7_TxoxeRtrRFUowOOPDvzOPKiKKYRsL4G2IjEw4R_468HljWfGLjqEn9d2gfBGOjMgKY-85_z7p4DDlbO2GIoQlaNP_tp7_Ym6Kjp16nuwPLGfhnV59rUbA8UNspY5fbqHu28Au3pF8uth0J88VJhLfwSLwqfsKR3L0OgMvCmKhH1pIJ2lqMp4ng4--tWZnYx8DXq8Borxk4Jz8y5yq_pXEafBrI36GOZ_D5xuXFRiE83m79PE0',_binary 'eyJhbGciOiJSUzI1NiIsImtpZCI6ImlmRkw3ZnNGQ2pQLUxnR1lSX0FDd2xFLXpxVDhGbE5Lc085S3NlcFBsRnMiLCJ0eXAiOiJKV1QifQ.eyJnbnQiOjEsInR0IjoxLCJleHAiOjE5MTQyMzU2MTMsImlhdCI6MTY1MTQzNTYxM30.GgNXUrPpKYCUcThjArMyNjYo8wZt_JT_zT79gqNcJPXSMy7KNB5dR8o_dfCZWWcwC5pWAhjHB8JCfrLcKmMfYDxaj1S1cooNypjh_mGd-BHmz5IbBPBaFV-xYj3mArKZ_d2avZBK0NxxWsdwTuOfXdmnCz68Z7C9K6NCN6hGm5ilb2XOAUtyVtCrgURLP6GB1cM7CAyFOiMWF3MAiXZcGkXP-lZ5v1_CEB8I8X8AinsdRsNhRXruJfAvFnNscHZmGIh0r4jJG-_Qe9tZGbnZsuFDQpNaGu4tlD1bOfEM4roEsUXPkceAE9fv9HWr1yhHi4amEgACR79RPccRC3huZUoO0wRGswAJEGdBLLDu6WlqTD9YevLEqveIBJf5zbQyycuVfHBjuSSNaJjowC40qlbkbi8aQZWOB94X45z7jhoaPernU6oMz9-CaoU1yP-qfhR-rJc5q2ItTHbjb8KZ2_D2hyLyNFyLj38fWnUrVVHxcZfgs1jUjg5o4_e5TiM9mQPl1ef3Y-QAr02_viJe7qOvrP2oRqDByJQSczkJsRK9OGK4d-heVh1za-3PsQAwAoAKizKfsIu8saMExG4WFz32VnBcp3dtRKmtNYc_iDyZg1-VJKEBHS3iuX6FNMFcH10oZceIEcG7X07j6owkCel27ptvZVe344S1-HfUFD8',1651439213,'5PVYqFHjdYWpzyOk6PVj9OUQULibBJeL');

TRUNCATE repos;

TRUNCATE secrets;

TRUNCATE cron;

TRUNCATE orgsecrets;

TRUNCATE templates;

SET foreign_key_checks = 1;