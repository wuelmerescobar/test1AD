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
('System Design','Martin Fowler','1010101010','Engineering'),
('Refactoring','Martin Fowler','1111111112','Programming'),
('The Pragmatic Programmer','Andrew Hunt','1212121212','Programming'),
('Designing Data-Intensive Applications','Martin Kleppmann','1313131313','Technology'),
('Computer Networks','Andrew Tanenbaum','1414141414','Networking'),
('Introduction to Algorithms','Thomas Cormen','1515151515','Computer Science'),
('Effective Java','Joshua Bloch','1616161616','Programming'),
('Web Application Security','Andrew Hoffman','1717171717','Security'),
('Operating System Concepts','Abraham Silberschatz','1818181818','Computer Science'),
('Building Microservices','Sam Newman','1919191919','Architecture'),
('Domain-Driven Design','Eric Evans','2020202020','Software Design');


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
(10,1,'available'),
(11,2,'available'),
(12,3,'available'),
(13,1,'available'),
(14,2,'available'),
(15,3,'available'),
(16,1,'available'),
(17,2,'available'),
(18,3,'available'),
(19,1,'available'),
(20,2,'available');


UPDATE book_copies
SET status = 'borrowed'
WHERE id IN (2, 5, 8);


INSERT INTO loans (member_id, book_copy_id, borrowed_at, due_at, returned_at, status) VALUES
(1,2,CURRENT_TIMESTAMP - INTERVAL '12 days',CURRENT_TIMESTAMP + INTERVAL '2 days',NULL,'borrowed'),
(3,5,CURRENT_TIMESTAMP - INTERVAL '21 days',CURRENT_TIMESTAMP - INTERVAL '7 days',NULL,'overdue'),
(6,8,CURRENT_TIMESTAMP - INTERVAL '18 days',CURRENT_TIMESTAMP - INTERVAL '4 days',CURRENT_TIMESTAMP - INTERVAL '1 day','returned');


INSERT INTO fines (loan_id, member_id, amount, reason, paid) VALUES
(2,3,17.50,'Book is overdue by 7 days',FALSE),
(3,6,10.00,'Book was returned 3 days late',TRUE);
