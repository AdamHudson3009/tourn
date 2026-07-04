-- MySQL dump 10.13  Distrib 8.0.41, for Win64 (x86_64)
--
-- Host: localhost    Database: league
-- ------------------------------------------------------
-- Server version	8.0.41

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!50503 SET NAMES utf8 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `divi`
--

DROP TABLE IF EXISTS `divi`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `divi` (
  `id` int NOT NULL AUTO_INCREMENT,
  `conf_id` int DEFAULT NULL,
  `description` varchar(100) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=63 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `divi`
--

LOCK TABLES `divi` WRITE;
/*!40000 ALTER TABLE `divi` DISABLE KEYS */;
INSERT INTO `divi` VALUES (3,14,'southern'),(6,15,'north'),(7,15,'west'),(9,14,'western'),(11,18,'red'),(12,18,'blue'),(16,19,'green'),(17,19,'blue'),(18,19,'red'),(21,18,'green'),(28,36,'Chinese'),(29,39,'blue'),(30,39,'green'),(31,39,'red'),(32,40,'blue'),(33,40,'green'),(34,40,'red'),(41,42,'East'),(42,42,'South'),(43,42,'North'),(44,42,'West'),(45,43,'East'),(46,43,'North'),(47,43,'South'),(48,43,'West'),(49,44,'blue'),(50,44,'green'),(51,44,'red'),(52,45,'blue'),(53,45,'green'),(54,45,'red'),(55,46,'waterhouse'),(56,46,'genaro'),(57,47,'waterhouse'),(58,47,'genaro'),(59,48,'1 1-5'),(60,48,'2 6-10'),(61,48,'3 11-15'),(62,48,'4 16-20');
/*!40000 ALTER TABLE `divi` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2026-07-04 16:50:45
