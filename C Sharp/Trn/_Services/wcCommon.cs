using Dapper;
using Trn.Data;
using Trn.Dtos;
using Microsoft.Extensions.Configuration;

namespace Trn.Services;

public static class WCCommon
{
    private static readonly DbConnectionFactory _db = new DbConnectionFactory(
        new ConfigurationBuilder()
            .SetBasePath(AppContext.BaseDirectory)
            .AddJsonFile("appsettings.json", optional: false)
            .Build());

    public static MatchResponse wc_round_setup(List<int> plyrsIds) {
        var plyr1 = new List<int>();
        var plyr2 = new List<int>();

        for (int indx = 0; indx < plyrsIds.Count/2; indx++) {
            plyr1.Add(plyrsIds[indx]);
            plyr2.Add(plyrsIds[plyrsIds.Count-1-indx]);
        }
        return new MatchResponse(plyr1, plyr2);
    }

    public static bool wc_insert(int trnId, int currRnd, List<int> plyr1, List<int> plyr2, int wnIncr = -1) {
        using var conn = _db.Create();
        bool inserted = true;
        int rec = 0;

        for (int indx = 0; indx < plyr1.Count; indx++) {
            int gmsId = conn.ExecuteScalar<int>(
            @"INSERT INTO gms (tourn_id, rnd, has_stnd) VALUES (@TrnId, @CurrRnd, 'N');
            SELECT LAST_INSERT_ID();",
            new { TrnId = trnId, CurrRnd = currRnd });

            if (wnIncr == indx) rec++;
            Console.WriteLine($"trnId: {trnId}, currRnd: {currRnd}, plyr1: {plyr1[indx]}, plyr2: {plyr2[indx]}, gmsId: {gmsId}");
            inserted &= insert_gms_plyrtm(gmsId, plyr1[indx], rec);
            inserted &= insert_gms_plyrtm(gmsId, plyr2[indx], rec);
        }

        return inserted;
    }

    private static bool insert_gms_plyrtm(int gmsId, int plyrtmId, int rec = 0) {
        using var conn = _db.Create();
        int result = conn.Execute(
            @"INSERT INTO gms_plyrtm (gms_id, plyrtm_id, rec) VALUES (@GmsId, @PlyrtmId, @Rec);",
            new { GmsId = gmsId, PlyrtmId = plyrtmId, Rec = rec });
        return result > 0;
    }

    public static (bool, List<int>) wc_validate(int trnId, int currRnd, List<int> plyr1, List<int> plyr2) {
        using var conn = _db.Create();

        List<int> plyrstm_wns = new List<int>();

        if (plyr1.Count != plyr2.Count) {
            return (false, plyrstm_wns);
        }

        for (int indx = 0; indx < plyr1.Count; indx++) {
            var rows = conn.Query(
                @"SELECT coalesce(plyrtm_wn, -1) plyrtm_wn FROM gms
                inner join gms_plyrtm gp1 on gp1.gms_id = gms.id and gp1.plyrtm_id = @PlyrId1
                inner join gms_plyrtm gp2 on gp2.gms_id = gms.id and gp2.plyrtm_id = @PlyrId2
                WHERE tourn_id = @TrnId AND rnd = @CurrRnd AND has_stnd = 'N'",
                new { TrnId = trnId, CurrRnd = currRnd, PlyrId1 = plyr1[indx], PlyrId2 = plyr2[indx] }).ToList();

            if (rows == null || rows.Count() == 0) {
                return (false, plyrstm_wns);
            }
            plyrstm_wns.Add(Convert.ToInt32(rows[0].plyrtm_wn));
        }

        return (true, plyrstm_wns);
    }
}