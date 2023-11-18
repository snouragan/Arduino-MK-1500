---
runme:
  id: 01HFJ06TK5X5CXJ0XN8DVQ72JW
  version: v2.0
---

mysql> CREATE TABLE SolSensor (
    ->     id INT AUTO_INCREMENT PRIMARY KEY,
    ->     Temp1 FLOAT,
    ->     Temp2 FLOAT,
    ->     Temp3 FLOAT,
    ->     Temp4 FLOAT,
    ->     Temp5 FLOAT,
    ->     Temp6 FLOAT
    -> );
Query OK, 0 rows affected (0.05 sec)

mysql> CREATE TABLE EnvSensor (
    ->     id INT AUTO_INCREMENT PRIMARY KEY,
    ->     EnvTemp FLOAT,
    ->     EnvHumid FLOAT,
    ->     EnvLux FLOAT
    -> );
Query OK, 0 rows affected (0.02 sec)

mysql> CREATE TABLE ShtSensor (
    ->     id INT AUTO_INCREMENT PRIMARY KEY,
    ->     ShtHumid FLOAT,
    ->     ShtTemp FLOAT
    -> );
Query OK, 0 rows affected (0.03 sec)

mysql> CREATE TABLE DataModel (
    ->     id INT AUTO_INCREMENT PRIMARY KEY,
    ->     Time DATETIME,
    ->     SolSensor_id INT,
    ->     EnvSensor_id INT,
    ->     ShtSensor_id INT,
    ->     FOREIGN KEY (SolSensor_id) REFERENCES SolSensor(id),
    ->     FOREIGN KEY (EnvSensor_id) REFERENCES EnvSensor(id),
    ->     FOREIGN KEY (ShtSensor_id) REFERENCES ShtSensor(id)
    -> );
Query OK, 0 rows affected (0.05 sec)

mysql>