PRAGMA foreign_keys = ON;

BEGIN TRANSACTION;
DROP TABLE IF EXISTS user;
CREATE TABLE user (userid INTEGER PRIMARY KEY NOT NULL, username TEXT, password TEXT, role INTEGER);

DROP TABLE IF EXISTS emppto;
CREATE TABLE emppto (empptoid INTEGER PRIMARY KEY NOT NULL, empid INTEGER, ptoid INTEGER, amt REAL, FOREIGN KEY(empid) REFERENCES emp, FOREIGN KEY(ptoid) REFERENCES pto);

DROP TABLE IF EXISTS pto;
CREATE TABLE pto (ptoid INTEGER PRIMARY KEY NOT NULL, ptoname TEXT);

DROP TABLE IF EXISTS emp;
CREATE TABLE emp (empid INTEGER PRIMARY KEY NOT NULL, userid INTEGER, empno TEXT, firstname TEXT, lastname TEXT);

INSERT INTO user (username, password, role) VALUES ('robdelacruz', 'password', 1);
INSERT INTO user (username, password, role) VALUES ('lky', 'password', 0);

INSERT INTO emp (userid, empno, firstname, lastname) VALUES (1, "101", "Rob", "de la Cruz");
INSERT INTO emp (userid, empno, firstname, lastname) VALUES (2, "102", "Kuan Yew", "Lee");

INSERT INTO pto (ptoname) VALUES ('Sick Leave');
INSERT INTO pto (ptoname) VALUES ('Vacation Leave');
INSERT INTO emppto (empid, ptoid, amt) VALUES (1, 1, 15);
INSERT INTO emppto (empid, ptoid, amt) VALUES (1, 2, 20);

END TRANSACTION;

SELECT empno, firstname, lastname
FROM emp


