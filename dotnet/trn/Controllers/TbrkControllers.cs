using System;
using Microsoft.AspNetCore.Mvc;
using MyRestApi.Models;
using MySql.Data.MySqlClient;
using System.Collections.Generic;

namespace MyRestApi.Controllers
{
    [ApiController]
    [Route("tbrk/[controller]")]
    public class TbrkController : ControllerBase
    {
        private readonly string connectionString = "Server=localhost;Database=league;User Id=root;Password=MySQL=1;";

        [HttpPost]
        [Route("/tbrk/GetCopyFrom")]
        public ActionResult<List<TournTbrk>> GetCopyFrom([FromBody] TournTbrk input)
        {
            List<TournTbrk> results = new List<TournTbrk>();

            try {
                var connection = new MySqlConnection(connectionString);
                connection.Open();

                string sql = @"SELECT distinct tr.id id, tr.description tourn 
                            from tbrk tb inner join tourn tr on tr.id = tb.tourn_id
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
                return StatusCode(500, "Data error 1: " + ex.Message);
            }
        }   
   
        [HttpPost]
        [Route("/tbrk/CopyFrom")]
        public IActionResult CopyFrom([FromBody] TournToFrom input)
        {
            Console.WriteLine("ABCCBA");
            try {
                var connection = new MySqlConnection(connectionString);
                connection.Open();

                string sql = @"SELECT 1 from tbrk where tourn_id = @tournId";
                
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

                List<string> dvCnOtherList = new List<string>();
                List<string> descList = new List<string>();
                List<int> rnkList = new List<int>();

                sql = @"SELECT dv_cn_other, description, rnk from tbrk where tourn_id = @tournId";

                command = new MySqlCommand(sql, connection);
                command.Parameters.AddWithValue("@tournId", input.TournIdFrom);
                reader = command.ExecuteReader();
                    
                while (reader.Read()) {
                    dvCnOtherList.Add(reader["dv_cn_other"].ToString());
                    descList.Add(reader["description"].ToString());
                    rnkList.Add(Convert.ToInt32(reader["rnk"]));
                }
                reader.Close();
                
                Console.WriteLine("Here");
                int rows = 0;
                for (int indx=0; indx<dvCnOtherList.Count; indx++) {
                    sql = @"insert into tbrk (tourn_id, dv_cn_other, description, rnk)
                        values (@tournId,@dvCnOther,@description,@rnk)";
                    command = new MySqlCommand(sql, connection);
                    command.Parameters.AddWithValue("@tournId", input.TournIdTo); 
                    command.Parameters.AddWithValue("@dvCnOther", dvCnOtherList[indx]); 
                    command.Parameters.AddWithValue("@description", descList[indx]); 
                    command.Parameters.AddWithValue("@rnk",rnkList[indx]);
                    rows += command.ExecuteNonQuery();
                }
                return Ok(new { inserted = rows });
            }
            catch (Exception ex) {
                return StatusCode(500, "Data error 2: " + ex.Message);
            }
        }   
    }

}