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
('Кузнецов Артём Викторович', '+7-913-501-77-11', 'пр. Строителей, д. 14, кв. 8', 2008),
('Григорьева Юлия Олеговна', '+7-913-502-88-22', 'ул. Машиностроителей, д. 27, кв. 3', 2012),
('Макаров Денис Игоревич', '+7-913-503-99-33', 'ул. Заводская, д. 5, кв. 12', 2014),
('Тарасова Светлана Борисовна', '+7-913-504-10-44', 'ул. Индустриальная, д. 8, кв. 22', 2010),
('Орлов Роман Сергеевич', '+7-913-505-21-55', 'ул. Кузнечная, д. 19, кв. 7', 2016),
('Беляева Наталья Юрьевна', '+7-913-506-32-66', 'ул. Литейная, д. 3, кв. 15', 2018),
('Захаров Вячеслав Петрович', '+7-913-507-43-77', 'пр. Металлургов, д. 22, кв. 9', 2013),
('Панова Оксана Владимировна', '+7-913-508-54-88', 'ул. Рабочая, д. 11, кв. 4', 2019),
('Дмитриев Константин Андреевич', '+7-913-509-65-99', 'ул. Цеховая, д. 6, кв. 18', 2007),
('Ершова Марина Александровна', '+7-913-510-76-10', 'ул. Проектная, д. 30, кв. 11', 2015),
('Фомин Григорий Николаевич', '+7-913-511-87-21', 'ул. Энергетиков, д. 4, кв. 6', 2011),
('Родионова Ирина Витальевна', '+7-913-512-98-32', 'ул. Снабженческая, д. 17, кв. 14', 2017);

INSERT INTO Qualifications (employee_id, position) VALUES
(1, 'Фрезеровщик'),
(2, 'Шлифовщик'),
(3, 'Инженер-механик'),
(4, 'Экономист'),
(5, 'Слесарь-ремонтник'),
(6, 'Кладовщик'),
(7, 'Оператор ЧПУ'),
(8, 'Контролёр ОТК'),
(9, 'Начальник цеха'),
(10, 'Инженер-конструктор'),
(11, 'Электрик'),
(12, 'Специалист по снабжению');

INSERT INTO Customers (full_name, address, regular_discount) VALUES
('ООО "УралПром"', 'г. Челябинск, ул. Танкистов, д. 22', 0.10),
('АО "СибирьМеталл"', 'г. Новосибирск, ул. Кирова, д. 113', 0.00),
('ООО "КузбассСтрой"', 'г. Кемерово, пр. Октябрьский, д. 45', 0.10),
('ИП Волков И.С.', 'г. Томск, ул. Пролетарская, д. 8', 0.00),
('ООО "АлтайТрансСервис"', 'г. Барнаул, ул. Попова, д. 179', 0.10),
('ООО "СибПромРесурс"', 'г. Красноярск, ул. Мира, д. 102', 0.10),
('ЗАО "МостСтрой-Регион"', 'г. Иркутск, ул. Лермонтова, д. 78', 0.00),
('ООО "ГорноДобыча"', 'г. Норильск, ул. Ленинградская, д. 12', 0.10),
('ООО "ТрансЛогистика"', 'г. Омск, пр. Маркса, д. 51', 0.00),
('ООО "ВостокМеталлургия"', 'г. Хабаровск, ул. Муравьёва-Амурского, д. 15', 0.10);

INSERT INTO Products (product_name, cost) VALUES
('Стальной швеллер', 6200),
('Лист горячекатаный', 4100),
('Труба бесшовная', 5500),
('Круг стальной', 1900),
('Полоса стальная', 2700),
('Кованая заготовка', 14000),
('Зубчатое колесо', 3500),
('Опорная рама', 21000),
('Металлический короб', 8500),
('Промышленные ворота', 32000),
('Технологическая дверь', 18500),
('Защитная решётка', 4200),
('Секция ограждения', 5100),
('Вал токарный', 2600),
('Корпус редуктора', 9800);

INSERT INTO Orders (customer_id, description, order_date) VALUES
(1, 'Партия швеллеров для цеха', '2026-01-18'),
(2, 'Металлоконструкции для моста', '2026-01-22'),
(3, 'Комплект ограждений', '2026-01-28');

INSERT INTO Production (order_id, employee_id, description, production_date) VALUES
(1, 1, 'Резка швеллера - первая партия', '2026-02-12'),
(1, 2, 'Сварка опорной рамы', '2026-02-17'),
(2, 2, 'Раскрой листового металла', '2026-02-03'),
(2, 1, 'Контрольный обмер конструкций', '2026-02-07'),
(3, 3, 'Проверка секций ограждения', '2026-02-22'),
(3, 1, 'Изготовление ограждений', '2026-02-27'),
(1, 3, 'Финишная обработка швеллеров', '2026-03-03'),
(2, 3, 'Антикоррозийная обработка', '2026-03-08'),
(3, 2, 'Сварка секций', '2026-03-12'),
(1, 1, 'Финальная сборка партии', '2026-03-18');

INSERT INTO ProductionProducts (production_id, product_id) VALUES
(1, 1),
(1, 4),
(2, 8),
(2, 3),
(3, 13),
(3, 1),
(4, 4),
(4, 1),
(5, 15),
(5, 13),
(6, 1),
(7, 7),
(8, 3),
(9, 8),
(10, 10);
