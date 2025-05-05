package db

const (
	SelectBlockedIPAddressesSQL = `
	SELECT address, verdict, created_at 
	FROM blocked_ip_address; 
`
	UpsertBlockedIPAddressesSQL = `
	INSERT INTO blocked_ip_address(address, verdict, created_at)
	VALUES (:address, :verdict, :created_at)
	ON CONFLICT DO NOTHING
`
)
