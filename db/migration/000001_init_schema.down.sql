-- Drop tables in reverse order of creation (respecting dependencies
DROP TABLE IF EXISTS article_comment;
DROP TABLE IF EXISTS article_like;
DROP TABLE IF EXISTS article_tag;
DROP TABLE IF EXISTS article_category;
DROP TABLE IF EXISTS tag;
DROP TABLE IF EXISTS category;
DROP TABLE IF EXISTS article;
DROP TABLE IF EXISTS follow_relationship;
DROP TABLE IF EXISTS users;


-- Drop indexes
DROP INDEX IF EXISTS idx_follow_rel_follower_followed;
DROP INDEX IF EXISTS idx_article_category_article_category;
DROP INDEX IF EXISTS idx_article_tag_article_tag;
DROP INDEX IF EXISTS idx_article_like_uid_article;