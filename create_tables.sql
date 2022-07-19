CREATE TABLE TSC_APPLICATIONS (
  APP_ID TEXT NOT NULL,
APP_NAME TEXT NOT NULL,
APP_URL TEXT,
APP_DESC TEXT,
APP_LOGO BLOB
);

CREATE OR REPLACE TABLE TSC_CLIENT_LOGS (
  LOG_ID INTEGER PRIMARY KEY AUTOINCREMENT,
  APP_ID TEXT,
SESSION_ID TEXT,
LOG_LEVEL TEXT,
URL TEXT,
MSG TEXT,
STACKTRACE TEXT,
TIMESTAMP TEXT,
USERAGENT TEXT,
CLIENT_IP TEXT,
REMOTE_IP TEXT
);

