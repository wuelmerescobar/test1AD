INSERT INTO branches (name, code, address) VALUES
('Central Library','CL01','Belize City'),
('West Branch','WB01','San Ignacio'),
('University Branch','UB01','Belmopan');


INSERT INTO books (title, author, isbn, genre) VALUES
('Cybersecurity Basics','John Smith','1111111111','Technology'),
('Learning Go','Alan Donovan','2222222222','Programming'),
('Database Design','Chris Date','3333333333','Technology'),
('Network Security','Bruce Schneier','4444444444','Security'),
('Linux Fundamentals','Brian Ward','5555555555','Technology'),
('API Development','Sam Newman','6666666666','Programming'),
('Cloud Computing','Thomas Erl','7777777777','Technology'),
('Software Engineering','Ian Sommerville','8888888888','Engineering'),
('Clean Code','Robert Martin','9999999999','Programming'),
('System Design','Martin Fowler','1010101010','Engineering');


INSERT INTO members (first_name,last_name,email,phone,branch_id) VALUES
('Juan','Perez','juan@test.com','6000001',1),
('Maria','Gomez','maria@test.com','6000002',2),
('Carlos','Lopez','carlos@test.com','6000003',3),
('Ana','Martinez','ana@test.com','6000004',1),
('Luis','Hernandez','luis@test.com','6000005',2),
('Elena','Ramirez','elena@test.com','6000006',3),
('Pedro','Torres','pedro@test.com','6000007',1),
('Lucia','Flores','lucia@test.com','6000008',2),
('Jorge','Castillo','jorge@test.com','6000009',3),
('Sofia','Mendoza','sofia@test.com','6000010',1);


INSERT INTO book_copies (book_id,branch_id,status) VALUES
(1,1,'available'),
(2,1,'available'),
(3,2,'available'),
(4,2,'available'),
(5,3,'available'),
(6,3,'available'),
(7,1,'available'),
(8,2,'available'),
(9,3,'available'),
(10,1,'available');