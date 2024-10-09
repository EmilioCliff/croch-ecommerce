-- name: GetBlog :one
SELECT * FROM blogs
WHERE id = ? LIMIT 1;

-- name: GetBlogByAuthor :one
SELECT * FROM blogs
WHERE author = ? LIMIT 1;

-- name: ListBlogs :many
SELECT * FROM blogs
ORDER BY created_at DESC;

-- name: CreateBlogs :execresult
INSERT INTO blogs (
  author, title, content, img_urls
) VALUES (
  ?, ?, ?, ?
);

-- name: DeleteAuthor :exec
DELETE FROM blogs
WHERE id = ?;

-- name: UpdateBlog :exec
UPDATE blogs
  set title = ?,
  content = ?,
  img_urls = ?
WHERE id = ?;