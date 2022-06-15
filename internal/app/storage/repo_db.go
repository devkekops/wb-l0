package storage

import (
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
)

var schema = `
CREATE TABLE IF NOT EXISTS delivery(
	delivery_id		SERIAL PRIMARY KEY,
	name			VARCHAR(32) NOT NULL,
	phone			VARCHAR(12),
	zip				VARCHAR(12),
	city			VARCHAR(32),
	address			VARCHAR(32),
	region			VARCHAR(32),
	email			VARCHAR(32),
	UNIQUE (name, phone, zip, city, address, region, email)
);
CREATE TABLE IF NOT EXISTS payment(
	transaction		VARCHAR(32) PRIMARY KEY,
	request_id		VARCHAR(32),
	currency		VARCHAR(12),
	provider		VARCHAR(12),
	amount 			INTEGER,
	payment_dt		INTEGER,
	bank			VARCHAR(12),
	delivery_cost	INTEGER,
	goods_total		INTEGER,
	custom_fee		INTEGER
);
CREATE TABLE IF NOT EXISTS item(
	chrt_id			INTEGER,
	track_number	VARCHAR(32),
	price			INTEGER,
	rid				VARCHAR(32) PRIMARY KEY,
	name			TEXT,
	sale			INTEGER,
	size			VARCHAR(12),
	total_price		INTEGER,
	nm_id			INTEGER,
	brand			VARCHAR(32),
	status			INTEGER
);
CREATE TABLE IF NOT EXISTS orders(
	order_uid			VARCHAR(32) NOT NULL PRIMARY KEY,
	delivery_id			INTEGER NOT NULL REFERENCES delivery(delivery_id),
	payment_id			VARCHAR(32) NOT NULL REFERENCES payment(transaction),	
	track_numer			VARCHAR(32),
	entry				VARCHAR(32),
	locale				VARCHAR(12),
	internal_signature	VARCHAR(12),
	customer_id			VARCHAR(32),
	delivery_service	VARCHAR(12),
	shardkey			VARCHAR(12),
	sm_id				INTEGER,
	date_created		TIMESTAMP,
	oof_shard			VARCHAR(12)
);
CREATE TABLE IF NOT EXISTS order_items(
	order_id		VARCHAR(32) NOT NULL REFERENCES orders(order_uid),
	item_id			VARCHAR(32) NOT NULL REFERENCES item(rid)
);`

type OrderRepoDB struct {
	db              *sqlx.DB
	orderRepoMemory *OrderRepoMemory
}

func NewOrderRepoDB(databaseURI string) (*OrderRepoDB, error) {
	db, err := sqlx.Connect("pgx", databaseURI)
	if err != nil {
		return nil, err
	}
	db.MustExec(schema)

	r := &OrderRepoDB{
		db:              db,
		orderRepoMemory: NewOderRepoMemory(),
	}

	return r, nil
}

func (r *OrderRepoDB) SaveOrder(order Order) error {
	/*tx, err := r.db.Beginx()
	if err != nil {
		return err
	}
	defer func(tx *sqlx.Tx) {
		err := tx.Rollback()
		if err != nil {
			log.Println(err)
		}
	}(tx)

	var deliveryID int
	querySaveDelivery := `INSERT INTO delivery (name, phone, zip, city, address, region, email) VALUES ($1, $2, $3, $4, $5, $6, $7) ON CONFLICT DO NOTHING RETURNING delivery_id`
	tx.QueryRow(querySaveDelivery, order.Delivery).Scan(&deliveryID)
	fmt.Println(deliveryID)

	querySavePayment := `INSERT INTO payment (transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee)
							VALUES (:transaction, :request_id, :currency, :provider, :amount, :payment_dt, :bank, :delivery_cost, :goods_total, :custom_fee) `
	_, err = tx.NamedExec(querySavePayment, order.Payment)
	if err != nil {
		return err
	}

	querySaveItems := `INSERT INTO item (chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status)
							VALUES (:chrt_id, :track_number, :price, :rid, :name, :sale, :size, :total_price, :nm_id, :brand, :status)`
	_, err = tx.NamedExec(querySaveItems, order.Items)
	if err != nil {
		return err
	}

	querySaveOrder := `INSERT INTO orders (order_uid, delivery_id, payment_id, track_number, entry, locale, internal_signture, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard
							VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $)`
	_, err = tx.Exec(querySaveOrder, order.Items)
	if err != nil {
		return err
	};

	err = tx.Commit()
	if err != nil {
		return err
	}*/

	err := r.orderRepoMemory.SaveOrder(order)
	if err != nil {
		return err
	}

	return nil
}

func (r *OrderRepoDB) GetOrderByID(orderID string) (Order, error) {
	return r.orderRepoMemory.GetOrderByID(orderID)
}
