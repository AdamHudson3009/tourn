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
-- Table structure for table `grammar`
--

DROP TABLE IF EXISTS `grammar`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `grammar` (
  `id` int NOT NULL AUTO_INCREMENT,
  `tourn_id` int DEFAULT NULL,
  `letters` varchar(5) DEFAULT NULL,
  `dv1_id` int DEFAULT NULL,
  `dv2_id` int DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=140 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `grammar`
--

LOCK TABLES `grammar` WRITE;
/*!40000 ALTER TABLE `grammar` DISABLE KEYS */;
INSERT INTO `grammar` VALUES (29,32,'gbr',17,18),(30,32,'bbr',12,11),(31,32,'grg',18,16),(32,32,'brg',11,21),(33,32,'ibgb',12,17),(34,32,'ibgg',21,16),(35,32,'gb',17,17),(36,32,'gg',16,16),(37,32,'gbg',17,16),(38,32,'bb',12,12),(39,32,'bg',21,21),(40,32,'bbg',12,21),(41,32,'ibgr',11,18),(42,32,'gr',18,18),(43,32,'br',11,11),(60,31,'bb',29,29),(61,31,'brg',31,30),(62,31,'bg',30,30),(63,31,'bbg',29,30),(64,31,'br',31,31),(65,31,'bbr',29,31),(66,31,'gb',32,32),(67,31,'ibgb',29,32),(68,31,'grg',34,33),(69,31,'gg',33,33),(70,31,'gbg',32,33),(71,31,'ibgg',30,33),(72,31,'gr',34,34),(73,31,'gbr',32,34),(74,31,'ibgr',31,34),(78,37,'y',40,40),(79,37,'w',37,37),(80,38,'NE-NE',45,45),(81,38,'AW-AW',44,44),(82,38,'AE-AS',41,42),(83,38,'NW-NS',48,47),(84,38,'NS-AS',47,42),(85,38,'NS-NS',47,47),(86,38,'AN-AE',43,41),(87,38,'AN-AN',43,43),(88,38,'AW-AE',44,41),(89,38,'AS-AW',42,44),(90,38,'NW-NW',48,48),(91,38,'NN-NN',46,46),(92,38,'AS-NW',42,48),(93,38,'NE-NN',45,46),(94,38,'NW-AN',48,43),(95,38,'AE-AE',41,41),(96,38,'AS-AN',42,43),(97,38,'NW-AS',48,42),(98,38,'AW-AS',44,42),(99,38,'NS-NW',47,48),(100,38,'NE-AW',45,44),(101,38,'NS-NN',47,46),(102,38,'AE-NS',41,47),(103,38,'AW-NE',44,45),(104,38,'AS-AS',42,42),(105,38,'NW-NE',48,45),(106,38,'AN-NN',43,46),(107,38,'NN-AN',46,43),(108,38,'NE-NS',45,47),(109,38,'NS-AE',47,41),(110,38,'NN-AW',46,44),(111,38,'AN-AW',43,44),(112,38,'NN-NE',46,45),(113,38,'NE-AE',45,41),(114,38,'AE-AN',41,43),(115,38,'NN-NW',46,48),(116,39,'bb',49,49),(117,39,'brg',51,50),(118,39,'bg',50,50),(119,39,'bbg',49,50),(120,39,'bbr',49,51),(121,39,'br',51,51),(122,39,'ibgb',49,52),(123,39,'gb',52,52),(124,39,'grg',54,53),(125,39,'ibgg',50,53),(126,39,'gg',53,53),(127,39,'gbg',52,53),(128,39,'gbr',52,54),(129,39,'ibgr',51,54),(130,39,'gr',54,54),(131,37,'bgbw',56,55),(132,37,'gggw',58,57),(133,37,'bggw',56,57),(135,37,'bggg',56,58),(136,37,'bwgg',55,58),(137,37,'bwgw',55,57),(138,37,'gwgw',57,57),(139,37,'gggg',58,58);
/*!40000 ALTER TABLE `grammar` ENABLE KEYS */;
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
