CREATE TABLE IF NOT EXISTS friends (
  friend_id SERIAL PRIMARY KEY NOT NULL,
	user_id_requester INTEGER NOT NULL,
	user_id_accepter INTEGER NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (user_id_requester) REFERENCES users (user_id) ON DELETE CASCADE,
  FOREIGN KEY (user_id_accepter) REFERENCES users (user_id) ON DELETE CASCADE
);

CREATE UNIQUE INDEX unique_requester_accepter_index ON friends (user_id_requester, user_id_accepter)