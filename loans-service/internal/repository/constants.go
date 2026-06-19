package repository

// Queries
const (
	QueryCreateLoan = `
		INSERT INTO loans (user_id, book_id, loaned_at, status)
		VALUES ($1, $2, $3, 'active')
		RETURNING id, user_id, book_id, loaned_at, returned_at, status`

	QueryUpdateStatus = `
    	UPDATE loans SET returned_at = $1, status = $2
    	WHERE id = $3 AND status = 'active'
    	RETURNING id, user_id, book_id, loaned_at, returned_at, status`

	QueryGetActiveByUser = `
		SELECT id, user_id, book_id, loaned_at, returned_at, status
		FROM loans WHERE user_id = $1 AND status = 'active'`

	QueryGetHistoryByUser = `
		SELECT id, user_id, book_id, loaned_at, returned_at, status
		FROM loans WHERE user_id = $1 ORDER BY loaned_at DESC`

	QueryGetByID = `
    	SELECT id, user_id, book_id, loaned_at, returned_at, status
    	FROM loans WHERE id = $1`
)
