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
-- Table structure for table `plyrtm`
--

DROP TABLE IF EXISTS `plyrtm`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `plyrtm` (
  `id` int NOT NULL AUTO_INCREMENT,
  `divi_id` int NOT NULL,
  `description` varchar(100) NOT NULL,
  `description2` varchar(100) DEFAULT NULL,
  `sds` int NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=263 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `plyrtm`
--

LOCK TABLES `plyrtm` WRITE;
/*!40000 ALTER TABLE `plyrtm` DISABLE KEYS */;
INSERT INTO `plyrtm` VALUES (14,7,'Arizona','Cardinal',1),(40,18,'elise',NULL,4),(48,18,'joana',NULL,1),(49,18,'amber',NULL,6),(50,18,'tara',NULL,3),(51,18,'mary',NULL,5),(52,18,'laiken',NULL,2),(53,17,'heatherf',NULL,1),(54,17,'deidra',NULL,2),(55,17,'victoria',NULL,3),(56,17,'lucy',NULL,4),(57,17,'kenzie',NULL,5),(59,16,'danielle',NULL,1),(60,16,'dawn',NULL,2),(61,16,'emma',NULL,3),(62,16,'kashani',NULL,4),(63,16,'zizi',NULL,6),(64,16,'siena',NULL,5),(65,11,'dean',NULL,1),(66,11,'chris',NULL,3),(67,11,'scott',NULL,2),(68,11,'trent',NULL,5),(69,11,'anson',NULL,4),(70,11,'kyle',NULL,6),(71,12,'gary',NULL,2),(73,12,'bob',NULL,3),(74,12,'jackson',NULL,4),(75,21,'kenny',NULL,5),(76,21,'cliff',NULL,6),(77,21,'kaleb',NULL,4),(78,21,'cecil',NULL,3),(79,21,'tyler',NULL,1),(81,21,'chase',NULL,2),(83,12,'jonathan',NULL,6),(84,12,'jacob',NULL,5),(85,12,'damian',NULL,1),(87,28,'li zhe',NULL,1),(88,17,'sara',NULL,6),(89,29,'licorice',NULL,1),(90,29,'august',NULL,2),(91,29,'conrad',NULL,3),(92,29,'james',NULL,4),(93,29,'able',NULL,5),(99,31,'greg',NULL,3),(100,31,'jack',NULL,4),(101,31,'carry',NULL,5),(102,31,'clarence',NULL,6),(103,30,'juan',NULL,2),(104,30,'ash',NULL,1),(105,32,'elana',NULL,1),(106,32,'jenna',NULL,2),(107,32,'nora',NULL,5),(108,32,'twes',NULL,3),(110,34,'cila',NULL,6),(111,30,'kit',NULL,3),(112,29,'johnb',NULL,6),(113,32,'heatherp',NULL,6),(114,32,'mariana',NULL,4),(115,33,'cali',NULL,1),(116,30,'charlie',NULL,6),(117,30,'jaylo',NULL,4),(119,31,'al',NULL,1),(120,31,'jeremy',NULL,2),(121,30,'crayton',NULL,5),(122,33,'carla',NULL,3),(123,33,'bridgette',NULL,2),(124,34,'jc',NULL,1),(125,34,'le',NULL,2),(126,33,'nancy',NULL,4),(127,33,'jen',NULL,6),(128,34,'ashley',NULL,3),(129,34,'katy',NULL,4),(130,33,'megan',NULL,5),(131,34,'lexi',NULL,5),(152,41,'Buffalo Bills',NULL,1),(153,41,'Miami Dolphins',NULL,2),(154,41,'New England Patriots',NULL,4),(155,41,'New York Jets',NULL,3),(156,42,'Houston Texans',NULL,1),(157,42,'Indianapolis Colts',NULL,2),(158,42,'Jacksonville Jaguars',NULL,3),(159,42,'Tennessee Titans',NULL,4),(160,43,'Baltimore Ravens',NULL,1),(161,43,'Cincinnati Bengals',NULL,3),(162,43,'Cleveland Browns',NULL,4),(163,43,'Pittsburgh Steelers',NULL,2),(164,44,'Denver Broncos',NULL,3),(165,44,'Kansas City Chiefs',NULL,1),(166,44,'Los Angeles Chargers',NULL,2),(167,44,'Las Vegas Raiders',NULL,4),(168,45,'Philadelphia Eagles',NULL,1),(169,45,'Washington Commanders',NULL,2),(170,45,'Dallas Cowboys',NULL,3),(171,45,'New York Giants',NULL,4),(172,46,'Chicago Bears',NULL,4),(173,46,'Detroit Lions',NULL,1),(174,46,'Green Bay Packers',NULL,3),(175,46,'Minnesota Vikings',NULL,2),(176,47,'Atlanta Falcons',NULL,2),(177,47,'Carolina Panthers',NULL,3),(178,47,'New Orleans Saints',NULL,4),(179,47,'Tampa Bay Buccaneers',NULL,1),(180,48,'Arizona Cardinals',NULL,3),(181,48,'Los Angeles Rams',NULL,1),(182,48,'San Francisco 49ers',NULL,4),(183,48,'Seattle Seahawks',NULL,2),(184,49,'damian',NULL,1),(185,49,'gary',NULL,2),(186,49,'bob',NULL,3),(187,49,'jackson',NULL,4),(188,49,'jacob',NULL,5),(189,49,'jonathan',NULL,6),(190,50,'tyler',NULL,1),(191,50,'chase',NULL,2),(192,50,'cecil',NULL,3),(193,50,'kaleb',NULL,4),(194,50,'kenny',NULL,5),(195,50,'cliff',NULL,6),(196,51,'dean',NULL,1),(197,51,'chris',NULL,3),(198,51,'anson',NULL,5),(199,51,'scott',NULL,4),(200,51,'trent',NULL,2),(201,51,'kev',NULL,6),(202,52,'heatherf',NULL,1),(203,52,'deidra',NULL,2),(204,52,'victoria',NULL,3),(205,52,'sara',NULL,5),(206,52,'marina',NULL,4),(207,52,'kenzie',NULL,6),(208,53,'danielle',NULL,1),(209,53,'dawn',NULL,2),(210,53,'emma',NULL,3),(211,53,'kashani',NULL,4),(212,53,'siena',NULL,5),(213,53,'zizi',NULL,6),(214,54,'joana',NULL,1),(215,54,'laiken',NULL,2),(216,54,'tara',NULL,3),(217,54,'elise',NULL,4),(218,54,'lucy',NULL,5),(219,54,'amber',NULL,6),(220,55,'cliff',NULL,1),(221,55,'kenny',NULL,2),(222,55,'dean',NULL,3),(223,55,'gary',NULL,4),(224,55,'damian',NULL,5),(225,57,'deidra',NULL,1),(226,57,'joana',NULL,2),(227,57,'heatherf',NULL,3),(228,57,'danielle',NULL,4),(229,57,'elise',NULL,5),(230,56,'carry',NULL,1),(231,56,'ash',NULL,2),(232,56,'licorice',NULL,3),(233,56,'juan',NULL,4),(234,56,'august',NULL,5),(235,58,'elana',NULL,1),(236,58,'cali',NULL,2),(237,58,'jc',NULL,3),(238,58,'le',NULL,4),(239,58,'carla',NULL,5),(240,59,'carry',NULL,1),(241,59,'deidra',NULL,2),(242,59,'elana',NULL,3),(243,59,'joana',NULL,4),(244,59,'cliff',NULL,5),(248,60,'ash',NULL,1),(249,60,'cali',NULL,2),(250,60,'licorice',NULL,3),(251,60,'kenny',NULL,4),(252,60,'heatherf',NULL,5),(253,61,'jc',NULL,1),(254,61,'dean',NULL,2),(255,61,'gary',NULL,3),(256,61,'le',NULL,4),(257,61,'juan',NULL,5),(258,62,'danielle',NULL,1),(259,62,'carla',NULL,2),(260,62,'damian',NULL,3),(261,62,'august',NULL,4),(262,62,'elise',NULL,5);
/*!40000 ALTER TABLE `plyrtm` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2026-07-04 16:50:46
