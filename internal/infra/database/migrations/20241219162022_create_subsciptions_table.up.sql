CREATE TABLE subscriptions (
                               event_id int NOT NULL,
                               user_id int NOT NULL,
                               PRIMARY KEY (event_id, user_id),
                               FOREIGN KEY (event_id) REFERENCES events(id),
                               FOREIGN KEY (user_id) REFERENCES users(id)
);