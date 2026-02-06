DROP TABLE IF EXISTS coin_transfer;
DROP TABLE IF EXISTS coin_entry;
DROP TABLE IF EXISTS coin;

ALTER TABLE follow_relationship
  DROP CONSTRAINT IF EXISTS follow_relationship_follower_uid_fkey;
ALTER TABLE follow_relationship
  DROP CONSTRAINT IF EXISTS follow_relationship_followed_uid_fkey;