package dbstore

const insertQuery string = `INSERT INTO links (short, original, correlation_id, user_id)
				VALUES ($1, $2, $3, $4)
				RETURNING short, original, correlation_id, user_id;
				`
