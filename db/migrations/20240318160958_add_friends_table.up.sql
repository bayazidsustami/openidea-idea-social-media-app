CREATE TABLE IF NOT EXISTS friends (
  friend_id SERIAL PRIMARY KEY NOT NULL,
	user_id_requester UUID NOT NULL,
	user_id_accepter UUID NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (user_id_requester) REFERENCES users (user_id) ON DELETE CASCADE,
  FOREIGN KEY (user_id_accepter) REFERENCES users (user_id) ON DELETE CASCADE,
  CONSTRAINT unique_friends_combination UNIQUE (user_id_requester, user_id_accepter)
);