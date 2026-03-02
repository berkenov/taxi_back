-- Тестовые пользователи: Samat (пассажир), Nurik (водитель)
-- Телефоны: 77000000000, 77000000001
-- Код для входа: 0000 (WhatsApp не вызывается)

INSERT INTO users (id, phone, name, role, is_active) VALUES
    ('a0000000-0000-0000-0000-000000000001', '77000000000', 'Samat', 'passenger', true),
    ('a0000000-0000-0000-0000-000000000002', '77000000001', 'Nurik', 'driver', true)
ON CONFLICT (phone) DO NOTHING;

-- Авто для Nurik (Toyota Camry, 777 ABC 01)
INSERT INTO cars (id, driver_id, model, number, color)
VALUES ('b0000000-0000-0000-0000-000000000001', 'a0000000-0000-0000-0000-000000000002', 'Toyota Camry', '777 ABC 01', 'Белый')
ON CONFLICT (id) DO NOTHING;
