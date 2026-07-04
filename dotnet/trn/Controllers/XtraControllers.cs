using System;
using Microsoft.AspNetCore.Mvc;
using MyRestApi.Models;
using MySql.Data.MySqlClient;
using System.Collections.Generic;

namespace MyRestApi.Controllers
{
    [ApiController]
    [Route("lllextra/[controller]")]
    public class XtraController : ControllerBase
    {
        private readonly string connectionString = "Server=localhost;Database=league;User Id=root;Password=MySQL=1;";

        [HttpPost]
        [Route("/lllextra/GetCopyFrom")]
        public ActionResult<List<TournTbrk>> GetCopyFrom([FromBody] TournTbrk input)
        {
            List<TournTbrk> results = new List<TournTbrk>();

            try {
                var connection = new MySqlConnection(connectionString);
                connection.Open();

                string sql = @"SELECT distinct tr.id id, tr.description tourn 
                            from lll_line l inner join tourn tr on tr.id = l.tourn_id
                            where tr.league_id in (select t1.league_id from tourn t1 where t1.id = @tournId)
                            order by tr.description ";
                
                MySqlCommand command = new MySqlCommand(sql, connection);
                command.Parameters.AddWithValue("@tournId", input.TournId);
                var reader = command.ExecuteReader();
                    
                while (reader.Read()) {
                    int tempId = Convert.ToInt32(reader["id"]);
                    if (input.TournId == tempId) {
                        results.Clear();
                        results.Add(new TournTbrk {TournId = 0, Tourn = ""});
                        break;
                    }
                    results.Add(new TournTbrk {
                        TournId = Convert.ToInt32(reader["id"]), 
                        Tourn = reader["tourn"].ToString()
                    });
                }
                reader.Close();
                return Ok(results);
            }
            catch (Exception ex) {
                return StatusCode(500, "Data error 3: " + ex.Message);
            }
        }   
   
        [HttpPost]
        [Route("/lllextra/CopyFrom")]
        public IActionResult CopyFrom([FromBody] TournToFrom input)
        {
            try {
                var connection = new MySqlConnection(connectionString);
                connection.Open();

                string sql = @"SELECT 1 from lll_line where tourn_id = @tournId";
                
                MySqlCommand command = new MySqlCommand(sql, connection);
                command.Parameters.AddWithValue("@tournId", input.TournIdTo);
                var reader = command.ExecuteReader();
                    
                bool found = false;
                while (reader.Read()) {
                    found = true;
                }
                reader.Close();
                if (found) {
                    return Ok("Already Exists");
                }

                List<string> descList = new List<string>();
                List<int> ordrList = new List<int>();

                sql = @"SELECT description, ordr from lll_line where tourn_id = @tournId";

                command = new MySqlCommand(sql, connection);
                command.Parameters.AddWithValue("@tournId", input.TournIdFrom);
                reader = command.ExecuteReader();
                    
                while (reader.Read()) {
                    descList.Add(reader["description"].ToString());
                    ordrList.Add(Convert.ToInt32(reader["ordr"]));
                }
                reader.Close();
                
                int rows = 0;
                for (int indx=0; indx<descList.Count; indx++) {
                    sql = @"insert into lll_line (tourn_id, description, ordr)
                        values (@tournId,@description,@ordr)";
                    command = new MySqlCommand(sql, connection);
                    command.Parameters.AddWithValue("@tournId", input.TournIdTo); 
                    command.Parameters.AddWithValue("@description", descList[indx]); 
                    command.Parameters.AddWithValue("@ordr",ordrList[indx]);
                    rows += command.ExecuteNonQuery();
                }
                return Ok(new { inserted = rows });
            }
            catch (Exception ex) {
                return StatusCode(500, "Data error 4: " + ex.Message);
            }
        }   
    }

}