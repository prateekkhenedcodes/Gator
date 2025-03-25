-- +goose up 
CREATE TABLE feed_follow (
    id UUID PRIMARY KEY, 
    created_at TIMESTAMP NOT NULL, 
    updated_at TIMESTAMP NOT NULL, 
    user_id UUID NOT NULL, 
    FOREIGN KEY (user_id)
    REFERENCES users(id) ON DELETE CASCADE, 
    feed_id UUID NOT NULL, 
    FOREIGN KEY (feed_id)
    REFERENCES feeds(id) ON DELETE CASCADE,
    UNIQUE(user_id, feed_id)
); 

-- +goose down 
DROP TABLE IF EXISTS feed_follows; 