package storage

import (
	"database/sql"
	"errors"

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
	track_number		VARCHAR(32),
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

	orderRepoMemory := NewOderRepoMemory()

	queryGetItems := `SELECT * FROM item`
	queryGetOrderItems := `SELECT * FROM order_items`
	queryGetOrders := `SELECT * FROM orders INNER JOIN delivery ON orders.delivery_id = delivery.delivery_id INNER JOIN payment ON orders.payment_id = payment.transaction`

	itemIDToItem := make(map[string]Item)
	var items []Item
	err = db.Select(&items, queryGetItems)
	if err != nil {
		return nil, err
	}
	for _, item := range items {
		itemIDToItem[item.RID] = item
	}

	orderIDToItems := make(map[string][]Item)
	var orderItems []OrderItem
	err = db.Select(&orderItems, queryGetOrderItems)
	if err != nil {
		return nil, err
	}
	for _, orderItem := range orderItems {
		orderIDToItems[orderItem.OrderID] = append(orderIDToItems[orderItem.OrderID], itemIDToItem[orderItem.ItemID])
	}

	var orders []Order
	rows, err := db.Query(queryGetOrders)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var o Order
		err = rows.Scan(&o.OrderUID, &o.DeliveryID, &o.PaymentID, &o.TrackNumber, &o.Entry, &o.Locale, &o.InternalSignature, &o.CustomerID, &o.DeliveryService, &o.ShardKey, &o.SmID, &o.DateCreated, &o.OofShard,
			&o.Delivery.DeliveryID, &o.Delivery.Name, &o.Delivery.Phone, &o.Delivery.Zip, &o.Delivery.City, &o.Delivery.Address, &o.Delivery.Region, &o.Delivery.Email,
			&o.Payment.Transaction, &o.Payment.RequestID, &o.Payment.Currency, &o.Payment.Provider, &o.Payment.Amount, &o.Payment.PaymentDt, &o.Payment.Bank, &o.Payment.DeliveryCost, &o.Payment.GoodsTotal, &o.Payment.CustomFee)
		if err != nil {
			return nil, err
		}
		o.Items = append(o.Items, orderIDToItems[o.OrderUID]...)

		orders = append(orders, o)
	}

	//fmt.Println(orders)

	for _, order := range orders {
		orderRepoMemory.idToOrderMap[order.OrderUID] = order
	}

	r := &OrderRepoDB{
		db:              db,
		orderRepoMemory: orderRepoMemory,
	}

	return r, nil
}

func (r *OrderRepoDB) SaveOrder(order Order) error {
	querySaveDelivery := `INSERT INTO delivery (name, phone, zip, city, address, region, email) VALUES ($1, $2, $3, $4, $5, $6, $7) ON CONFLICT DO NOTHING RETURNING delivery_id`
	err := r.db.Get(&order.Delivery.DeliveryID, querySaveDelivery, order.Delivery.Name, order.Delivery.Phone, order.Delivery.Zip, order.Delivery.City, order.Delivery.Address, order.Delivery.Region, order.Delivery.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			queryGetDeliveryID := `SELECT delivery_id FROM delivery WHERE name = $1 and phone = $2 and zip = $3 and city = $4 and address = $5 and region = $6 and email = $7`
			err := r.db.Get(&order.Delivery.DeliveryID, queryGetDeliveryID, order.Delivery.Name, order.Delivery.Phone, order.Delivery.Zip, order.Delivery.City, order.Delivery.Address, order.Delivery.Region, order.Delivery.Email)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}
	order.DeliveryID = order.Delivery.DeliveryID

	querySavePayment := `INSERT INTO payment (transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee)
							VALUES (:transaction, :request_id, :currency, :provider, :amount, :payment_dt, :bank, :delivery_cost, :goods_total, :custom_fee)`
	_, err = r.db.NamedExec(querySavePayment, order.Payment)
	if err != nil {
		return err
	}
	order.PaymentID = order.Payment.Transaction

	querySaveItems := `INSERT INTO item (chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status)
							VALUES (:chrt_id, :track_number, :price, :rid, :name, :sale, :size, :total_price, :nm_id, :brand, :status) ON CONFLICT DO NOTHING`
	_, err = r.db.NamedExec(querySaveItems, order.Items)
	if err != nil {
		return err
	}

	querySaveOrder := `INSERT INTO orders (order_uid, delivery_id, payment_id, track_number, entry, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard)
							VALUES (:order_uid, :delivery_id, :payment_id, :track_number, :entry, :locale, :internal_signature, :customer_id, :delivery_service, :shardkey, :sm_id, :date_created, :oof_shard)`
	_, err = r.db.NamedExec(querySaveOrder, order)
	if err != nil {
		return err
	}

	querySaveOrderItems := `INSERT INTO order_items (order_id, item_id) VALUES ($1, $2)`
	for _, item := range order.Items {
		_, err := r.db.Exec(querySaveOrderItems, order.OrderUID, item.RID)
		if err != nil {
			return err
		}
	}

	err = r.orderRepoMemory.SaveOrder(order)
	if err != nil {
		return err
	}

	return nil
}

func (r *OrderRepoDB) GetOrderByID(orderID string) (Order, error) {
	return r.orderRepoMemory.GetOrderByID(orderID)
}
