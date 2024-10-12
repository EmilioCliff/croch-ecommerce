-- name: GetBlog :one
SELECT * FROM blogs
WHERE id = ? LIMIT 1;

-- name: GetBlogsByAuthor :many
SELECT * FROM blogs
WHERE author = ?;

-- name: ListBlogs :many
SELECT * FROM blogs
ORDER BY created_at DESC;

-- name: CreateBlog :execresult
INSERT INTO blogs (
  author, title, content, img_urls
) VALUES (
  ?, ?, ?, ?
);

-- name: DeleteBlog :exec
DELETE FROM blogs
WHERE id = ?;

-- name: UpdateBlog :exec
UPDATE blogs
  set title = coalesce(sqlc.narg('title'), title),
  content = coalesce(sqlc.narg('content'), content),
  img_urls = coalesce(sqlc.narg('img_urls'), img_urls)
WHERE id = sqlc.arg('id');

