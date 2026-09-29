CREATE DATABASE IF NOT EXISTS factory;
USE factory;

CREATE TABLE IF NOT EXISTS Employees (
    id INT PRIMARY KEY AUTO_INCREMENT,
    full_name VARCHAR(255) NOT NULL,
    phone VARCHAR(20) NOT NULL UNIQUE,
    address VARCHAR(255) NOT NULL,
    graduation_year INT NOT NULL
);

CREATE TABLE IF NOT EXISTS Qualifications (
    employee_id INT PRIMARY KEY,
    position VARCHAR(255) NOT NULL,
    FOREIGN KEY (employee_id) REFERENCES Employees(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS Customers (
    id INT PRIMARY KEY AUTO_INCREMENT,
    full_name VARCHAR(255) NOT NULL,
    address VARCHAR(255) NOT NULL,
    regular_discount DECIMAL(5,2) DEFAULT 0.00
);

CREATE TABLE IF NOT EXISTS Products (
    id INT PRIMARY KEY AUTO_INCREMENT,
    product_name VARCHAR(255) NOT NULL,
    cost INT NOT NULL
);

CREATE TABLE IF NOT EXISTS Orders (
    id INT PRIMARY KEY AUTO_INCREMENT,
    customer_id INT NOT NULL,
    description VARCHAR(255) NOT NULL,
    order_date DATE NOT NULL,
    FOREIGN KEY (customer_id) REFERENCES Customers(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS Production (
    id INT PRIMARY KEY AUTO_INCREMENT,
    order_id INT NOT NULL,
    employee_id INT NOT NULL,
    description VARCHAR(255) NOT NULL,
    production_date DATE NOT NULL,
    FOREIGN KEY (order_id) REFERENCES Orders(id) ON DELETE CASCADE,
    FOREIGN KEY (employee_id) REFERENCES Employees(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS ProductionProducts (
    production_id INT,
    product_id INT,
    PRIMARY KEY (production_id, product_id),
    FOREIGN KEY (production_id) REFERENCES Production(id) ON DELETE CASCADE,
    FOREIGN KEY (product_id) REFERENCES Products(id) ON DELETE CASCADE
);

INSERT INTO Employees (full_name, phone, address, graduation_year) VALUES
('Иванов Иван Иванович', '+7-999-111-22-33', 'ул. Ленина, д. 1, кв. 5', 2010),
('Петров Пётр Сергеевич', '+7-999-222-33-44', 'ул. Пушкина, д. 10, кв. 15', 2015),
('Сидоров Сидор Николаевич', '+7-999-111-11-11', 'ул. Ленина, д. 25, кв. 4', 2015),
('Соколов Ирина Дмитриевна', '+7-999-222-22-22', 'ул. Мира, д. 12, кв. 7', 2013),
('Волков Сергей Александрович', '+7-999-333-33-33', 'ул. Пушкина, д. 8, кв. 15', 2016),
('Морозова Наталья Петровна', '+7-999-444-44-44', 'ул. Советская, д. 3, кв. 9', 2019),
('Новиков Алексей Иванович', '+7-999-555-55-55', 'ул. Лесная, д. 7, кв. 12', 2017),
('Фёдорова Ольга Сергеевна', '+7-999-666-66-66', 'ул. Речная, д. 5, кв. 6', 2014),
('Попов Дмитрий Владимирович', '+7-999-777-77-77', 'ул. Парковая, д. 10, кв. 3', 2021),
('Лебедева Анна Михайловна', '+7-999-888-88-88', 'ул. Садовая, д. 2, кв. 8', 2018),
('Громов Игорь Петрович', '+7-999-999-99-99', 'ул. Северная, д. 15, кв. 1', 2011),
('Крылова Татьяна Сергеевна', '+7-999-000-00-00', 'ул. Южная, д. 9, кв. 5', 2020);

INSERT INTO Qualifications (employee_id, position) VALUES
(1, 'Токарь'),
(2, 'Сварщик'),
(3, 'Инженер-технолог'),
(4, 'Бухгалтер'),
(5, 'Механик'),
(6, 'Кладовщик'),
(7, 'Оператор ЧПУ'),
(8, 'Контролёр ОТК'),
(9, 'Начальник цеха'),
(10, 'Инженер-конструктор'),
(11, 'Электрик'),
(12, 'Специалист по снабжению');

INSERT INTO Customers (full_name, address, regular_discount) VALUES
('ООО "СтройТехника"', 'г. Барнаул, ул. Промышленная, д. 8', 0.10),
('ИП Смирнов А.П.', 'г. Бийск, ул. Советская, д. 15', 0.00),
('ООО "АгроМаш"', 'г. Рубцовск, ул. Победы, д. 3', 0.10),
('ЗАО "Мебель-Сервис"', 'г. Новоалтайск, ул. Лесная, д. 5', 0.00),
('ООО "МеталлПро"', 'г. Барнаул, ул. Садовая, д. 10', 0.10),
('ООО "СибСтрой"', 'г. Камень-на-Оби, ул. Парковая, д. 5', 0.10),
('ИП Денисов В.П.', 'г. Заринск, ул. Центральная, д. 12', 0.00),
('ООО "ТехноСервис"', 'г. Алейск, ул. Молодёжная, д. 7', 0.10),
('ООО "ПромСнаб"', 'г. Барнаул, ул. Новая, д. 1', 0.00),
('ООО "УралТрейд"', 'г. Славгород, ул. Весенняя, д. 8', 0.10);

INSERT INTO Products (product_name, cost) VALUES
('Металлическая балка', 5000),
('Труба профильная', 3500),
('Лист стальной', 2500),
('Уголок стальной', 1800),
('Арматура', 4000),
('Сварная конструкция', 12000),
('Кронштейн', 1500),
('Металлический каркас', 18000),
('Ограждение', 7000),
('Ворота металлические', 25000),
('Дверь металлическая', 15000),
('Решётка', 3000),
('Заборная секция', 4500),
('Токарная деталь', 2200),
('Фрезерная деталь', 2800);

INSERT INTO Orders (customer_id, description, order_date) VALUES
(1, 'Металлоконструкции для склада', '2026-01-15'),
(2, 'Изготовление ограждения', '2026-01-20'),
(3, 'Партия арматуры', '2026-01-25');

INSERT INTO Production (order_id, employee_id, description, production_date) VALUES
(1, 1, 'Изготовление балок - первая партия', '2026-02-10'),
(1, 2, 'Сварка каркаса', '2026-02-15'),
(2, 2, 'Раскрой металла для ограждения', '2026-02-01'),
(2, 1, 'Контрольный осмотр', '2026-02-05'),
(3, 3, 'Проверка арматуры', '2026-02-20'),
(3, 1, 'Изготовление арматуры', '2026-02-25'),
(1, 3, 'Финишная обработка', '2026-03-01'),
(2, 3, 'Покраска ограждения', '2026-03-05'),
(3, 2, 'Сварка арматуры', '2026-03-10'),
(1, 1, 'Финальная сборка', '2026-03-15');

INSERT INTO ProductionProducts (production_id, product_id) VALUES
(1, 1),
(1, 3),
(2, 6),
(2, 2),
(3, 5),
(3, 1),
(4, 3),
(4, 1),
(5, 8),
(5, 5),
(6, 1),
(7, 7),
(8, 2),
(9, 6),
(10, 10);
