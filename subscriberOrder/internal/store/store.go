package store

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"

	"awesomeProject5/subscriberOrder/internal/models"
)

type Store struct {
	db *sql.DB
}

func NewStore() *Store {
	return &Store{}
}

func (s *Store) OpenDB() error {
	url := os.Getenv("POSTGRES_URL")

	db, err := sql.Open("postgres", url)
	if err != nil {
		return err
	}
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}
	s.db = db
	return err
}

func (s *Store) CloseConnectionDB() {
	if err := s.db.Close(); err != nil {
		log.Fatal(err)
	}
}

func (s *Store) AddOrder(m *models.Order) error {
	order := &models.Order{
		OrderUid:          m.OrderUid,
		TrackNumber:       m.TrackNumber,
		Entry:             m.Entry,
		Locale:            m.Locale,
		InternalSignature: m.InternalSignature,
		CustomerId:        m.CustomerId,
		DeliveryService:   m.DeliveryService,
		Shardkey:          m.Shardkey,
		SmId:              m.SmId,
		DateCreated:       m.DateCreated,
		OofShard:          m.OofShard,
	}
	query := `insert into orders (order_uid, track_number, entry, locale, internal_signature, customer_id, delivery_service, shardkey,sm_id,date_created, oof_shard)
				values ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)
				RETURNING order_id`

	err := s.db.QueryRow(query, order.OrderUid, order.TrackNumber, order.Entry, order.Locale, order.InternalSignature, order.CustomerId, order.DeliveryService, order.Shardkey, order.SmId, order.DateCreated, order.OofShard).Scan(&models.OrderID)

	if err != nil {
		return err
	}
	return nil
}

func (s *Store) AddDelivery(m *models.Delivery) error {
	delivery := &models.Delivery{
		Name:    m.Name,
		Phone:   m.Phone,
		Zip:     m.Zip,
		City:    m.City,
		Address: m.Address,
		Region:  m.Region,
		Email:   m.Email,
	}

	query := `insert into deliveries (name,phone,zip,city,address,region,email, order_id)
				values ($1,$2,$3,$4,$5,$6,$7,$8)
`

	_, err := s.db.Exec(query, delivery.Name, delivery.Phone, delivery.Zip, delivery.City, delivery.Address, delivery.Region, delivery.Email, models.OrderID)
	if err != nil {
		return err
	}
	return nil

}

func (s *Store) AddPayments(m *models.Payment) error {
	payments := &models.Payment{
		Transaction:  m.Transaction,
		RequestID:    m.RequestID,
		Currency:     m.Currency,
		Provider:     m.Provider,
		Amount:       m.Amount,
		PaymentDT:    m.PaymentDT,
		Bank:         m.Bank,
		DeliveryCost: m.DeliveryCost,
		GoodsTotal:   m.GoodsTotal,
		CustomFee:    m.CustomFee,
	}

	query := `insert into payments (transaction, request_id, currency, provider, amount, payment_dt,bank, delivery_cost,goods_total,custom_fee,order_id)
				values ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)`
	_, err := s.db.Exec(query, payments.Transaction, payments.RequestID, payments.Currency, payments.Provider, payments.Amount, payments.PaymentDT, payments.Bank, payments.DeliveryCost, payments.GoodsTotal, payments.CustomFee, models.OrderID)
	if err != nil {
		return err
	}
	return nil
}

func (s *Store) AddItems(items *[]models.Item) error {

	query := `insert into items (chrt_id, track_number, price, rid, name, sale, size, total_price,nm_id,brand, status, order_id)
				values ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)
`
	for _, item := range *items {
		_, err := s.db.Exec(query, item.ChrtId, item.TrackNumber, item.Price, item.Rid, item.Name, item.Sale, item.Size, item.TotalPrice, item.NmId, item.Brand, item.Status, models.OrderID)
		if err != nil {
			return err
		}
	}
	return nil

}

func (s *Store) GetByID(id int) (models.Order, error) {
	query := fmt.Sprintf(`SELECT order_uid, "orders".track_number, entry,locale,internal_signature,customer_id,delivery_service,shardkey,sm_id,date_created,oof_shard,
       chrt_id, "items".track_number,price,rid,"items".name,sale,size,total_price,nm_id,brand,status,  "deliveries".name,phone,zip,city,address,region,email,
       transaction, request_id,currency,provider,amount,payment_dt,bank,delivery_cost,goods_total,custom_fee
	FROM orders
	JOIN items ON orders.order_id = items.order_id
	JOIN deliveries ON orders.order_id = deliveries.order_id
	JOIN payments ON orders.order_id = payments.order_id
	WHERE orders.order_id = %d
	`, id)

	rows, err := s.db.Query(query)
	if err != nil {
		return models.Order{}, err
	}
	defer rows.Close()
	var order models.Order
	var delivery models.Delivery
	var payment models.Payment
	var items []models.Item

	for rows.Next() {
		var item models.Item
		err = rows.Scan(
			&order.OrderUid, &order.TrackNumber, &order.Entry, &order.Locale, &order.InternalSignature,
			&order.CustomerId, &order.DeliveryService, &order.Shardkey, &order.SmId, &order.DateCreated, &order.OofShard,
			&item.ChrtId, &item.TrackNumber, &item.Price, &item.Rid,
			&item.Name, &item.Sale, &item.Size, &item.TotalPrice,
			&item.NmId, &item.Brand, &item.Status,
			&delivery.Name, &delivery.Phone, &delivery.Zip, &delivery.City, &delivery.Address, &delivery.Region, &delivery.Email,
			&payment.Transaction, &payment.RequestID, &payment.Currency, &payment.Provider, &payment.Amount, &payment.PaymentDT, &payment.Bank,
			&payment.DeliveryCost, &payment.GoodsTotal, &payment.CustomFee,
		)
		if err != nil {
			return models.Order{}, err
		}
		items = append(items, item)
	}

	order.Delivery = &delivery
	order.Payment = &payment
	order.Items = &items

	return order, nil
}

func (s *Store) CounterIDs() (int, error) {
	query := `SELECT COUNT(*) FROM orders`

	var count int

	err := s.db.QueryRow(query).Scan(count)
	if err != nil {
		return 0, err
	}

	return count, nil
}
