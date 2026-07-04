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
-- Table structure for table `tbrk`
--

DROP TABLE IF EXISTS `tbrk`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `tbrk` (
  `id` int NOT NULL AUTO_INCREMENT,
  `tourn_id` int DEFAULT NULL,
  `league_id` int DEFAULT NULL,
  `dv_cn_other` varchar(20) DEFAULT NULL,
  `description` varchar(20) DEFAULT NULL,
  `rnk` int DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=164 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `tbrk`
--

LOCK TABLES `tbrk` WRITE;
/*!40000 ALTER TABLE `tbrk` DISABLE KEYS */;
INSERT INTO `tbrk` VALUES (13,32,NULL,'dv','arch',2),(14,32,NULL,'dv','hhx',3),(15,32,NULL,'dv','cm',4),(16,32,NULL,'dv','dv',5),(17,32,NULL,'dv','sht-pt',6),(18,32,NULL,'dv','pts-pt',7),(19,32,NULL,'dv','strk',8),(20,32,NULL,'dv','clsgme',9),(21,32,NULL,'cn','hhxa',2),(22,32,NULL,'cn','cm4',3),(23,32,NULL,'cn','sht',4),(24,32,NULL,'cn','pts',5),(25,32,NULL,'cn','strk',6),(26,32,NULL,'cn','clsgme',7),(27,32,NULL,'hdr-spr-ln','T1',1),(28,32,NULL,'hdr-spr-ln','T2',2),(29,32,NULL,'hdr-spr-ln','W',3),(30,32,NULL,'hdr-spr-ln','Pt',4),(31,32,NULL,'hdr-spr-ln','Pf',5),(32,32,NULL,'hdr-spr-ln','Pa',6),(33,32,NULL,'hdr-spr-ln','Rec',7),(34,32,NULL,'wc','cn',1),(35,32,NULL,'wc','4',2),(36,32,NULL,'cn','o',1),(37,32,NULL,'dv','o',1),(63,31,NULL,'dv','arch',2),(64,31,NULL,'dv','hhx',3),(65,31,NULL,'dv','cm',4),(66,31,NULL,'dv','dv',5),(67,31,NULL,'dv','sht-pt',6),(68,31,NULL,'dv','pts-pt',7),(69,31,NULL,'dv','strk',8),(70,31,NULL,'dv','clsgme',9),(71,31,NULL,'cn','hhxa',2),(72,31,NULL,'cn','cm4',3),(73,31,NULL,'cn','sht',4),(74,31,NULL,'cn','pts',5),(75,31,NULL,'cn','strk',6),(76,31,NULL,'cn','clsgme',7),(77,31,NULL,'hdr-spr-ln','T1',1),(78,31,NULL,'hdr-spr-ln','T2',2),(79,31,NULL,'hdr-spr-ln','W',3),(80,31,NULL,'hdr-spr-ln','Pt',4),(81,31,NULL,'hdr-spr-ln','Pf',5),(82,31,NULL,'hdr-spr-ln','Pa',6),(83,31,NULL,'hdr-spr-ln','Rec',7),(84,31,NULL,'wc','cn',1),(85,31,NULL,'wc','4',2),(86,31,NULL,'cn','o',1),(87,31,NULL,'dv','o',1),(88,32,NULL,'wc format','rhombus',1),(89,32,NULL,'playoff format','external',1),(90,31,NULL,'wc format','rhombus',1),(91,31,NULL,'playoff format','external',1),(92,37,NULL,'dv','hhx',2),(93,37,NULL,'dv','sht',3),(94,37,NULL,'dv','pts',5),(95,37,NULL,'dv','hh-sht',4),(96,37,NULL,'dv','hh-pts',6),(97,37,NULL,'dv','hh',7),(98,37,NULL,'cn','o',1),(99,37,NULL,'dv','o',1),(100,37,NULL,'hdr-spr-ln','T1',1),(101,37,NULL,'hdr-spr-ln','T2',2),(102,37,NULL,'hdr-spr-ln','W',3),(103,37,NULL,'hdr-spr-ln','Pt',4),(104,37,NULL,'hdr-spr-ln','Pf',5),(105,37,NULL,'hdr-spr-ln','Pa',6),(106,37,NULL,'hdr-spr-ln','Rec',7),(107,38,NULL,'dv','hh',2),(108,38,NULL,'dv','dv',3),(109,38,NULL,'dv','cm',4),(110,38,NULL,'dv','cn',5),(111,38,NULL,'dv','sov',6),(112,38,NULL,'dv','sos',7),(113,38,NULL,'dv','rnk-cn',8),(114,38,NULL,'dv','rnk',9),(115,38,NULL,'dv','pts-cn',10),(116,38,NULL,'dv','pts',11),(118,38,NULL,'cn','hhb',2),(119,38,NULL,'cn','cm4',4),(120,38,NULL,'cn','sov',5),(121,38,NULL,'cn','sos',6),(122,38,NULL,'cn','cn',3),(123,38,NULL,'cn','rnk-cn',7),(124,38,NULL,'cn','rnk',8),(125,38,NULL,'cn','pts-cn',9),(126,38,NULL,'cn','pts',10),(127,38,NULL,'dv','o',1),(128,38,NULL,'cn','o',1),(129,38,NULL,'wc','cn',1),(130,38,NULL,'hdr-spr-ln','f',1),(131,38,NULL,'wc','3',2),(132,38,NULL,'hdr-spr-ln','T1',2),(133,38,NULL,'hdr-spr-ln','T2',3),(134,38,NULL,'hdr-spr-ln','W',4),(135,38,NULL,'hdr-spr-ln','Pt',5),(136,38,NULL,'hdr-spr-ln','Pf',6),(137,38,NULL,'hdr-spr-ln','Pa',7),(138,38,NULL,'hdr-spr-ln','Rec',8),(139,39,NULL,'cn','o',1),(140,39,NULL,'dv','o',1),(141,39,NULL,'cn','hhxa',2),(142,39,NULL,'cn','cm4',3),(143,39,NULL,'cn','sht',4),(144,39,NULL,'cn','pts',5),(145,39,NULL,'cn','strk',6),(146,39,NULL,'cn','clsgme',7),(147,39,NULL,'dv','arch',2),(148,39,NULL,'dv','hhx',3),(149,39,NULL,'dv','cm',4),(150,39,NULL,'dv','dv',5),(151,39,NULL,'dv','sht-pt',6),(152,39,NULL,'dv','pts-pt',7),(153,39,NULL,'dv','strk',8),(154,39,NULL,'dv','clsgme',9),(155,39,NULL,'hdr-spr-ln','T1',1),(156,39,NULL,'wc','cn',1),(157,39,NULL,'wc','4',2),(158,39,NULL,'hdr-spr-ln','T2',2),(159,39,NULL,'hdr-spr-ln','W',3),(160,39,NULL,'hdr-spr-ln','Pt',4),(161,39,NULL,'hdr-spr-ln','Pf',5),(162,39,NULL,'hdr-spr-ln','Pa',6),(163,39,NULL,'hdr-spr-ln','Rec',7);
/*!40000 ALTER TABLE `tbrk` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2026-07-04 16:50:47
