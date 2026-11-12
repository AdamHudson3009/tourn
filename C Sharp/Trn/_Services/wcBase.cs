using Dapper;
using Trn.Data;
using Microsoft.Extensions.Configuration;


namespace Trn.Services;

public static class WCBase
{
    //private static readonly DbConnectionFactory _db = new DbConnectionFactory();
    private static readonly DbConnectionFactory _db = new DbConnectionFactory(
        new ConfigurationBuilder()
            .SetBasePath(AppContext.BaseDirectory)
            .AddJsonFile("appsettings.json", optional: false)
            .Build());

    public static int GetWCBase(int trnId) {
        using var conn = _db.Create();
        var row = conn.Query(
            @"SELECT description AS description
            FROM tbrk
            WHERE tourn_id = @TrnId AND dv_cn_other = 'wc' AND rnk = 2",
            new { TrnId = trnId }).ToList();
        var description = row.FirstOrDefault()?.description;
        return int.Parse(description ?? "2"); //2 wc as default
    }
}